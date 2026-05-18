# radarctl — CLI Design Spec

**Date:** 2026-05-18  
**Status:** Approved  
**Scope:** v1 — init wizard, deploy, status

---

## Overview

`radarctl` is a Go CLI tool that improves the deployment experience for the RADAR-Kubernetes stack. It lives inside the RADAR-Kubernetes repository at `cli/` and shells out to existing tools (kubectl, helm, helmfile) rather than reimplementing their logic.

**Target users:** Cluster operators, researchers setting up new studies, and RADAR-base developers — each served by different modes and commands.

**v1 scope:** `radarctl init`, `radarctl deploy`, `radarctl status`, `radarctl diagnose`, `radarctl validate`

---

## Architecture

### Directory structure

```
cli/
├── main.go
├── cmd/
│   ├── root.go        # root command, global flags (--output, --context, --yes)
│   ├── init.go        # radarctl init (wizard, interactive, expert modes)
│   ├── deploy.go      # radarctl deploy
│   ├── status.go      # radarctl status
│   ├── diagnose.go    # radarctl diagnose (full diagnostic snapshot)
│   └── validate.go    # radarctl validate (config validation only)
├── pkg/
│   ├── config/
│   │   ├── loader.go      # read/write base.yaml, production.yaml, secrets.yaml
│   │   ├── validator.go   # validate completeness and consistency
│   │   └── features.go    # feature flag → config expansion
│   ├── wizard/
│   │   ├── wizard.go      # orchestrates wizard flow and mode selection
│   │   ├── questions.go   # question definitions and branching logic
│   │   └── writer.go      # writes collected answers to config files
│   ├── helmfile/
│   │   └── runner.go      # shells out to helmfile (sync, diff, template)
│   └── kubectl/
│       └── runner.go      # shells out to kubectl (get pods, logs, describe, ingress)
└── go.mod
```

### Key dependencies

| Dependency | Purpose |
|------------|---------|
| [Cobra](https://github.com/spf13/cobra) | Command structure and flags |
| [Huh](https://github.com/charmbracelet/huh) | Interactive terminal prompts and wizard forms |
| [Viper](https://github.com/spf13/viper) | Config file reading/writing |
| [go-yaml](https://github.com/goccy/go-yaml) | YAML manipulation for config generation |
| [pterm](https://github.com/pterm/pterm) | Progress bars, spinners, status tables |

### Design principles

- **Shell out, don't reimplement** — use kubectl, helm, helmfile for all cluster operations
- **Structured output everywhere** — every command supports `-o json` for agent/CI consumption
- **Fail loudly with context** — errors include the release name, pod name, and relevant logs
- **Progressive disclosure** — wizard mode for beginners, expert mode for power users
- **Resumable** — wizard state saved to `.radarctl-state.yaml` (gitignored)

---

## Commands

### `radarctl init`

Guided setup with three modes selectable at the start:

- **Wizard** — high-level questions that expand into full config (recommended for first-time users)
- **Interactive** — field-by-field prompts for all config values
- **Expert** — skip setup, validate existing config files only

#### Phase 0: Prerequisites check

Runs before any wizard prompts. Checks:

- Tool presence and minimum version: `kubectl`, `helm` (v3), `helmfile` (v0.169.1), `helm-diff` (v3.9.12), `yq` (v4.44.3), `java`, `openssl`, `git`
- Kubernetes connectivity and version (warn if outside v1.30–v1.33)
- Helm repositories present and reachable
- Cluster resources (node count, allocatable CPU/memory — warn if below recommended)

Example output:
```
Checking prerequisites...

  ✓ kubectl        v1.30.2    (context: my-cluster, nodes: 3 Ready)
  ✓ helm           v3.15.1
  ✓ helmfile       v0.169.1
  ✓ helm-diff      v3.9.12
  ✓ yq             v4.44.3
  ✗ java           not found  (required for keystore generation)
  ✗ openssl        not found  (required for secret generation)

2 prerequisites missing. Install them before continuing?
→ Show install instructions  /  Exit
```

On failure: show per-OS install instructions. Bypassable with `--skip-prereqs` for advanced users.

#### Wizard mode flow

```
1. Cluster basics
   - Server hostname              → server_name
   - Maintainer email             → maintainer_email
   - Kubernetes context           → kubeContext

2. Deployment profile
   - Production / Staging / Local dev
   → dev: auto-applies mods/minimal + mods/localdev + mods/disable_tls + mods/fast_deploy
   → staging: applies mods/minimal
   → production: no mods applied by default

3. Kafka
   - Local Kafka or Confluent Cloud?
   → Confluent: prompt for bootstrap URL, API key, API secret
   → Local: prompt for broker count and replica factors (or accept defaults)

4. Data sources  (yes/no toggle per source)
   - Fitbit        → client ID + secret
   - Garmin        → consumer key + secret
   - REDCap        → API URL + token
   - Upload portal → no extra config needed

5. Storage
   - Local Minio or external S3?
   → External S3: prompt for bucket, region, access key, secret key

6. Authentication
   - Enable Ory Hydra/Kratos (OAuth2/OIDC)? → yes/no

7. Monitoring & logging
   - Enable Prometheus + Grafana?       → yes/no
   - Enable Graylog + Elasticsearch?    → yes/no

8. Review & confirm
   - Display summary of all choices
   - Show which files will be written (production.yaml, secrets.yaml, environments.yaml)
   - Confirm before writing
```

#### What gets written

- `etc/production.yaml` — all non-secret values from wizard answers
- `etc/secrets.yaml` — auto-generated passwords + user-provided credentials
- `environments.yaml` — rendered from `environments.yaml.tmpl` with selected mods
- Calls existing `bin/keystore-init` for Java keystores

Re-running `radarctl init` on an existing deployment pre-fills all prompts from current config values.

#### Interactive mode

Presents every config field one by one with its current value pre-filled. No branching logic — the user sees and can set all options explicitly. Best for operators who know the stack well and want full control.

#### Expert mode

Skips all prompts. Runs the prerequisites check and config validator only, then exits. Equivalent to running `radarctl validate` directly.

#### Resumability

Wizard progress is saved to `.radarctl-state.yaml` in the repo root (added to `.gitignore`). If the wizard is interrupted, re-running `radarctl init` offers to resume from where it left off.

---

### `radarctl deploy`

Wraps `helmfile sync` with validation, a change preview, live progress, and post-deploy health checks.

#### Flags

```bash
radarctl deploy                        # full sync
radarctl deploy --diff                 # preview changes only (helmfile diff)
radarctl deploy --dry-run              # render templates, no apply
radarctl deploy --selector cert-manager  # deploy specific release(s)
radarctl deploy --yes                  # skip confirmation prompt
radarctl deploy --no-atomic            # disable rollback on failure
radarctl deploy -o json                # structured output
```

#### Deploy flow

```
1. Pre-deploy validation
   - Run config validator; abort on errors, warn on suspicious values
   - Warn if secrets.yaml contains placeholder values ("change_me", "secret")

2. Change preview
   - Run helmfile diff, summarise: "3 releases will be updated, 1 installed, 0 removed"
   - Prompt for confirmation (skipped with --yes)

3. Live progress display
   ┌─────────────────────────────────────────────┐
   │  Deploying RADAR stack...                   │
   │                                             │
   │  ✓ cert-manager          installed  (12s)   │
   │  ✓ kube-prometheus-stack installed  (45s)   │
   │  ⠸ mongodb               syncing...         │
   │  ○ kafka                 waiting            │
   │  ○ radar-appserver       waiting            │
   └─────────────────────────────────────────────┘

4. Post-deploy health check
   - Poll kubectl rollout status for each deployed release
   - Aggregate failures and surface with pod name + error reason
   - On failure: print last 20 lines of logs from failing pod automatically
   - Final summary: "22 releases healthy, 1 degraded → radar-appserver"
```

#### Exit codes

| Code | Meaning |
|------|---------|
| 0 | All releases healthy |
| 1 | One or more releases failed (fixable) |
| 2 | Config invalid (needs human input) |
| 3 | Prerequisites missing |
| 4 | Cluster unreachable |

`--atomic` is on by default (rolls back failed releases).

---

### `radarctl status`

Live health dashboard for all deployed releases.

#### Default output

```
RADAR Stack Status — my-cluster (production)
Last updated: 2026-05-18 14:32:01

INFRASTRUCTURE
  ✓ cert-manager            healthy    1/1 pods    v1.14.2
  ✓ kube-prometheus-stack   healthy    3/3 pods    v55.5.0
  ✓ nginx-ingress           healthy    2/2 pods    v4.10.0

KAFKA
  ✓ zookeeper               healthy    3/3 pods
  ✓ kafka                   healthy    3/3 pods
  ✓ schema-registry         healthy    1/1 pods
  ✗ ksql-server             degraded   0/1 pods    CrashLoopBackOff
    └─ radar-ksql-0: OOMKilled — last 3 restarts in 10m

STORAGE
  ✓ mongodb                 healthy    3/3 pods
  ✓ postgresql              healthy    2/2 pods
  ✓ redis                   healthy    1/1 pods
  ✓ minio                   healthy    1/1 pods

RADAR SERVICES
  ✓ radar-appserver         healthy    2/2 pods
  ✓ management-portal       healthy    1/1 pods
  ⚠ radar-fitbit-connector  warning    1/2 pods    1 pod pending (resource pressure)

IDENTITY
  ✓ kratos                  healthy    1/1 pods
  ✓ hydra                   healthy    1/1 pods

Summary: 19 healthy  1 degraded  1 warning
```

#### Flags

```bash
radarctl status --watch              # refresh every 10s
radarctl status --component kafka    # filter to one group
radarctl status --show-urls          # print ingress URLs for each service
radarctl status -o json              # structured output for agents/CI
```

#### Data sources (all via kubectl)

- `kubectl get pods` — pod counts and states
- `kubectl describe pod` — events for degraded pods
- `kubectl logs --tail=20` — recent logs for failing pods
- `kubectl get ingress` — service URLs

---

### `radarctl diagnose`

Collects a complete diagnostic snapshot in one command — designed for agentic loops and CI debugging.

```bash
radarctl diagnose -o json
```

Output includes: config validation results, pod states for all releases, recent Kubernetes events, and log tails for any failing pods. An AI agent can consume this single JSON blob and propose targeted fixes.

---

### `radarctl validate`

Runs config validation standalone, without deploying.

```bash
radarctl validate              # validate production.yaml + secrets.yaml
radarctl validate -o json      # structured output
```

Checks: required fields present, no placeholder values, feature dependencies satisfied (e.g. Fitbit enabled but client ID empty), mod compatibility (e.g. confluent_kafka + local kafka both enabled).

---

## Agent-friendliness

All commands output structured JSON with `--output json` / `-o json`. Consistent exit codes allow branching. The intended agentic loop:

```
radarctl deploy -o json
  → agent reads result
  → if failed: radarctl diagnose -o json
  → agent reads logs/events, proposes fix
  → agent applies fix
  → radarctl deploy --yes -o json
  → repeat until healthy or escalate to human
```

JSON schema for deploy result:
```json
{
  "status": "degraded",
  "releases": [
    { "name": "mongodb", "status": "healthy", "duration_s": 34 },
    { "name": "radar-appserver", "status": "failed", "error": "CrashLoopBackOff",
      "pod": "radar-appserver-6d4f9b-xkp2q", "logs": "..." }
  ],
  "summary": { "healthy": 21, "failed": 1, "pending": 0 }
}
```

---

## Out of scope for v1

- Upgrade orchestration (`radarctl upgrade`) — future version
- Secret manager integration (Vault, AWS Secrets Manager) — future version
- Plugin architecture — monolithic binary for v1
- Direct Kubernetes API calls — kubectl shell-outs only for v1
- Windows support — macOS and Linux only for v1
