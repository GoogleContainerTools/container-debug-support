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
import java.util.Optional;
import javax.annotation.Nullable;

/** Extracts configuration from Jib. */
public interface JibAdapter {

  /** Resolved information about Jib. */
  class ResolvedJib {

    @FunctionalInterface
    interface ImageReferenceResolver {

      Optional<ImageReference> getImageReference() throws InvalidImageReferenceException;
    }

    /**
     * Instantiates for no Jib configuration found.
     *
     * @return a new {@link ResolvedJib}
     */
    static ResolvedJib jibNotFound() {
      return new ResolvedJib(false, false, null);
    }

    /**
     * Instantiates for Jib configuration found and the version is supported.
     *
     * @return a new {@link ResolvedJib}
     */
    static ResolvedJib supportedVersion(ImageReferenceResolver imageReferenceResolver) {
      return new ResolvedJib(true, true, imageReferenceResolver);
    }

    /**
     * Instantiates for Jib configuration found, but the version is not supported.
     *
     * @return a new {@link ResolvedJib}
     */
    static ResolvedJib unsupportedVersion(ImageReferenceResolver imageReferenceResolver) {
      return new ResolvedJib(true, false, imageReferenceResolver);
    }

    private final boolean hasJib;
    private final boolean isVersionSupported;
    @Nullable private final ImageReferenceResolver imageReferenceResolver;

    private ResolvedJib(
        boolean hasJib,
        boolean isVersionSupported,
        @Nullable ImageReferenceResolver imageReferenceResolver) {
      this.hasJib = hasJib;
      this.isVersionSupported = isVersionSupported;
      this.imageReferenceResolver = imageReferenceResolver;
    }

    /**
     * Returns {@code true} if Jib was resolved; {@code false} otherwise.
     *
     * @return {@code true} if Jib was resolved; {@code false} otherwise
     */
    public boolean hasJib() {
      return hasJib;
    }

    /**
     * Returns {@code true} if the current version of Jib is supported; {@code false} otherwise. If
     * the current version of Jib is not supported, {@link #getImageReference} can still be
     * attempted.
     *
     * @return {@code true} if the current version is supported; {@code false} otherwise
     */
    public boolean isVersionSupported() {
      return isVersionSupported;
    }

    /**
     * Gets the target image reference defined in the Jib plugin configuration, or {@link
     * Optional#empty} if a valid image reference configuration cannot be found.
     *
     * @return the optional target image reference
     */
    Optional<ImageReference> getImageReference() throws InvalidImageReferenceException {
      if (imageReferenceResolver == null) {
        return Optional.empty();
      }
      return imageReferenceResolver.getImageReference();
    }
  }

  /**
   * Resolves Jib configuration.
   *
   * @return the {@link ResolvedJib} with the resolved information
   */
  ResolvedJib resolveJib();
}
