<!-- omit in toc -->
# Contributing to k8gb

- [Getting started](#getting-started)
- [Getting help](#getting-help)
- [Reporting issues](#reporting-issues)
- [Contribution flow](#contribution-flow)
- [Local setup](#local-setup)
  - [Deploy k8gb locally](#deploy-k8gb-locally)
  - [Upgrade k8gb to local candidate](#upgrade-k8gb-to-local-candidate)
- [Overriding dev environment settings](#overriding-dev-environment-settings)
- [Testing](#testing)
  - [Unit testing](#unit-testing)
  - [Testing against real k8s clusters](#testing-against-real-k8s-clusters)
- [Debugging](#debugging)
- [Metrics](#metrics)
- [Code style](#code-style)
  - [Logging](#logging)
  - [Error handling](#error-handling)
- [Commit and Pull Request message](#commit-and-pull-request-message)
  - [Signature](#signature)
  - [Changelog](#changelog)
- [Documentation](#documentation)
- [k8gb.io website](#k8gbio-website)
  - [Local website authoring and testing](#local-website-authoring-and-testing)
    - [Using custom SSL CA certificate](#using-custom-ssl-ca-certificate)
    - [.env support](#env-support)
- [End-to-end demo helper](#end-to-end-demo-helper)
- [Release process](#release-process)
  - [Signed releases](#signed-releases)
  - [Software bill of materials](#software-bill-of-materials)

k8gb is licensed under [Apache 2 License](https://github.com/k8gb-io/k8gb/blob/master/LICENSE) and accepts contributions via GitHub pull requests.
This document outlines the resources and guidelines necessary to follow by contributors to the k8gb project.

## Getting started

- Fork the repository on GitHub
- See the [local playground guide](local.md) for local dev environment setup

## Getting help

Feel free to ask for help and join the discussions at [k8gb community discussions forum](https://github.com/k8gb-io/k8gb/discussions).
We have [dedicated `#k8gb` channel on Cloud Native Computing Foundation (CNCF) Slack](https://cloud-native.slack.com/archives/C021P656HGB),
and we can also actively monitoring [`#sig-multicluster` channel on Kubernetes Slack](https://kubernetes.slack.com/archives/C09R1PJR3).

## Reporting issues

Reporting bugs is one of the best ways to contribute.
Feel free to open an issue describing your problem or question.

## Contribution flow

Following is a rough outline for the contributor's workflow:

- Create a topic branch from where to base the contribution.
- Make commits of logical units.
- Make sure your code is clean and follows the [code style and logging guidelines](#code-style).
- Make sure the commit messages are in the [proper format](#commit-and-pull-request-message).
- Make sure the changes are covered by [reasonable amount of testing](#testing).
- Push changes in a topic branch to a personal fork of the repository.
- Submit a pull request to k8gb-io/k8gb GitHub repository.
- Resolve review comments.
- PR must receive an "LGTM" approval from at least one maintainer listed in the `CODEOWNERS` file.

## Local setup

### Deploy k8gb locally

```sh
make deploy-full-local-setup
```
deploys k8gb from scratch, including:

* 2 local clusters
* Stable k8gb helm chart
* Gslb Custom Resources examples
* Test applications

### Upgrade k8gb to local candidate
```sh
make upgrade-candidate
```
performs upgrade of k8gb helm chart and controller to the testing version built from your current development tree.

## Overriding dev environment settings

Sometimes there is a need to override environment variables used by `make` targets for local k8gb development.
This can be easily achieved by providing the list of environment variables with respective values in the `.env` file at the local repository root:
```sh
cat .env

# .env:
LOG_LEVEL=info
LOG_FORMAT=json
```
Overrides done this way can persist between terminal sessions and can be used as a single point of configuration for development in terminal and IDE of choice.

## Testing

- Any functional GSLB controller code change should be secured by the corresponding [unit tests](https://github.com/k8gb-io/k8gb/tree/master/controllers/gslb_controller_reconciliation_test.go).
- Integration terratest suite is located [here](https://github.com/k8gb-io/k8gb/tree/master/terratest).
  These tests are updated only if the change is substantial enough to affect the main end-to-end flow.
- See the [local playground guide](https://github.com/k8gb-io/k8gb/blob/master/docs/local.md) for local testing environment setup and integration test execution.

### Unit testing
- Include unit tests when you contribute new features, as they help to a) prove that your code works correctly, and b) guard against future breaking changes to lower the maintenance cost.
- Bug fixes also generally require unit tests, because the presence of bugs usually indicates insufficient test coverage.

Use `make test` to check your implementation changes.

### Testing against real k8s clusters

There is a possibility to execute the integration terratest suite over the real clusters.
For this, you need to override the set of test settings as in the example below.
```sh
PRIMARY_GEO_TAG=af-south-1 \
SECONDARY_GEO_TAG=eu-west-1 \
DNS_SERVER1=a377095726f1845fb85b95c2afef8ac0-9a1a10f24e634e28.elb.af-south-1.amazonaws.com \
DNS_SERVER1_PORT=53 \
DNS_SERVER2=a873f5c83be624a0a84c05a743d598a8-443f7e0285e4a28f.elb.eu-west-1.amazonaws.com \
DNS_SERVER2_PORT=53 \
GSLB_DOMAIN=test.k8gb.io \
K8GB_CLUSTER1=arn:aws:eks:af-south-1:<aws-account-id>:cluster/k8gb-cluster-af-south-1 \
K8GB_CLUSTER2=arn:aws:eks:eu-west-1:<aws-account-id>:cluster/k8gb-cluster-eu-west-1 \
make terratest
```

## Debugging

1. Install Delve debugger first. Follow the [installation instructions](https://github.com/go-delve/delve/tree/master/Documentation/installation) for specific platforms from Delve's website.

2. Run delve with options specific to IDE of choice.
There is a dedicated make target available for Goland:

    ```sh
    make debug-idea
    ```
    [This article](https://dev4devs.com/2019/05/04/operator-framework-how-to-debug-golang-operator-projects/) describes possible option examples for Goland and VS Code.

3. Attach debugger of your IDE to port `2345`.

## Metrics
More info about k8gb metrics can be found in the [metrics.md](metrics.md) document.
If you need to check and query the k8gb metrics locally, you can install a Prometheus in the local clusters using the `make deploy-prometheus` command.

The deployed Prometheus scrapes metrics from the dedicated k8gb operator endpoint and makes them accessible via Prometheus web UI:

- http://127.0.0.1:9080
- http://127.0.0.1:9081

All the metric data is ephemeral and will be lost with pod restarts.
To uninstall Prometheus, run `make uninstall-prometheus`

Optionally, you can also install Grafana that will have the datasources configured and example dashboard ready using `make deploy-grafana`

## Code style

k8gb project is using the coding style suggested by the Golang community. See the [golang-style-doc](https://github.com/golang/go/wiki/CodeReviewComments) for details.

Please follow this style to make k8gb easy to review, maintain and develop.
Run `make check` to automatically check if your code is compliant.

### Logging

k8gb project is using the [zerolog](https://github.com/rs/zerolog) library for logging.

- Please make sure to follow the zerolog library concepts and conventions in the code.
- Try to use [contextual logging](https://github.com/rs/zerolog#contextual-logging) whenever possible.
- Pay attention to [error logging](https://github.com/rs/zerolog#error-logging) recommendations.

### Error handling
See [effective go errors](https://golang.org/doc/effective_go.html#errors) first. Do not discard errors using `_` variables except
tests, or, in truly exceptional situations. If a function returns an error, check it to make sure the function succeeded.
If the function fails, consider logging the error and recording errors in metrics (see: [logging recommendations](#logging)).

The following example demonstrates error handling inside the reconciliation loop:
```go
	var log = logging.Logger()

	var m = metrics.Metrics()
	...
	err = r.DNSProvider.CreateZoneDelegationForExternalDNS(gslb)
	if err != nil {
		log.Err(err).Msg("Unable to create zone delegation")
		m.ErrorIncrement(gslb)
		return result.Requeue()
	}
```

## Commit and Pull Request message

We follow a rough convention for PR and commit messages, which is designed to answer two questions: what changed and why.
The subject line should feature the what, and the body of the message should describe the why.
The format can be described more formally as follows:

```
<what was changed>

<why this change was made>

<footer>
```

The first line is the subject and should be no longer than 70 characters.
The second line is always blank.
Consequent lines should be wrapped at 80 characters.
This way, the message is easier to read on GitHub as well as in various git tools.

```
scripts: add the test-cluster command

This command uses "k3d" to set up a test cluster for debugging.

Fixes #38
```

Commit message can be made lightweight unless it is the only commit forming the PR.
In that case, the message can follow the simplified convention:

```
<what was changed and why>
```
This convention is useful when several minimalistic commit messages are going to form PR descriptions as bullet points of what was done during the final squash and merge for PR.

### Signature

As a CNCF project, k8gb must comply with [Developer Certificate of Origin (DCO)](https://developercertificate.org/) requirement.
[DCO GitHub Check](https://github.com/apps/dco) automatically enforces DCO for all commits.
Contributors are required to ensure that every commit message contains the following signature:
```txt
Signed-off-by: NAME SURNAME <email@address.example.org>
```
The best way to achieve this automatically for local development is to create the following alias in the `~/.gitconfig` file:
```.gitconfig
[alias]
ci = commit -s
```
When a commit is created in GitHub UI as a result of [accepted suggested change](https://docs.github.com/en/github/collaborating-with-issues-and-pull-requests/incorporating-feedback-in-your-pull-request#applying-suggested-changes), the signature should be manually added to the "optional extended description" field.

### Changelog

The [CHANGELOG](https://github.com/k8gb-io/k8gb/blob/master/CHANGELOG.md) is automatically generated from Github PRs and Issues during release.
Use dedicated [keywords](https://docs.github.com/en/github/managing-your-work-on-github/linking-a-pull-request-to-an-issue#linking-a-pull-request-to-an-issue-using-a-keyword) in PR message or [manual PR and Issue linking](https://docs.github.com/en/github/managing-your-work-on-github/linking-a-pull-request-to-an-issue#manually-linking-a-pull-request-to-an-issue) for clean changelog generation.
Issues and PRs should be also properly tagged with valid project tags ("bug", "enhancement", "wontfix", etc )

## Documentation

If contribution changes the existing APIs or user interface, it must include sufficient documentation explaining the use of the new or updated feature.

## k8gb.io website

k8gb.io website is a Jekyll-based static website generated from project markdown documentation and hosted by GitHub Pages.
`gh-pages` branch contains the website source, including configuration, website layout, and styling.
Markdown documents are automatically populated to `gh-pages` from the main branch and should be authored there.
Changes to the k8gb.io website layout and styling should be checked out from the `gh-pages` branch and  PRs should be created against `gh-pages`.

### Local website authoring and testing

These instructions will help you to set up and use local website authoring and testing environment:

- Check-out the [`gh-pages`](https://github.com/k8gb-io/k8gb/tree/gh-pages) branch
- Create dedicated [GitHub Personal Access Token](https://github.com/settings/tokens/new) with `public_repo` permission and assign it to the `JEKYLL_GITHUB_TOKEN` environment variable:

  ```sh
  export JEKYLL_GITHUB_TOKEN=<your-github-token>
  ```

- Run `make` and serve the local copy of the k8gb.io website.

  *The target utilizes the `jekyll/jekyll` docker container to avoid local environment pollution from Jekyll gem dependencies.*

- Open the `http://localhost:4000/` page in your browser.
- Website will get automatically rebuilt and refreshed in the browser to accommodate the related code changes.

#### Using custom SSL CA certificate

When the dev environment is behind a proxy, it might be required to use a custom SSL CA certificate,
otherwise, Jekyll is not able to update its plugins and dependencies.
Local website authoring target automatically detects and wires the SSL certificate to its docker container from the path provided by the `CUSTOM_CERT_PATH` environment variable:

```sh
export CUSTOM_CERT_PATH=/path/to/custom/certificate.pem
```

> NOTE: Certificate should be in a PEM format

#### .env support

Environment variables can be stored in a local `.env` file in the project's root to retain their persistency between the authoring sessions:

```.env
JEKYLL_GITHUB_TOKEN=<your-github-token>
CUSTOM_CERT_PATH=<your-custom-ssl-certificate>
```

## End-to-end demo helper

The demo helper is designed to work with `podinfo` that was deployed by

```sh
make deploy-test-apps
```

It will configure `podinfo` to expose geotag as part of an HTTP response.

To test and/or demonstrate continuous query to GSLB enabled endpoint execute

```sh
make demo DEMO_URL=https://failover.test.exampledns.tk
```

The happy path will look like:

```sh
[Thu May 27 15:35:26 UTC 2021] ...

200  "message": "eu-west-1",

[Thu May 27 15:35:31 UTC 2021] ...

200
  "message": "eu-west-1",
[Thu May 27 15:35:36 UTC 2021] ...
```

The sources for demo helper images can be found [here](https://github.com/k8gb-io/k8gb/tree/master/deploy/test-apps/curldemo/)

To enable verbose debug output declare `DEMO_DEBUG=1` like
```sh
make demo DEMO_URL=https://failover.test.exampledns.tk DEMO_DEBUG=1
```

## Release process

* Bump the version in `Chart.yaml`, see [example PR](https://github.com/k8gb-io/k8gb/pull/749). Make sure the
commit message starts with `RELEASE:`.
* In addition, bump the version in `chart/k8gb/README.md` to reflect correct version tags in artifacthub after
release, e.g.
```
chart/k8gb/README.md:![Version: v0.13.0](https://img.shields.io/badge/Version-v0.13.0-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: v0.13.0](https://img.shields.io/badge/AppVersion-v0.13.0-informational?style=flat-square)
```
* Ensure `stable` become `next`. Update `deploy/helm/stable.yaml` with `deploy/helm/next.yaml`.
* Merge the Pull Request after the review approval (make sure the squash or rebase is used, merge commit will not trigger the release pipeline)
* At this point a DRAFT release will be created on GitHub. After the [automatic tag](https://github.com/k8gb-io/k8gb/actions/workflows/cut_release.yaml) & [release pipeline](https://github.com/k8gb-io/k8gb/actions/workflows/release.yaml)
have been successfully completed, you check the [release DRAFT](https://github.com/k8gb-io/k8gb/releases) and if it is OK, you click on the **"Publish release"** button.

* Check the [helm publish pipeline](https://github.com/k8gb-io/k8gb/actions/workflows/helm_publish.yaml) status
* Check the [offline changelog](https://github.com/k8gb-io/k8gb/actions/workflows/changelog_pr.yaml) status. This pipeline creates
a pull request with an offline changelog. Do a review and if everything is ok, merge it.

Congratulations, the release is complete!

### Signed releases

During the release process we generate SLSA (Supply-chain Levels for Software Artifacts) Level 3 provenance
attestations using the official [slsa-framework/slsa-github-generator](https://github.com/slsa-framework/slsa-github-generator).

This provenance file is signed itself and attached with the signature to the release artifacts. For signing
the artifacts we use [`cosign`](https://github.com/sigstore/cosign) tool and private key stored as the
repository secret. Public key is available in the repository itself in file [`cosign.pub`](https://github.com/k8gb-io/k8gb/blob/master/cosign.pub).
This way anybody can verify the origin of arbitrary artifact. In order to regenerate the keys for cosign,
one can run `cosign generate-key-pair`, use some passphrase and update the `COSIGN_{PRIVATE,PUBLIC}_KEY` &
`COSIGN_PASSWORD` repo secret and also the content of [`cosign.pub`](https://github.com/k8gb-io/k8gb/blob/master/cosign.pub) file.

#### SLSA Provenance

The release artifacts include a `multiple.intoto.jsonl` file that contains SLSA provenance attestations
compliant with the [in-toto Attestation Framework](https://in-toto.io/Statement/v0.1). This file provides:

- **Build metadata**: Information about the GitHub Actions workflow that built the artifacts
- **Material provenance**: Source code repository, commit SHA, and build environment details
- **Artifact attestations**: Cryptographic signatures for all release binaries and archives
- **Supply chain security**: Verifiable proof that artifacts weren't tampered with during build

#### Artifact Signing

For signing the artifacts we use [`cosign`](https://github.com/sigstore/cosign) tool with:
- **Private key**: Stored as repository secret (`COSIGN_PRIVATE_KEY`)
- **Public key**: Available in [`cosign.pub`](./cosign.pub) file
- **Passphrase**: Stored as repository secret (`COSIGN_PASSWORD`)

To regenerate signing keys:
```bash
cosign generate-key-pair
# Update COSIGN_{PRIVATE,PUBLIC}_KEY & COSIGN_PASSWORD secrets
# Update ./cosign.pub file content
```

#### Container Image Signatures

All container images are signed with `cosign` and signatures are pushed to registries:
- **Signature location**: OCI format under predictable names
- **Discovery**: `cosign triangulate $IMAGE`
- **Verification**: `cosign verify --key cosign.pub $IMAGE`

#### Verification Instructions

To verify release artifacts using the new SLSA provenance:

**1. Install slsa-verifier:**
```bash
# Download from https://github.com/slsa-framework/slsa-verifier/releases
wget https://github.com/slsa-framework/slsa-verifier/releases/latest/download/slsa-verifier-linux-amd64
chmod +x slsa-verifier-linux-amd64
```

**2. Download artifacts:**
```bash
# Download the binary and provenance from GitHub release
wget https://github.com/k8gb-io/k8gb/releases/download/v0.15.0/k8gb_0.15.0_linux_amd64.tar.gz
wget https://github.com/k8gb-io/k8gb/releases/download/v0.15.0/multiple.intoto.jsonl
```

**3. Verify SLSA provenance:**
```bash
./slsa-verifier-linux-amd64 verify-artifact \
  --provenance-path multiple.intoto.jsonl \
  --source-uri github.com/k8gb-io/k8gb \
  k8gb_0.15.0_linux_amd64.tar.gz
```

**4. Verify container images:**
```bash
cosign verify --key cosign.pub docker.io/absaoss/k8gb:v0.15.0
```

**5. Verify container provenance:**
```bash
cosign verify-attestation \
  --key cosign.pub \
  --type slsaprovenance \
  docker.io/absaoss/k8gb@sha256:digest
```

The verification confirms that artifacts were:
- Built by the official GitHub Actions workflow
- From the expected source repository (github.com/k8gb-io/k8gb)
- Without tampering during the build process

### Software bill of materials

For each container image we also create Software bill of materials (SBOM) file + its signature that ends up
as part of the release. These files follows this naming pattern:
`k8gb_{version}_{os}_{arch}.tar.gz.sbom.json` and are generated using [Syft](https://github.com/anchore/syft)tool.

---
Thanks for contributing!
