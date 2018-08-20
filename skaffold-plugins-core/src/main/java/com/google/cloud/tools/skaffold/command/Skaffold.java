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
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.function.Function;
import java.util.function.Supplier;
import javax.annotation.Nullable;

/** Runs {@code skaffold} commands. */
public class Skaffold {

  @VisibleForTesting
  static Supplier<ExecutorService> executorServiceSupplier = Executors::newCachedThreadPool;

  /**
   * Initializes {@link Skaffold} with a custom path to the {@code skaffold} executable.
   *
   * @param executablePath the path to {@code skaffold}
   * @return a new {@link Skaffold}
   */
  public static Skaffold atPath(Path executablePath) {
    return new Skaffold(getListeningExecutorService(), executablePath.toString());
  }

  /**
   * Sets the {@link ExecutorService} to handle the {@code skaffold} process. Uses {@link
   * Executors#newCachedThreadPool} by default.
   *
   * @param executorService the executor
   */
  public static void setExecutorService(ExecutorService executorService) {
    Skaffold.executorServiceSupplier = () -> executorService;
  }

  @VisibleForTesting
  static ListeningExecutorService getListeningExecutorService() {
    return MoreExecutors.listeningDecorator(Skaffold.executorServiceSupplier.get());
  }

  private static Callable<Void> redirect(InputStream inputStream, OutputStream outputStream) {
    return () -> {
      ByteStreams.copy(inputStream, outputStream);
      return null;
    };
  }

  private final ListeningExecutorService listeningExecutorService;
  private final List<String> initialTokens;

  private Function<List<String>, ProcessBuilder> processBuilderFactory = ProcessBuilder::new;
  @Nullable private Path skaffoldYaml;
  @Nullable private String profile;
  @Nullable private InputStream stdinInputStream;
  @Nullable private OutputStream stdoutOutputStream;
  @Nullable private OutputStream stderrOutputStream;

  @VisibleForTesting
  Skaffold(ListeningExecutorService listeningExecutorService, String... initialTokens) {
    this.listeningExecutorService = listeningExecutorService;
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
  public Skaffold redirectToStdin(InputStream stdinInputStream) {
    this.stdinInputStream = stdinInputStream;
    return this;
  }

  /**
   * Sets the {@link OutputStream} to receive the stdout.
   *
   * @param stdoutOutputStream receives the stdout
   * @return this
   */
  public Skaffold redirectStdoutTo(OutputStream stdoutOutputStream) {
    this.stdoutOutputStream = stdoutOutputStream;
    return this;
  }

  /**
   * Sets the {@link OutputStream} to receive the stderr.
   *
   * @param stderrOutputStream receives the stderr
   * @return this
   */
  public Skaffold redirectStderrTo(OutputStream stderrOutputStream) {
    this.stderrOutputStream = stderrOutputStream;
    return this;
  }

  /**
   * Calls {@code skaffold deploy}.
   *
   * @return the process exit code
   * @throws ExecutionException if an error occurred while handling the process I/O
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

    List<ListenableFuture<Void>> listenableFutures = new ArrayList<>();

    try (InputStream stdout = skaffoldProcess.getInputStream();
        InputStream stderr = skaffoldProcess.getErrorStream()) {
      if (stdoutOutputStream != null) {
        listenableFutures.add(
            listeningExecutorService.submit(redirect(stdout, stdoutOutputStream)));
      }
      if (stderrOutputStream != null) {
        listenableFutures.add(
            listeningExecutorService.submit(redirect(stderr, stderrOutputStream)));
      }

      for (ListenableFuture<Void> listenableFuture : listenableFutures) {
        listenableFuture.get();
      }
    }

    return skaffoldProcess.waitFor();
  }

  @VisibleForTesting
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
