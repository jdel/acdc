package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

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
	var output []byte
	var incomingPayload dockerHubPayload
	hookName := r.URL.Query().Get("hook")
	apiKey := r.URL.Query().Get("apiKey")
	tag := r.URL.Query().Get("tag")

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

	if tag != incomingPayload.PushData.Tag {
		jsonOutput(w, http.StatusNotFound,
			fmt.Sprintln("Ignoring tag", incomingPayload.PushData.Tag, "does not match", tag))
		return
	}

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

	output, err = hook.Pull().CombinedOutput()
	if err != nil {
		logRoute.WithField("route", "RouteDockerHub").Error(string(output), err)
		jsonOutput(w, http.StatusNotFound,
			outputDockerHubPayload("Could not pull images for hook", hook.Name))
		return
	}

	output, err = hook.Up().CombinedOutput()
	if err != nil {
		logRoute.WithField("route", "RouteDockerHub").Error(string(output), err)
		jsonOutput(w, http.StatusInternalServerError,
			outputDockerHubPayload("Could not bring hook up", hook.Name))
		return
	}

	jsonOutput(w, http.StatusOK,
		outputDockerHubPayload("Upgraded hook", hook.Name))
}
