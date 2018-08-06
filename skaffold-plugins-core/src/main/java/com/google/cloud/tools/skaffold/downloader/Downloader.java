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
import java.net.URL;
import java.nio.channels.Channels;
import java.nio.channels.FileChannel;
import java.nio.file.Path;
import java.nio.file.StandardOpenOption;

/** Downloads files. */
public class Downloader {

  private final URL url;

  /**
   * Instantiates with a {@link URL} to download from.
   *
   * @param url the {@link URL} to download
   */
  public Downloader(URL url) {
    this.url = url;
  }

  /**
   * Downloads to a destination file.
   *
   * @param destination the destination file to download to
   * @return the size of the downloaded contents
   * @throws IOException if an I/O exception occurred during the download process
   */
  public long download(Path destination) throws IOException {
    return FileChannel.open(destination, StandardOpenOption.WRITE, StandardOpenOption.CREATE)
        .transferFrom(Channels.newChannel(url.openStream()), 0, Long.MAX_VALUE);
  }
}
