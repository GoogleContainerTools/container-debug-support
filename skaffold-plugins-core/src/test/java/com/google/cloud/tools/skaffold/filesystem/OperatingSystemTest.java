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

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.fail;

import java.util.Properties;
import org.junit.jupiter.api.Test;

/** Tests for {@link OperatingSystem}. */
class OperatingSystemTest {

  @Test
  void testLinux() {
    Properties fakeProperties = new Properties();
    fakeProperties.setProperty("os.name", "os is LiNuX");
    assertEquals(OperatingSystem.LINUX, OperatingSystem.resolve(fakeProperties));
  }

  @Test
  void testMacOs_mac() {
    Properties fakeProperties = new Properties();
    fakeProperties.setProperty("os.name", "os is mAc");
    assertEquals(OperatingSystem.MAC_OS, OperatingSystem.resolve(fakeProperties));
  }

  @Test
  void testMacOs_darwin() {
    Properties fakeProperties = new Properties();
    fakeProperties.setProperty("os.name", "os is DaRwIn");
    assertEquals(OperatingSystem.MAC_OS, OperatingSystem.resolve(fakeProperties));
  }

  @Test
  void testWindows() {
    Properties fakeProperties = new Properties();
    fakeProperties.setProperty("os.name", "os is WiNdOwS");
    assertEquals(OperatingSystem.WINDOWS, OperatingSystem.resolve(fakeProperties));
  }

  @Test
  void testUnknown() {
    Properties fakeProperties = new Properties();
    fakeProperties.setProperty("os.name", "UnKnOwN");
    try {
      OperatingSystem.resolve(fakeProperties);
      fail("Resolve should have failed");

    } catch (IllegalStateException ex) {
      assertEquals("Unknown OS: UnKnOwN", ex.getMessage());
    }
  }
}
