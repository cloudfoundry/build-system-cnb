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
	"os"
	"path/filepath"

	"github.com/buildpack/libbuildpack"
	"github.com/cloudfoundry/libjavabuildpack"
)

// MavenCache represents the location that Maven caches its downloaded artifacts for reuse.
type MavenCache struct {
	layer  libbuildpack.CacheLayer
	logger libjavabuildpack.Logger
}

// Contribute links the cache layer to $HOME/.m2.
func (m MavenCache) Contribute() error {
	m2 := m.m2()

	exists, err := libjavabuildpack.FileExists(m2)
	if err != nil {
		return err
	}

	if exists {
		m.logger.Debug("Maven cache already exists")
		return nil
	}

	m.logger.SubsequentLine("Linking Maven Cache to %s", m2)

	m.logger.Debug("Creating cache directory %s", m.layer.Root)
	if err := os.MkdirAll(m.layer.Root, 0755); err != nil {
		return err
	}

	m.logger.Debug("Linking %s => %s", m.layer.Root, m2)
	return os.Symlink(m.layer.Root, m2)
}

// String makes MavenCache satisfy the Stringer interface.
func (m MavenCache) String() string {
	return fmt.Sprintf("MavenCache{ layer :%s , logger: %s}", m.layer, m.logger)
}

func (m MavenCache) m2() string {
	return filepath.Join(os.Getenv("HOME"), ".m2")
}

// NewMavenCache creates a new MavenCache instance.
func NewMavenCache(build libjavabuildpack.Build) MavenCache {
	return MavenCache{
		build.Cache.Layer("maven-cache"),
		build.Logger,
	}
}
