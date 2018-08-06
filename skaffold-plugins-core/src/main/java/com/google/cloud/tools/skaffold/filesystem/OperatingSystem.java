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

import java.util.Locale;
import java.util.Properties;

/** Gets information about the operating system. */
public enum OperatingSystem {
  LINUX,
  MAC_OS,
  WINDOWS;

  /**
   * Resolves the operating system type. *
   *
   * @return the {@link OperatingSystem}
   */
  public static OperatingSystem resolve() {
    return resolve(System.getProperties());
  }

  static OperatingSystem resolve(Properties properties) {
    String rawOsName = properties.getProperty("os.name");
    String osName = rawOsName.toLowerCase(Locale.ENGLISH);

    if (osName.contains("linux")) {
      return LINUX;
    }
    if (osName.contains("windows")) {
      return WINDOWS;
    }
    if (osName.contains("mac") || osName.contains("darwin")) {
      return MAC_OS;
    }

    throw new IllegalStateException("Unknown OS: " + rawOsName);
  }
}
