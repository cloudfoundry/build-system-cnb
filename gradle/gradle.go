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

package gradle

import (
	"fmt"
	"path/filepath"

	"github.com/buildpack/libbuildpack"
	"github.com/cloudfoundry/libjavabuildpack"
	"github.com/cloudfoundry/openjdk-buildpack"
)

// GradleDependency is the key identifying the Gradle build system in the buildpack plan.
const GradleDependency = "gradle"

// Gradle represents the Gradle executable contributed by the buildpack.
type Gradle struct {
	application libbuildpack.Application
	logger      libjavabuildpack.Logger
	layer       libjavabuildpack.DependencyCacheLayer
}

// Contribute makes the contribution to the cache layer
func (g Gradle) Contribute() error {
	if g.hasWrapper() {
		g.logger.SubsequentLine("Using Gradle wrapper")
		return nil
	}

	return g.layer.Contribute(func(artifact string, layer libjavabuildpack.DependencyCacheLayer) error {
		layer.Logger.SubsequentLine("Expanding to %s", layer.Root)
		return libjavabuildpack.ExtractZip(artifact, layer.Root, 1)
	})
}

// Executable returns the path to the executable that should be used.  Will be the wrapper if it exists, the downloaded
// Gradle distribution otherwise.
func (g Gradle) Executable() string {
	if g.hasWrapper() {
		return g.wrapper()
	}

	return g.gradle()
}

// String makes Gradle satisfy the Stringer interface.
func (g Gradle) String() string {
	return fmt.Sprintf("Gradle{ application: %s, logger: %s, layer :%s }", g.application, g.logger, g.layer)
}

func (g Gradle) hasWrapper() bool {
	exists, err := libjavabuildpack.FileExists(g.wrapper())
	if err != nil {
		return false
	}

	return exists
}

func (g Gradle) gradle() string {
	return filepath.Join(g.layer.Root, "bin", "gradle")
}

func (g Gradle) wrapper() string {
	return filepath.Join(g.application.Root, "gradlew")
}

// BuildPlanContribution returns the BuildPlan with requirements for Gradle.
func BuildPlanContribution() libbuildpack.BuildPlan {
	return libbuildpack.BuildPlan{
		GradleDependency:  libbuildpack.BuildPlanDependency{},
		"jvm-application": libbuildpack.BuildPlanDependency{}, // TODO use constants for jvm-application
		openjdk_buildpack.JDKDependency: libbuildpack.BuildPlanDependency{
			Version: "1.*",
		},
	}
}

// NewGradle creates a new Gradle instance. OK is true if build plan contains "gradle" dependency, otherwise false.
func NewGradle(build libjavabuildpack.Build) (Gradle, bool, error) {
	bp, ok := build.BuildPlan[GradleDependency]
	if !ok {
		return Gradle{}, false, nil
	}

	deps, err := build.Buildpack.Dependencies()
	if err != nil {
		return Gradle{}, false, err
	}

	dep, err := deps.Best(GradleDependency, bp.Version, build.Stack)
	if err != nil {
		return Gradle{}, false, err
	}

	return Gradle{
		build.Application,
		build.Logger,
		build.Cache.DependencyLayer(dep),
	}, true, nil
}