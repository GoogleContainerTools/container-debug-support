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

import java.io.ByteArrayInputStream;
import java.io.IOException;
import java.io.InputStream;
import java.net.URL;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import org.junit.Assert;
import org.junit.Rule;
import org.junit.Test;
import org.junit.rules.TemporaryFolder;
import org.junit.runner.RunWith;
import org.mockito.Mock;
import org.mockito.Mockito;
import org.mockito.junit.MockitoJUnitRunner;

/** Test for {@link Downloader}. */
@RunWith(MockitoJUnitRunner.class)
public class DownloaderTest {

  @Rule public TemporaryFolder temporaryFolder = new TemporaryFolder();

  @Mock private URL mockURL;

  @Test
  public void testDownload_newFile() throws IOException {
    downloadToFile(temporaryFolder.newFolder().toPath().resolve("nonexistent"));
  }

  @Test
  public void testDownload_overwrite() throws IOException {
    downloadToFile(temporaryFolder.newFile().toPath());
  }

  private void downloadToFile(Path destination) throws IOException {
    String expectedDownloadContents = "downloaded file contents";
    InputStream fakeInputStream =
        new ByteArrayInputStream(expectedDownloadContents.getBytes(StandardCharsets.UTF_8));
    Mockito.when(mockURL.openStream()).thenReturn(fakeInputStream);

    long size = new Downloader(mockURL).download(destination);
    Assert.assertEquals(expectedDownloadContents.getBytes(StandardCharsets.UTF_8).length, size);
    Assert.assertEquals(
        expectedDownloadContents,
        new String(Files.readAllBytes(destination), StandardCharsets.UTF_8));
  }
}
