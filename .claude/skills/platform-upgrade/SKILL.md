---
name: platform-upgrade
description: Automated security-patch sweep across the RADAR-base microservice fleet — checks out each component repo (shared libraries first), scans dependencies with Trivy for HIGH/CRITICAL CVEs, auto-fixes same-major-version bumps, rebuilds, and stages a release branch per repo for review. Use when the user asks to run a security patch sweep, vulnerability scan, or dependency upgrade across RADAR-base services, or invokes /platform-upgrade.
---

# Platform security upgrade sweep

Formalizes the "Automated security patching" workflow described in this repo's `AGENTS.md`. Read
`AGENTS.md`'s **Component inventory & local checkouts** section first, every run — it is the single source of
truth for the repo→local-path mapping and build tooling per component; do not hardcode a copy of it here, it will
drift.

## Scope

- If the user names specific components, scope the run to those repos plus any shared library they depend on
  (see ordering below).
- Otherwise, run the full fleet from the component inventory table, in the order below.
- This skill only patches **dependency CVEs** (library/package versions). Dockerfile base-image CVEs and IaC
  misconfigurations are out of scope unless the user explicitly asks for them.
- This skill never modifies `RADAR-Kubernetes` itself (this repo), except as an explicit, separately-confirmed
  follow-up if a patch requires a chart version or image tag bump in `etc/base.yaml`.

## Ordering — shared libraries first

A shared library's own vulnerabilities should be patched before scanning its consumers, so process in this fixed
order first, then the remaining service repos from the component table in any order:

1. `RADAR-Schemas` — no dependency on the others below.
2. `radar-jersey` — independent.
3. `radar-commons` — depends on `RADAR-Schemas` (`radar-schemas-commons`).
4. `radar-commons-android` — depends on `RADAR-Schemas` (`radar-schemas-commons`).

Caveat: patching and releasing a shared library does **not** automatically bump the version pinned in a
consumer's `build.gradle(.kts)` — that pin is just another dependency the consumer's own Trivy pass may or may not
flag (it will only show up once the new library version is published and resolvable, and only if it's actually
outdated/vulnerable per Trivy). Don't assume a consumer picked up the library fix; check its own scan results.

## Prerequisites

Before starting, confirm the tools this skill needs are available: `trivy`, `gh` (authenticated), `git`. Confirm
each target repo is checked out as documented in AGENTS.md's component table — verify the directory exists and
`git remote get-url origin` matches the expected GitHub repo before touching anything (that parent directory has
many unrelated repos; never operate on a directory that doesn't match).

## Per-repo procedure

Repeat the following for each repo in the resolved order. Report progress as you go — don't batch silent work
across many repos before saying anything.

### 1. Resolve & verify checkout

Resolve the local path from the AGENTS.md table. `cd` there and confirm `git remote get-url origin` matches the
expected `RADAR-base/<repo>` (or, for `RADAR-JDBC-Connector`, the `origin` remote specifically — it also has an
`upstream` remote pointing at `confluentinc/kafka-connect-jdbc`, don't confuse the two). Run `git status`; if
there's uncommitted work already sitting in the checkout, stop and flag it to the user rather than branching over
it.

### 2. Base the release branch on the most current upstream ref

1. `git fetch origin` (and `dev`/`develop` if present).
2. Determine the default branch: `gh repo view RADAR-base/<repo> --json defaultBranchRef -q .defaultBranchRef.name`.
3. Check for a `dev`/`develop` branch: `git ls-remote --heads origin dev develop`.
4. If one exists, compare it against the default branch: `git rev-list --count origin/<default>..origin/dev`. If
   that count is > 0, `dev` has commits not yet on the default branch — base the release branch on `dev`.
   Otherwise base it on the default branch.
5. Create the release branch off the chosen base: `security/patch-YYYY-MM-DD` (today's date; if a branch with that
   name already exists from an earlier run today, suffix `-2`, `-3`, …).

### 3. Make the real dependency graph visible to Trivy

- **Gradle repos**: Trivy's Java/Gradle scanner reads `gradle.lockfile`. If the repo doesn't already commit
  lockfiles (check for `dependencyLocking` in `build.gradle(.kts)`/`settings.gradle(.kts)` or existing
  `*.lockfile` files), generate them temporarily without touching the repo's own build files, using a Gradle init
  script:
  ```
  echo 'allprojects { dependencyLocking { lockAllConfigurations() } }' > /tmp/lock-init.gradle
  ./gradlew -I /tmp/lock-init.gradle dependencies --write-locks
  ```
  Regenerate these lockfiles after every dependency bump in step 6 (so re-scans in step 7 see the new resolved
  graph), and delete the generated `*.lockfile` files at the end of the repo's run if the project doesn't
  otherwise commit them — they're a scanning aid, not a deliberate project change.
- **npm repos** (`radar-home`, `radar-self-enrolment-ui`, the Angular side of `ManagementPortal`): a
  `package-lock.json` should already exist; if not, `npm install --package-lock-only`.
- Run the scan: `trivy fs --scanners vuln --severity HIGH,CRITICAL -f json -o trivy-report.json .`

### 4. Classify every HIGH/CRITICAL finding

For each finding in `trivy-report.json`, determine:

- **Fixable?** `FixedVersion` is present and non-empty. If absent, bucket as `NO-FIX-AVAILABLE`.
- **Direct or transitive?**
  - npm: trust Trivy's package relationship info from the lockfile (direct vs indirect).
  - Gradle: a `gradle.lockfile` is a flat resolved list with no parent/child structure, so Trivy can't reliably
    tell you this for Gradle. Cross-check with `./gradlew <module>:dependencies --configuration runtimeClasspath`
    (or `compileClasspath` if that's where it resolves): a top-level entry in the tree is direct; a nested entry
    is transitive, and its immediate parent(s) in the tree is what you'd need to bump instead.
- **Target version**: among Trivy's `FixedVersion` candidate(s) — sometimes a comma/space-separated list spanning
  multiple release lines — pick the lowest one that shares the **same major version** as the currently resolved
  (installed) version.
  - Candidate found → bucket `AUTO-FIX`.
  - Only higher-major candidates exist → bucket `NEEDS-MAJOR-BUMP` (report only, never auto-applied — a major
    bump needs a human compatibility judgment call).

Show the user the classified table (CVE, package, severity, direct/transitive, current → target version, bucket)
for this repo before proceeding to step 5.

### 5. Apply `AUTO-FIX` bumps

- **Direct dependency**: edit the version wherever it's actually declared — `build.gradle(.kts)`,
  `gradle.properties`, or a version catalog (`gradle/libs.versions.toml`) for Gradle; `package.json` (plus
  `npm install <pkg>@<target>` to update the lockfile) for npm.
- **Transitive dependency**: prefer bumping the direct parent dependency identified in step 4 to a version whose
  own manifest pulls in a fixed transitive version — verify with
  `./gradlew <module>:dependencyInsight --dependency <transitive-pkg> --configuration runtimeClasspath` (Gradle)
  after the trial bump, or `npm ls <pkg>` (npm) — rather than force-pinning the transitive version directly. Only
  fall back to an explicit override (Gradle `constraints {}` / `resolutionStrategy.force`, npm
  `overrides`/`resolutions`) if no parent bump resolves it, and flag that repo's report entry as
  `forced transitive override — verify compatibility` so the user knows it wasn't a clean parent-version bump.
- Apply all `AUTO-FIX` bumps for the repo together as one batch — failures are isolated in step 6, not avoided by
  going one at a time.

### 6. Build, verify, and re-scan

1. Build: `./gradlew build` (Gradle — full build including tests, since a dependency bump can break behavior that
   only tests catch) or `npm ci && npm run build && npm test --if-present` (npm).
2. **Build succeeds** → regenerate the lockfile (per step 3) and re-run the Trivy scan on the patched tree.
   Confirm (a) the CVEs you targeted are gone, and (b) no *new* HIGH/CRITICAL finding appeared — a bump can pull
   in a newer major of some other transitive dependency that's itself vulnerable. Treat any new HIGH/CRITICAL
   finding as a failure of that specific bump and handle it via the bisection step below.
3. **Build fails** (or step 2's re-scan surfaces a new HIGH/CRITICAL) → bisect: revert the batched bumps one at a
   time and rebuild after each revert, until the build passes, to isolate exactly which bump(s) are responsible.
   Keep every bump that wasn't implicated (build passes without it reverted). For the implicated bump(s): **stop
   and ask the user for advice** — show the compile/test error (or the new CVE), which CVE the bump was meant to
   fix, and the version(s) involved. Do not guess at a workaround or silently drop the bump; it must remain
   visible in the final report as `BUILD-BLOCKED — awaiting user decision`, not disappear.

### 7. Commit — do not push without confirmation

Commit the surviving bumps on the release branch with a message listing the CVEs fixed. Pushing the branch and
opening a PR (`gh pr create`) are visible, hard-to-reverse actions against a shared repo — do not do either
without showing the user this repo's diff and final vulnerability table first and getting explicit confirmation
to proceed, per repo.

### 8. Repo-level report

Emit the step-4 table with final status per finding: `AUTO-FIXED` / `NEEDS-MAJOR-BUMP (manual)` /
`BUILD-BLOCKED (awaiting advice)` / `NO-FIX-AVAILABLE`.

## Cross-repo final summary

After the ordered sweep (or the user-scoped subset) completes, produce one consolidated table across all repos
processed, so the user can see the overall sweep outcome — and everything still awaiting a push/PR decision or
user advice — at a glance.

## Safety notes

- Never push or open a PR without explicit per-repo confirmation (step 7).
- Never touch `RADAR-Kubernetes` itself as part of the automated sweep; a chart-version/image-tag bump here is a
  separate, explicitly-confirmed follow-up.
- Run `git status` before any branch/checkout operation in a sibling repo, exactly as you would in this one —
  never discard uncommitted work that predates this skill run without flagging it to the user first.
- If a repo's checkout doesn't match its expected `origin` remote, or the sibling directory is missing entirely,
  skip it and report that rather than guessing at a different local path.
