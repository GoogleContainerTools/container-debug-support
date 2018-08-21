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

import com.google.cloud.tools.skaffold.filesystem.UserCacheHome;
import com.google.common.annotations.VisibleForTesting;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.Arrays;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/** Static helpers to manage a cached {@code skaffold} executable. */
public class CachedSkaffoldManager {

  private static final Logger logger = LoggerFactory.getLogger(CachedSkaffoldManager.class);

  /** The location to store {@code skaffold} if auto-downloading it. */
  private static final Path CACHED_SKAFFOLD_LOCATION =
      UserCacheHome.getCacheHome().resolve("skaffold");

  /** The location to store the digest for {@code skaffold} if auto-downloading it. */
  private static final Path CACHED_SKAFFOLD_DIGEST_LOCATION =
      CACHED_SKAFFOLD_LOCATION.resolveSibling("skaffold.sha256");

  /**
   * Checks if the cached {@code skaffold} executable is the latest version.
   *
   * @return {@code true} if is latest; {@code false} otherwise
   * @throws IOException if an I/O exception occurs
   */
  public static boolean checkIsLatest() throws IOException {
    return checkIsLatest(CACHED_SKAFFOLD_LOCATION, CACHED_SKAFFOLD_DIGEST_LOCATION);
  }

  /**
   * Updates the cached {@code skaffold} executable to the latest version. Check with {@link
   * #checkIsLatest} beforehand to avoid updating when already at latest version.
   *
   * @throws IOException if an I/O exception occurs
   */
  public static void updateToLatest() throws IOException {
    updateToLatest(CACHED_SKAFFOLD_LOCATION, CACHED_SKAFFOLD_DIGEST_LOCATION);
  }

  /**
   * Gets the cached {@code skaffold} executable.
   *
   * @return the path to {@code skaffold}
   */
  public static Path getCachedSkaffold() {
    return CACHED_SKAFFOLD_LOCATION;
  }

  @VisibleForTesting
  static boolean checkIsLatest(Path cachedSkaffoldLocation, Path cachedSkaffoldDigestLocation)
      throws IOException {
    if (!Files.exists(cachedSkaffoldLocation)) {
      return false;
    }
    // Checks if the digest is up-to-date and redownloads skaffold if not.
    if (!Files.exists(cachedSkaffoldDigestLocation)) {
      return false;
    }
    byte[] storedDigest = Files.readAllBytes(cachedSkaffoldDigestLocation);
    byte[] latestDigest = Files.readAllBytes(downloadLatestDigest());
    if (!Arrays.equals(storedDigest, latestDigest)) {
      return false;
    }
    logger.debug("Cached skaffold is latest version");
    return true;
  }

  @VisibleForTesting
  static void updateToLatest(Path cachedSkaffoldLocation, Path cachedSkaffoldDigestLocation)
      throws IOException {
    Files.copy(downloadLatestDigest(), cachedSkaffoldDigestLocation);
    SkaffoldDownloader.downloadLatest(cachedSkaffoldLocation);
  }

  private static Path downloadLatestDigest() throws IOException {
    Path temporaryDigestFile = Files.createTempFile("", "");
    temporaryDigestFile.toFile().deleteOnExit();
    logger.debug("Downloading latest skaffold release digest");
    SkaffoldDownloader.downloadLatestDigest(temporaryDigestFile);
    return temporaryDigestFile;
  }

  private CachedSkaffoldManager() {}
}
