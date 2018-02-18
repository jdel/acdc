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
	apiKey := r.FormValue("token")
	args := slackParseArgs(r.FormValue("text"))
	hookName := args[1]

	if len(args) < 3 {
		jsonOutput(w, http.StatusOK,
			slackCallbackPayload("This Slack hook requires 2 arguments", "acdc"))
		return
	}

	key := api.FindKey(apiKey)
	if key == nil {
		jsonOutput(w, http.StatusOK,
			slackCallbackPayload("Could not find key", "acdc"))
		return
	}

	hook := key.GetHook(hookName)
	if hook == nil {
		jsonOutput(w, http.StatusOK,
			slackCallbackPayload("Could not find hook", hookName))
		return
	}

	output, err := hook.ExecuteSequentially(strings.Split(args[2], "+")...)
	if err != nil {
		jsonOutput(w, http.StatusInternalServerError,
			slackCallbackPayload(err.Error(), hook.Name))
		return
	}

	anonymizedOutput := fmt.Sprintf(
		"```%s```",
		strings.Replace(output, strings.ToLower(key.Unique), "", -1),
	)

	jsonOutput(w, http.StatusOK,
		slackCallbackPayload(anonymizedOutput, hook.Name))
	return
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
