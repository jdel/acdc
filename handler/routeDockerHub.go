package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
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

// RouteDockerHub handles incoming docker hub hooks
func RouteDockerHub(w http.ResponseWriter, r *http.Request) {
	var output []byte
	var incomingPayload dockerHubPayload
	apiKey := mux.Vars(r)["apiKey"]
	tag := mux.Vars(r)["tag"]

	key := api.FindKey(apiKey)
	if key.Unique == "" {
		jsonOutput(w, http.StatusInternalServerError,
			outputKey("Could not get key", apiKey))
		return
	}

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&incomingPayload)
	if err != nil {
		logRoute.Error(err)
	}
	defer r.Body.Close()

	if tag != incomingPayload.PushData.Tag {
		jsonOutput(w, http.StatusInternalServerError,
			fmt.Sprintln("Ignoring tag", incomingPayload.PushData.Tag, "does not match", tag))
		return
	}

	hook := key.GetHook(incomingPayload.Repository.Name)

	output, err = hook.Pull().CombinedOutput()
	if err != nil {
		logRoute.Error(string(output), err)
		jsonOutput(w, http.StatusInternalServerError,
			outputHook("Could not pull images for hook", hook.Name))
		return
	}

	output, err = hook.Down().CombinedOutput()
	if err != nil {
		logRoute.Error(string(output), err)
		jsonOutput(w, http.StatusInternalServerError,
			outputHook("Could not bring hook down", hook.Name))
		return
	}

	output, err = hook.Up().CombinedOutput()
	if err != nil {
		logRoute.Error(string(output), err)
		jsonOutput(w, http.StatusInternalServerError,
			outputHook("Could not bring hook up", hook.Name))
		return
	}

	jsonOutput(w, http.StatusOK,
		outputHook("Upgraded hook", hook.Name))
}
