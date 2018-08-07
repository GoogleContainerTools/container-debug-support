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

import com.google.common.annotations.VisibleForTesting;
import java.io.BufferedInputStream;
import java.io.IOException;
import java.io.InputStream;
import java.net.URL;
import java.net.URLConnection;
import java.nio.channels.Channels;
import java.nio.channels.FileChannel;
import java.nio.channels.ReadableByteChannel;
import java.nio.file.Path;
import java.nio.file.StandardOpenOption;

/** Downloads files. */
public class Downloader {

  /**
   * Downloads to a destination file.
   *
   * @param url the {@link URL} to download
   * @param destination the destination file to download to
   * @return the size of the downloaded contents, or -1 if {@code Content-Length} is unknown and
   *     thus nothing downloaded
   * @throws IOException if an I/O exception occurred during the download process
   */
  public static long download(URL url, Path destination) throws IOException {
    return download(url, destination, Long.MAX_VALUE);
  }

  @VisibleForTesting
  static long download(URL url, Path destination, long chunkSize) throws IOException {
    URLConnection connection = url.openConnection();
    try (FileChannel fileChannel =
            FileChannel.open(destination, StandardOpenOption.WRITE, StandardOpenOption.CREATE);
        InputStream connectionInputStream = new BufferedInputStream(connection.getInputStream());
        ReadableByteChannel urlChannel = Channels.newChannel(connectionInputStream)) {
      long totalSize = connection.getContentLengthLong();
      while (fileChannel.position() < totalSize) {
        fileChannel.position(
            fileChannel.position()
                + fileChannel.transferFrom(urlChannel, fileChannel.position(), chunkSize));
      }
      return totalSize;
    }
  }

  private Downloader() {}
}
