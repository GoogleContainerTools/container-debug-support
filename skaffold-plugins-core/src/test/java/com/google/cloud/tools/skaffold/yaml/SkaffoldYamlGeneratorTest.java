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
import com.google.common.io.Resources;
import java.io.IOException;
import java.net.URISyntaxException;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import org.junit.Assert;
import org.junit.Test;

/** Tests for {@link SkaffoldYamlGenerator}. */
public class SkaffoldYamlGeneratorTest {

  @Test
  public void testGenerate() throws URISyntaxException, IOException {
    String expected =
        new String(
            Files.readAllBytes(Paths.get(Resources.getResource("yaml/test-skaffold.yaml").toURI())),
            StandardCharsets.UTF_8);
    ImmutableList<Path> paths =
        ImmutableList.of(Paths.get("MANIFEST_PATH_1"), Paths.get("MANIFEST_PATH_2"));
    SkaffoldYamlGenerator generator = new SkaffoldYamlGenerator(paths);
    Assert.assertEquals(expected, generator.generate());
  }
}
