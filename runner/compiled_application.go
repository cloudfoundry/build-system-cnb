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

package runner

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/buildpacks/libbuildpack/v2/application"
	"github.com/cloudfoundry/libcfbuildpack/v2/logger"
	"github.com/cloudfoundry/libcfbuildpack/v2/runner"
)

// CompiledApplication represents metadata about a compiled application.
type CompiledApplication struct {
	// JavaVersion is the version of Java used to compile the application.
	JavaVersion string `toml:"java-version"`

	// Sources is metadata about the source files used to compile the application.
	Sources Sources `toml:"sources"`
}

// Identity makes CompiledApplication satisfy the Identifiable interface.
func (c CompiledApplication) Identity() (string, string) {
	return "Compiled Application", fmt.Sprintf("(%d files)", len(c.Sources))
}

func NewCompiledApplication(application application.Application, runner runner.Runner, logger logger.Logger) (CompiledApplication, error) {
	v, err := javaVersion(application, runner)
	if err != nil {
		return CompiledApplication{}, err
	}

	s, err := sources(application, logger)
	if err != nil {
		return CompiledApplication{}, err
	}

	return CompiledApplication{
		v,
		s,
	}, nil
}

func javaVersion(application application.Application, runner runner.Runner) (string, error) {
	v, err := runner.RunWithOutput("javac", application.Root, "-version")
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

type result struct {
	err   error
	value Source
}

func sources(application application.Application, logger logger.Logger) (Sources, error) {
	ch := make(chan result)
	var wg sync.WaitGroup

	if err := filepath.Walk(application.Root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		wg.Add(1)
		go func() {
			defer wg.Done()

			source, err := NewSource(path, info, logger)
			if err != nil {
				ch <- result{err: err}
				return
			}

			ch <- result{value: source}
		}()

		return nil
	}); err != nil {
		return nil, err
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	var s Sources
	for r := range ch {
		if r.err != nil {
			return Sources{}, r.err
		}

		s = append(s, r.value)
	}
	sort.Sort(s)

	return s, nil
}
