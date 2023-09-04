package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/google/go-github/v54/github"
)

const (
	skew = 10 * time.Second
)

var repoList = []string{
	"kubernetes-sigs/cluster-api",
	"kubernetes-sigs/cluster-api-provider-aws",
	"kubernetes-sigs/cluster-api-provider-azure",
	"kubernetes-sigs/cluster-api-provider-gcp",
}

func main() {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("missing runtime info")
	}
	dir := filepath.Join(file, "../../../files")
	fmt.Println("Using dir =", dir)

	ctx := context.Background()
	client := github.NewClient(nil)

	for _, fullname := range repoList {
		fmt.Printf("Refreshing %s\n", fullname)
		idx := strings.IndexRune(fullname, '/')
		if idx == -1 {
			panic("expected owner/repo, found " + fullname)
		}
		owner, repo := fullname[:idx], fullname[idx+1:]
		if err := DownloadReleaseAssets(ctx, client, owner, repo, dir); err != nil {
			panic(err)
		}
	}
}

func DownloadReleaseAssets(ctx context.Context, client *github.Client, owner, repo, dir string) error {
	releases, err := ListReleases(ctx, client, owner, repo)
	if err != nil {
		return err
	}
	for _, r := range releases {
		subDir := filepath.Join(dir, repo, r.GetName())
		if err := os.MkdirAll(subDir, 0o755); err != nil {
			return err
		}
		for _, a := range r.Assets {
			if strings.HasSuffix(a.GetName(), ".yaml") ||
				strings.HasSuffix(a.GetName(), ".yml") ||
				strings.HasSuffix(a.GetName(), ".json") {

				filename := filepath.Join(subDir, a.GetName())
				if exists(filename) {
					continue
				}

				fmt.Println("Downloading", a.GetBrowserDownloadURL())
				resp, err := http.DefaultClient.Get(a.GetBrowserDownloadURL())
				if err != nil {
					return err
				}
				defer resp.Body.Close()

				data, err := io.ReadAll(resp.Body)
				if err != nil {
					return err
				}
				err = os.WriteFile(filename, data, 0o644)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func ListReleases(ctx context.Context, client *github.Client, owner, repo string) ([]*github.RepositoryRelease, error) {
	opt := &github.ListOptions{
		PerPage: 100,
	}

	var result []*github.RepositoryRelease
	for {
		branch, resp, err := client.Repositories.ListReleases(ctx, owner, repo, opt)
		switch e := err.(type) {
		case *github.RateLimitError:
			time.Sleep(time.Until(e.Rate.Reset.Time.Add(skew)))
			continue
		case *github.AbuseRateLimitError:
			time.Sleep(e.GetRetryAfter())
			continue
		case *github.ErrorResponse:
			if e.Response.StatusCode == http.StatusNotFound {
				log.Println(err)
				break
			} else {
				return nil, err
			}
		default:
			if e != nil {
				return nil, err
			}
		}

		result = append(result, branch...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return result, nil
}

func exists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}
