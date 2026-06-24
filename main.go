package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"

	"net/http"
	"os"
	"strings"
	"time"
)

// Github URL
const (
	Github    = "https://api.github.com/repos/betorvs"
	layoutISO = "2006-01-02"
	path      = "README.md"
	MyRepos   = "https://github.com/betorvs"
)

// Repository struct
type Repository struct {
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
	repositories := []string{}
	repo, ok := os.LookupEnv("REPOSITORY_LIST")
	if ok {
		list := strings.Split(repo, ",")
		for _, v := range list {
			repositories = append(repositories, strings.TrimSpace(v))
		}
	}
	var list []string
	client := newClientWeb()
	if len(repositories) != 0 {
		for _, v := range repositories {
			result := ""
			perPage := 1
			if strings.Contains(v, "modules") {
				perPage = 5
			}
			github := fmt.Sprintf("%s/%s/releases?per_page=%d", Github, v, perPage)
			repos, err := client.getRepositories(github)
			if err == nil {
				head := ""
				lines := []string{}
				for _, repo := range repos {
					if head == "" {
						head = fmt.Sprintf("[%s](%s/%s) ", v, MyRepos, v)
					}
					r := fmt.Sprintf("[%s](%s) - %s", repo.TagName, repo.HTMLURL, repo.PublishedAt.Format(layoutISO))
					lines = append(lines, r)
				}
				result = fmt.Sprintf("%s: %v", head, lines)

			}
			list = append(list, result)
		}
	} else {
		fmt.Println("Nothing to do")
	}

	f, err := os.Create(path)
	if err != nil {
		fmt.Println("create file: ", err)
	}
	t := template.New("README.tpl") // Create a template.
	t, _ = t.ParseFiles("./README.tpl")
	err = t.Execute(f, list)
	if err != nil {
		fmt.Println("executing template:", err)
	}
	f.Close()
}

type clientWeb struct {
	web *http.Client
}

func newClientWeb() *clientWeb {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	return &clientWeb{web: client}
}

func (client *clientWeb) getRepositories(github string) ([]Repository, error) {
	result := []Repository{}
	req, err := http.NewRequest(http.MethodGet, github, nil)
	if err != nil {
		fmt.Printf("[ERROR]  GET %s", err)
		return result, err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	token := fmt.Sprintf("token %s", os.Getenv("GITHUB_TOKEN"))
	req.Header.Add("Authorization", token)
	resp, err := client.web.Do(req)
	if err != nil {
		fmt.Printf("[ERROR] client %s", err)
		return result, err
	}
	res, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("[ERROR] ReadAll %s", err)
		return result, err
	}
	_ = json.Unmarshal(res, &result)

	defer resp.Body.Close()
	return result, nil
}
