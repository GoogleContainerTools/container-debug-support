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
import com.google.cloud.tools.skaffold.jib.JibAdapter.ResolvedJib.ImageReferenceResolver;
import java.util.Optional;
import java.util.Properties;
import org.apache.maven.model.Plugin;
import org.apache.maven.project.MavenProject;
import org.codehaus.plexus.util.xml.Xpp3Dom;

/** Fetches the target image reference from {@code jib-maven-plugin}. */
public class JibMavenPluginAdapter implements JibAdapter {

  /** Adapter for {@code jib-maven-plugin} versions {@code 0.9.x}. */
  private static class AdapterBeta9 implements ImageReferenceResolver {

    private final Properties mavenProperties;
    private final Xpp3Dom jibConfiguration;

    private AdapterBeta9(Properties mavenProperties, Xpp3Dom jibConfiguration) {
      this.mavenProperties = mavenProperties;
      this.jibConfiguration = jibConfiguration;
    }

    @Override
    public Optional<ImageReference> getImageReference() throws InvalidImageReferenceException {
      String imageProperty = mavenProperties.getProperty("image");
      if (imageProperty != null) {
        return Optional.of(ImageReference.parse(imageProperty));
      }

      Xpp3Dom to = jibConfiguration.getChild("to");
      if (to == null) {
        return Optional.empty();
      }
      Xpp3Dom image = to.getChild("image");
      if (image == null) {
        return Optional.empty();
      }
      return Optional.of(ImageReference.parse(image.getValue()));
    }
  }

  private static final String JIB_MAVEN_PLUGIN_KEY = "com.google.cloud.tools:jib-maven-plugin";

  private final MavenProject mavenProject;

  public JibMavenPluginAdapter(MavenProject mavenProject) {
    this.mavenProject = mavenProject;
  }

  @Override
  public ResolvedJib resolveJib() {
    Plugin jibPlugin = mavenProject.getPlugin(JIB_MAVEN_PLUGIN_KEY);
    if (jibPlugin == null) {
      return ResolvedJib.jibNotFound();
    }

    Xpp3Dom jibConfiguration = (Xpp3Dom) jibPlugin.getConfiguration();
    if (jibConfiguration == null) {
      return ResolvedJib.jibNotFound();
    }

    String version = jibPlugin.getVersion();
    if (version.startsWith("0.9.")) {
      return ResolvedJib.supportedVersion(
          new AdapterBeta9(mavenProject.getProperties(), jibConfiguration));
    }

    // Defaulting to latest adapter.
    return ResolvedJib.unsupportedVersion(
        new AdapterBeta9(mavenProject.getProperties(), jibConfiguration));
  }
}
