package lib

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
	"log"
	"strings"
)

type JiraLib struct {
	JiraClient *jira.Client
}

func NewJira(baseUrl string, email string, token string) JiraLib {
	j := JiraLib{}
	tp := jira.BasicAuthTransport{
		Username: email,
		Password: token,
	}
	c, err := jira.NewClient(tp.Client(), baseUrl)
	if err != nil {
		log.Fatal(err)
	}
	j.JiraClient = c

	return j
}

func (j JiraLib) GetMultipleTickets(ticketids []string) ([]jira.Issue, error) {
	opts := jira.SearchOptions{}

	jql := fmt.Sprintf("key IN (\"%s\") ORDER BY key", strings.Join(ticketids[:], "\",\""))

	issues, _, err := j.JiraClient.Issue.Search(jql, &opts)

	if err != nil {
		return nil, err
	}

	return issues, nil
}
