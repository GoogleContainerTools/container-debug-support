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

package com.google.cloud.tools.skaffold.downloader;

import com.google.cloud.tools.skaffold.filesystem.OperatingSystem;
import com.google.common.io.CharStreams;
import com.google.common.io.Resources;
import java.io.IOException;
import java.io.InputStreamReader;
import java.net.URL;
import java.nio.charset.StandardCharsets;
import java.nio.file.Path;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;
import org.junit.Assert;
import org.junit.Assume;
import org.junit.Rule;
import org.junit.Test;
import org.junit.rules.TemporaryFolder;

/** Integration tests for {@link Downloader}. */
public class DownloaderIntegationTest {

  private static String downloadAndRun(URL url, Path destination, String... command)
      throws IOException, InterruptedException {
    // Downloads a script that says "hello".
    Downloader.download(url, destination, 1);
    Assert.assertTrue(destination.toFile().setExecutable(true));

    // Runs the downloaded script.
    List<String> commandList = new ArrayList<>(Arrays.asList(command));
    commandList.add(destination.toString());
    Process process = new ProcessBuilder(commandList).start();
    String stdout =
        CharStreams.toString(
            new InputStreamReader(process.getInputStream(), StandardCharsets.UTF_8));
    String stderr =
        CharStreams.toString(
            new InputStreamReader(process.getErrorStream(), StandardCharsets.UTF_8));
    Assert.assertEquals("", stderr);
    Assert.assertEquals(0, process.waitFor());
    return stdout;
  }

  @Rule public TemporaryFolder temporaryFolder = new TemporaryFolder();

  @Test
  public void testDownload() throws IOException, InterruptedException {
    Assume.assumeTrue("non-Windows test", OperatingSystem.resolve() != OperatingSystem.WINDOWS);

    Assert.assertEquals(
        "hello\n",
        downloadAndRun(
            Resources.getResource("helloScript.sh"),
            temporaryFolder.newFolder().toPath().resolve("hello.sh"),
            System.getenv("SHELL")));
  }

  @Test
  public void testDownload_windows() throws IOException, InterruptedException {
    Assume.assumeTrue("Windows test", OperatingSystem.resolve() == OperatingSystem.WINDOWS);

    Assert.assertEquals(
        "hello\n",
        downloadAndRun(
            Resources.getResource("helloScript.bat"),
            temporaryFolder.newFolder().toPath().resolve("hello.bat"),
            "cmd",
            "/c",
            "/q"));
  }
}
