package renovate

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

const (
	dependencyDashboard = "https://github.com/redwoodjs/redwood/issues/3795"
	renovatePRs         = "https://github.com/redwoodjs/redwood/pulls/app%2Frenovate"
)

func PrintErrMsg() {
	fmt.Println("Expected one of: open, update")
}

func Open() {
	exec.Command("open", dependencyDashboard).Run()
	exec.Command("open", renovatePRs).Run()
}

func Update() {
	// client
	// ------------------------

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	client := githubv4.NewClient(httpClient)

	// query
	// ------------------------

	var query SearchQuery
	queryVars := map[string]interface{}{
		"query": githubv4.String("repo:redwoodjs/redwood is:pr is:open author:app/renovate -label:release:feature-breaking"),
	}

	queryErr := client.Query(context.Background(), &query, queryVars)
	if queryErr != nil {
		log.Fatal(queryErr)
	}

	var wg sync.WaitGroup
	for _, prNode := range query.Search.Nodes {
		if prNode.PullRequest.Milestone.Title != "" {
			fmt.Printf("Skipped %s %s\n", prNode.PullRequest.URL, prNode.PullRequest.Title)
			continue
		}

		wg.Add(1)

		go func(prNode Node) {
			defer wg.Done()

			// mutation
			// ------------------------

			// release:chore
			labelIDs := []githubv4.ID{"LA_kwDOC2M2f87afZ1K"}
			// next-release
			milestoneID := githubv4.ID("MI_kwDOC2M2f84Aa82f")

			var mutation UpdatePullRequestMutation
			input := githubv4.UpdatePullRequestInput{
				PullRequestID: prNode.PullRequest.ID,
				LabelIDs:      &labelIDs,
				MilestoneID:   &milestoneID,
			}

			err := client.Mutate(context.Background(), &mutation, input, nil)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("Updated % s%s\n", prNode.PullRequest.URL, prNode.PullRequest.Title)
		}(prNode)
	}

	wg.Wait()
}

type LabelConnection struct {
	Nodes []struct {
		Name string
	}
}

type PullRequest struct {
	ID        string
	URL       string
	Title     string
	Labels    LabelConnection `graphql:"labels(first: 10)"`
	Milestone struct {
		Title string
	}
}

type Node struct {
	PullRequest PullRequest `graphql:"... on PullRequest"`
}

type SearchQuery struct {
	Search struct {
		Nodes []Node
	} `graphql:"search(query: $query, type: ISSUE, first: 100)"`
}

type UpdatePullRequestPayload struct {
	ClientMutationID string
	PullRequest      struct {
		Title string
		URL   string
	}
}

type UpdatePullRequestMutation struct {
	UpdatePullRequest UpdatePullRequestPayload `graphql:"updatePullRequest(input: $input)"`
}
