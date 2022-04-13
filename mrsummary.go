package main

import (
	"MRSummary/lib"
	"github.com/jessevdk/go-flags"
	"log"
	"os"
)

var opts struct {
	JiraUrl        string `long:"jiraurl" description:"Jira instance URL" env:"JIRA_URL"`
	JiraEmail      string `long:"jiraemail" description:"Jira API Token" env:"JIRA_EMAIL"`
	JiraToken      string `long:"jiratoken" description:"Jira API Token" env:"JIRA_TOKEN"`
	GitlabProject  string `long:"project" description:"Gitlab Project ID" env:"CI_MERGE_REQUEST_PROJECT_ID"`
	GitlabMergeIID int    `long:"mergeiid" description:"Gitlab Merge Internal ID" env:"CI_MERGE_REQUEST_IID"`
	GitlabToken    string `long:"gitlabtoken" env:"GITLAB_TOKEN"`
}

func main() {
	var err error
	_, err = flags.Parse(&opts)
	if err != nil {
		log.Print(err)
		os.Exit(-2)
	}

	git, err := lib.NewGit(opts.GitlabToken)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	log.Println("Parsing git commit logs for this MR")
	tickets, err := git.Logs(opts.GitlabProject, opts.GitlabMergeIID)

	if err != nil {
		log.Print(err)
		os.Exit(-3)

	}

	j := lib.NewJira(opts.JiraUrl, opts.JiraEmail, opts.JiraToken)

	log.Println("Getting Jira ticket info")
	parsedTickets, err := j.GetMultipleTickets(tickets)
	if err != nil {
		log.Print(err)
		os.Exit(-4)
	}

	err = git.SetComment(opts.GitlabProject, opts.GitlabMergeIID, opts.JiraUrl, parsedTickets)
	if err != nil {
		log.Print(err)
		os.Exit(-5)
	}
}
