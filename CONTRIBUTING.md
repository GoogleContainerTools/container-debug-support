# Contributing to Container Debug Runtime Support

We'd love to accept your patches and contributions to this project. There are
just a few small guidelines you need to follow.

## Contributor License Agreement

Contributions to this project must be accompanied by a Contributor License
Agreement. You (or your employer) retain the copyright to your contribution; 
this simply gives us permission to use and redistribute your contributions as
part of the project. Head over to <https://cla.developers.google.com/> to see
your current agreements on file or to sign a new one.

You generally only need to submit a CLA once, so if you've already submitted one
(even if it was for a different project), you probably don't need to do it
again.

## Code Reviews

1. Set your git user.email property to the address used for step 1. E.g.
   ```
   git config --global user.email "janedoe@google.com"
   ```
   If you're a Googler or other corporate contributor,
   use your corporate email address here, not your personal address.
2. Fork the repository into your own Github account.
3. Please include unit tests (and integration tests if applicable) for all new code.
4. Make sure all existing tests pass.
5. Associate the change with an existing issue or file a [new issue](../../issues).
6. Create a pull request!


# Development

This project uses Skaffold's multiple config support to allow
developing for each language runtime separately.

Each language runtime is broken out to a separate directory
with a `skaffold.yaml` for development of the `duct-tape` initContainer
image.  Each image is expected to be standalone and should not require
downloading additional content across the network.  To add support for a new
language runtime, an image definition should download the necessary
files into the container image.  The image's entrypoint should then
copy those files into place at `/dbg/<runtime>`.  The image should
be added to the respective `skaffold.yaml` and referenced within
`test/k8s-*-installation.yaml`.

We currently build language support images for both `linux/amd64` and
`linux/arm64`.

Images are currently published with two names: the short-form (like
`go`) and a longer-form (`skaffold-debug-go`).  The short-forms are
marked as deprecated as we intend to move away from, but they are
still used by Skaffold.

## Testing

Integration tests are found under each language runtime directory in
`test/`.  These tests build and launch applications as pods that
are similar to the transformed form produced by `skaffold debug`.
To run:

```sh
sh run-its.sh
```

You can run this script from a language-runtime location too to
test only that runtime's support:

```sh
cd nodejs
sh ../run-its.sh
```

# Staging and Deploying

To stage a set of images for testing use the following command,
where `$REPO` is the image repository where you plan to host the
images.
```sh
skaffold build -p release,deprecated-names --default-repo $REPO
```

The `release` profile causes the images to be pushed to the specified
repository and also enables multi-arch builds using buildx (default
`linux/amd64` and `linux/arm64`).  The `deprecated-names` profile enables
using the short-form image names (`go`, `netcore`, `nodejs`, `python`)
that are currently used by `skaffold debug`.

Then configure Skaffold to point to that location:
```sh
skaffold config set --global debug-helpers-registry $REPO
```

You should then be able to use `skaffold debug` with the
staged images.
