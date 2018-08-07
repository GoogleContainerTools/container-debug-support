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

package com.google.cloud.tools.skaffold.command;

import com.google.common.io.Resources;
import java.io.ByteArrayInputStream;
import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.io.InputStream;
import java.net.URISyntaxException;
import java.nio.charset.StandardCharsets;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.Arrays;
import java.util.concurrent.ExecutionException;
import org.junit.Assert;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.mockito.junit.MockitoJUnitRunner;

/** Tests for {@link Skaffold}. */
@RunWith(MockitoJUnitRunner.class)
public class SkaffoldTest {

  @Before
  public void setUp() {}

  @Test
  public void testDeploy()
      throws URISyntaxException, IOException, InterruptedException, ExecutionException {
    Path executable = Paths.get(Resources.getResource("command.sh").toURI());
    InputStream stdinInputStream =
        new ByteArrayInputStream("input".getBytes(StandardCharsets.UTF_8));
    ByteArrayOutputStream stdoutOutputStream = new ByteArrayOutputStream();
    ByteArrayOutputStream stderrOutputStream = new ByteArrayOutputStream();

    int exitCode =
        new Skaffold(executable.toString())
            .setProcessBuilderFactory(
                command -> {
                  Assert.assertEquals(
                      Arrays.asList(
                          executable.toString(),
                          "--filename",
                          "skaffoldYaml",
                          "--profile",
                          "profile",
                          "deploy"),
                      command);
                  return new ProcessBuilder(command);
                })
            .setSkaffoldYaml(Paths.get("skaffoldYaml"))
            .setProfile("profile")
            .setStdinInputStream(stdinInputStream)
            .setStdoutOutputStream(stdoutOutputStream)
            .setStderrOutputStream(stderrOutputStream)
            .deploy();

    Assert.assertEquals(0, exitCode);
    Assert.assertEquals(
        "input\noutput\n", new String(stdoutOutputStream.toByteArray(), StandardCharsets.UTF_8));
    Assert.assertEquals(
        "error\n", new String(stderrOutputStream.toByteArray(), StandardCharsets.UTF_8));
  }
}
