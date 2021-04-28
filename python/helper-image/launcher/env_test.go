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
	"path/filepath"
	"sort"
	"testing"
)

func TestEnvAsPairs(t *testing.T) {
	tests := []struct {
		description string
		env         map[string]string
		expected    []string
	}{
		{"nil", nil, nil},
		{"empty", map[string]string{}, nil},
		{"single", map[string]string{"a": "b"}, []string{"a=b"}},
		{"multiple", map[string]string{"a": "b", "c": "d"}, []string{"a=b", "c=d"}},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			result := env(test.env).AsPairs()
			sort.Strings(result)
			if len(result) != len(test.expected) {
				t.Errorf("expected %v but got %v", test.expected, result)
			} else {
				for i := 0; i < len(result); i++ {
					if result[i] != test.expected[i] {
						t.Errorf("expected %v but got %v", test.expected[i], result[i])
					}
				}
			}
		})
	}
}

func TestEnvFromPairs(t *testing.T) {
	tests := []struct {
		description string
		env         []string
		expected    map[string]string
	}{
		{"nil", nil, nil},
		{"empty", []string{}, nil},
		{"single", []string{"a=b"}, map[string]string{"a": "b"}},
		{"multiple", []string{"a=b", "c=d"}, map[string]string{"a": "b", "c": "d"}},
		{"collisions", []string{"a=b", "a=d"}, map[string]string{"a": "d"}},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			result := EnvFromPairs(test.env)
			if len(result) != len(test.expected) {
				t.Errorf("expected %v but got %v", test.expected, result)
			} else {
				for k, v := range result {
					if v != test.expected[k] {
						t.Errorf("for %v expected %v but got %v", k, test.expected[k], v)
					}
				}
			}
		})
	}
}

func TestEnvPrependFilepath(t *testing.T) {
	tests := []struct {
		description string
		env         env
		key         string
		value       string
		expected    map[string]string
	}{
		{"empty", env{}, "PATH", "value", env{"PATH": "value"}},
		{"existing value", env{"PATH": "other"}, "PATH", "value", env{"PATH": "value" + string(filepath.ListSeparator) + "other"}},
		{"other value unchanged", env{"PYTHONPATH": "other"}, "PATH", "value", env{"PATH": "value", "PYTHONPATH": "other"}},
		{"existing value with other value unchanged", env{"PATH": "other", "PYTHONPATH": "other"}, "PATH", "value", env{"PATH": "value" + string(filepath.ListSeparator) + "other", "PYTHONPATH": "other"}},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			result := test.env // not a copy but that's ok for this test
			result.PrependFilepath(test.key, test.value)
			if len(result) != len(test.expected) {
				t.Errorf("expected %v but got %v", test.expected, result)
			} else {
				for k, v := range result {
					if v != test.expected[k] {
						t.Errorf("for %v expected %v but got %v", k, test.expected[k], v)
					}
				}
			}
		})
	}
}
