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

// Code generated. DO NOT EDIT

package role

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/perses/perses/internal/api/authorization"
	"github.com/perses/perses/internal/api/interface/v1/role"
	"github.com/perses/perses/internal/api/route"
	"github.com/perses/perses/internal/api/toolbox"
	"github.com/perses/perses/internal/api/utils"
	v1 "github.com/perses/perses/pkg/model/api/v1"
)

type endpoint struct {
	toolbox  toolbox.Toolbox[*v1.Role, *role.Query]
	readonly bool
}

func NewEndpoint(service role.Service, authz authorization.Authorization, readonly bool, caseSensitive bool) route.Endpoint {
	return &endpoint{
		toolbox:  toolbox.New[*v1.Role, *v1.Role, *role.Query](service, authz, v1.KindRole, caseSensitive),
		readonly: readonly,
	}
}

func (e *endpoint) CollectRoutes(g *route.Group) {
	group := g.Group(fmt.Sprintf("/%s", utils.PathRole))
	subGroup := g.Group(fmt.Sprintf("/%s/:%s/%s", utils.PathProject, utils.ParamProject, utils.PathRole))
	if !e.readonly {
		group.POST("", e.Create, false)
		subGroup.POST("", e.Create, false)
		subGroup.PUT(fmt.Sprintf("/:%s", utils.ParamName), e.Update, false)
		subGroup.DELETE(fmt.Sprintf("/:%s", utils.ParamName), e.Delete, false)
	}
	group.GET("", e.List, false)
	subGroup.GET("", e.List, false)
	subGroup.GET(fmt.Sprintf("/:%s", utils.ParamName), e.Get, false)
}

func (e *endpoint) Create(ctx echo.Context) error {
	entity := &v1.Role{}
	return e.toolbox.Create(ctx, entity)
}

func (e *endpoint) Update(ctx echo.Context) error {
	entity := &v1.Role{}
	return e.toolbox.Update(ctx, entity)
}

func (e *endpoint) Delete(ctx echo.Context) error {
	return e.toolbox.Delete(ctx)
}

func (e *endpoint) Get(ctx echo.Context) error {
	return e.toolbox.Get(ctx)
}

func (e *endpoint) List(ctx echo.Context) error {
	q := &role.Query{}
	return e.toolbox.List(ctx, q)
}
