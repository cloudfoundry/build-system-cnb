/*
 * Copyright 2018 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package maven

import (
	"fmt"
	"path/filepath"

	"github.com/buildpack/libbuildpack"
	"github.com/cloudfoundry/jvm-application-buildpack"
	"github.com/cloudfoundry/libjavabuildpack"
	"github.com/cloudfoundry/openjdk-buildpack"
)

// MavenDependency is the key identifying the Maven build system in the buildpack plan.
const MavenDependency = "maven"

// Maven represents the Maven executable contributed by the buildpack.
type Maven struct {
	application libbuildpack.Application
	layer       libjavabuildpack.DependencyCacheLayer
	logger      libjavabuildpack.Logger
}

// Contribute makes the contribution to the cache layer.
func (m Maven) Contribute() error {
	if m.hasWrapper() {
		m.logger.FirstLine("Using Maven wrapper")
		return nil
	}

	return m.layer.Contribute(func(artifact string, layer libjavabuildpack.DependencyCacheLayer) error {
		layer.Logger.SubsequentLine("Expanding to %s", layer.Root)
		return libjavabuildpack.ExtractTarGz(artifact, layer.Root, 1)
	})
}

// Executable returns the path to the executable that should be used.  Will be the wrapper if it exists, the downloaded
// Maven distribution otherwise.
func (m Maven) Executable() string {
	if m.hasWrapper() {
		return m.wrapper()
	}

	return m.maven()
}

// String makes Maven satisfy the Stringer interface.
func (m Maven) String() string {
	return fmt.Sprintf("Maven{ application: %s, layer :%s , logger: %s}", m.application, m.layer, m.logger)
}

func (m Maven) hasWrapper() bool {
	exists, err := libjavabuildpack.FileExists(m.wrapper())
	if err != nil {
		return false
	}

	return exists
}

func (m Maven) maven() string {
	return filepath.Join(m.layer.Root, "bin", "mvn")
}

func (m Maven) wrapper() string {
	return filepath.Join(m.application.Root, "mvnw")
}

// BuildPlanContribution returns the BuildPlan with requirements for Maven.
func BuildPlanContribution() libbuildpack.BuildPlan {
	return libbuildpack.BuildPlan{
		MavenDependency:                          libbuildpack.BuildPlanDependency{},
		jvm_application_buildpack.JVMApplication: libbuildpack.BuildPlanDependency{},
		openjdk_buildpack.JDKDependency: libbuildpack.BuildPlanDependency{
			Version: "1.*",
		},
	}
}

// NewMaven creates a new Maven instance. OK is true if build plan contains "maven" dependency, otherwise false.
func NewMaven(build libjavabuildpack.Build) (Maven, bool, error) {
	bp, ok := build.BuildPlan[MavenDependency]
	if !ok {
		return Maven{}, false, nil
	}

	deps, err := build.Buildpack.Dependencies()
	if err != nil {
		return Maven{}, false, err
	}

	dep, err := deps.Best(MavenDependency, bp.Version, build.Stack)
	if err != nil {
		return Maven{}, false, err
	}

	return Maven{
		build.Application,
		build.Cache.DependencyLayer(dep),
		build.Logger,
	}, true, nil
}
