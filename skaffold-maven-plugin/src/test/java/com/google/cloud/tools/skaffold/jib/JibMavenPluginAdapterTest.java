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
import java.io.ByteArrayInputStream;
import java.io.IOException;
import java.io.InputStreamReader;
import java.nio.charset.StandardCharsets;
import java.util.Optional;
import java.util.Properties;
import org.apache.maven.model.Plugin;
import org.apache.maven.project.MavenProject;
import org.codehaus.plexus.util.xml.Xpp3Dom;
import org.codehaus.plexus.util.xml.Xpp3DomBuilder;
import org.codehaus.plexus.util.xml.pull.XmlPullParserException;
import org.junit.Assert;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.mockito.Mock;
import org.mockito.Mockito;
import org.mockito.junit.MockitoJUnitRunner;

/** Tests for {@link JibMavenPluginAdapter}. */
@RunWith(MockitoJUnitRunner.class)
public class JibMavenPluginAdapterTest {

  private static Xpp3Dom buildXpp3Dom(String xml) throws IOException, XmlPullParserException {
    return Xpp3DomBuilder.build(
        new InputStreamReader(
            new ByteArrayInputStream(xml.getBytes(StandardCharsets.UTF_8)),
            StandardCharsets.UTF_8));
  }

  @Mock private MavenProject mockMavenProject;
  @Mock private Plugin mockJibPlugin;

  private JibAdapter testJibAdapter;

  @Before
  public void setUp() {
    testJibAdapter = new JibMavenPluginAdapter(mockMavenProject);
  }

  @Test
  public void testResolveJib_notFound() throws InvalidImageReferenceException {
    ResolvedJib resolvedJib = testJibAdapter.resolveJib();

    Assert.assertFalse(resolvedJib.hasJib());
    Assert.assertFalse(resolvedJib.isVersionSupported());
    Assert.assertFalse(resolvedJib.getImageReference().isPresent());
  }

  @Test
  public void testResolveJib_configurationNotFound() throws InvalidImageReferenceException {
    Mockito.when(mockMavenProject.getPlugin("com.google.cloud.tools:jib-maven-plugin"))
        .thenReturn(mockJibPlugin);

    ResolvedJib resolvedJib = testJibAdapter.resolveJib();
    Mockito.verify(mockJibPlugin).getConfiguration();

    Assert.assertFalse(resolvedJib.hasJib());
    Assert.assertFalse(resolvedJib.isVersionSupported());
    Assert.assertFalse(resolvedJib.getImageReference().isPresent());
  }

  @Test
  public void testResolveJib_beta9_imageProperty() throws InvalidImageReferenceException {
    Mockito.when(mockMavenProject.getPlugin("com.google.cloud.tools:jib-maven-plugin"))
        .thenReturn(mockJibPlugin);

    Mockito.when(mockJibPlugin.getConfiguration()).thenReturn(Mockito.mock(Xpp3Dom.class));
    Mockito.when(mockJibPlugin.getVersion()).thenReturn("0.9.0");

    Properties mavenProperties = new Properties();
    mavenProperties.setProperty("image", "image");
    Mockito.when(mockMavenProject.getProperties()).thenReturn(mavenProperties);

    ResolvedJib resolvedJib = testJibAdapter.resolveJib();
    Assert.assertTrue(resolvedJib.hasJib());
    Assert.assertTrue(resolvedJib.isVersionSupported());
    Optional<ImageReference> optionalImageReference = resolvedJib.getImageReference();
    Assert.assertTrue(optionalImageReference.isPresent());
    Assert.assertEquals("image", optionalImageReference.get().toString());
  }

  @Test
  public void testResolveJib_beta9_noTo()
      throws InvalidImageReferenceException, IOException, XmlPullParserException {
    Mockito.when(mockMavenProject.getPlugin("com.google.cloud.tools:jib-maven-plugin"))
        .thenReturn(mockJibPlugin);
    Mockito.when(mockMavenProject.getProperties()).thenReturn(new Properties());

    Mockito.when(mockJibPlugin.getConfiguration())
        .thenReturn(buildXpp3Dom("<configuration></configuration>"));
    Mockito.when(mockJibPlugin.getVersion()).thenReturn("0.9.0");

    ResolvedJib resolvedJib = testJibAdapter.resolveJib();
    Assert.assertTrue(resolvedJib.hasJib());
    Assert.assertTrue(resolvedJib.isVersionSupported());
    Assert.assertFalse(resolvedJib.getImageReference().isPresent());
  }

  @Test
  public void testResolveJib_beta9_hasTo_noImage()
      throws InvalidImageReferenceException, IOException, XmlPullParserException {
    Mockito.when(mockMavenProject.getPlugin("com.google.cloud.tools:jib-maven-plugin"))
        .thenReturn(mockJibPlugin);
    Mockito.when(mockMavenProject.getProperties()).thenReturn(new Properties());

    Mockito.when(mockJibPlugin.getConfiguration())
        .thenReturn(buildXpp3Dom("<configuration><to></to></configuration>"));
    Mockito.when(mockJibPlugin.getVersion()).thenReturn("0.9.0");

    ResolvedJib resolvedJib = testJibAdapter.resolveJib();
    Assert.assertTrue(resolvedJib.hasJib());
    Assert.assertTrue(resolvedJib.isVersionSupported());
    Assert.assertFalse(resolvedJib.getImageReference().isPresent());
  }

  @Test
  public void testResolveJib_beta9_hasImage()
      throws IOException, XmlPullParserException, InvalidImageReferenceException {
    Mockito.when(mockMavenProject.getPlugin("com.google.cloud.tools:jib-maven-plugin"))
        .thenReturn(mockJibPlugin);
    Mockito.when(mockMavenProject.getProperties()).thenReturn(new Properties());

    Mockito.when(mockJibPlugin.getConfiguration())
        .thenReturn(buildXpp3Dom("<configuration><to><image>image</image></to></configuration>"));
    Mockito.when(mockJibPlugin.getVersion()).thenReturn("0.9.0");

    ResolvedJib resolvedJib = testJibAdapter.resolveJib();
    Assert.assertTrue(resolvedJib.hasJib());
    Assert.assertTrue(resolvedJib.isVersionSupported());
    Optional<ImageReference> optionalImageReference = resolvedJib.getImageReference();
    Assert.assertTrue(optionalImageReference.isPresent());
    Assert.assertEquals("image", optionalImageReference.get().toString());
  }

  @Test
  public void testResolveJib_unsupportedVersion()
      throws IOException, XmlPullParserException, InvalidImageReferenceException {
    Mockito.when(mockMavenProject.getPlugin("com.google.cloud.tools:jib-maven-plugin"))
        .thenReturn(mockJibPlugin);
    Mockito.when(mockMavenProject.getProperties()).thenReturn(new Properties());

    Mockito.when(mockJibPlugin.getConfiguration())
        .thenReturn(buildXpp3Dom("<configuration><to><image>image</image></to></configuration>"));
    Mockito.when(mockJibPlugin.getVersion()).thenReturn("0.1.0");

    ResolvedJib resolvedJib = testJibAdapter.resolveJib();
    Assert.assertTrue(resolvedJib.hasJib());
    Assert.assertFalse(resolvedJib.isVersionSupported());
    Optional<ImageReference> optionalImageReference = resolvedJib.getImageReference();
    Assert.assertTrue(optionalImageReference.isPresent());
    Assert.assertEquals("image", optionalImageReference.get().toString());
  }
}
