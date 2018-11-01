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

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cloudfoundry/build-system-buildpack/gradle"
	"github.com/cloudfoundry/build-system-buildpack/maven"
	"github.com/cloudfoundry/libjavabuildpack"
)

func main() {
	detect, err := libjavabuildpack.DefaultDetect()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize Detect: %s\n", err.Error())
		os.Exit(101)
	}

	if isGradle(detect) {
		detect.Logger.Debug("Gradle application")
		detect.Pass(gradle.BuildPlanContribution())
		return
	}

	if isMaven(detect) {
		detect.Logger.Debug("Maven application")
		detect.Pass(maven.BuildPlanContribution())
		return
	}

	detect.Fail()
	return
}

func isGradle(detect libjavabuildpack.Detect) bool {
	build := filepath.Join(detect.Application.Root, "build.gradle")

	exists, err := libjavabuildpack.FileExists(build)
	if err != nil {
		return false
	}

	return exists
}

func isMaven(detect libjavabuildpack.Detect) bool {
	pom := filepath.Join(detect.Application.Root, "pom.xml")

	exists, err := libjavabuildpack.FileExists(pom)
	if err != nil {
		return false
	}

	return exists
}
