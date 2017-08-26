package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"hash"
	"io/ioutil"
	"net/http"
	"time"
)

/*
{
   "zen":"Avoid administrative distraction.",
   "hook_id":15681425,
   "hook":{
      "type":"Repository",
      "id":15681425,
      "name":"web",
      "active":true,
      "events":[
         "push"
      ],
      "config":{
         "content_type":"json",
         "insecure_ssl":"0",
         "secret":"********",
         "url":"https://requestb.in/q0n3guq0"
      },
      "updated_at":"2017-08-22T21:43:59Z",
      "created_at":"2017-08-22T21:43:59Z",
      "url":"https://api.github.com/repos/jdel/acdc/hooks/15681425",
      "test_url":"https://api.github.com/repos/jdel/acdc/hooks/15681425/test",
      "ping_url":"https://api.github.com/repos/jdel/acdc/hooks/15681425/pings",
      "last_response":{
         "code":null,
         "status":"unused",
         "message":null
      }
   },
   "repository":{
      "id":100611332,
      "name":"acdc",
      "full_name":"jdel/acdc",
      "owner":{
         "login":"jdel",
         "id":1107511,
         "avatar_url":"https://avatars2.githubusercontent.com/u/1107511?v=4",
         "gravatar_id":"",
         "url":"https://api.github.com/users/jdel",
         "html_url":"https://github.com/jdel",
         "followers_url":"https://api.github.com/users/jdel/followers",
         "following_url":"https://api.github.com/users/jdel/following{/other_user}",
         "gists_url":"https://api.github.com/users/jdel/gists{/gist_id}",
         "starred_url":"https://api.github.com/users/jdel/starred{/owner}{/repo}",
         "subscriptions_url":"https://api.github.com/users/jdel/subscriptions",
         "organizations_url":"https://api.github.com/users/jdel/orgs",
         "repos_url":"https://api.github.com/users/jdel/repos",
         "events_url":"https://api.github.com/users/jdel/events{/privacy}",
         "received_events_url":"https://api.github.com/users/jdel/received_events",
         "type":"User",
         "site_admin":false
      },
      "private":false,
      "html_url":"https://github.com/jdel/acdc",
      "description":null,
      "fork":false,
      "url":"https://api.github.com/repos/jdel/acdc",
      "forks_url":"https://api.github.com/repos/jdel/acdc/forks",
      "keys_url":"https://api.github.com/repos/jdel/acdc/keys{/key_id}",
      "collaborators_url":"https://api.github.com/repos/jdel/acdc/collaborators{/collaborator}",
      "teams_url":"https://api.github.com/repos/jdel/acdc/teams",
      "hooks_url":"https://api.github.com/repos/jdel/acdc/hooks",
      "issue_events_url":"https://api.github.com/repos/jdel/acdc/issues/events{/number}",
      "events_url":"https://api.github.com/repos/jdel/acdc/events",
      "assignees_url":"https://api.github.com/repos/jdel/acdc/assignees{/user}",
      "branches_url":"https://api.github.com/repos/jdel/acdc/branches{/branch}",
      "tags_url":"https://api.github.com/repos/jdel/acdc/tags",
      "blobs_url":"https://api.github.com/repos/jdel/acdc/git/blobs{/sha}",
      "git_tags_url":"https://api.github.com/repos/jdel/acdc/git/tags{/sha}",
      "git_refs_url":"https://api.github.com/repos/jdel/acdc/git/refs{/sha}",
      "trees_url":"https://api.github.com/repos/jdel/acdc/git/trees{/sha}",
      "statuses_url":"https://api.github.com/repos/jdel/acdc/statuses/{sha}",
      "languages_url":"https://api.github.com/repos/jdel/acdc/languages",
      "stargazers_url":"https://api.github.com/repos/jdel/acdc/stargazers",
      "contributors_url":"https://api.github.com/repos/jdel/acdc/contributors",
      "subscribers_url":"https://api.github.com/repos/jdel/acdc/subscribers",
      "subscription_url":"https://api.github.com/repos/jdel/acdc/subscription",
      "commits_url":"https://api.github.com/repos/jdel/acdc/commits{/sha}",
      "git_commits_url":"https://api.github.com/repos/jdel/acdc/git/commits{/sha}",
      "comments_url":"https://api.github.com/repos/jdel/acdc/comments{/number}",
      "issue_comment_url":"https://api.github.com/repos/jdel/acdc/issues/comments{/number}",
      "contents_url":"https://api.github.com/repos/jdel/acdc/contents/{+path}",
      "compare_url":"https://api.github.com/repos/jdel/acdc/compare/{base}...{head}",
      "merges_url":"https://api.github.com/repos/jdel/acdc/merges",
      "archive_url":"https://api.github.com/repos/jdel/acdc/{archive_format}{/ref}",
      "downloads_url":"https://api.github.com/repos/jdel/acdc/downloads",
      "issues_url":"https://api.github.com/repos/jdel/acdc/issues{/number}",
      "pulls_url":"https://api.github.com/repos/jdel/acdc/pulls{/number}",
      "milestones_url":"https://api.github.com/repos/jdel/acdc/milestones{/number}",
      "notifications_url":"https://api.github.com/repos/jdel/acdc/notifications{?since,all,participating}",
      "labels_url":"https://api.github.com/repos/jdel/acdc/labels{/name}",
      "releases_url":"https://api.github.com/repos/jdel/acdc/releases{/id}",
      "deployments_url":"https://api.github.com/repos/jdel/acdc/deployments",
      "created_at":"2017-08-17T14:19:02Z",
      "updated_at":"2017-08-17T14:21:30Z",
      "pushed_at":"2017-08-22T20:49:51Z",
      "git_url":"git://github.com/jdel/acdc.git",
      "ssh_url":"git@github.com:jdel/acdc.git",
      "clone_url":"https://github.com/jdel/acdc.git",
      "svn_url":"https://github.com/jdel/acdc",
      "homepage":null,
      "size":89,
      "stargazers_count":0,
      "watchers_count":0,
      "language":"Go",
      "has_issues":true,
      "has_projects":true,
      "has_downloads":true,
      "has_wiki":true,
      "has_pages":false,
      "forks_count":0,
      "mirror_url":null,
      "open_issues_count":0,
      "forks":0,
      "open_issues":0,
      "watchers":0,
      "default_branch":"master"
   },
   "sender":{
      "login":"jdel",
      "id":1107511,
      "avatar_url":"https://avatars2.githubusercontent.com/u/1107511?v=4",
      "gravatar_id":"",
      "url":"https://api.github.com/users/jdel",
      "html_url":"https://github.com/jdel",
      "followers_url":"https://api.github.com/users/jdel/followers",
      "following_url":"https://api.github.com/users/jdel/following{/other_user}",
      "gists_url":"https://api.github.com/users/jdel/gists{/gist_id}",
      "starred_url":"https://api.github.com/users/jdel/starred{/owner}{/repo}",
      "subscriptions_url":"https://api.github.com/users/jdel/subscriptions",
      "organizations_url":"https://api.github.com/users/jdel/orgs",
      "repos_url":"https://api.github.com/users/jdel/repos",
      "events_url":"https://api.github.com/users/jdel/events{/privacy}",
      "received_events_url":"https://api.github.com/users/jdel/received_events",
      "type":"User",
      "site_admin":false
   }
}
*/
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

func checkMAC(message, messageMAC, key []byte) bool {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	fmt.Println(string(expectedMAC))
	fmt.Println(string(messageMAC))
	return hmac.Equal(messageMAC, expectedMAC)
}

// veiryfyHMAC generated the HMAC signature with Hash and secret and compare with the provided digest
func verifyHMAC(h func() hash.Hash, secret string, payload string, digest []byte) bool {
	hHash := hmac.New(h, []byte(secret))
	_, _ = hHash.Write([]byte(payload)) // assignations are required not to get an errcheck issue (linter)
	computedDigest := hHash.Sum(nil)

	return hmac.Equal(computedDigest, digest)
}

// RouteGithub handles incoming github pushes
func RouteGithub(w http.ResponseWriter, r *http.Request) {
	var output []byte
	// var incomingPayload githubPayload

	// event := r.Header["X-GitHub-Event"]
	// signature := r.Header["X-Hub-Signature"]

	if r.Header.Get("X-GitHub-Event") != "push" {
		return
	}

	body, _ := ioutil.ReadAll(r.Body)

	if !checkMAC(body, []byte(r.Header.Get("X-Hub-Signature")), []byte("cakecakecake")) {
		jsonOutput(w, http.StatusOK,
			outputKey("Invalid signature", ""))
		return
	}

	// apiKey := mux.Vars(r)["apiKey"]
	// tag := mux.Vars(r)["tag"]

	// key := api.FindKey(apiKey)
	// if key.Unique == "" {
	// 	jsonOutput(w, http.StatusInternalServerError,
	// 		outputKey("Could not get key", apiKey))
	// 	return
	// }

	// decoder := json.NewDecoder(r.Body)
	// err := decoder.Decode(&incomingPayload)
	// if err != nil {
	// 	logRoute.Error(err)
	// }
	// defer r.Body.Close()

	// hook := key.GetHook("incomingPayload.Repository.Name")

	// output, err = hook.Pull().CombinedOutput()
	// if err != nil {
	// 	logRoute.Error(string(output), err)
	// 	jsonOutput(w, http.StatusInternalServerError,
	// 		outputHook("Could not pull images for hook", hook.Name))
	// 	return
	// }

	// output, err = hook.Down().CombinedOutput()
	// if err != nil {
	// 	logRoute.Error(string(output), err)
	// 	jsonOutput(w, http.StatusInternalServerError,
	// 		outputHook("Could not bring hook down", hook.Name))
	// 	return
	// }

	// output, err = hook.Up().CombinedOutput()
	// if err != nil {
	// 	logRoute.Error(string(output), err)
	// 	jsonOutput(w, http.StatusInternalServerError,
	// 		outputHook("Could not bring hook up", hook.Name))
	// 	return
	// }

	jsonOutput(w, http.StatusOK,
		outputHook(string(output), ""))
}
