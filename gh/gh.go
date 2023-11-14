package gh

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
)

const url = "https://api.github.com/repos/redwoodjs/redwood/actions/caches"

func PrintErrMsg() {
	fmt.Println("Expected one of: cache-clean")
}

func CacheClean() {
	client := &http.Client{}

	req, err := http.NewRequest("GET", fmt.Sprintf("%v?per_page=100", url), nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ghTok := os.Getenv("GITHUB_TOKEN")

	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("Authorization", fmt.Sprintf("token %v", ghTok))
	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var cachesResponse CachesResponse
	err = json.Unmarshal(body, &cachesResponse)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Found %v cache(s)\n", len(cachesResponse.ActionsCaches))

	var wg sync.WaitGroup
	for _, cache := range cachesResponse.ActionsCaches {
		wg.Add(1)

		go func(cache ActionsCache) {
			defer wg.Done()

			req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/%v", url, cache.ID), nil)

			req.Header.Add("Accept", "application/vnd.github+json")
			req.Header.Add("Authorization", fmt.Sprintf("token %v", ghTok))
			req.Header.Add("X-GitHub-Api-Version", "2022-11-28")

			_, err = client.Do(req)

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}(cache)
	}

	wg.Wait()

	fmt.Printf("Deleted %v cache(s)\n", len(cachesResponse.ActionsCaches))
}

type ActionsCache struct {
	ID             int    `json:"id"`
	Ref            string `json:"ref"`
	Key            string `json:"key"`
	Version        string `json:"version"`
	LastAccessedAt string `json:"last_accessed_at"`
	CreatedAt      string `json:"created_at"`
	SizeInBytes    int    `json:"size_in_bytes"`
}

type CachesResponse struct {
	TotalCount    int            `json:"total_count"`
	ActionsCaches []ActionsCache `json:"actions_caches"`
}
