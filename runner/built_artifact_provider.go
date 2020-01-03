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
	"archive/zip"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"github.com/buildpack/libbuildpack/application"
	"github.com/magiconair/properties"
)

// BuiltArtifactProvider returns the artifact built as part of running a build system.
type BuiltArtifactProvider struct {
	target string
}

// Get returns the built artifact if exactly one exists.  If less than or more than one exists, returns an error.
func (b BuiltArtifactProvider) Get(application application.Application) (string, error) {
	candidates, err := filepath.Glob(filepath.Join(application.Root, b.target))
	if err != nil {
		return "", err
	}

	var artifacts []string

	if len(candidates) == 1 {
		artifacts = candidates
	} else {
		for _, c := range candidates {
			if i, err := b.isInterestingFile(c); err != nil {
				return "", err
			} else if i {
				artifacts = append(artifacts, c)
			}
		}
	}

	if len(artifacts) != 1 {
		sort.Strings(candidates)
		return "", fmt.Errorf("unable to find built artifact (executable JAR or WAR) in %s, candidates: %s", b.target, candidates)
	}

	return artifacts[0], nil
}

func (BuiltArtifactProvider) isInterestingEntry(f *zip.File) (bool, error) {
	if f.Name == "WEB-INF/" && f.FileInfo().IsDir() {
		return true, nil
	}

	if f.Name == "META-INF/MANIFEST.MF" {
		m, err := f.Open()
		if err != nil {
			return false, err
		}
		defer m.Close()

		b, err := ioutil.ReadAll(m)
		if err != nil {
			return false, err
		}

		p, err := properties.Load(b, properties.UTF8)
		if err != nil {
			return false, nil
		}

		if _, ok := p.Get("Main-Class"); ok {
			return true, nil
		}
	}

	return false, nil
}

func (b BuiltArtifactProvider) isInterestingFile(f string) (bool, error) {
	z, err := zip.OpenReader(f)
	if err != nil {
		return false, err
	}
	defer z.Close()

	for _, f := range z.File {
		if i, err := b.isInterestingEntry(f); err != nil {
			return false, err
		} else if i {
			return true, nil
		}
	}

	return false, nil
}

// NewBuiltArtifactProvider creates a new instance using the default target if not otherwise configured.
func NewBuiltArtifactProvider(defaultTarget ...string) BuiltArtifactProvider {
	if target, ok := os.LookupEnv("BP_BUILT_ARTIFACT"); ok {
		return BuiltArtifactProvider{target}
	}

	target := filepath.Join(defaultTarget...)

	if module, ok := os.LookupEnv("BP_BUILT_MODULE"); ok {
		target = filepath.Join(module, target)
	}

	return BuiltArtifactProvider{target}
}
