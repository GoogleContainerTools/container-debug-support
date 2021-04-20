/*
Copyright 2021 The Skaffold Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

type env map[string]string

// EnvFromPairs turns a set of VAR=VALUE strings to a map.
func EnvFromPairs(entries []string) env {
	e := make(env)
	for _, entry := range entries {
		kv := strings.SplitN(entry, "=", 2)
		e[kv[0]] = kv[1]
	}
	return e
}

// AsPairs turns a map of variable:value pairs into a set of VAR=VALUE string pairs.
func (e env) AsPairs() []string {
	var m []string
	for k, v := range e {
		m = append(m, k+"="+v)
	}
	return m
}

// PrependFilepath prepands a path to a environment variable.
func (e env) PrependFilepath(key string, path string) {
	v := e[key]
	if v != "" {
		v = path + string(filepath.ListSeparator) + v
	} else {
		v = path
	}
	e[key] = v
}

func (e env) String() string {
	var keys []string
	for k, _ := range e {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var s string
	for _, k := range keys {
		s += fmt.Sprintf("%s=%q ", k, e[k])
	}
	return s
}
