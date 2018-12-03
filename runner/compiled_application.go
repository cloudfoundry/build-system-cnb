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

package runner

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/buildpack/libbuildpack/application"
	"github.com/cloudfoundry/libcfbuildpack/logger"
)

// CompiledApplication represents metadata about a compiled application.
type CompiledApplication struct {
	// JavaVersion is the version of Java used to compile the application.
	JavaVersion string `toml:"java-version"`
}

func (c CompiledApplication) Identity() (string, string) {
	return "Compiled Application", ""
}

// String makes CompiledApplication satisfy the Stringer interface.
func (c CompiledApplication) String() string {
	return fmt.Sprintf("CompiledApplication{ JavaVersion: %s }", c.JavaVersion)
}

func NewCompiledApplication(application application.Application, executor Executor, logger logger.Logger) (CompiledApplication, error) {
	v, err := javaVersion(application, executor, logger)
	if err != nil {
		return CompiledApplication{}, err
	}

	return CompiledApplication{
		v,
	}, nil
}

func javaVersion(application application.Application, executor Executor, logger logger.Logger) (string, error) {
	v, err := executor.ExecuteWithOutput(application, exec.Command("javac", "-version"), logger)
	if err != nil {
		return "", err
	}

	s := strings.Split(strings.TrimSpace(string(v)), " ")
	switch len(s) {
	case 2:
		return s[1], nil
	case 1:
		return s[0], nil
	default:
		return "unknown", nil
	}
}
