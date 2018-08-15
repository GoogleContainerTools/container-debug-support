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

import static java.nio.file.attribute.PosixFilePermission.GROUP_EXECUTE;
import static java.nio.file.attribute.PosixFilePermission.OTHERS_EXECUTE;
import static java.nio.file.attribute.PosixFilePermission.OWNER_EXECUTE;

import com.google.common.annotations.VisibleForTesting;
import java.io.IOException;
import java.nio.file.FileSystem;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.attribute.PosixFilePermission;
import java.util.EnumSet;
import java.util.Set;

/** Static helpers for modifying file permissions. */
public class FilePermissions {

  /**
   * Makes a file executable (same as {@code chmod a+x}).
   *
   * @param file the file to make executable
   * @return {@code true} if successful, {@code false} otherwise (such as if the file system does
   *     not support POSIX permissions)
   * @throws IOException if an I/O exception occurred
   */
  public static boolean makeExecutable(Path file) throws IOException {
    try {
      Set<PosixFilePermission> executableFilePermissions = Files.getPosixFilePermissions(file);
      executableFilePermissions.addAll(EnumSet.of(OWNER_EXECUTE, GROUP_EXECUTE, OTHERS_EXECUTE));
      Files.setPosixFilePermissions(file, executableFilePermissions);
      return true;

    } catch (UnsupportedOperationException ex) {
      // File system does not support POSIX.
      return false;
    }
  }

  @VisibleForTesting
  public static boolean isFilesystemPosix(FileSystem filesystem) {
    return filesystem.supportedFileAttributeViews().contains("posix");
  }

  private FilePermissions() {}
}
