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

	"github.com/buildpack/libbuildpack/application"
	"github.com/buildpack/libbuildpack/buildplan"
	"github.com/cloudfoundry/jvm-application-buildpack/jvmapplication"
	"github.com/cloudfoundry/libcfbuildpack/build"
	"github.com/cloudfoundry/libcfbuildpack/layers"
	"github.com/cloudfoundry/libcfbuildpack/logger"
	"github.com/cloudfoundry/openjdk-buildpack/jdk"
)

// Dependency is the key identifying the Gradle build system in the buildpack plan.
const Dependency = "gradle"

// Gradle represents the Gradle executable contributed by the buildpack.
type Gradle struct {
	application application.Application
	layer       layers.DependencyLayer
	logger      logger.Logger
}

// Contribute makes the contribution to the cache layer
func (g Gradle) Contribute() error {
	if g.hasWrapper() {
		g.logger.FirstLine("Using Gradle wrapper")
		return nil
	}

	return g.layer.Contribute(func(artifact string, layer layers.DependencyLayer) error {
		layer.Logger.SubsequentLine("Expanding to %s", layer.Root)
		return layers.ExtractZip(artifact, layer.Root, 1)
	}, layers.Build, layers.Cache)
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
	return fmt.Sprintf("Gradle{ application: %s, layer :%s, logger: %s }", g.application, g.layer, g.logger)
}

func (g Gradle) hasWrapper() bool {
	exists, err := layers.FileExists(g.wrapper())
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
func BuildPlanContribution() buildplan.BuildPlan {
	return buildplan.BuildPlan{
		Dependency:                buildplan.Dependency{},
		jvmapplication.Dependency: buildplan.Dependency{},
		jdk.Dependency:            buildplan.Dependency{Version: "1.*"},
	}
}

// IsGradle returns whether this application is built using Gradle.
func IsGradle(application application.Application) bool {
	exists, err := layers.FileExists(filepath.Join(application.Root, "build.gradle"))
	if err != nil {
		return false
	}

	return exists
}

// NewGradle creates a new Gradle instance. OK is true if build plan contains "gradle" dependency, otherwise false.
func NewGradle(build build.Build) (Gradle, bool, error) {
	bp, ok := build.BuildPlan[Dependency]
	if !ok {
		return Gradle{}, false, nil
	}

	deps, err := build.Buildpack.Dependencies()
	if err != nil {
		return Gradle{}, false, err
	}

	dep, err := deps.Best(Dependency, bp.Version, build.Stack)
	if err != nil {
		return Gradle{}, false, err
	}

	return Gradle{
		build.Application,
		build.Layers.DependencyLayer(dep),
		build.Logger,
	}, true, nil
}
