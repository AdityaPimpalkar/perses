// Copyright 2021 The Perses Authors
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

package role

import (
	"encoding/json"
	"fmt"

	"github.com/brunoga/deep"
	"github.com/labstack/echo/v4"
	"github.com/perses/perses/internal/api/authorization"
	apiInterface "github.com/perses/perses/internal/api/interface"
	"github.com/perses/perses/internal/api/interface/v1/role"
	"github.com/perses/perses/internal/api/plugin/schema"
	"github.com/perses/perses/pkg/model/api"
	v1 "github.com/perses/perses/pkg/model/api/v1"
	"github.com/sirupsen/logrus"
)

type service struct {
	role.Service
	dao   role.DAO
	authz authorization.Authorization
	sch   schema.Schema
}

func NewService(dao role.DAO, authz authorization.Authorization, sch schema.Schema) role.Service {
	return &service{
		dao:   dao,
		authz: authz,
		sch:   sch,
	}
}

func (s *service) Create(_ echo.Context, entity *v1.Role) (*v1.Role, error) {
	copyEntity, err := deep.Copy(entity)
	if err != nil {
		return nil, fmt.Errorf("failed to copy entity: %w", err)
	}
	return s.create(copyEntity)
}

func (s *service) create(entity *v1.Role) (*v1.Role, error) {
	// Update the time contains in the entity
	entity.Metadata.CreateNow()
	if err := s.dao.Create(entity); err != nil {
		return nil, err
	}
	// Refreshing RBAC cache as the role can add or remove new permissions to users
	if err := s.authz.RefreshPermissions(); err != nil {
		logrus.WithError(err).Error("failed to refresh RBAC cache")
	}
	return entity, nil
}

func (s *service) Update(_ echo.Context, entity *v1.Role, parameters apiInterface.Parameters) (*v1.Role, error) {
	copyEntity, err := deep.Copy(entity)
	if err != nil {
		return nil, fmt.Errorf("failed to copy entity: %w", err)
	}
	return s.update(copyEntity, parameters)
}

func (s *service) update(entity *v1.Role, parameters apiInterface.Parameters) (*v1.Role, error) {
	if entity.Metadata.Name != parameters.Name {
		logrus.Debugf("name in role %q and name from the http request %q don't match", entity.Metadata.Name, parameters.Name)
		return nil, apiInterface.HandleBadRequestError("metadata.name and the name in the http path request don't match")
	}
	if len(entity.Metadata.Project) == 0 {
		entity.Metadata.Project = parameters.Project
	} else if entity.Metadata.Project != parameters.Project {
		logrus.Debugf("project in role %q and project from the http request %q don't match", entity.Metadata.Project, parameters.Project)
		return nil, apiInterface.HandleBadRequestError("metadata.project and the project name in the http path request don't match")
	}

	// find the previous version of the role
	oldEntity, err := s.dao.Get(parameters.Project, parameters.Name)
	if err != nil {
		return nil, err
	}
	entity.Metadata.Update(oldEntity.Metadata)
	if updateErr := s.dao.Update(entity); updateErr != nil {
		logrus.WithError(updateErr).Errorf("unable to perform the update of the role %q, something wrong with the database", entity.Metadata.Name)
		return nil, updateErr
	}
	// Refreshing RBAC cache as the role can add or remove new permissions to users
	if err := s.authz.RefreshPermissions(); err != nil {
		logrus.WithError(err).Error("failed to refresh RBAC cache")
	}
	return entity, nil
}

func (s *service) Delete(_ echo.Context, parameters apiInterface.Parameters) error {
	if err := s.dao.Delete(parameters.Project, parameters.Name); err != nil {
		return err
	}
	// Refreshing RBAC cache as the role can add or remove new permissions to users
	if err := s.authz.RefreshPermissions(); err != nil {
		logrus.WithError(err).Error("failed to refresh RBAC cache")
	}
	return nil
}

func (s *service) Get(parameters apiInterface.Parameters) (*v1.Role, error) {
	return s.dao.Get(parameters.Project, parameters.Name)
}

func (s *service) List(q *role.Query, params apiInterface.Parameters) ([]*v1.Role, error) {
	query, err := manageQuery(q, params)
	if err != nil {
		return nil, err
	}
	return s.dao.List(query)
}

func (s *service) RawList(q *role.Query, params apiInterface.Parameters) ([]json.RawMessage, error) {
	query, err := manageQuery(q, params)
	if err != nil {
		return nil, err
	}
	return s.dao.RawList(query)
}

func (s *service) MetadataList(q *role.Query, params apiInterface.Parameters) ([]api.Entity, error) {
	query, err := manageQuery(q, params)
	if err != nil {
		return nil, err
	}
	return s.dao.MetadataList(query)
}

func (s *service) RawMetadataList(q *role.Query, params apiInterface.Parameters) ([]json.RawMessage, error) {
	query, err := manageQuery(q, params)
	if err != nil {
		return nil, err
	}
	return s.dao.RawMetadataList(query)
}

func manageQuery(q *role.Query, params apiInterface.Parameters) (*role.Query, error) {
	// Query is copied because it can be modified by the toolbox.go: listWhenPermissionIsActivated(...) and need to `q` need to keep initial value
	query, err := deep.Copy(q)
	if err != nil {
		return nil, fmt.Errorf("unable to copy the query: %w", err)
	}
	if len(query.Project) == 0 {
		query.Project = params.Project
	}
	return query, nil
}
