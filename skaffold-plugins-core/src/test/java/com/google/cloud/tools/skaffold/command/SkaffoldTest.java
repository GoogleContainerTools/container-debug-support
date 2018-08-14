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

import com.google.cloud.tools.skaffold.filesystem.FilePermissions;
import com.google.cloud.tools.skaffold.filesystem.OperatingSystem;
import com.google.common.io.Resources;
import java.io.ByteArrayInputStream;
import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.io.InputStream;
import java.net.URISyntaxException;
import java.nio.charset.StandardCharsets;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;
import java.util.concurrent.ExecutionException;
import org.junit.Assert;
import org.junit.Assume;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.mockito.junit.MockitoJUnitRunner;

/** Tests for {@link Skaffold}. */
@RunWith(MockitoJUnitRunner.class)
public class SkaffoldTest {

  /**
   * Tests that {@link Skaffold} uses the correct parameters to executor {@code skaffold}. Executes
   * a script to mock the {@code skaffold} executable.
   */
  private static void verifyDeploy(String expectedStdout, String expectedStderr, String... command)
      throws InterruptedException, ExecutionException, IOException {
    InputStream stdinInputStream =
        new ByteArrayInputStream("input".getBytes(StandardCharsets.UTF_8));
    ByteArrayOutputStream stdoutOutputStream = new ByteArrayOutputStream();
    ByteArrayOutputStream stderrOutputStream = new ByteArrayOutputStream();

    List<String> expectedCommandList = new ArrayList<>();
    expectedCommandList.addAll(Arrays.asList(command));
    expectedCommandList.addAll(
        Arrays.asList("--filename", "skaffoldYaml", "--profile", "profile", "deploy"));

    int exitCode =
        new Skaffold(Skaffold.getListeningExecutorService(), command)
            .setProcessBuilderFactory(
                commandList -> {
                  Assert.assertEquals(expectedCommandList, commandList);
                  return new ProcessBuilder(commandList);
                })
            .setSkaffoldYaml(Paths.get("skaffoldYaml"))
            .setProfile("profile")
            .redirectToStdin(stdinInputStream)
            .redirectStdoutTo(stdoutOutputStream)
            .redirectStderrTo(stderrOutputStream)
            .deploy();

    Assert.assertEquals(0, exitCode);
    Assert.assertEquals(
        expectedStdout, new String(stdoutOutputStream.toByteArray(), StandardCharsets.UTF_8));
    Assert.assertEquals(
        expectedStderr, new String(stderrOutputStream.toByteArray(), StandardCharsets.UTF_8));
  }

  @Test
  public void testDeploy()
      throws URISyntaxException, IOException, InterruptedException, ExecutionException {
    Assume.assumeTrue("non-Windows test", OperatingSystem.resolve() != OperatingSystem.WINDOWS);
    Path commandSh = Paths.get(Resources.getResource("command.sh").toURI());
    Assert.assertTrue(FilePermissions.makeExecutable(commandSh));
    verifyDeploy("input\noutput\n", "error\n", commandSh.toString());
  }

  @Test
  public void testDeploy_windows()
      throws URISyntaxException, IOException, InterruptedException, ExecutionException {
    Assume.assumeTrue("Windows test", OperatingSystem.resolve() == OperatingSystem.WINDOWS);
    verifyDeploy(
        "input\r\noutput\r\n",
        "error\r\n",
        "cmd",
        "/c",
        Paths.get(Resources.getResource("command.bat").toURI()).toString());
  }
}
