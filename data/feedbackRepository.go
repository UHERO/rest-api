package data

import (
	"fmt"
	"github.com/UHERO/rest-api/models"
	"github.com/andygrunwald/go-jira"
	"log"
	"os"
)

type FeedbackRepository struct {
}

func (r *FeedbackRepository) CreateFeedback(feedback models.Feedback) (err error) {
	jiraClient, err := jira.NewClient(nil, "https://uhero-analytics.atlassian.net")
	if err != nil {
		log.Printf("JIRA connection error: %s", err)
		return
	}

	res, err := jiraClient.Authentication.AcquireSessionCookie(os.Getenv("JIRA_USER"), os.Getenv("JIRA_PASSWORD"))
	if err != nil || res == false {
		log.Printf("JIRA Result: %v\n", res)
		log.Printf("JIRA authentication error: %s", err)
		return
	}

	i := jira.Issue{
		Fields: &jira.IssueFields{
			Assignee: &jira.User{
				Name: "admin",
			},
			Reporter: &jira.User{
				Name: "admin",
			},
			Summary: "Data Portal Feedback",
			Type: jira.IssueType{
				Name: "Bug",
			},
			Project: jira.Project{
				Key: "UA",
			},
			Description: fmt.Sprintf(
				"name: %s email: %s feedback: %s",
				feedback.Name,
				feedback.Email,
				feedback.Feedback,
			),
		},
	}
	issue, _, err := jiraClient.Issue.Create(&i)
	if err != nil {
		log.Printf("JIRA issue not saved: %s", err)
		return
	}

	log.Printf("JIRA issue saved: %s\n", issue)
	return
}
