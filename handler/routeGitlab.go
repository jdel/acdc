package handler // import jdel.org/acdc/handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"jdel.org/acdc/api"
)

type gitlabPushPayload struct {
	ObjectKind   string      `json:"object_kind"`
	EventName    string      `json:"event_name"`
	Before       string      `json:"before"`
	After        string      `json:"after"`
	Ref          string      `json:"ref"`
	CheckoutSha  string      `json:"checkout_sha"`
	Message      interface{} `json:"message"`
	UserID       int         `json:"user_id"`
	UserName     string      `json:"user_name"`
	UserUsername string      `json:"user_username"`
	UserEmail    string      `json:"user_email"`
	UserAvatar   string      `json:"user_avatar"`
	ProjectID    int         `json:"project_id"`
	Project      struct {
		ID                int         `json:"id"`
		Name              string      `json:"name"`
		Description       string      `json:"description"`
		WebURL            string      `json:"web_url"`
		AvatarURL         interface{} `json:"avatar_url"`
		GitSSHURL         string      `json:"git_ssh_url"`
		GitHTTPURL        string      `json:"git_http_url"`
		Namespace         string      `json:"namespace"`
		VisibilityLevel   int         `json:"visibility_level"`
		PathWithNamespace string      `json:"path_with_namespace"`
		DefaultBranch     string      `json:"default_branch"`
		CiConfigPath      string      `json:"ci_config_path"`
		Homepage          string      `json:"homepage"`
		URL               string      `json:"url"`
		SSHURL            string      `json:"ssh_url"`
		HTTPURL           string      `json:"http_url"`
	} `json:"project"`
	Commits []struct {
		ID        string    `json:"id"`
		Message   string    `json:"message"`
		Timestamp time.Time `json:"timestamp"`
		URL       string    `json:"url"`
		Author    struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"author"`
		Added    []string      `json:"added"`
		Modified []interface{} `json:"modified"`
		Removed  []interface{} `json:"removed"`
	} `json:"commits"`
	TotalCommitsCount int `json:"total_commits_count"`
	Repository        struct {
		Name            string `json:"name"`
		URL             string `json:"url"`
		Description     string `json:"description"`
		Homepage        string `json:"homepage"`
		GitHTTPURL      string `json:"git_http_url"`
		GitSSHURL       string `json:"git_ssh_url"`
		VisibilityLevel int    `json:"visibility_level"`
	} `json:"repository"`
}

type gitlabPipelinePayload struct {
	ObjectKind       string `json:"object_kind"`
	ObjectAttributes struct {
		ID             int      `json:"id"`
		Ref            string   `json:"ref"`
		Tag            bool     `json:"tag"`
		Sha            string   `json:"sha"`
		BeforeSha      string   `json:"before_sha"`
		Status         string   `json:"status"`
		DetailedStatus string   `json:"detailed_status"`
		Stages         []string `json:"stages"`
		CreatedAt      string   `json:"created_at"`
		FinishedAt     string   `json:"finished_at"`
		Duration       int      `json:"duration"`
	} `json:"object_attributes"`
	User struct {
		Name      string `json:"name"`
		Username  string `json:"username"`
		AvatarURL string `json:"avatar_url"`
	} `json:"user"`
	Project struct {
		ID                int         `json:"id"`
		Name              string      `json:"name"`
		Description       string      `json:"description"`
		WebURL            string      `json:"web_url"`
		AvatarURL         interface{} `json:"avatar_url"`
		GitSSHURL         string      `json:"git_ssh_url"`
		GitHTTPURL        string      `json:"git_http_url"`
		Namespace         string      `json:"namespace"`
		VisibilityLevel   int         `json:"visibility_level"`
		PathWithNamespace string      `json:"path_with_namespace"`
		DefaultBranch     string      `json:"default_branch"`
		CiConfigPath      string      `json:"ci_config_path"`
	} `json:"project"`
	Commit struct {
		ID        string    `json:"id"`
		Message   string    `json:"message"`
		Timestamp time.Time `json:"timestamp"`
		URL       string    `json:"url"`
		Author    struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"author"`
	} `json:"commit"`
	Builds []struct {
		ID         int         `json:"id"`
		Stage      string      `json:"stage"`
		Name       string      `json:"name"`
		Status     string      `json:"status"`
		CreatedAt  string      `json:"created_at"`
		StartedAt  interface{} `json:"started_at"`
		FinishedAt interface{} `json:"finished_at"`
		When       string      `json:"when"`
		Manual     bool        `json:"manual"`
		User       struct {
			Name      string `json:"name"`
			Username  string `json:"username"`
			AvatarURL string `json:"avatar_url"`
		} `json:"user"`
		Runner        interface{} `json:"runner"`
		ArtifactsFile struct {
			Filename interface{} `json:"filename"`
			Size     int         `json:"size"`
		} `json:"artifacts_file"`
	} `json:"builds"`
}

type gitlabCallbackPayload struct {
	Message []string `json:"message"`
	Context string   `json:"context,omitempty"`
}

func outputGitlabPayload(message, context string) gitlabCallbackPayload {
	return gitlabCallbackPayload{
		Message: strings.Split(message, "\n"),
		Context: context,
	}
}

func findGitlabMatchingKey(key string) *api.Key {
	return api.FindKey(key)
}

// RouteGitlab handles incoming Gitlab pushes
func RouteGitlab(w http.ResponseWriter, r *http.Request) {
	key := findGitlabMatchingKey(r.Header.Get("X-Gitlab-Token"))

	if key == nil {
		jsonOutput(w, http.StatusNotFound,
			outputGitlabPayload("Could not find a matching key", ""))
		return
	}

	if r.Header.Get("X-Gitlab-Event") != r.URL.Query().Get("event") {
		jsonOutput(w, http.StatusMethodNotAllowed,
			outputGithubPayload("wrong event", "error"))
		return
	}

	switch event := r.Header.Get("X-Gitlab-Event"); event {
	case "Pipeline Hook":
		var incomingPayload gitlabPipelinePayload
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		decoder.Decode(&incomingPayload)
		if incomingPayload.ObjectAttributes.Status != r.URL.Query().Get("status") {
			jsonOutput(w, http.StatusMethodNotAllowed,
				outputGithubPayload("wrong status", "error"))
			return
		}
	}

	hookName := r.URL.Query().Get("hook")

	hook := key.GetHook(hookName)
	if hook == nil {
		logRoute.Error("Cannot find hook ", hookName)
		jsonOutput(w, http.StatusInternalServerError,
			outputGitlabPayload("Could not find hook", hookName))
		return
	}

	if key.IsRemote() {
		key.Pull()
	}

	actions := strings.Split(r.URL.Query().Get("actions"), " ")
	tickets, _ := hook.ExecuteSequentially(actions...)

	jsonOutput(w, http.StatusOK,
		outputGithubPayload(tickets[:len(tickets)-1], "queued"))
}
