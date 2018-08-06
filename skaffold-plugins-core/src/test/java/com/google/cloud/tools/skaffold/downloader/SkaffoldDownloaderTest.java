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
import org.junit.Assert;
import org.junit.Test;

/** Test for {@link SkaffoldDownloader}. */
public class SkaffoldDownloaderTest {

  @Test
  public void testGetUrl() {
    Assert.assertEquals(
        "https://storage.googleapis.com/skaffold/releases/version/skaffold-linux-amd64",
        SkaffoldDownloader.getUrl("version", OperatingSystem.LINUX));
    Assert.assertEquals(
        "https://storage.googleapis.com/skaffold/releases/someversion/skaffold-darwin-amd64",
        SkaffoldDownloader.getUrl("someversion", OperatingSystem.MAC_OS));
    Assert.assertEquals(
        "https://storage.googleapis.com/skaffold/releases/latest/skaffold-windows-amd64.exe",
        SkaffoldDownloader.getUrl("latest", OperatingSystem.WINDOWS));
  }
}
