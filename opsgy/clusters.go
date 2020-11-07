package opsgy

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Cluster struct {
	ClusterID string         `json:"clusterId,omitempty"`
	ProjectID string         `json:"projectId,omitempty"`
	Name      string         `json:"name,omitempty"`
	Region    string         `json:"region,omitempty"`
	Status    *ClusterStatus `json:"status,omitempty"`
}

type ClusterStatus struct {
	KubeApiURL string `json:"kubeApiUrl,omitempty"`
}

type ClusterPage struct {
	NextPage string     `json:"nextPage,omitempty"`
	PageSize int        `json:"pageSize,omitempty"`
	Items    []*Cluster `json:"items,omitempty"`
}

func GetClusters(client *http.Client, projectID string) ([]*Cluster, error) {
	list := make([]*Cluster, 0)

	var lastItem *string
	for {
		var startFrom = ""
		if lastItem != nil {
			startFrom = "&startFrom=" + *lastItem
		}
		url := OpsgyApiUrl + "/v1/projects/" + projectID + "/clusters?pageSize=50" + startFrom

		resp, err := client.Get(url)
		if err != nil {
			return nil, err
		}

		defer resp.Body.Close()

		var page ClusterPage

		err = json.NewDecoder(resp.Body).Decode(&page)
		if err != nil {
			return nil, err
		}

		for _, item := range page.Items {
			list = append(list, item)
		}
		lastItem = &page.NextPage
		if len(page.Items) < 50 {
			break
		}
	}
	return list, nil
}

func GetClusterByName(client *http.Client, projectID string, clusterName string) (*Cluster, error) {
	clusters, err := GetClusters(client, projectID)
	if err != nil {
		return nil, err
	}

	for _, cluster := range clusters {
		if cluster.Name == clusterName {
			return cluster, nil
		}
	}
	return nil, fmt.Errorf("Could not find cluster with name: %s", clusterName)
}
