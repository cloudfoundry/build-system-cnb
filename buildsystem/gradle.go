/*
 * Copyright 2018-2020 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package buildsystem

import (
	"path/filepath"

	"github.com/buildpack/libbuildpack/application"
	"github.com/buildpack/libbuildpack/buildplan"
	"github.com/cloudfoundry/libcfbuildpack/build"
	"github.com/cloudfoundry/libcfbuildpack/helper"
	"github.com/cloudfoundry/libcfbuildpack/layers"
)

// GradleDependency is the key identifying the Gradle build system in the buildpack plan.
const GradleDependency = "gradle"

// GradlePlan returns the Plan with requirements for Gradle.
func GradlePlan() buildplan.Plan {
	return buildplan.Plan{
		Provides: []buildplan.Provided{
			{Name: GradleDependency},
			{Name: "jvm-application"},
		},
		Requires: []buildplan.Required{
			{Name: GradleDependency},
			{Name: "openjdk-jdk"},
		},
	}
}

// IsGradle returns whether this application is built using Gradle.
func IsGradle(application application.Application) bool {
	e1, err := helper.FileExists(filepath.Join(application.Root, "build.gradle"))
	if err != nil {
		return false
	}

	e2, err := helper.FileExists(filepath.Join(application.Root, "build.gradle.kts"))
	if err != nil {
		return false
	}

	return e1 || e2
}

// NewGradleBuildSystem creates a new Gradle BuildSystem instance. OK is true if build plan contains "gradle"
// dependency, otherwise false.
func NewGradleBuildSystem(build build.Build) (BuildSystem, bool, error) {
	p, ok, err := build.Plans.GetShallowMerged(GradleDependency)
	if err != nil {
		return BuildSystem{}, false, err
	} else if !ok {
		return BuildSystem{}, false, nil
	}

	deps, err := build.Buildpack.Dependencies()
	if err != nil {
		return BuildSystem{}, false, err
	}

	dep, err := deps.Best(GradleDependency, p.Version, build.Stack)
	if err != nil {
		return BuildSystem{}, false, err
	}

	layer := build.Layers.DependencyLayer(dep)
	distribution := filepath.Join(layer.Root, "bin", "gradle")
	wrapper := filepath.Join(build.Application.Root, "gradlew")

	return BuildSystem{
		contributeGradleDistribution,
		distribution,
		layer,
		build.Logger,
		wrapper,
	}, true, nil
}

func contributeGradleDistribution(artifact string, layer layers.DependencyLayer) error {
	layer.Logger.Body("Expanding to %s", layer.Root)
	return helper.ExtractZip(artifact, layer.Root, 1)
}
