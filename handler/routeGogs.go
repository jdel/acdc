package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

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

func checkHexHMacSha256Signature(secret, message, expectedSum []byte) bool {
	hash := hmac.New(sha256.New, secret)
	hash.Write(message)
	return subtle.ConstantTimeCompare([]byte(hex.EncodeToString(hash.Sum(nil))), expectedSum) == 1
}

func findGogsMatchingKey(expectedSum, message []byte) *api.Key {
	keys, _ := api.AllAPIKeys()
	for _, key := range keys {
		if checkHexHMacSha256Signature([]byte(key.Unique), message, expectedSum) {
			return api.FindKey(key.Unique)
		}
	}
	return nil
}

// RouteGogs handles incoming gogs pushes
func RouteGogs(w http.ResponseWriter, r *http.Request) {
	var incomingPayload gogsPayload
	if r.Header.Get("X-Gogs-Event") != "push" {
		return
	}

	body, _ := ioutil.ReadAll(r.Body)
	key := findGogsMatchingKey([]byte(r.Header.Get("X-Gogs-Signature")), body)

	if key.Unique == "" {
		jsonOutput(w, http.StatusNotFound,
			outputGogsPayload("Could not find a matching key", ""))
		return
	}

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	decoder.Decode(&incomingPayload)

	hookName := r.URL.Query().Get("hook")

	hook := key.GetHook(hookName)
	if hook == nil {
		logRoute.Error("Cannot find hook", hookName)
		jsonOutput(w, http.StatusInternalServerError,
			outputGogsPayload("Could not find hook", hookName))
		return
	}

	if key.IsRemote() {
		key.Pull()
	}

	actions := strings.Split(r.URL.Query().Get("actions"), " ")
	hook.ExecuteSequentially(actions...)

	jsonOutput(w, http.StatusOK,
		outputGogsPayload("ok", ""))
}
