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

import com.google.common.io.CharStreams;
import com.google.common.io.Resources;
import java.io.IOException;
import java.io.InputStreamReader;
import java.nio.charset.StandardCharsets;
import java.nio.file.Path;
import org.junit.Assert;
import org.junit.Rule;
import org.junit.Test;
import org.junit.rules.TemporaryFolder;

/** Integration tests for {@link Downloader}. */
public class DownloaderIntegationTest {

  @Rule public TemporaryFolder temporaryFolder = new TemporaryFolder();

  @Test
  public void testDownload() throws IOException, InterruptedException {
    // Downloads a script that says "hello".
    Path temporaryFile = temporaryFolder.newFolder().toPath().resolve("hello.sh");
    new Downloader(Resources.getResource("helloScript.sh")).download(temporaryFile);
    Assert.assertTrue(temporaryFile.toFile().setExecutable(true));

    // Runs the downloaded script.
    Process helloProcess = new ProcessBuilder(temporaryFile.toString()).start();
    String stdout =
        CharStreams.toString(
            new InputStreamReader(helloProcess.getInputStream(), StandardCharsets.UTF_8));
    Assert.assertEquals(0, helloProcess.waitFor());
    Assert.assertEquals("hello\n", stdout);
  }
}
