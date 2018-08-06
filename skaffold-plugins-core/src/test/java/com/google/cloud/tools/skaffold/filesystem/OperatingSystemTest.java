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

import java.util.Properties;
import org.junit.Assert;
import org.junit.Test;

/** Tests for {@link OperatingSystem}. */
public class OperatingSystemTest {

  @Test
  public void testLinux() {
    Properties fakeProperties = new Properties();
    fakeProperties.setProperty("os.name", "os is LiNuX");
    Assert.assertEquals(OperatingSystem.LINUX, OperatingSystem.resolve(fakeProperties));
  }

  @Test
  public void testMacOs() {
    Properties fakeProperties = new Properties();
    fakeProperties.setProperty("os.name", "os is mAc or DaRwIn");
    Assert.assertEquals(OperatingSystem.MAC_OS, OperatingSystem.resolve(fakeProperties));
  }

  @Test
  public void testWindows() {
    Properties fakeProperties = new Properties();
    fakeProperties.setProperty("os.name", "os is WiNdOwS");
    Assert.assertEquals(OperatingSystem.WINDOWS, OperatingSystem.resolve(fakeProperties));
  }

  @Test
  public void testUnknown() {
    Properties fakeProperties = new Properties();
    fakeProperties.setProperty("os.name", "UnKnOwN");
    try {
      OperatingSystem.resolve(fakeProperties);
      Assert.fail("Resolve should have failed");

    } catch (IllegalStateException ex) {
      Assert.assertEquals("Unknown OS: UnKnOwN", ex.getMessage());
    }
  }
}
