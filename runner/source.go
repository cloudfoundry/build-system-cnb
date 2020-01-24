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
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"

	"github.com/cloudfoundry/libcfbuildpack/v2/logger"
)

// Source is metadata about a source file.
type Source struct {
	// Path is the path of the source file.
	Path string `toml:"path"`

	// Mode is the human-readable mode of the source file.
	Mode string `toml:"mode"`

	// SHA256 is the hash of the source file.
	SHA256 string `toml:"sha256"`
}

func NewSource(path string, info os.FileInfo, logger logger.Logger) (Source, error) {
	if info.IsDir() {
		return Source{
			Path: path,
			Mode: info.Mode().String(),
		}, nil
	}

	h, err := hash(path)
	if err != nil {
		return Source{}, nil
	}

	s := Source{
		path,
		info.Mode().String(),
		h,
	}

	logger.Debug("Source: %s", s)
	return s, nil
}

func hash(file string) (string, error) {
	s := sha256.New()

	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = io.Copy(s, f)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(s.Sum(nil)), nil
}
