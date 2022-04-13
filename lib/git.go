package lib

import (
	"bytes"
	"fmt"
	"github.com/andygrunwald/go-jira"
	"github.com/xanzy/go-gitlab"
	"html/template"
	"log"
	"regexp"
	"strings"
	"time"
)

type Git struct {
	gitlabClient *gitlab.Client
}

type Ticket struct {
	Key     string
	Url     string
	Summary string
}

type Summary struct {
	Date    string
	Tickets []Ticket
}

const CommentTemplate string = `# [MRSummary]

## Tickets in this MR
{{ with .Tickets}}
{{ range .}}- [**[{{ .Key }}]**]({{ .Url }}) {{ .Summary }}
{{ end }}
{{ end }}
*Updated: {{ .Date }} | Created by [MRSummary](https://github.com/simbachain/mrsummary)*`

func unique(intSlice []string) []string {
	keys := make(map[string]bool)
	var list []string
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func NewGit(gitlabToken string) (Git, error) {
	g := Git{}
	var err error
	g.gitlabClient, err = gitlab.NewClient(gitlabToken)

	if err != nil {
		log.Print(err)
		return g, err
	}

	return g, nil
}

func (g Git) Logs(project string, merge int) ([]string, error) {
	var mergeCommits []*gitlab.Commit

	page := 1

	for {
		optsmrc := &gitlab.GetMergeRequestCommitsOptions{
			Page:    page,
			PerPage: 5,
		}
		mrc, resp, err := g.gitlabClient.MergeRequests.GetMergeRequestCommits(project, merge, optsmrc)
		if err != nil {
			return nil, err
		}

		mergeCommits = append(mergeCommits, mrc...)

		if resp.CurrentPage >= resp.TotalPages {
			break
		}

		page++
	}

	re, _ := regexp.Compile("\\[(?P<ticket>(?P<code>[A-Z]{2,})-?(?P<number>[0-9]+))\\]")

	var tickets []string

	for _, commit := range mergeCommits {
		match := re.FindStringSubmatch(commit.Message)
		if match != nil {
			// ensure we always use the ABC-123 format, not ABC123
			code := match[re.SubexpIndex("code")]
			number := match[re.SubexpIndex("number")]
			tickets = append(tickets, fmt.Sprintf("%s-%s", code, number))
		}
	}

	tickets = unique(tickets)

	return tickets, nil
}

func (g Git) SetComment(project string, merge int, jiraUrl string, tickets []jira.Issue) error {

	summary := Summary{
		Date:    time.Now().Format(time.RFC3339),
		Tickets: []Ticket{},
	}

	for _, jiraTicket := range tickets {
		ticket := Ticket{
			Key:     jiraTicket.Key,
			Summary: jiraTicket.Fields.Summary,
			Url:     fmt.Sprintf("%s/browse/%s", jiraUrl, jiraTicket.Key),
		}

		summary.Tickets = append(summary.Tickets, ticket)
	}

	tmpl, err := template.New("comment").Parse(CommentTemplate)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)

	err = tmpl.Execute(buf, summary)

	if err != nil {
		return err
	}

	comment := buf.String()

	var comments []*gitlab.Note

	page := 1

	log.Println("Checking for existing Jira Comment")
	for {
		optsmrc := &gitlab.ListMergeRequestNotesOptions{
			ListOptions: gitlab.ListOptions{
				PerPage: 5,
				Page:    page,
			},
		}

		mrc, resp, err := g.gitlabClient.Notes.ListMergeRequestNotes(project, merge, optsmrc)
		if err != nil {
			return err
		}

		comments = append(comments, mrc...)

		if resp.CurrentPage >= resp.TotalPages {
			break
		}

		page++
	}

	var targetComment *gitlab.Note

	for _, comment := range comments {
		if strings.HasPrefix(comment.Body, "# [MRSummary]") {
			targetComment = comment
			break
		}
	}

	if targetComment != nil {
		log.Println("Existing comment found, updating")
		opts := &gitlab.UpdateMergeRequestNoteOptions{
			Body: &comment,
		}
		note, _, err := g.gitlabClient.Notes.UpdateMergeRequestNote(project, merge, targetComment.ID, opts)
		if err != nil {
			return err
		}
		log.Println(note.Body)
	} else {
		log.Println("No comment found, creating a new one")
		opt := &gitlab.CreateMergeRequestNoteOptions{
			Body: &comment,
		}
		note, _, err := g.gitlabClient.Notes.CreateMergeRequestNote(project, merge, opt)
		if err != nil {
			return err
		}
		log.Println(note.Body)
	}

	return nil
}
