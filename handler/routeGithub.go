package handler

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/jdel/acdc/api"
)

type githubPayload struct {
	Ref     string      `json:"ref"`
	Before  string      `json:"before"`
	After   string      `json:"after"`
	Created bool        `json:"created"`
	Deleted bool        `json:"deleted"`
	Forced  bool        `json:"forced"`
	BaseRef interface{} `json:"base_ref"`
	Compare string      `json:"compare"`
	Commits []struct {
		ID        string `json:"id"`
		TreeID    string `json:"tree_id"`
		Distinct  bool   `json:"distinct"`
		Message   string `json:"message"`
		Timestamp string `json:"timestamp"`
		URL       string `json:"url"`
		Author    struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Username string `json:"username"`
		} `json:"author"`
		Committer struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Username string `json:"username"`
		} `json:"committer"`
		Added    []interface{} `json:"added"`
		Removed  []interface{} `json:"removed"`
		Modified []string      `json:"modified"`
	} `json:"commits"`
	HeadCommit struct {
		ID        string `json:"id"`
		TreeID    string `json:"tree_id"`
		Distinct  bool   `json:"distinct"`
		Message   string `json:"message"`
		Timestamp string `json:"timestamp"`
		URL       string `json:"url"`
		Author    struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Username string `json:"username"`
		} `json:"author"`
		Committer struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Username string `json:"username"`
		} `json:"committer"`
		Added    []interface{} `json:"added"`
		Removed  []interface{} `json:"removed"`
		Modified []string      `json:"modified"`
	} `json:"head_commit"`
	Repository struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		FullName string `json:"full_name"`
		Owner    struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"owner"`
		Private          bool        `json:"private"`
		HTMLURL          string      `json:"html_url"`
		Description      string      `json:"description"`
		Fork             bool        `json:"fork"`
		URL              string      `json:"url"`
		ForksURL         string      `json:"forks_url"`
		KeysURL          string      `json:"keys_url"`
		CollaboratorsURL string      `json:"collaborators_url"`
		TeamsURL         string      `json:"teams_url"`
		HooksURL         string      `json:"hooks_url"`
		IssueEventsURL   string      `json:"issue_events_url"`
		EventsURL        string      `json:"events_url"`
		AssigneesURL     string      `json:"assignees_url"`
		BranchesURL      string      `json:"branches_url"`
		TagsURL          string      `json:"tags_url"`
		BlobsURL         string      `json:"blobs_url"`
		GitTagsURL       string      `json:"git_tags_url"`
		GitRefsURL       string      `json:"git_refs_url"`
		TreesURL         string      `json:"trees_url"`
		StatusesURL      string      `json:"statuses_url"`
		LanguagesURL     string      `json:"languages_url"`
		StargazersURL    string      `json:"stargazers_url"`
		ContributorsURL  string      `json:"contributors_url"`
		SubscribersURL   string      `json:"subscribers_url"`
		SubscriptionURL  string      `json:"subscription_url"`
		CommitsURL       string      `json:"commits_url"`
		GitCommitsURL    string      `json:"git_commits_url"`
		CommentsURL      string      `json:"comments_url"`
		IssueCommentURL  string      `json:"issue_comment_url"`
		ContentsURL      string      `json:"contents_url"`
		CompareURL       string      `json:"compare_url"`
		MergesURL        string      `json:"merges_url"`
		ArchiveURL       string      `json:"archive_url"`
		DownloadsURL     string      `json:"downloads_url"`
		IssuesURL        string      `json:"issues_url"`
		PullsURL         string      `json:"pulls_url"`
		MilestonesURL    string      `json:"milestones_url"`
		NotificationsURL string      `json:"notifications_url"`
		LabelsURL        string      `json:"labels_url"`
		ReleasesURL      string      `json:"releases_url"`
		CreatedAt        int         `json:"created_at"`
		UpdatedAt        time.Time   `json:"updated_at"`
		PushedAt         int         `json:"pushed_at"`
		GitURL           string      `json:"git_url"`
		SSHURL           string      `json:"ssh_url"`
		CloneURL         string      `json:"clone_url"`
		SvnURL           string      `json:"svn_url"`
		Homepage         interface{} `json:"homepage"`
		Size             int         `json:"size"`
		StargazersCount  int         `json:"stargazers_count"`
		WatchersCount    int         `json:"watchers_count"`
		Language         interface{} `json:"language"`
		HasIssues        bool        `json:"has_issues"`
		HasDownloads     bool        `json:"has_downloads"`
		HasWiki          bool        `json:"has_wiki"`
		HasPages         bool        `json:"has_pages"`
		ForksCount       int         `json:"forks_count"`
		MirrorURL        interface{} `json:"mirror_url"`
		OpenIssuesCount  int         `json:"open_issues_count"`
		Forks            int         `json:"forks"`
		OpenIssues       int         `json:"open_issues"`
		Watchers         int         `json:"watchers"`
		DefaultBranch    string      `json:"default_branch"`
		Stargazers       int         `json:"stargazers"`
		MasterBranch     string      `json:"master_branch"`
	} `json:"repository"`
	Pusher struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"pusher"`
	Sender struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"sender"`
}

type githubCallbackPayload struct {
	Message []string `json:"message"`
	Context string   `json:"context,omitempty"`
}

func outputGithubPayload(message, context string) githubCallbackPayload {
	return githubCallbackPayload{
		Message: strings.Split(message, "\n"),
		Context: context,
	}
}

func checkHexHMacSha1Signature(secret, message, expectedSum []byte) bool {
	hash := hmac.New(sha1.New, secret)
	hash.Write(message)
	return subtle.ConstantTimeCompare([]byte(hex.EncodeToString(hash.Sum(nil))), expectedSum) == 1
}

func findGithubMatchingKey(expectedSum, message []byte) *api.Key {
	keys, _ := api.AllAPIKeys()
	for _, key := range keys {
		if checkHexHMacSha1Signature([]byte(key.Unique), message, expectedSum) {
			return api.FindKey(key.Unique)
		}
	}
	return nil
}

// RouteGithub handles incoming github pushes
func RouteGithub(w http.ResponseWriter, r *http.Request) {
	var incomingPayload githubPayload
	if r.Header.Get("X-Github-Event") != "push" {
		return
	}

	body, _ := ioutil.ReadAll(r.Body)
	key := findGithubMatchingKey([]byte(r.Header.Get("X-Hub-Signature")[5:]), body)
	if key.Unique == "" {
		jsonOutput(w, http.StatusNotFound,
			outputGithubPayload("Could not find a matching key", ""))
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
			outputGithubPayload("Could not find hook", hookName))
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
