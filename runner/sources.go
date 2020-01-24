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

// Sources a collection of Source in order to add additional functionality.
type Sources []Source

// Len makes Sources satisfy the sort.Interface interface.
func (s Sources) Len() int {
	return len(s)
}

// Less makes Sources satisfy the sort.Interface interface.
func (s Sources) Less(i, j int) bool {
	return s[i].Path < s[j].Path
}

// Swap makes Sources satisfy the sort.Interface interface.
func (s Sources) Swap(i int, j int) {
	s[i], s[j] = s[j], s[i]
}
