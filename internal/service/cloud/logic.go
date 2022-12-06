// Copyright 2022 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cloud

import (
	"fmt"
	"math"
	"os"

	"tidbcloud-cli/internal/ui"

	clusterApi "github.com/c4pt0r/go-tidbcloud-sdk-v1/client/cluster"
	projectApi "github.com/c4pt0r/go-tidbcloud-sdk-v1/client/project"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/juju/errors"
)

type Project struct {
	ID   string
	Name string
}

func (p Project) String() string {
	return fmt.Sprintf("%s(%s)", p.Name, p.ID)
}

type Cluster struct {
	ID   string
	Name string
}

func (c Cluster) String() string {
	return fmt.Sprintf("%s(%s)", c.Name, c.ID)
}

func GetSelectedProject(pageSize int64, client TiDBCloudClient) (*Project, error) {
	_, projectItems, err := RetrieveProjects(pageSize, client)
	if err != nil {
		return nil, err
	}
	set := hashset.New()
	for _, item := range projectItems {
		set.Add(&Project{
			ID:   item.ID,
			Name: item.Name,
		})
	}
	model, err := ui.InitialSelectModel(set.Values(), "Choose the project:")
	if err != nil {
		return nil, err
	}
	p := tea.NewProgram(model)
	projectModel, err := p.StartReturningModel()
	if err != nil {
		return nil, err
	}
	if m, _ := projectModel.(ui.SelectModel); m.Interrupted {
		os.Exit(130)
	}
	res := projectModel.(ui.SelectModel).Choices[projectModel.(ui.SelectModel).Selected].(*Project)
	return res, nil
}

func GetSelectedCluster(projectID string, pageSize int64, client TiDBCloudClient) (*Cluster, error) {
	_, clusterItems, err := RetrieveClusters(projectID, pageSize, client)
	if err != nil {
		return nil, err
	}
	set := hashset.New()
	for _, item := range clusterItems {
		set.Add(&Cluster{
			ID:   *(item.ID),
			Name: item.Name,
		})
	}
	clusters := set.Values()
	if len(clusters) == 0 {
		return nil, fmt.Errorf("no available clusters found")
	}

	model, err := ui.InitialSelectModel(clusters, "Choose the cluster")
	if err != nil {
		return nil, errors.Trace(err)
	}
	p := tea.NewProgram(model)
	clusterModel, err := p.StartReturningModel()
	if err != nil {
		return nil, errors.Trace(err)
	}
	if m, _ := clusterModel.(ui.SelectModel); m.Interrupted {
		os.Exit(130)
	}
	cluster := clusterModel.(ui.SelectModel).Choices[clusterModel.(ui.SelectModel).Selected].(*Cluster)
	return cluster, nil
}

func RetrieveProjects(size int64, d TiDBCloudClient) (int64, []*projectApi.ListProjectsOKBodyItemsItems0, error) {
	params := projectApi.NewListProjectsParams()
	var total int64 = math.MaxInt64
	var page int64 = 1
	var pageSize = size
	var items []*projectApi.ListProjectsOKBodyItemsItems0
	for (page-1)*pageSize < total {
		projects, err := d.ListProjects(params.WithPage(&page).WithPageSize(&pageSize))
		if err != nil {
			return 0, nil, errors.Trace(err)
		}

		total = *projects.Payload.Total
		page += 1
		items = append(items, projects.Payload.Items...)
	}
	return total, items, nil
}

func RetrieveClusters(pID string, pageSize int64, d TiDBCloudClient) (int64, []*clusterApi.ListClustersOfProjectOKBodyItemsItems0, error) {
	params := clusterApi.NewListClustersOfProjectParams().WithProjectID(pID)
	var total int64 = math.MaxInt64
	var page int64 = 1
	var items []*clusterApi.ListClustersOfProjectOKBodyItemsItems0
	// loop to get all clusters
	for (page-1)*pageSize < total {
		clusters, err := d.ListClustersOfProject(params.WithPage(&page).WithPageSize(&pageSize))
		if err != nil {
			return 0, nil, errors.Trace(err)
		}

		total = *clusters.Payload.Total
		page += 1
		items = append(items, clusters.Payload.Items...)
	}
	return total, items, nil
}
