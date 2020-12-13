package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// Github URL
const (
	Github    = "https://api.github.com/repos/betorvs"
	layoutISO = "2006-01-02"
	path      = "README.md"
)

var (
	// Repositories list of
	Repositories []string
)

// LatestRepository struct
type LatestRepository struct {
	URL             string    `json:"url"`
	HTMLURL         string    `json:"html_url"`
	AssetsURL       string    `json:"assets_url"`
	UploadURL       string    `json:"upload_url"`
	TarballURL      string    `json:"tarball_url"`
	ZipballURL      string    `json:"zipball_url"`
	ID              int       `json:"id"`
	NodeID          string    `json:"node_id"`
	TagName         string    `json:"tag_name"`
	TargetCommitish string    `json:"target_commitish"`
	Name            string    `json:"name"`
	Body            string    `json:"body"`
	Draft           bool      `json:"draft"`
	Prerelease      bool      `json:"prerelease"`
	CreatedAt       time.Time `json:"created_at"`
	PublishedAt     time.Time `json:"published_at"`
	Author          struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		NodeID            string `json:"node_id"`
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
	} `json:"author"`
	Assets []struct {
		URL                string    `json:"url"`
		BrowserDownloadURL string    `json:"browser_download_url"`
		ID                 int       `json:"id"`
		NodeID             string    `json:"node_id"`
		Name               string    `json:"name"`
		Label              string    `json:"label"`
		State              string    `json:"state"`
		ContentType        string    `json:"content_type"`
		Size               int       `json:"size"`
		DownloadCount      int       `json:"download_count"`
		CreatedAt          time.Time `json:"created_at"`
		UpdatedAt          time.Time `json:"updated_at"`
		Uploader           struct {
			Login             string `json:"login"`
			ID                int    `json:"id"`
			NodeID            string `json:"node_id"`
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
		} `json:"uploader"`
	} `json:"assets"`
}

func main() {
	Repositories = []string{"sensu-opsgenie-handler", "sensu-hangouts-chat-handler", "sensu-alertmanager-events", "sensu-grafana-mutator", "sensu-dynamic-check-mutator", "sensu-kubernetes-events"}
	var repos []string
	for _, v := range Repositories {
		github := fmt.Sprintf("%s/%s/releases/latest", Github, v)
		repo, err := getRepositories(github)
		if err == nil {
			r := fmt.Sprintf("[%s](%s) %s - %s   \n", v, repo.HTMLURL, repo.TagName, repo.PublishedAt.Format(layoutISO))
			// fmt.Println(r)
			repos = append(repos, r)
		}
	}
	f, err := os.Create(path)
	if err != nil {
		fmt.Println("create file: ", err)
	}
	t := template.New("README.tpl") // Create a template.
	t, _ = t.ParseFiles("./README.tpl")
	err = t.Execute(f, repos)
	if err != nil {
		fmt.Println("executing template:", err)
	}
	f.Close()
	// fmt.Println(repos)
}

func getRepositories(github string) (LatestRepository, error) {
	result := LatestRepository{}
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest(http.MethodGet, github, nil)
	if err != nil {
		fmt.Printf("[ERROR]  GET %s", err)
		return result, err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	token := fmt.Sprintf("token %s", os.Getenv("GITHUB_TOKEN"))
	req.Header.Add("Authorization", token)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("[ERROR] client %s", err)
		return result, err
	}
	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("[ERROR] ReadAll %s", err)
		return result, err
	}
	// fmt.Println(resp.Status)
	_ = json.Unmarshal(res, &result)

	defer resp.Body.Close()
	return result, nil
}
