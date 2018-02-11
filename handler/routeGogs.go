package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/jdel/acdc/util"

	"github.com/jdel/acdc/api"
)

type gogsPayload struct {
	Ref        string `json:"ref"`
	Before     string `json:"before"`
	After      string `json:"after"`
	CompareURL string `json:"compare_url"`
	Commits    []struct {
		ID      string `json:"id"`
		Message string `json:"message"`
		URL     string `json:"url"`
		Author  struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Username string `json:"username"`
		} `json:"author"`
		Committer struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Username string `json:"username"`
		} `json:"committer"`
		Added     []string      `json:"added"`
		Removed   []interface{} `json:"removed"`
		Modified  []interface{} `json:"modified"`
		Timestamp time.Time     `json:"timestamp"`
	} `json:"commits"`
	Repository struct {
		ID    int `json:"id"`
		Owner struct {
			ID        int    `json:"id"`
			Login     string `json:"login"`
			FullName  string `json:"full_name"`
			Email     string `json:"email"`
			AvatarURL string `json:"avatar_url"`
			Username  string `json:"username"`
		} `json:"owner"`
		Name            string      `json:"name"`
		FullName        string      `json:"full_name"`
		Description     string      `json:"description"`
		Private         bool        `json:"private"`
		Fork            bool        `json:"fork"`
		Parent          interface{} `json:"parent"`
		Empty           bool        `json:"empty"`
		Mirror          bool        `json:"mirror"`
		Size            int         `json:"size"`
		HTMLURL         string      `json:"html_url"`
		SSHURL          string      `json:"ssh_url"`
		CloneURL        string      `json:"clone_url"`
		Website         string      `json:"website"`
		StarsCount      int         `json:"stars_count"`
		ForksCount      int         `json:"forks_count"`
		WatchersCount   int         `json:"watchers_count"`
		OpenIssuesCount int         `json:"open_issues_count"`
		DefaultBranch   string      `json:"default_branch"`
		CreatedAt       time.Time   `json:"created_at"`
		UpdatedAt       time.Time   `json:"updated_at"`
	} `json:"repository"`
	Pusher struct {
		ID        int    `json:"id"`
		Login     string `json:"login"`
		FullName  string `json:"full_name"`
		Email     string `json:"email"`
		AvatarURL string `json:"avatar_url"`
		Username  string `json:"username"`
	} `json:"pusher"`
	Sender struct {
		ID        int    `json:"id"`
		Login     string `json:"login"`
		FullName  string `json:"full_name"`
		Email     string `json:"email"`
		AvatarURL string `json:"avatar_url"`
		Username  string `json:"username"`
	} `json:"sender"`
}

type gogsCallbackPayload struct {
	Message string `json:"message"`
	Context string `json:"context,omitempty"`
}

func outputGogsPayload(message, context string) gogsCallbackPayload {
	return gogsCallbackPayload{
		Message: message,
		Context: context,
	}
}

func checkHexHMacSignature(secret, message, expectedSum []byte) bool {
	hash := hmac.New(sha256.New, secret)
	hash.Write(message)
	if hex.EncodeToString(hash.Sum(nil)) != string(expectedSum) {
		return false
	}
	return true
}

func findMatchingKey(expectedSum, message []byte) api.Key {
	keys, _ := api.AllAPIKeys()
	for _, key := range keys {
		if checkHexHMacSignature([]byte(key.Unique), message, expectedSum) {
			return api.FindKey(key.Unique)
		}
	}
	return api.Key{}
}

// RouteGogs handles incoming gogs pushes
func RouteGogs(w http.ResponseWriter, r *http.Request) {
	var incomingPayload gogsPayload
	var output []byte

	if r.Header.Get("X-Gogs-Event") != "push" {
		return
	}

	actions := strings.Split(r.URL.Query().Get("actions"), " ")
	logRoute.Debugf("Actions: %+v", actions)
	body, _ := ioutil.ReadAll(r.Body)
	key := findMatchingKey([]byte(r.Header.Get("X-Gogs-Signature")), body)

	if key.Unique == "" {
		jsonOutput(w, http.StatusNotFound,
			outputGogsPayload("Could not get a matching key", ""))
		return
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&incomingPayload)
	if err != nil {
		logRoute.Error(err)
	}
	defer r.Body.Close()

	hookName := r.URL.Query().Get("hook")

	hook := key.GetHook(hookName)
	if hook.Name == "" {
		logRoute.Error("Cannot find hook", err)
		jsonOutput(w, http.StatusInternalServerError,
			outputGogsPayload("Could not pull images for hook", hookName))
		return
	}

	if key.IsRemote() {
		key.Pull()
	}

	if util.IsStringInSlice("pull", actions) {
		logRoute.Debugf("Pulling %s", hook.Name)
		output, err = hook.Pull().CombinedOutput()
		if err != nil {
			logRoute.Error(string(output), err)
			jsonOutput(w, http.StatusInternalServerError,
				outputGogsPayload("Could not pull images for hook", hook.Name))
			return
		}
	}

	if util.IsStringInSlice("down", actions) {
		logRoute.Debugf("Bringing %s down", hook.Name)
		output, err = hook.Down().CombinedOutput()
		if err != nil {
			logRoute.Error(string(output), err)
			jsonOutput(w, http.StatusInternalServerError,
				outputGogsPayload("Could not bring hook down", hook.Name))
			return
		}
	}

	if util.IsStringInSlice("build", actions) {
		logRoute.Debugf("Building %s", hook.Name)
		output, err = hook.Build().CombinedOutput()
		if err != nil {
			logRoute.Error(string(output), err)
			jsonOutput(w, http.StatusInternalServerError,
				outputGogsPayload("Could not build hook", hook.Name))
			return
		}
	}

	if util.IsStringInSlice("up", actions) {
		logRoute.Debugf("Bringing %s up", hook.Name)
		output, err = hook.Up().CombinedOutput()
		if err != nil {
			logRoute.Error(string(output), err)
			jsonOutput(w, http.StatusInternalServerError,
				outputGogsPayload("Could not bring hook up", hook.Name))
			return
		}
	}

	jsonOutput(w, http.StatusOK,
		outputGogsPayload("ok", ""))
}
