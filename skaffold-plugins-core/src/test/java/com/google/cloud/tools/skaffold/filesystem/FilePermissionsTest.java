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

package com.google.cloud.tools.skaffold.filesystem;

import com.google.common.io.CharStreams;
import com.google.common.io.Resources;
import java.io.IOException;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.net.URISyntaxException;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import org.junit.Assert;
import org.junit.Assume;
import org.junit.Rule;
import org.junit.Test;
import org.junit.rules.TemporaryFolder;

/** Tests for {@link FilePermissions}. */
public class FilePermissionsTest {

  @Rule public final TemporaryFolder temporaryFolder = new TemporaryFolder();

  @Test
  public void testMakeExecutable() throws URISyntaxException, IOException, InterruptedException {
    Assume.assumeTrue(
        "only for posix filesystems",
        FilePermissions.isFilesystemPosix(temporaryFolder.getRoot().toPath().getFileSystem()));

    Path nonExecutableSh = Paths.get(Resources.getResource("non-executable.sh").toURI());

    try {
      new ProcessBuilder(nonExecutableSh.toString()).start();
      Assert.fail("executing non-executable should fail");

    } catch (IOException ex) {
      // pass
    }

    Path executableSh = temporaryFolder.newFolder().toPath().resolve("executable.sh");
    Files.copy(nonExecutableSh, executableSh);
    FilePermissions.makeExecutable(executableSh);

    Process process = new ProcessBuilder(executableSh.toString()).start();
    try (InputStream stdout = process.getInputStream();
        InputStreamReader stdoutReader = new InputStreamReader(stdout, StandardCharsets.UTF_8)) {
      Assert.assertEquals("hello\n", CharStreams.toString(stdoutReader));
    }
    Assert.assertEquals(0, process.waitFor());
  }
}
