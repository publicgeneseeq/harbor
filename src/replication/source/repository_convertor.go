// Copyright (c) 2017 VMware, Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package source

import (
	"github.com/vmware/harbor/src/common/utils/log"
	"github.com/vmware/harbor/src/replication"
	"github.com/vmware/harbor/src/replication/models"
	"github.com/vmware/harbor/src/replication/registry"
)

// RepositoryConvertor implement Convertor interface, convert projects to repositories
type RepositoryConvertor struct {
	registry registry.Adaptor
}

// NewRepositoryConvertor returns an instance of RepositoryConvertor
func NewRepositoryConvertor(registry registry.Adaptor) *RepositoryConvertor {
	return &RepositoryConvertor{
		registry: registry,
	}
}

// Convert projects to repositories
func (r *RepositoryConvertor) Convert(items []models.FilterItem) []models.FilterItem {
	result := []models.FilterItem{}
	for _, item := range items {
		if item.Kind != replication.FilterItemKindProject {
			log.Warningf("unexpected filter item kind for repository convertor, expected %s got %s, skip",
				replication.FilterItemKindProject, item.Kind)
			continue
		}

		repositories := r.registry.GetRepositories(item.Value)
		for _, repository := range repositories {
			result = append(result, models.FilterItem{
				Kind:  replication.FilterItemKindRepository,
				Value: repository.Name,
				// public is used to create project if it does not exist when replicating
				Metadata: map[string]interface{}{
					"public": item.Metadata["public"],
				},
			})
		}
	}
	return result
}