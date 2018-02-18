package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/jdel/acdc/api"
)

type dockerHubPayload struct {
	PushData struct {
		PushedAt int           `json:"pushed_at"`
		Images   []interface{} `json:"images"`
		Tag      string        `json:"tag"`
		Pusher   string        `json:"pusher"`
	} `json:"push_data"`
	CallbackURL string `json:"callback_url"`
	Repository  struct {
		Status          string `json:"status"`
		Description     string `json:"description"`
		IsTrusted       bool   `json:"is_trusted"`
		FullDescription string `json:"full_description"`
		RepoURL         string `json:"repo_url"`
		Owner           string `json:"owner"`
		IsOfficial      bool   `json:"is_official"`
		IsPrivate       bool   `json:"is_private"`
		Name            string `json:"name"`
		Namespace       string `json:"namespace"`
		StarCount       int    `json:"star_count"`
		CommentCount    int    `json:"comment_count"`
		DateCreated     int    `json:"date_created"`
		Dockerfile      string `json:"dockerfile"`
		RepoName        string `json:"repo_name"`
	} `json:"repository"`
}

type dockerHubCallbackPayload struct {
	Message string `json:"message"`
	Context string `json:"context,omitempty"`
}

func outputDockerHubPayload(message, context string) dockerHubCallbackPayload {
	return dockerHubCallbackPayload{
		Message: message,
		Context: context,
	}
}

// RouteDockerHub handles incoming docker hub hooks
func RouteDockerHub(w http.ResponseWriter, r *http.Request) {
	var incomingPayload dockerHubPayload
	apiKey := r.URL.Query().Get("apiKey")
	key := api.FindKey(apiKey)
	if key == nil {
		jsonOutput(w, http.StatusUnauthorized,
			outputDockerHubPayload("Could not find key", apiKey))
		return
	}

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&incomingPayload)
	if err != nil {
		logRoute.WithField("route", "RouteDockerHub").Error(err)
	}
	defer r.Body.Close()

	tag := r.URL.Query().Get("tag")
	if tag != incomingPayload.PushData.Tag {
		jsonOutput(w, http.StatusNotFound,
			fmt.Sprintln("Ignoring tag", incomingPayload.PushData.Tag, "does not match", tag))
		return
	}

	hookName := r.URL.Query().Get("hook")
	if hookName == "" {
		hookName = incomingPayload.Repository.Name
	}

	hook := key.GetHook(hookName)
	if hook == nil {
		logRoute.Error("Cannot find hook", hookName)
		jsonOutput(w, http.StatusInternalServerError,
			outputGogsPayload("Could not find hook", hookName))
		return
	}

	actions := strings.Split(r.URL.Query().Get("actions"), " ")
	ticket, _ := hook.ExecuteSequentially(actions...)

	jsonOutput(w, http.StatusOK,
		outputGithubPayload("queued", ticket))
}
