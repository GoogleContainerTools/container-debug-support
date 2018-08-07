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
import com.google.common.io.ByteStreams;
import com.google.common.util.concurrent.Futures;
import com.google.common.util.concurrent.ListeningExecutorService;
import com.google.common.util.concurrent.MoreExecutors;
import java.io.BufferedInputStream;
import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.nio.file.Path;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;
import java.util.concurrent.Executor;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import javax.annotation.Nullable;

/** Runs {@code skaffold} commands. */
public class Skaffold {

  public static Skaffold atPath(Path executablePath) {
    return new Skaffold(executablePath);
  }

  public static Skaffold init() {
    SkaffoldDownloader
    return new Skaffold(executablePath);
  }

  private final Path executablePath;

  @Nullable
  private Path skaffoldYaml;
  @Nullable
  private String profile;
  @Nullable
  private InputStream stdinInputStream;
  @Nullable
  private OutputStream stdoutOutputStream;
  @Nullable
  private OutputStream stderrOutputStream;

  private Skaffold(Path executablePath) {
    this.executablePath = executablePath;
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

  public int deploy() throws InterruptedException, IOException {
    List<String> command = new ArrayList<>();
    command.add(executablePath.toString());
    command.addAll(getFlags());
    command.add("deploy");

    Process skaffoldProcess = new ProcessBuilder(command).start();

    if (stdinInputStream != null) {
      try (OutputStream stdin = skaffoldProcess.getOutputStream()) {
        ByteStreams.copy(stdinInputStream, stdin);
      }
    }

    ListeningExecutorService listeningExecutorService = MoreExecutors.listeningDecorator(Executors.newCachedThreadPool());
    if (stdoutOutputStream != null) {
      listeningExecutorService.submit(
          () -> {
            try (InputStream stdout = skaffoldProcess.getInputStream()) {
              ByteStreams.copy(stdout, stdoutOutputStream);
            }
            return null;
          });
    }
    if (stderrOutputStream != null) {
      listeningExecutorService.submit(
          () -> {
            try (InputStream stderr = skaffoldProcess.getErrorStream()) {
              ByteStreams.copy(stderr, stderrOutputStream);
            }
            return null;
          });
    }

    return skaffoldProcess.waitFor();
  }

  private List<String> getFlags() {
    List<String> flags = new ArrayList<>();
    if (skaffoldYaml != null) {
      flags.add("--filename");
      flags.add(skaffoldYaml.toString());
    }
    return flags;
  }
}
