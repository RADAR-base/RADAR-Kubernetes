# AGENTS.md

Guidance for AI coding agents (and human contributors) working in this repository.

## About this repository

`RADAR-Kubernetes` is the Infrastructure-as-Code (IaC) project for deploying the RADAR-base platform onto a
Kubernetes cluster. It does not contain application source code; instead it wires together and configures the
Helm charts of the various RADAR-base components using [Helm](https://helm.sh) and
[helmfile](https://github.com/helmfile/helmfile).

Key locations:

- `helmfile.d/` — helmfile release definitions, split into numbered files (`00-…`, `10-…`, `20-…`). Lower numbers
  are installed first; a component that other components depend on (e.g. Kafka, PostgreSQL) should live in a
  lower-numbered file.
- `etc/` — default and per-component values files (`etc/base.yaml`, `etc/base-secrets.yaml`,
  `etc/production.yaml`, `etc/production.yaml.gotmpl`, and per-component subfolders such as `etc/mongodb`,
  `etc/radar-kratos`, etc.).
- `bin/` — operational scripts (`bin/init`, `bin/generate-secrets`, `bin/chart-updates`, `bin/version_scan`, …).
- `docs/` — supplementary docs, notably `docs/development_guide.md` (how to add a component, test locally with
  k3d) and `docs/upgrade.md`.
- `dev/` — local development cluster configuration (k3d).
- `test/` — behave-based feature tests for the deployment.

Standard workflow: configure `etc/production.yaml`, then run `helmfile sync` (or `helmfile apply`/`helmfile diff`
scoped with `--selector name=<component>`) against the target `kubeContext`.

## Component inventory & local checkouts

This repo only wires together Helm charts; the actual application source code for each RADAR-base microservice
lives in separate GitHub repositories under the `RADAR-base` org. On this machine, those repos are expected to be
checked out as **siblings of this repo's parent directory** (i.e. `../<repo-name>` relative to `RADAR-Kubernetes`,
same level as `radar-helm-charts`). That parent directory also contains many unrelated repos (other projects,
forks, scratch clones) — always match on exact repo name, don't assume a directory is a RADAR-base component just
because it's nearby.

The table below maps each deployed component to its helmfile release name(s) in `helmfile.d/10-services.yaml`,
its source repo, local checkout path, and build tooling. Docker image `repository:` values (from
`radar-helm-charts/charts/<chart>/values.yaml`) are namespaced as `radar-base/<github-repo>/<image-name>` on
ghcr.io/Docker Hub, which is how a given chart was traced back to its source repo below.

| Component | Helmfile release(s) | GitHub repo | Local path | Build tool |
|---|---|---|---|---|
| Management Portal (auth/identity, legacy OAuth2 authserver) | `management-portal` | `RADAR-base/ManagementPortal` | `../ManagementPortal` | Gradle (Java) + npm (Angular frontend) |
| App Config service | `app-config` | `RADAR-base/radar-app-config` | `../radar-app-config` | Gradle Kotlin DSL |
| App Config frontend | `app-config-frontend` | `RADAR-base/radar-app-config` (same repo, separate image) | `../radar-app-config` | Gradle Kotlin DSL |
| Self-enrolment UI (Next.js) | `radar-self-enrolment-ui` | `RADAR-base/radar-self-enrolment-ui` | `../radar-self-enrolment-ui` | npm (Next.js) |
| App server (push/questionnaire scheduling) | `radar-appserver` | `RADAR-base/RADAR-Appserver` | `../RADAR-Appserver` | Gradle (Java) |
| Data dashboard backend | `data-dashboard-backend` | `RADAR-base/radar-data-dashboard-backend` | `../radar-data-dashboard-backend` | Gradle Kotlin DSL |
| REST API Gateway | `radar-gateway` | `RADAR-base/RADAR-Gateway` | `../RADAR-Gateway` | Gradle Kotlin DSL |
| REDCap integration | `radar-integration` | `RADAR-base/RADAR-RedcapIntegration` | `../RADAR-RedcapIntegration` | Gradle (Java) |
| JDBC output connector (Kafka Connect) | `radar-jdbc-connector-grafana`, `radar-jdbc-connector-data-dashboard`, `radar-jdbc-connector-realtime-dashboard` | `RADAR-base/RADAR-JDBC-Connector` (fork of `confluentinc/kafka-connect-jdbc`) | `../RADAR-JDBC-Connector` | fork remote is `origin`; `upstream`/`kafka-connect-jdbc` remotes point at Confluent's repo — check for upstream security fixes there too |
| Fitbit / Oura REST source connectors | `radar-fitbit-connector`, `radar-oura-connector` | `RADAR-base/RADAR-REST-Connector` | `../RADAR-REST-Connector` | Gradle Kotlin DSL |
| REST sources auth backend + authorizer UI | `radar-rest-sources-backend`, `radar-rest-sources-authorizer` | `RADAR-base/RADAR-Rest-Source-Auth` | `../RADAR-Rest-Source-Auth` | Gradle Kotlin DSL |
| S3 sink connector (Kafka Connect transform) | `radar-s3-connector` | `RADAR-base/kafka-connect-transform-keyvalue` | `../kafka-connect-transform-keyvalue` | Gradle (Java) |
| Output restructure (S3 → flat file) | `radar-output` | `RADAR-base/radar-output-restructure` | `../radar-output-restructure` | Gradle Kotlin DSL |
| Upload connector backend + frontend | `radar-upload-connect-backend`, `radar-upload-connect-frontend` | `RADAR-base/radar-upload-source-connector` | `../radar-upload-source-connector` | Gradle Kotlin DSL |
| Upload source connector (Kafka Connect) | `radar-upload-source-connector` | `RADAR-base/radar-upload-source-connector` (same repo) | `../radar-upload-source-connector` | Gradle Kotlin DSL |
| Push notification endpoint | `radar-push-endpoint` | `RADAR-base/RADAR-PushEndpoint` | `../RADAR-PushEndpoint` | Gradle Kotlin DSL |
| Landing/portal page | `radar-home` | `RADAR-base/radar-home` | `../radar-home` | npm |
| Helm charts for all of the above | (n/a — chart source, not deployed as a release) | `RADAR-base/radar-helm-charts` | `../radar-helm-charts` | Helm |

Shared libraries — not deployed directly, but pulled in as dependencies by several of the Gradle-based services
above. A vulnerable transitive dependency is often best fixed by bumping the library version here and then bumping
the library version in each consumer's `build.gradle(.kts)`:

| Library | GitHub repo | Local path |
|---|---|---|
| Jersey/REST framework wrapper used by most Java backends | `RADAR-base/radar-jersey` | `../radar-jersey` |
| Shared Kafka/Avro producer-consumer utilities | `RADAR-base/radar-commons` | `../radar-commons` |
| Android client SDK (mobile app, not deployed to k8s but same org/security process) | `RADAR-base/radar-commons-android` | `../radar-commons-android` |
| Avro schema definitions used across the Kafka pipeline | `RADAR-base/RADAR-Schemas` | `../RADAR-Schemas` |

Not tracked as separate RADAR-base repos (no local checkout expected): `radar-kratos` and `radar-hydra` deploy
upstream Ory images directly (no RADAR-specific source/patching needed beyond the chart itself); infra components
like `mongodb`, `elasticsearch`, `cp-kafka`, `postgresql`, `minio`, `redis`, `radar-grafana` are third-party charts
with upstream-maintained images.

## Conventions to follow

- Match existing YAML formatting/indentation in `helmfile.d/*.yaml` and `etc/*.yaml` files.
- New components: add a release entry in the appropriate `helmfile.d` file, a default config block in
  `etc/base.yaml` (`_install`, `_chart_version`, `_extra_timeout`), and secrets in `etc/base-secrets.yaml` /
  `bin/generate-secrets` if applicable — see `docs/development_guide.md` for the full checklist.
- Prefer testing changes against a local k3d cluster (`docs/development_guide.md`) and scoped helmfile commands
  (`helmfile apply --file helmfile.d/<file>.yaml --selector name=<component>`) rather than a full `helmfile sync`.
- Do not commit real secrets/credentials; use `bin/generate-secrets` patterns and placeholder/example values.
- Keep documentation (`README.md`, `docs/`) in sync with structural changes, especially when components move
  between "hot-wired local path" and "published chart" states.

## Features under development

The sections below document work-in-progress efforts that don't reflect the steady-state of this repo — treat
them as current context, not as permanent architecture.

### Ory auth integration (`feature/ory-auth-base`)

This project is migrating authentication (authn) and authorization (authz) to
[Ory Kratos](https://www.ory.sh/kratos/) and [Ory Hydra](https://www.ory.sh/hydra/), replacing/augmenting the
existing Management Portal OAuth2 authserver. Relevant branches/repos:

- **This repo, branch `feature/ory-auth-base`** — the integration branch for this work. Relevant releases live in
  `helmfile.d/10-services.yaml`: `radar-kratos`, `radar-hydra`, `radar-self-enrolment-ui`, `app-config`,
  `app-config-frontend`, and `management_portal` (with `authserver.internal` toggling the legacy authserver).
  Default values for these are under `etc/radar-kratos/` (add `etc/radar-hydra/` similarly if/when needed) and in
  `etc/base.yaml` / `etc/production.yaml`.
- **`../radar-helm-charts`, branch `feature/sep-oauth2`** — sibling repo (expected checked out alongside this
  repo, i.e. at `../radar-helm-charts` relative to this project root) containing the actual Helm charts for
  `radar-kratos`, `radar-hydra`, `app-config`, `app-config-frontend`, `radar-self-enrolment-ui`, etc. This repo's
  `feature/ory-auth-base` branch depends on chart changes made there.
- **`../radar.connectdigitalstudy.com`, `main` branch** — sibling repo with a partial Kratos implementation that
  uses MFA. Use it as a reference/source of configuration patterns (e.g. identity schemas, MFA/TOTP settings,
  courier/SMTP config) when completing the Kratos setup here.
- **`../radar-self-enrolment-ui`, branch `feature/mfa-flow`** — sibling repo (expected checked out alongside this
  repo, i.e. at `../radar-self-enrolment-ui` relative to this project root) containing the actual frontend source
  for the `radar-self-enrolment-ui` Next.js app deployed via `../radar-helm-charts/charts/radar-self-enrolment-ui`.
  The `feature/mfa-flow` branch holds the login-time TOTP challenge and settings-flow TOTP enrollment UI changes
  (`app/_ui/auth/totpChallenge.tsx`, `app/_ui/auth/totpEnrollment.tsx`, `app/_ui/auth/mfaRequired.tsx`, plus the
  corresponding `kratos.ts`/API route changes) that pair with the MFA/AAL2 config work in this repo and in
  `radar-helm-charts`. Publishing/building a new container image from that branch and bumping the chart's image tag
  is still a manual follow-up.

#### Important: hot-wired local chart paths

While this feature is in progress, `helmfile.d/10-services.yaml` does **not** pull the Ory-related charts from the
published `radar` chart repository. Instead it points directly at local file paths into the sibling
`radar-helm-charts` checkout, e.g.:

```yaml
- name: radar-kratos
#    chart: radar/radar-kratos
    chart: ../../radar-helm-charts/charts/radar-kratos
    ...
- name: radar-hydra
#    chart: radar/radar-hydra
    chart: ../../radar-helm-charts/charts/radar-hydra
    ...
```

The commented-out `chart: radar/...` line above each hot-wired entry shows the eventual, published-chart form.
Implications for agents:

- These paths assume `radar-helm-charts` is checked out as a sibling directory to `RADAR-Kubernetes`
  (`../radar-helm-charts`) and is on (or compatible with) branch `feature/sep-oauth2`. Local runs/tests of
  `helmfile template`/`helmfile diff`/`helmfile apply` for these releases will fail if that sibling checkout is
  missing or out of date.
- This is a **temporary, intentional hack** for development. Before this feature branch can be merged, these
  `chart:` entries must be reverted back to versioned chart-repo references (the commented `radar/...` lines) once
  the corresponding charts are published to the `radar` chart repository, and a proper `_chart_version` set in
  `etc/base.yaml`.
- When adding new Ory-related config, follow the same conventions as other components (see
  `docs/development_guide.md`): install toggle, chart version and extra timeout in `etc/base.yaml`; secrets in
  `etc/base-secrets.yaml` and `bin/generate-secrets` if auto-generated; component defaults in an `etc/<name>/`
  folder if needed.

### Automated security patching

Implemented as the `platform-upgrade` skill (`.claude/skills/platform-upgrade/SKILL.md`): an ordered sweep over
the RADAR-base microservice fleet from
[Component inventory & local checkouts](#component-inventory--local-checkouts) above — shared libraries first,
then services — that scans each repo's dependencies with [Trivy](https://trivy.dev), auto-fixes same-major-version
CVE bumps, rebuilds to verify, and stages (but never pushes without confirmation) a release branch per repo. See
that file for the full, current procedure — don't duplicate it here; this section just points at it so the two
don't drift out of sync.
