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
import com.google.common.annotations.VisibleForTesting;
import java.io.IOException;
import java.net.MalformedURLException;
import java.net.URL;
import java.nio.file.Path;

/** Downloads {@code skaffold} executable. */
public class SkaffoldDownloader {

  /**
   * Downloads the latest {@code skaffold} release.
   *
   * @return {@code true} if the destination file could be set to executable; {@code false}
   *     otherwise
   * @throws MalformedURLException if the URL to download from is malformed
   */
  public static boolean downloadLatest(Path destination) throws IOException {
    return download(new URL(getUrl("latest", OperatingSystem.resolve())), destination);
  }

  /**
   * Downloads to the {@code destination}.
   *
   * @param destination the destination file to download {@code skaffold} to
   * @return {@code true} if the destination file could be set to executable; {@code false}
   *     otherwise
   * @throws IOException if an I/O exception occurs during download
   */
  @VisibleForTesting
  static boolean download(URL url, Path destination) throws IOException {
    if (Downloader.download(url, destination) == -1) {
      throw new IOException("Could not get size of skaffold binary to download");
    }
    return destination.toFile().setExecutable(true);
  }

  /**
   * Resolves the correct URL to download {@code skaffold} from based on the version to download and
   * the operating system.
   *
   * @param version the version to download (use {@code latest} for the latest version)
   * @param operatingSystem the {@link OperatingSystem}
   * @return the URL to download from
   * @see <a
   *     href="https://github.com/GoogleContainerTools/skaffold/releases">https://github.com/GoogleContainerTools/skaffold/releases</a>
   */
  @VisibleForTesting
  static String getUrl(String version, OperatingSystem operatingSystem) {
    String base = "https://storage.googleapis.com/skaffold/releases/" + version + "/";

    switch (operatingSystem) {
      case LINUX:
        return base + "skaffold-linux-amd64";

      case MAC_OS:
        return base + "skaffold-darwin-amd64";

      case WINDOWS:
        return base + "skaffold-windows-amd64.exe";

      default:
        throw new IllegalStateException("unreachable");
    }
  }

  private SkaffoldDownloader() {}
}
