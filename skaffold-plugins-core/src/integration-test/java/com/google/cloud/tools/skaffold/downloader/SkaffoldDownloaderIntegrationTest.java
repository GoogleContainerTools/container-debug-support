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
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;
import javax.xml.bind.DatatypeConverter;
import org.junit.Assert;
import org.junit.Rule;
import org.junit.Test;
import org.junit.rules.TemporaryFolder;

/** Integration tests for {@link SkaffoldDownloader}. */
public class SkaffoldDownloaderIntegrationTest {

  @Rule public TemporaryFolder temporaryFolder = new TemporaryFolder();

  @Test
  public void testDownloadLatest()
      throws IOException, InterruptedException, NoSuchAlgorithmException {
    Path temporarySkaffoldExecutable = temporaryFolder.newFile().toPath();
    Assert.assertTrue(SkaffoldDownloader.downloadLatest(temporarySkaffoldExecutable));
    Process skaffoldProcess = new ProcessBuilder(temporarySkaffoldExecutable.toString()).start();
    Assert.assertEquals(0, skaffoldProcess.waitFor());

    // Downloads and checks that the digest matches.
    Path temporarySkaffoldExecutableDigest = temporaryFolder.newFile().toPath();
    SkaffoldDownloader.downloadLatestDigest(temporarySkaffoldExecutableDigest);

    MessageDigest messageDigest = MessageDigest.getInstance("SHA-256");
    byte[] expectedDigest = messageDigest.digest(Files.readAllBytes(temporarySkaffoldExecutable));

    String receivedDigestHex =
        new String(Files.readAllBytes(temporarySkaffoldExecutableDigest), StandardCharsets.UTF_8)
            .substring(0, 64);
    byte[] receivedDigest = DatatypeConverter.parseHexBinary(receivedDigestHex);

    Assert.assertArrayEquals(expectedDigest, receivedDigest);
  }
}
