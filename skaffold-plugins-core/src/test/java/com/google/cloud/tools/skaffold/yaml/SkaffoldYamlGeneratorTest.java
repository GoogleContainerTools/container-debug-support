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

import com.google.common.collect.ImmutableList;
import java.nio.file.Path;
import java.nio.file.Paths;
import org.junit.Assert;
import org.junit.Test;

/** Tests for {@link SkaffoldYamlGenerator}. */
public class SkaffoldYamlGeneratorTest {

  private static final String EXPECTED_OUTPUT =
      "apiVersion: skaffold/v1alpha2\n"
          + "kind: Config\n"
          + "deploy:\n"
          + "  kubectl:\n"
          + "    manifests:\n"
          + "    - MANIFEST_PATH_1\n"
          + "    - MANIFEST_PATH_2\n";

  @Test
  public void testGenerate() {
    ImmutableList<Path> paths =
        ImmutableList.of(Paths.get("MANIFEST_PATH_1"), Paths.get("MANIFEST_PATH_2"));
    SkaffoldYamlGenerator generator = new SkaffoldYamlGenerator(paths);
    Assert.assertEquals(EXPECTED_OUTPUT, generator.generate());
  }
}
