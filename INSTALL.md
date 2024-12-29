# Install

yodoc is written in Go. So you only have to install a binary in your `PATH`.

There are some ways to install yodoc.

1. [Homebrew](#homebrew)
1. [Scoop](#scoop)
1. [aqua](#aqua)
1. [GitHub Releases](#github-releases)
1. [Build an executable binary from source code yourself using Go](#build-an-executable-binary-from-source-code-yourself-using-go)

## Homebrew

You can install yodoc using [Homebrew](https://brew.sh/).

```sh
brew install suzuki-shunsuke/yodoc/yodoc
```

## Scoop

You can install yodoc using [Scoop](https://scoop.sh/).

```sh
scoop bucket add suzuki-shunsuke https://github.com/suzuki-shunsuke/scoop-bucket
scoop install yodoc
```

## aqua

You can install yodoc using [aqua](https://aquaproj.github.io/).

```sh
aqua g -i suzuki-shunsuke/yodoc
```

## Build an executable binary from source code yourself using Go

```sh
go install github.com/suzuki-shunsuke/yodoc/cmd/yodoc@latest
```

## GitHub Releases

You can download an asset from [GitHub Releases](https://github.com/suzuki-shunsuke/yodoc/releases).
Please unarchive it and install a pre built binary into `$PATH`. 

### Verify downloaded assets from GitHub Releases

You can verify downloaded assets using some tools.

1. [GitHub CLI](https://cli.github.com/)
1. [slsa-verifier](https://github.com/slsa-framework/slsa-verifier)
1. [Cosign](https://github.com/sigstore/cosign)

### 1. GitHub CLI

You can install GitHub CLI by aqua.

```sh
aqua g -i cli/cli
```

```sh
version=v0.1.2
asset=yodoc_darwin_arm64.tar.gz
gh release download -R suzuki-shunsuke/yodoc "$version" -p "$asset"
gh attestation verify "$asset" \
  -R suzuki-shunsuke/yodoc \
  --signer-workflow suzuki-shunsuke/go-release-workflow/.github/workflows/release.yaml
```

### 2. slsa-verifier

You can install slsa-verifier by aqua.

```sh
aqua g -i slsa-framework/slsa-verifier
```

```sh
version=v0.1.2
asset=yodoc_darwin_arm64.tar.gz
gh release download -R suzuki-shunsuke/yodoc "$version" -p "$asset" -p multiple.intoto.jsonl
slsa-verifier verify-artifact "$asset" \
  --provenance-path multiple.intoto.jsonl \
  --source-uri github.com/suzuki-shunsuke/yodoc \
  --source-tag "$version"
```

### 3. Cosign

You can install Cosign by aqua.

```sh
aqua g -i sigstore/cosign
```

```sh
version=v0.1.2
checksum_file="yodoc_${version#v}_checksums.txt"
asset=yodoc_darwin_arm64.tar.gz
gh release download "$version" \
  -R suzuki-shunsuke/yodoc \
  -p "$asset" \
  -p "$checksum_file" \
  -p "${checksum_file}.pem" \
  -p "${checksum_file}.sig"
cosign verify-blob \
  --signature "${checksum_file}.sig" \
  --certificate "${checksum_file}.pem" \
  --certificate-identity-regexp 'https://github\.com/suzuki-shunsuke/go-release-workflow/\.github/workflows/release\.yaml@.*' \
  --certificate-oidc-issuer "https://token.actions.githubusercontent.com" \
  "$checksum_file"
cat "$checksum_file" | sha256sum -c --ignore-missing
```