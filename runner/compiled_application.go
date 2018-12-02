/*
 * Copyright 2018, Pivotal Software, Inc. All Rights Reserved.
 * Proprietary and Confidential.
 * Unauthorized use, copying or distribution of this source code via any medium is
 * strictly prohibited without the express written consent of Pivotal Software,
 * Inc.
 */

package runner

import (
	"os/exec"
	"strings"

	"github.com/buildpack/libbuildpack/application"
	"github.com/cloudfoundry/libcfbuildpack/logger"
)

// CompiledApplication represents metadata about a compiled application.
type CompiledApplication struct {
	JavaVersion string `toml:"java-version"`
}

func (c CompiledApplication) Identity() (string, string) {
	return "Compiled Application", ""
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

	return strings.TrimSpace(string(v)), nil
}
