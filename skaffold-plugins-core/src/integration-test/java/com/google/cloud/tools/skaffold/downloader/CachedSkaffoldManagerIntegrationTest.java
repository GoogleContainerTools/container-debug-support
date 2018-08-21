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

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.StandardCopyOption;
import org.junit.Assert;
import org.junit.Before;
import org.junit.Rule;
import org.junit.Test;
import org.junit.rules.TemporaryFolder;

/** Integration tests for {@link CachedSkaffoldManager}. */
public class CachedSkaffoldManagerIntegrationTest {

  @Rule public final TemporaryFolder temporaryFolder = new TemporaryFolder();

  private Path fakeSkaffoldLocation;
  private Path fakeDigestLocation;
  private Path latestDigestLocation;

  @Before
  public void setUp() throws IOException {
    Path temporaryDirectory = temporaryFolder.newFolder().toPath();
    fakeSkaffoldLocation = temporaryDirectory.resolve("skaffold");
    fakeDigestLocation = temporaryDirectory.resolve("digest");
    latestDigestLocation = temporaryDirectory.resolve("latestdigest");
  }

  @Test
  public void testCheckIsLatest() throws IOException {
    Assert.assertFalse(
        CachedSkaffoldManager.checkIsLatest(fakeSkaffoldLocation, fakeDigestLocation));

    // Creating the skaffold executable file is not enough.
    Files.createFile(fakeSkaffoldLocation);
    Assert.assertFalse(
        CachedSkaffoldManager.checkIsLatest(fakeSkaffoldLocation, fakeDigestLocation));

    // Creating the digest file is not enough.
    Files.createFile(fakeDigestLocation);
    Assert.assertFalse(
        CachedSkaffoldManager.checkIsLatest(fakeSkaffoldLocation, fakeDigestLocation));

    // Populating the digest file with an arbitrary byte is not enough.
    byte[] arbitraryByteArray = new byte[] {0x10};
    Files.write(fakeDigestLocation, arbitraryByteArray);
    Assert.assertFalse(
        CachedSkaffoldManager.checkIsLatest(fakeSkaffoldLocation, fakeDigestLocation));

    // Downloading the latest digest is enough.
    SkaffoldDownloader.downloadLatestDigest(latestDigestLocation);
    Files.move(latestDigestLocation, fakeDigestLocation, StandardCopyOption.REPLACE_EXISTING);
    Assert.assertTrue(
        CachedSkaffoldManager.checkIsLatest(fakeSkaffoldLocation, fakeDigestLocation));
  }

  @Test
  public void testUpdateToLatest() throws IOException, InterruptedException {
    Assert.assertFalse(
        CachedSkaffoldManager.checkIsLatest(fakeSkaffoldLocation, fakeDigestLocation));
    CachedSkaffoldManager.updateToLatest(fakeSkaffoldLocation, fakeDigestLocation);
    Assert.assertTrue(
        CachedSkaffoldManager.checkIsLatest(fakeSkaffoldLocation, fakeDigestLocation));
    Assert.assertEquals(0, new ProcessBuilder(fakeSkaffoldLocation.toString()).start().waitFor());
  }
}
