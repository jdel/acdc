package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/jdel/acdc/api"
)

/*
 * token=GSukJsasLYROypkh1nowEHrX
 * team_id=T0001
 * team_domain=example
 * channel_id=C2147483705
 * channel_name=test
 * timestamp=1355517523.000005
 * user_id=U2147483697
 * user_name=Steve
 * text=googlebot: What is the air-speed velocity of an unladen swallow?
 * trigger_word=googlebot:
 */

// RouteSlack handles incoming slack hooks
func RouteSlack(w http.ResponseWriter, r *http.Request) {
	var output []byte
	var err error
	apiKey := r.FormValue("token")
	args := slackParseArgs(r.FormValue("text"))
	hookName := args[1]

	if len(args) < 3 {
		jsonOutput(w, http.StatusOK,
			slackCallbackPayload("This Slack hook requires 2 arguments", "acdc"))
		return
	}

	key := api.FindKey(apiKey)
	if key.Unique == "" {
		jsonOutput(w, http.StatusOK,
			slackCallbackPayload("Could not get key", "acdc"))
		return
	}

	hook := key.GetHook(hookName)
	if hook.Name == "" {
		jsonOutput(w, http.StatusOK,
			slackCallbackPayload("Could not get hook", hookName))
		return
	}

	switch args[2] {
	case "pull":
		output, err = hook.Pull().CombinedOutput()
		if err != nil {
			jsonOutput(w, http.StatusOK,
				slackCallbackPayload("Could not pull images for hook", hook.Name))
			return
		}
	case "ps":
		output, err = hook.Ps().CombinedOutput()
		if err != nil {
			jsonOutput(w, http.StatusOK,
				slackCallbackPayload("Could not get hook", hook.Name))
			return
		}
	case "logs":
		output, err = hook.Logs().CombinedOutput()
		if err != nil {
			jsonOutput(w, http.StatusOK,
				slackCallbackPayload("Could not get hook logs", hook.Name))
			return
		}
	case "up":
		output, err = hook.Up().CombinedOutput()
		if err != nil {
			logRoute.Error(hook.Up(), err)
			jsonOutput(w, http.StatusOK,
				slackCallbackPayload("Could not bring hook up", hook.Name))
			return
		}
	case "upgrade":
		output, err = hook.Pull().CombinedOutput()
		if err != nil {
			jsonOutput(w, http.StatusOK,
				slackCallbackPayload("Could not pull images for hook", hook.Name))
			return
		}
		output, err = hook.Down().CombinedOutput()
		if err != nil {
			logRoute.WithField("key", apiKey).Error(err)
			jsonOutput(w, http.StatusOK,
				slackCallbackPayload("Could not bring hook down", hook.Name))
			return
		}
		output, err = hook.Up().CombinedOutput()
		if err != nil {
			logRoute.Error(hook.Up(), err)
			jsonOutput(w, http.StatusOK,
				slackCallbackPayload("Could not bring hook up", hook.Name))
			return
		}
	case "down":
		output, err = hook.Down().CombinedOutput()
		if err != nil {
			logRoute.WithField("key", apiKey).Error(err)
			jsonOutput(w, http.StatusOK,
				slackCallbackPayload("Could not bring hook down", hook.Name))
			return
		}
	case "start":
		output, err = hook.Start().CombinedOutput()
		if err != nil {
			logRoute.WithField("key", apiKey).Error(err)
			jsonOutput(w, http.StatusOK,
				slackCallbackPayload("Could not start hook", hook.Name))
			return
		}
	case "stop":
		output, err = hook.Stop().CombinedOutput()
		if err != nil {
			logRoute.WithField("key", apiKey).Error(err)
			jsonOutput(w, http.StatusOK,
				slackCallbackPayload("Could not stop hook", hook.Name))
			return
		}
	case "restart":
		output, err = hook.Restart().CombinedOutput()
		if err != nil {
			logRoute.WithField("key", apiKey).Error(err)
			jsonOutput(w, http.StatusOK,
				slackCallbackPayload("Could not restart hook", hook.Name))
			return
		}
	default:
		jsonOutput(w, http.StatusOK,
			slackCallbackPayload("Valid actions are: ps, pull, upgrade (pull+down+up), up, down, start, stop, restart, logs", hook.Name))
		return
	}

	anonymizedOutput := fmt.Sprintf(
		"```%s```",
		strings.Replace(string(output), strings.ToLower(key.Unique), "", -1),
	)

	jsonOutput(w, http.StatusOK,
		slackCallbackPayload(anonymizedOutput, hook.Name))
}

func slackParseArgs(text string) []string {
	return strings.Split(text, " ")
}

func slackCallbackPayload(text, username string) interface{} {
	return struct {
		Text     string `json:"text"`
		Username string `json:"username,omitempty"`
	}{
		Text:     text,
		Username: username,
	}
}
