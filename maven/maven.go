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
	"github.com/buildpack/libbuildpack"
	"github.com/cloudfoundry/openjdk-buildpack"
)

// MavenDependency is the key identifying the Maven build system in the buildpack plan.
const MavenDependency = "maven"

// BuildPlanContribution returns the BuildPlan with requirements for Maven.
func BuildPlanContribution() libbuildpack.BuildPlan {
	return libbuildpack.BuildPlan{
		MavenDependency:  libbuildpack.BuildPlanDependency{},
		"jvm-application": libbuildpack.BuildPlanDependency{}, // TODO use constants for jvm-application
		openjdk_buildpack.JDKDependency: libbuildpack.BuildPlanDependency{
			Version: "1.*",
		},
	}
}
