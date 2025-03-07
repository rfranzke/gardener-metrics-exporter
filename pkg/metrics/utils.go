// Copyright (c) 2018 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
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

package metrics

import (
	"fmt"

	gardenv1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	seedmanagementv1alpha1 "github.com/gardener/gardener/pkg/apis/seedmanagement/v1alpha1"

	constants "github.com/gardener/gardener/pkg/apis/core/v1beta1/constants"

	"github.com/prometheus/client_golang/prometheus"
)

const unknown = "unknown"

var (
	// ScrapeFailures is a metric, which counts the amount scrape issues grouped by kind.
	ScrapeFailures = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "garden_scrape_failure_total",
		Help: "Total count of scraping failures, grouped by kind/group of metric(s)",
	}, []string{"kind"})
)

func mapConditionStatus(status gardenv1beta1.ConditionStatus) float64 {
	switch status {
	case gardenv1beta1.ConditionTrue:
		return 1
	case gardenv1beta1.ConditionFalse:
		return 0
	case gardenv1beta1.ConditionProgressing:
		return 2
	default:
		return -1
	}
}

func usedAsSeed(shoot *gardenv1beta1.Shoot, managedSeeds []*seedmanagementv1alpha1.ManagedSeed) bool {
	if shoot.Namespace != constants.GardenNamespace {
		return false
	}
	for _, ms := range managedSeeds {
		if ms.Spec.Shoot.Name == shoot.Name && ms.Namespace == shoot.Namespace {
			return true
		}
	}

	return false
}

func findProject(projects []*gardenv1beta1.Project, match string) (*string, error) {
	var projectName string
	for _, project := range projects {
		if project.Spec.Namespace != nil && *project.Spec.Namespace == match {
			projectName = project.Name
			break
		}
	}
	if projectName == "" {
		return nil, fmt.Errorf("no project found for shoot %s", match)
	}
	return &projectName, nil
}
