package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Activity struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	Actor struct {
		ID           int    `json:"id"`
		DisplayLogin string `json:"display_login"`
	} `json:"actor"`
	Repo struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"repo"`
	Payload struct {
		Action  string `json:"action"`
		Commits []struct {
			Author struct {
				Name string `json:"name"`
			} `json:"author"`
			Message string `json:"message"`
			URL     string `json:"url"`
		} `json:"commits"`
	} `json:"payload"`
}

func main() {

	q := "kamranahmedse"

	if len(os.Args) >= 2 {
		q = os.Args[1]
	}
	res, err := http.Get("https://api.github.com/users/" + q + "/events")

	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		panic("either user or activities not available")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	activities := make([]Activity, 0)

	err = json.Unmarshal(body, &activities)
	if err != nil {
		panic(err)
	}

	for _, activity := range activities {
		fmt.Print("- ")

		switch activity.Type {
		case "PullRequestReviewEvent":
			fmt.Println("reviewed a PR for", activity.Repo.Name)
		case "PushEvent":
			fmt.Printf("pushed %d commit(s) to %s\n", len(activity.Payload.Commits), activity.Repo.Name)
		case "IssuesEvent":
			fmt.Printf("%s an issue on %s\n", activity.Payload.Action, activity.Repo.Name)
		case "PullRequestEvent":
			fmt.Printf("%s a PR at %s\n", activity.Payload.Action, activity.Repo.Name)
		case "IssueCommentEvent":
			fmt.Printf("%s on an Issue in %s\n", activity.Payload.Action, activity.Repo.Name)
		case "ForkEvent":
			fmt.Printf("Forked %s\n", activity.Repo.Name)
		case "CreateEvent":
			fmt.Printf("Created repository %s\n", activity.Repo.Name)
		}

	}

}
