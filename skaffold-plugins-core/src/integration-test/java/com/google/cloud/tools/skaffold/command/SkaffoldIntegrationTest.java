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

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.attribute.FileTime;
import org.junit.Assert;
import org.junit.Rule;
import org.junit.Test;
import org.junit.rules.TemporaryFolder;

/** Integration tests for {@link Skaffold}. */
public class SkaffoldIntegrationTest {

  @Rule public final TemporaryFolder temporaryFolder = new TemporaryFolder();

  @Test
  public void testEnsureSkaffoldIsLatestVersion() throws IOException {
    Path cache = temporaryFolder.newFolder().toPath();
    Path cachedSkaffoldLocation = cache.resolve("skaffold");
    Path cachedSkaffoldDigestLocation = cache.resolve("skaffold.sha256");

    Skaffold.ensureSkaffoldIsLatestVersion(cachedSkaffoldLocation, cachedSkaffoldDigestLocation);

    Assert.assertTrue(Files.exists(cachedSkaffoldLocation));

    // Checks that not having an up-to-date digest causes a re-download.
    byte[] arbitraryByteArray = new byte[] {0x10};
    Files.write(cachedSkaffoldDigestLocation, arbitraryByteArray);
    FileTime originalLastModifiedTime = Files.getLastModifiedTime(cachedSkaffoldLocation);
    Skaffold.ensureSkaffoldIsLatestVersion(cachedSkaffoldLocation, cachedSkaffoldDigestLocation);
    FileTime newLastModifiedTime = Files.getLastModifiedTime(cachedSkaffoldLocation);
    Assert.assertTrue(newLastModifiedTime.compareTo(originalLastModifiedTime) > 0);

    // Checks that will not re-download if has up-to-date digest.
    Skaffold.ensureSkaffoldIsLatestVersion(cachedSkaffoldLocation, cachedSkaffoldDigestLocation);
    FileTime anotherNewLastModifiedTime = Files.getLastModifiedTime(cachedSkaffoldLocation);
    Assert.assertEquals(newLastModifiedTime, anotherNewLastModifiedTime);
  }
}
