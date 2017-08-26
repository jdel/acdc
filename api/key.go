package api

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/jdel/acdc/cfg"
	"github.com/jdel/acdc/util"
	log "github.com/sirupsen/logrus"
	git "gopkg.in/src-d/go-git.v4"
)

// Key represents an API key
type Key struct {
	Unique string
	Remote string
}

// NewKey creates a new key with a known unique or generates a random one
// iif unique is "", repote param can be "" for a local key
func NewKey(unique, remote string) (Key, error) {
	var err error
	if unique == "" {
		unique, err = util.GenerateRandomString(12)
		if err != nil {
			logAPI.Error("Could not generate an API Key:", err)
			return Key{}, err
		}
	}
	key := newKey(unique, remote)
	if err != nil {
		logAPI.Error("Error creating new API Key:", err)
	}
	return key, err
}

// Delete deletes the specified key
func (key Key) Delete() error {
	var err error
	var keyPath = filepath.Join(cfg.GetComposeDir(), key.Unique)
	if err = os.RemoveAll(keyPath); err != nil {
		logAPI.
			WithField("key", key.Unique).
			WithField("err", err).
			Error("Couldn't delete path ", keyPath)
		return err
	}
	return err
}

// newKey creates and returns a new key
func newKey(unique, remote string) Key {
	apiKey := Key{
		Unique: unique,
		Remote: remote,
	}
	apiKey.create()
	return apiKey
}

// FindKey returns a Key from the unique string
func FindKey(unique string) Key {
	apiKey := Key{
		Unique: unique,
		Remote: "",
	}

	if !apiKey.Exists() {
		return Key{"", ""}
	}

	remote, err := apiKey.getRemote()
	if err == nil {
		origin, err := remote.Remote("origin")
		if err == nil {
			apiKey.Remote = origin.Config().URL
		}
	}

	return apiKey
}

// Create creates the local API Key
// (creates a directory under compose-dir/key)
func (key Key) create() {
	var dir string
	var err error

	dir, err = util.CreateDir(filepath.Join(cfg.GetComposeDir(), key.Unique))
	if err != nil {
		logAPI.WithFields(log.Fields{
			"key": key.Unique,
		}).Error("Could not create api-key directory ", dir, ": ", err)
	}

	if key.Remote != "" {
		if _, err = git.PlainClone(filepath.Join(cfg.GetComposeDir(), key.Unique), false, &git.CloneOptions{
			URL: key.Remote,
		}); err != nil {
			logAPI.WithFields(log.Fields{
				"key": key.Unique,
			}).Error("Could not git clone ", key.Remote, ": ", err)
			key.Delete()
		}
	}
}

// Pull does a git pull of the remote associated with the key
func (key Key) Pull() (string, error) {
	var headShortHash string

	repo, err := key.getRemote()
	if err != nil {
		logAPI.WithFields(log.Fields{
			"key": key.Unique,
		}).Error("Not a remote key: ", err)
		return "", err
	}

	worktree, err := repo.Worktree()
	if err != nil {
		logAPI.WithFields(log.Fields{
			"key": key.Unique,
		}).Error("Could not get worktree: ", err)
	}

	err = worktree.Pull(&git.PullOptions{})
	if err != nil {
		switch err {
		case git.NoErrAlreadyUpToDate:
			logAPI.WithFields(log.Fields{
				"key": key.Unique,
			}).Info(err)
			err = nil
		default:
			logAPI.WithFields(log.Fields{
				"key": key.Unique,
			}).Error("Could not git pull: ", err)

		}
	}

	head, _ := repo.Head()

	headShortHash = head.Hash().String()[:6]

	return headShortHash, err
}

// IsRemote returns true if the key is linked to git repo
func (key Key) getRemote() (*git.Repository, error) {
	var repoPath = filepath.Join(cfg.GetComposeDir(), key.Unique)
	return git.PlainOpen(repoPath)
}

// Exists returns true if the key is linked to git repo
func (key Key) Exists() bool {
	return util.FileExists(filepath.Join(cfg.GetComposeDir(), key.Unique))
}

// IsRemote returns true if the key is linked to git repo
func (key Key) IsRemote() bool {
	if key.Remote != "" {
		return true
	}
	return false
}

// AllHooks returns all hooks (.yml files) associated with the key
func (key Key) AllHooks() []Hook {
	var hooks []Hook
	var err error
	var children []os.FileInfo

	children, err = ioutil.ReadDir(filepath.Join(cfg.GetComposeDir(), key.Unique))
	if err != nil {
		return nil
	}

	for _, child := range children {
		childName := child.Name()
		// TODO: Handle yaml files too !
		if childNameNoExt := strings.TrimSuffix(childName, ".yml"); child.Mode().IsRegular() && childNameNoExt != childName {
			hooks = append(hooks, key.GetHook(childNameNoExt))
		} else {
			logAPI.Debug("Ignoring non-yml file ", childName)
		}

	}
	return hooks
}

// AllAPIKeys lists all API Keys
func AllAPIKeys() (map[string]Key, error) {
	var keys = make(map[string]Key)
	var err error
	var children []os.FileInfo

	children, err = ioutil.ReadDir(cfg.GetComposeDir())
	if err != nil {
		return nil, err
	}

	for _, child := range children {
		childName := child.Name()
		if child.Mode().IsDir() {
			keys[childName] = FindKey(childName)
		} else {
			logAPI.Debug("Ignoring regular file ", childName)
		}

	}
	return keys, err
}
