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

package com.google.cloud.tools.skaffold.jib;

import com.google.cloud.tools.skaffold.image.ImageReference;
import com.google.cloud.tools.skaffold.image.InvalidImageReferenceException;
import com.google.cloud.tools.skaffold.jib.JibAdapter.ResolvedJib;
import com.google.cloud.tools.skaffold.jib.JibAdapter.ResolvedJib.ImageReferenceResolver;
import java.util.Optional;
import org.junit.Assert;
import org.junit.Test;

/** Tests for {@link JibAdapter}. */
public class JibAdapterTest {

  private static final ImageReferenceResolver IMAGE_REFERENCE_RESOLVER =
      () -> Optional.of(ImageReference.of("registry", "repository", "tag"));

  @Test
  public void testResolvedJib_jibNotFound() throws InvalidImageReferenceException {
    ResolvedJib resolvedjib = ResolvedJib.jibNotFound();
    Assert.assertFalse(resolvedjib.hasJib());
    Assert.assertFalse(resolvedjib.isVersionSupported());
    Assert.assertFalse(resolvedjib.getImageReference().isPresent());
  }

  @Test
  public void testResolvedJib_supportedVersion() throws InvalidImageReferenceException {
    ResolvedJib resolvedjib = ResolvedJib.supportedVersion(IMAGE_REFERENCE_RESOLVER);
    Assert.assertTrue(resolvedjib.hasJib());
    Assert.assertTrue(resolvedjib.isVersionSupported());
    Optional<ImageReference> optionalImageReference = resolvedjib.getImageReference();
    Assert.assertTrue(optionalImageReference.isPresent());
    Assert.assertEquals("registry/repository:tag", optionalImageReference.get().toString());
  }

  @Test
  public void testResolvedJib_unsupportedVersion() throws InvalidImageReferenceException {
    ResolvedJib resolvedjib = ResolvedJib.unsupportedVersion(IMAGE_REFERENCE_RESOLVER);
    Assert.assertTrue(resolvedjib.hasJib());
    Assert.assertFalse(resolvedjib.isVersionSupported());
    Optional<ImageReference> optionalImageReference = resolvedjib.getImageReference();
    Assert.assertTrue(optionalImageReference.isPresent());
    Assert.assertEquals("registry/repository:tag", optionalImageReference.get().toString());
  }
}
