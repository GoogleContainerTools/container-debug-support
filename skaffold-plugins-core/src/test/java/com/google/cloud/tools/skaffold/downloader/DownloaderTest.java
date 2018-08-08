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
import java.net.URLConnection;
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

  @Rule public final TemporaryFolder temporaryFolder = new TemporaryFolder();

  @Mock private URL mockURL;
  @Mock private URLConnection mockURLConnection;

  @Test
  public void testDownload_newFile() throws IOException {
    downloadToFile(temporaryFolder.newFolder().toPath().resolve("nonexistent"));
  }

  @Test
  public void testDownload_overwrite() throws IOException {
    downloadToFile(temporaryFolder.newFile().toPath());
  }

  private void downloadToFile(Path destination) throws IOException {
    byte[] expectedDownloadContents = new byte[] {0x11, 0x22, 0x33, 0x44, 0x55};
    long expectedSize = expectedDownloadContents.length;
    InputStream fakeInputStream = new ByteArrayInputStream(expectedDownloadContents);
    Mockito.when(mockURL.openConnection()).thenReturn(mockURLConnection);
    Mockito.when(mockURLConnection.getInputStream()).thenReturn(fakeInputStream);
    Mockito.when(mockURLConnection.getContentLengthLong()).thenReturn(expectedSize);

    long size = Downloader.download(mockURL, destination, 1);
    Assert.assertEquals(expectedSize, size);
    Assert.assertArrayEquals(expectedDownloadContents, Files.readAllBytes(destination));
  }
}
