/*
 * Copyright 2018 Google LLC. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not
 * use this file except in compliance with the License. You may obtain a copy of
 * the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
 * License for the specific language governing permissions and limitations under
 * the License.
 */

package com.google.cloud.tools.skaffold.yaml;

import com.google.common.base.Preconditions;
import com.google.common.collect.ImmutableList;

/**
 * Automatically generates contents of a skaffold.yaml from a provided set of Kubernetes manifests.
 */
public class SkaffoldYamlGenerator {

  private final ImmutableList<String> manifestPaths;

  /**
   * Creates a new {@link SkaffoldYamlGenerator}.
   *
   * @param manifestPaths a non-empty list of paths to Kubernetes yamls (may include glob patterns)
   */
  public SkaffoldYamlGenerator(ImmutableList<String> manifestPaths) {
    Preconditions.checkArgument(manifestPaths.size() > 0);
    this.manifestPaths = manifestPaths;
  }

  /**
   * Generates the skaffold.yaml contents.
   *
   * @return the skaffold.yaml contents as a string
   */
  public String generate() {
    StringBuilder output = new StringBuilder();
    output.append("apiVersion: skaffold/v1alpha2\n");
    output.append("kind: Config\n");
    output.append("deploy:\n");
    output.append("  kubectl:\n");

    // Add manifests
    output.append("    manifests:\n");
    for (String manifestPath : manifestPaths) {
      output.append("    - ");
      output.append(manifestPath);
      output.append("\n");
    }

    return output.toString();
  }
}
