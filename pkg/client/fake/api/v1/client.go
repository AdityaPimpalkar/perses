// Copyright 2023 The Perses Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fakev1

import (
	v1 "github.com/perses/perses/pkg/client/api/v1"
	"github.com/perses/perses/pkg/client/perseshttp"
)

type client struct {
	v1.ClientInterface
	restClient *perseshttp.RESTClient
}

func New(restClient *perseshttp.RESTClient) v1.ClientInterface {
	return &client{restClient: restClient}
}

func (c *client) RESTClient() *perseshttp.RESTClient {
	return c.restClient
}

func (c *client) Dashboard(_ string) v1.DashboardInterface {
	return &dashboard{}
}

func (c *client) EphemeralDashboard(_ string) v1.EphemeralDashboardInterface {
	return &ephemeralDashboard{}
}

func (c *client) Folder(project string) v1.FolderInterface {
	return &folder{
		project: project,
	}
}

func (c *client) GlobalDatasource() v1.GlobalDatasourceInterface {
	return &globalDatasource{}
}

func (c *client) Health() v1.HealthInterface {
	return &health{}
}

func (c *client) Plugin() v1.PluginInterface {
	return &plugin{}
}

func (c *client) Project() v1.ProjectInterface {
	return &project{}
}
