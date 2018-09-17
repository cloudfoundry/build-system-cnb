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
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/buildpack/libbuildpack"
	"github.com/cloudfoundry/libjavabuildpack"
	"github.com/fatih/color"
)

// Runner represents the behavior of running the maven command to build an application.
type Runner struct {
	application libbuildpack.Application
	logger      libjavabuildpack.Logger
	mvn         string
}

// Contributes builds the application from source code, removes the source, and expands the built artifact into
// $APPLICATION_ROOT.
func (r Runner) Contribute() error {
	r.logger.FirstLine("%s application", color.YellowString("Building"))

	if err := r.command().Run(); err != nil {
		return err
	}

	a, err := r.builtArtifact()
	if err != nil {
		return err
	}

	tmp, err := r.preserveBuiltArtifact(a)
	if err != nil {
		return err
	}

	r.logger.SubsequentLine("Removing source code")
	if err := os.RemoveAll(r.application.Root); err != nil {
		return err
	}

	r.logger.Debug("Expanding %s to %s", tmp, r.application.Root)
	return libjavabuildpack.ExtractZip(tmp, r.application.Root, 0)

}

// String makes Runner satisfy the Stringer interface.
func (r Runner) String() string {
	return fmt.Sprintf("Runner{ application: %s, logger: %s, mvn: %s}", r.application, r.logger, r.mvn)
}

func (r Runner) builtArtifact() (string, error) {
	target := filepath.Join(r.application.Root, "target")

	candidates, err := filepath.Glob(filepath.Join(target, "*.jar"))
	if err != nil {
		return "", err
	}

	if len(candidates) == 0 {
		return "", fmt.Errorf("unable to find built artifact in %s", target)
	}

	artifact := candidates[0]
	r.logger.Debug("Built artifact: %s", artifact)
	return artifact, nil
}

func (r Runner) command() *exec.Cmd {
	cmd := exec.Command(r.mvn, "-Dmaven.test.skip=true", "package")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	cmd.Dir = r.application.Root

	r.logger.SubsequentLine("Running %s %s", cmd.Path, strings.Join(cmd.Args, " "))
	return cmd
}

func (r Runner) preserveBuiltArtifact(artifact string) (string, error) {
	tmp, err := ioutil.TempFile("", "runner")
	if err != nil {
		return "", err
	}

	r.logger.Debug("Copying %s to %s", artifact, tmp.Name())
	libjavabuildpack.CopyFile(artifact, tmp.Name())

	return tmp.Name(), nil
}

// NewRunner creates a new Runner instance.
func NewRunner(build libjavabuildpack.Build, maven Maven) Runner {
	return Runner{
		build.Application,
		build.Logger,
		maven.Executable(),
	}
}
