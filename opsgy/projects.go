package opsgy

import (
	"fmt"
	"encoding/json"
	"net/http"
)

type Project struct {
	ProjectID   string `json:"projectId,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	IsDefault   bool   `json:"isDefault,omitempty"`
}

type ProjectPage struct {
	NextPage string    `json:"nextPage,omitempty"`
	PageSize int       `json:"pageSize,omitempty"`
	Items    []*Project `json:"items,omitempty"`
}

func GetProjects(client *http.Client) ([]*Project, error) {
	projects := make([]*Project, 0)

	var lastItem *string
	for {
		var startFrom = ""
		if lastItem != nil {
			startFrom = "&startFrom=" + *lastItem
		}
		url := OpsgyApiUrl + "/v1/projects?pageSize=50" + startFrom

		resp, err := client.Get(url)
		if err != nil {
			return nil, err
		}

		defer resp.Body.Close()

		var projectPage ProjectPage

		err = json.NewDecoder(resp.Body).Decode(&projectPage)
		if err != nil {
			return nil, err
		}

		for _, project := range projectPage.Items {
			projects = append(projects, project)
		}
		lastItem = &projectPage.NextPage
		if len(projectPage.Items) < 50 {
			break
		}
	}
	return projects, nil
}

func GetProjectByName(client *http.Client, projectName string) (*Project, error) {
	projects, err := GetProjects(client);
	if err != nil {
		return nil, err
	}

	for _, project := range projects {
		if project.Name == projectName {
			return project, nil
		}
	}
	return nil, fmt.Errorf("Could not find project with name: %s", projectName)
}