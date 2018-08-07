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

package com.google.cloud.tools.skaffold.command;

import com.google.cloud.tools.skaffold.downloader.SkaffoldDownloader;
import com.google.cloud.tools.skaffold.filesystem.UserCacheHome;
import com.google.common.annotations.VisibleForTesting;
import com.google.common.io.ByteStreams;
import com.google.common.util.concurrent.ListenableFuture;
import com.google.common.util.concurrent.ListeningExecutorService;
import com.google.common.util.concurrent.MoreExecutors;
import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.nio.file.Path;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;
import java.util.concurrent.Callable;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.Executors;
import java.util.function.Function;
import javax.annotation.Nullable;

/** Runs {@code skaffold} commands. */
public class Skaffold {

  /** The location to store {@code skaffold} if auto-downloading it. */
  private static final Path CACHED_SKAFFOLD_LOCATION =
      UserCacheHome.getCacheHome().resolve("skaffold");

  /** The location to store the digest for {@code skaffold} if auto-downloading it. */
  private static final Path CACHED_SKAFFOLD_DIGEST_LOCATION =
      CACHED_SKAFFOLD_LOCATION.resolveSibling("skaffold.sha256");

  /**
   * Initializes {@link Skaffold} with a custom path to the {@code skaffold} executable.
   *
   * @param executablePath the path to {@code skaffold}
   * @return a new {@link Skaffold}
   */
  public static Skaffold atPath(Path executablePath) {
    return new Skaffold(executablePath.toString());
  }

  /**
   * Initializes {@link Skaffold} with a managed {@code skaffold} executable.
   *
   * @return a new {@link Skaffold}
   * @throws IOException
   */
  public static Skaffold init() throws IOException {
    SkaffoldDownloader.downloadLatestDigest(Files.tem);

    if ()

    SkaffoldDownloader.downloadLatest(CACHED_SKAFFOLD_LOCATION);
    return new Skaffold(CACHED_SKAFFOLD_LOCATION.toString());
  }

  private static Callable<Void> redirect(InputStream inputStream, OutputStream outputStream) {
    return () -> {
      try (InputStream inputStream1 = inputStream) {
        ByteStreams.copy(inputStream1, outputStream);
      }
      return null;
    };
  }

  private final List<String> initialTokens;

  private Function<List<String>, ProcessBuilder> processBuilderFactory = ProcessBuilder::new;
  @Nullable private Path skaffoldYaml;
  @Nullable private String profile;
  @Nullable private InputStream stdinInputStream;
  @Nullable private OutputStream stdoutOutputStream;
  @Nullable private OutputStream stderrOutputStream;

  @VisibleForTesting
  Skaffold(String... initialTokens) {
    this.initialTokens = Arrays.asList(initialTokens);
  }

  public Skaffold setSkaffoldYaml(Path skaffoldYaml) {
    this.skaffoldYaml = skaffoldYaml;
    return this;
  }

  public Skaffold setProfile(String profile) {
    this.profile = profile;
    return this;
  }

  /**
   * Sets the {@link InputStream} to provide to {@code skaffold} as the stdin.
   *
   * @param stdinInputStream provides the stdin
   * @return this
   */
  public Skaffold setStdinInputStream(InputStream stdinInputStream) {
    this.stdinInputStream = stdinInputStream;
    return this;
  }

  /**
   * Sets the {@link OutputStream} to receive the stdout.
   *
   * @param stdoutOutputStream receives the stdout
   * @return this
   */
  public Skaffold setStdoutOutputStream(OutputStream stdoutOutputStream) {
    this.stdoutOutputStream = stdoutOutputStream;
    return this;
  }

  /**
   * Sets the {@link OutputStream} to receive the stderr.
   *
   * @param stderrOutputStream receives the stderr
   * @return this
   */
  public Skaffold setStderrOutputStream(OutputStream stderrOutputStream) {
    this.stderrOutputStream = stderrOutputStream;
    return this;
  }

  /**
   * Calls {@code skaffold deploy}.
   *
   * @return the process exit code
   * @throws InterruptedException if the process was interrupted during execution
   * @throws IOException if an I/O exception occurred
   */
  public int deploy() throws InterruptedException, IOException, ExecutionException {
    List<String> command = new ArrayList<>();
    command.addAll(initialTokens);
    command.addAll(getFlags());
    command.add("deploy");

    Process skaffoldProcess = processBuilderFactory.apply(command).start();

    if (stdinInputStream != null) {
      try (OutputStream stdin = skaffoldProcess.getOutputStream()) {
        ByteStreams.copy(stdinInputStream, stdin);
      }
    }

    ListeningExecutorService listeningExecutorService =
        MoreExecutors.listeningDecorator(Executors.newCachedThreadPool());
    List<ListenableFuture<Void>> listenableFutures = new ArrayList<>();

    if (stdoutOutputStream != null) {
      listenableFutures.add(
          listeningExecutorService.submit(
              redirect(skaffoldProcess.getInputStream(), stdoutOutputStream)));
    }
    if (stderrOutputStream != null) {
      listenableFutures.add(
          listeningExecutorService.submit(
              redirect(skaffoldProcess.getErrorStream(), stderrOutputStream)));
    }

    for (ListenableFuture<Void> listenableFuture : listenableFutures) {
      listenableFuture.get();
    }

    return skaffoldProcess.waitFor();
  }

  Skaffold setProcessBuilderFactory(Function<List<String>, ProcessBuilder> processBuilderFactory) {
    this.processBuilderFactory = processBuilderFactory;
    return this;
  }

  private List<String> getFlags() {
    List<String> flags = new ArrayList<>();
    if (skaffoldYaml != null) {
      flags.add("--filename");
      flags.add(skaffoldYaml.toString());
    }
    if (profile != null) {
      flags.add("--profile");
      flags.add(profile);
    }
    return flags;
  }
}
