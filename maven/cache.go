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
	"os/user"
	"path/filepath"

	"github.com/buildpack/libbuildpack"
	"github.com/cloudfoundry/libjavabuildpack"
)

// Cache represents the location that Maven caches its downloaded artifacts for reuse.
type Cache struct {
	layer  libbuildpack.CacheLayer
	logger libjavabuildpack.Logger
}

// Contribute links the cache layer to $HOME/.m2.
func (c Cache) Contribute() error {
	m2, err := c.m2()
	if err != nil {
		return err
	}

	exists, err := libjavabuildpack.FileExists(m2)
	if err != nil {
		return err
	}

	if exists {
		c.logger.Debug("Maven cache already exists")
		return nil
	}

	c.logger.SubsequentLine("Linking Maven Cache to %s", m2)

	c.logger.Debug("Creating cache directory %s", c.layer.Root)
	if err := os.MkdirAll(c.layer.Root, 0755); err != nil {
		return err
	}

	c.logger.Debug("Linking %s => %s", c.layer.Root, m2)
	return os.Symlink(c.layer.Root, m2)
}

// String makes Cache satisfy the Stringer interface.
func (c Cache) String() string {
	return fmt.Sprintf("Cache{ layer :%s , logger: %s}", c.layer, c.logger)
}

func (c Cache) m2() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", err
	}

	return filepath.Join(u.HomeDir, ".m2"), nil
}

// NewCache creates a new Cache instance.
func NewCache(build libjavabuildpack.Build) Cache {
	return Cache{
		build.Cache.Layer("maven-cache"),
		build.Logger,
	}
}
