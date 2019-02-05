package main

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/google/go-github/v21/github"
)

// QueryOpenedPullRequests queries opened pull requests on specified term (fromto)
func QueryOpenedPullRequests(c *Client, fromto string) (map[string]int, error) {
	name := Username(c)
	page := 1
	pulls := make(map[string]int)
	for {
		result, resp, err := c.Search.Issues(
			context.Background(),
			fmt.Sprintf("type:pr author:%s created:%s", name, fromto),
			&github.SearchOptions{
				ListOptions: github.ListOptions{
					PerPage: 100,
					Page:    page,
				}})
		if err != nil {
			panic(err)
		}

		page = resp.NextPage

		for _, v := range result.Issues {
			ss := strings.Split(*v.RepositoryURL, "/")
			repo := strings.Join(ss[len(ss)-2:], "/")
			pulls[repo]++
		}

		if resp.NextPage == 0 {
			break
		}
	}
	return pulls, nil
}

// ShowOpenedPullRequests shows opened pull requests
func ShowOpenedPullRequests(pulls map[string]int) {
	keys := make([]string, len(pulls))
	index := 0
	for k := range pulls {
		keys[index] = k
		index++
	}
	sort.Strings(keys)

	total := 0
	for _, v := range keys {
		total += pulls[v]
	}
	fmt.Printf("Opened %d pull requests in %d repositories\n", total, len(pulls))
	for _, v := range keys {
		fmt.Printf("%d\t%s\n", pulls[v], v)
	}
}

// QueryReviewedPullRequests queries reviewed pull requests on specified term (fromto)
// TODO: it's not sure how to filter "pull requests by reviewed date".
// https://stackoverflow.com/questions/54396853/is-there-a-way-to-query-when-i-contributed-to-a-pull-request-with-submitting-rev/54441897
func QueryReviewedPullRequests(c *Client, fromto string) (map[string]int, error) {
	name := Username(c)
	page := 1
	pulls := make(map[string]int)
	for {
		// TODO: how to retrieve "reviewed date" ?
		result, resp, err := c.Search.Issues(
			context.Background(),
			fmt.Sprintf("type:pr reviewed-by:%s created:%s -author:%s", name, fromto, name),
			&github.SearchOptions{
				ListOptions: github.ListOptions{
					PerPage: 100,
					Page:    page,
				}})
		if err != nil {
			panic(err)
		}

		page = resp.NextPage

		for _, v := range result.Issues {
			ss := strings.Split(*v.RepositoryURL, "/")
			repo := strings.Join(ss[len(ss)-2:], "/")
			pulls[repo]++
		}

		if resp.NextPage == 0 {
			break
		}
	}
	return pulls, nil
}

// ShowReviewedPullRequests shows reviewed pull requests
func ShowReviewedPullRequests(pulls map[string]int) {
	keys := make([]string, len(pulls))
	index := 0
	for k := range pulls {
		keys[index] = k
		index++
	}
	sort.Strings(keys)

	total := 0
	for _, v := range keys {
		total += pulls[v]
	}
	fmt.Printf("Reviewed %d pull requests in %d repositories\n", total, len(pulls))
	for _, v := range keys {
		fmt.Printf("%d\t%s\n", pulls[v], v)
	}
}
