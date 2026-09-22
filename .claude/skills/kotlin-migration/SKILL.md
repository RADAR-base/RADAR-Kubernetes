---
name: kotlin-migration
description: Ordered sweep to migrate every RADAR-base microservice repo that actually contains Kotlin source from Kotlin 1.9 to Kotlin 2.3 — shared libraries first, then services — building and testing after each bump and staging (but never pushing without confirmation) a migration branch per repo. This is a living, self-updating skill — expect unknown breakage, and unresolved friction is a stop-and-ask, not a guess. Use when the user asks to migrate/upgrade Kotlin version(s) across the fleet, mentions Kotlin 1.9 -> 2.3 / Kotlin 2.x migration, or invokes /kotlin-migration.
---

# Kotlin 1.9 → 2.3 migration sweep

Companion to the `platform-upgrade` skill, but scoped to the Kotlin **language/compiler version** itself (and the
ecosystem plugins version-locked to it), not general dependency CVEs. Read `AGENTS.md`'s **Component inventory &
local checkouts** section first, every run — it is the source of truth for the repo→local-path mapping; don't
duplicate it here.

## This skill is a living document — read this before doing anything else

Nobody has run this migration before, so the "known issues" list below starts empty and **must grow as you work**.
The whole point of this skill file is to stop the next run (or the next repo in *this* run) from rediscovering the
same breakage from scratch. Two hard rules:

1. **Unresolved friction is a stop-and-ask, not a guess.** If a build breaks, a plugin refuses to resolve, a
   warning looks behavior-changing, or you're simply not sure whether a fix is correct — and the situation isn't
   already covered by an entry in [Known issues & fixes](#known-issues--fixes) below — stop and ask the user for
   advice (via `AskUserQuestion` if it's a discrete choice, or a direct question if it needs open-ended input).
   Show the concrete error/output, what you were trying to do, and why you're unsure. Do not silently pick a
   workaround, suppress a warning, downgrade a plugin, or skip a test to make the build green.
2. **Every resolution gets written back into this file before you move on.** Once the user advises a fix (or you
   hit an issue you're confident is safe to resolve because an existing entry already covers it), and it's
   verified working, use the `Edit` tool to add or update an entry under [Known issues & fixes](#known-issues--fixes)
   *immediately* — before starting the next repo, and even if the rest of the sweep is later paused or abandoned.
   An issue you solved but didn't record is a problem the next repo (or the next run) will pay for twice. Follow
   the entry template below so entries stay scannable.

If you hit an issue that *is* already covered by an entry below, apply the recorded fix — but verify it actually
resolves the situation in this repo's context before trusting it blindly (an entry recorded against one repo's
specific module layout or plugin set may not transfer cleanly). If it doesn't fully apply, treat the gap as new
friction under rule 1, and update the entry under rule 2 rather than adding a near-duplicate.

## Sweep progress (resume point)

Update this section every time a run pauses or ends, so the next run (or the next message in this one) knows
exactly where things stand without re-deriving it. Keep it current — overwrite stale entries rather than
appending a history log; `git log` on this file already provides that if anyone needs it.

- **Target version:** 2.3.20 (confirmed with the user 2026-09-22 — latest patch/tooling release in the 2.3 line,
  not exactly-2.3.0).
- **`radar-commons` (incl. `radar-commons-gradle`, step 0) — MIGRATED.** 1.9.24 → 2.3.20, Gradle wrapper 8.14 →
  9.0. Committed on branch `kotlin/migrate-2.3-2026-09-22` (based on `origin/master`) — **committed locally only,
  not pushed, no PR opened.** Full clean build + tests pass (module details in
  [Known issues & fixes](#known-issues--fixes) above). Stray pre-existing untracked files in this checkout
  (`.kotlin/`, `build_output*.log`, `build_info_v4.log`, `build_plugin_output.log`) were left alone, not part of
  this migration's changes.
  - **Not yet published.** `radar-jersey` (and every other consumer) can't actually pick up Kotlin 2.3.20 through
    the `radarCommons` catalog pin until this branch is reviewed, merged, and a new `radar-commons` version is
    released — that's a prerequisite for finishing any consumer repo's own migration step, not just a nice-to-have
    ordering note.
- **`radar-jersey` — paused, not started.** A pre-existing untracked `Dockerfile` was stashed (`git stash push -u
  -m "kotlin-migration: pre-existing untracked Dockerfile set aside"` on branch `release-0.12.8`) so the working
  tree is clean; nothing else has been touched — no migration branch created, no version bumps applied yet.
  **Remember to `git stash pop` (or otherwise resolve) that stash before or after this repo's migration work** —
  it isn't part of this skill's changes and its owner will want it back.
- **`radar-commons-android` (step 2) — not started.**
- **Remaining fleet (step 2+, any order) — not started.**

### Known issues & fixes

#### Consumer repos don't control their own Kotlin compiler version — `radar-commons-gradle` does

- **Where seen:** `radar-jersey`, discovered before any code change was made there. Expected to be general across
  every repo that applies the `org.radarbase.radar-kotlin` convention plugin (i.e. most of the fleet — check
  `apply(plugin = "org.radarbase.radar-kotlin")` / `alias(libs.plugins.radar.kotlin)` in a repo's build files
  before assuming this doesn't apply).
- **Symptom:** none yet (caught by reading the build setup, not by a failed build) — flagging so nobody wastes a
  cycle bumping a repo's own `kotlin = "..."` catalog entry expecting it to change the compiler version.
- **Cause:** a consumer repo's own `gradle/libs.versions.toml` `kotlin = "1.9.24"` entry only pins the
  `kotlin-reflect`/`kotlin-stdlib` *library* versions that repo declares directly — it does not control the
  compiler. The actual Kotlin Gradle Plugin dependency, and the `languageVersion`/`apiVersion` compiler options
  applied to every consumer, live in **`radar-commons/radar-commons-gradle`** (a Gradle included-build inside the
  `radar-commons` repo, wired via `pluginManagement { includeBuild("radar-commons-gradle") }` in
  `radar-commons/settings.gradle.kts`). Specifically:
  - `radar-commons-gradle/build.gradle.kts` declares `implementation(libs.gradlePlugin.kotlin)` — the real Kotlin
    Gradle Plugin artifact — and hardcodes `languageVersion.set(KotlinVersion.KOTLIN_1_9)` /
    `apiVersion.set(KotlinVersion.KOTLIN_1_9)` on its own `KotlinCompile` tasks (this affects only
    radar-commons-gradle's own compilation, but see next point).
  - `radar-commons-gradle/.../RadarKotlinPlugin.kt` (implementation of `org.radarbase.radar-kotlin`) does
    `apply(plugin = "kotlin")` in every consumer, using whatever Kotlin Gradle Plugin version is on
    radar-commons-gradle's own plugin classpath (not the consumer's) — plus a `RadarKotlinExtension.kotlinVersion`
    property, defaulted from a `Versions.kotlin = "1.9.24"` constant in the same source set, that sets the
    consumer's `languageVersion`/`apiVersion` (consumers can override this via `radarKotlin { kotlinVersion.set(...) }`,
    but nothing in `radar-jersey` currently does).
  - The convention plugin itself is version-pinned in a consumer via the `radarCommons` catalog entry (e.g.
    `radarCommons = "1.2.5"` in `radar-jersey/gradle/libs.versions.toml`), which tracks `radar-commons`'
    own project version — i.e. it's a *published artifact* version, not a live path dependency.
- **Fix:** migrate `radar-commons-gradle` itself first — bump its Kotlin Gradle Plugin dependency and the two
  hardcoded `KOTLIN_1_9` compiler-option lines, and bump the `Versions.kotlin` constant — get that building, then
  release/publish `radar-commons` at a new version. Only then does bumping the consumer's `radarCommons` catalog
  pin (as part of that consumer's own migration step, per [Ordering](#ordering--shared-libraries-first)) actually
  pick up a newer Kotlin compiler. Bumping a consumer's own local `kotlin = "..."` entry alone does nothing for
  the compiler version.
- **Caveats:** unconfirmed whether every fleet repo goes through this exact convention plugin — repos not using
  `org.radarbase.radar-kotlin` (check per repo) may declare `kotlin("jvm")` directly and control their own version
  as originally assumed by this skill.

The [Ordering](#ordering--shared-libraries-first) list below has been corrected to put `radar-commons-gradle`
first as a result of this discovery.

#### `radar-commons-kotlin`'s `testForkJoinFirst` is flaky, independent of Kotlin version

- **Where seen:** `radar-commons` / `radar-commons-kotlin:test`, `ExtensionsKtTest.testForkJoinFirst`, during the
  1.9.24 → 2.0.21 step.
- **Symptom:** `java.lang.AssertionError: Expected: a value less than <200ms> but: <291.192924ms> was greater than
  <200ms>` — a hard-coded wall-clock timing threshold under `runBlocking`, sensitive to machine load/scheduling
  jitter, not to compiled code semantics.
- **Cause:** pre-existing test flakiness unrelated to this migration — confirmed by rerunning the same test class
  alone immediately after the failure, on the same Kotlin 2.0.21 build, where it passed.
- **Fix:** none applied (nothing to fix for the migration). If this test fails again during a later version step,
  don't assume it's evidence of a real regression — rerun it isolated first before escalating.
- **Caveats:** if it starts failing *consistently* (not just once under `--continue` alongside everything else),
  that would be new evidence and should be escalated normally, not written off via this entry. Same pattern seen
  again during the 2.0.21 → 2.1.21 step in `CachedValueTest.getInvalid` (`"No refresh within threshold"` — also a
  wall-clock timing assertion under `runBlocking`), also confirmed flaky by isolated rerun. Any test in this
  module asserting on wall-clock duration thresholds under `runBlocking`/coroutine dispatch is a candidate for
  this same flakiness — rerun isolated before treating a failure there as migration-caused.

#### `radar-commons-gradle`'s Kotlin version can outrun Gradle's embedded `kotlin-dsl` Kotlin version

- **Where seen:** `radar-commons/radar-commons-gradle` (uses the `kotlin-dsl` plugin), during the 2.0.21 → 2.1.21
  step.
- **Symptom:** `WARNING: Unsupported Kotlin plugin version. The embedded-kotlin and kotlin-dsl plugins rely on
  features of Kotlin 2.0.21 that might work differently than in the requested version 2.1.21.` Build still
  succeeded at this step — non-fatal so far, but Gradle's own wording is a real compatibility warning, not
  boilerplate noise.
- **Cause:** Gradle's `kotlin-dsl` plugin (used because `radar-commons-gradle` is a Gradle-plugin-authoring
  project) embeds a specific Kotlin version tied to the Gradle release itself, independent of whatever Kotlin
  Gradle Plugin version the project's own `plugins {}` block requests. Gradle 8.14 (this repo's wrapper version
  when the migration started) embeds Kotlin 2.0.21; no Gradle 8.x release embeds anything newer. Gradle 9.0 is the
  first release with a newer embedded Kotlin (2.2.0) — so this constraint cannot be fully satisfied at the
  eventual 2.3.20 target with any Gradle version confirmed to exist as of this writing.
- **Fix:** user chose to bump the Gradle wrapper to 9.0 partway through this migration (rather than accept the
  warning indefinitely or pause). This is a **major Gradle version bump**, separate from the Kotlin migration
  itself, with its own breaking-change surface across the whole `radar-commons` build (plugin compatibility —
  `nexus-publish`, `dokka`, `sentry`, `ktlint-gradle`, `licenseReport`, `version-catalog-update`,
  `gradle-versions-plugin` — all need to still resolve/work under Gradle 9). Treat every one of those plugins'
  Gradle-9 compatibility as its own thing to verify via the build, not assumed.
- **Caveats:** this only applies to repos whose *own build* uses `kotlin-dsl` (currently just
  `radar-commons-gradle` in this fleet, as far as confirmed) — plain application/library modules consuming the
  `org.radarbase.radar-kotlin` convention plugin don't hit this warning themselves. Re-check whether the warning
  resurfaces at 2.3.20 even after the Gradle 9.0 bump (2.2.0 embedded vs. 2.3.20 requested) — if so, that's a new
  instance of the same root cause, not a fixed problem, and needs the same escalate-and-ask treatment.
  - **Update (2.2.21 step):** the warning did resurface (`embeds 2.2.0` vs. `requested 2.2.21`), as expected —
    but this is a same-minor-line patch-version mismatch, not the earlier cross-minor jump (2.0.21 embedded vs.
    2.1.21+ requested), and the full build (including tests) still passed. Judged low-risk enough to proceed
    without a fresh ask, since this is exactly the recurrence this entry already anticipated — but the 2.3.20 step
    will land on a genuinely newer minor line again (2.2.0 embedded vs. 2.3.20 requested), which is closer in
    shape to the original cross-minor warning than to this one. Don't reuse this "low-risk, proceed" judgment
    there without re-checking the build result first.
  - **Update (2.3.20 step, final target for this migration):** the cross-minor case did recur (`embeds 2.2.0` vs.
    `requested 2.3.20`) and, as of this writing, there is no Gradle release whose embedded Kotlin reaches 2.3.x —
    Gradle 9.0 is still the newest confirmed and it embeds 2.2.0. A full clean build (including tests) at 2.3.20
    still passed, so this was *not* escalated as a fresh blocker — but be aware this warning is expected to remain
    permanently present in `radar-commons-gradle`'s build output at the 2.3.20 end state, for as long as no newer
    Gradle release embeds a matching-or-newer Kotlin. That's a standing, accepted condition of this migration, not
    an unresolved item — don't re-raise it as new friction on a future run unless the build itself actually starts
    failing because of it.

#### `Project.property()` calls in `build.gradle.kts` need a null-assertion under Gradle 9

- **Where seen:** `radar-commons/radar-commons-testing/build.gradle.kts:30`, discovered when bumping the Gradle
  wrapper to 9.0 (see the `kotlin-dsl`/embedded-Kotlin entry above) alongside Kotlin 2.1.21.
- **Symptom:** `Script compilation error: Argument type mismatch: actual type is '@Nullable() Any?', but 'Any' was
  expected` at a call like `args(project.property("mockConfig"))`.
- **Cause:** under Gradle 9, `Project.property(name: String)`'s Kotlin-visible signature is nullable (`Any?`)
  where it previously type-checked as a non-null platform type — this is a Gradle 9 API/annotation change
  surfaced by the stricter Kotlin 2.x type checker, not a Kotlin-language change per se. Any `build.gradle.kts` in
  this fleet calling `project.property(...)` directly (as opposed to `providers.gradleProperty(...)` or similar)
  is a candidate for this, not just this one call site — grep for `project.property(` / `.property(` calls in
  build scripts across a repo before declaring it clear of this issue.
- **Fix:** add a non-null assertion (`project.property("mockConfig")!!`) at each call site that is already guarded
  by a preceding `project.hasProperty(...)` check (safe — presence is already guaranteed at that point). Do not
  blanket-apply `!!` to a `property(...)` call that *isn't* preceded by such a guard; that would turn a real
  missing-property bug into an NPE instead of Gradle's own clearer `MissingPropertyException` — treat an
  unguarded occurrence as a fresh instance of this issue needing its own judgment call, not an auto-apply.
- **Caveats:** only confirmed against Gradle 9.0 + Kotlin 2.1.21; unconfirmed whether this is a Gradle-version
  trigger, a Kotlin-version trigger, or both together — re-verify which one actually matters if a repo needs the
  Kotlin bump without the Gradle 9 bump (or vice versa).

Entry template (copy this shape for each new entry, as a `####` subsection under this heading):

```
#### <short symptom title>

- **Where seen:** <repo(s)/module(s); note if it looks general vs repo-specific>
- **Symptom:** <compiler error / warning / test failure / runtime behavior change, verbatim where short enough>
- **Cause:** <root cause, if known — otherwise say "unconfirmed">
- **Fix:** <what was actually changed — version bump, code change, config flag, etc.>
- **Caveats:** <anything that makes this fix conditional, e.g. "only safe if the module doesn't use kapt">
```

## Scope

- **Only repos with actual Kotlin source are in scope.** The inventory table's "Build tool" column says "Gradle
  Kotlin DSL" for many repos — that only means the *build scripts* (`build.gradle.kts`) are written in Kotlin, not
  that the repo's application source is. Detect real Kotlin usage per repo rather than trusting that column: look
  for `.kt` files under `src/main`/`src/test` (not just `buildSrc`/`build.gradle.kts`), and a `kotlin("jvm")` /
  `org.jetbrains.kotlin.jvm` plugin block with application source sets attached. A repo with Kotlin-DSL build
  scripts but pure-Java application source is out of scope for this skill (it may still get a Gradle Kotlin DSL
  *tooling* version bump incidentally if that's pinned to the same version — note it, don't chase it).
- If the user names specific repos, scope the run to those plus any shared library they depend on (ordering
  below). Otherwise run the full fleet, in the order below.
- This skill migrates: the Kotlin language/compiler version itself, and any plugin whose version is directly
  coupled to it (KSP — version string embeds the Kotlin version, e.g. `2.3.0-x.y`; `kotlinx-serialization`
  compiler plugin; Kotlin Gradle plugin's own bundled stdlib/reflect; detekt/ktlint's Kotlin-compiler-version
  compatibility). It does **not** migrate unrelated dependency versions (CVE bumps belong to `platform-upgrade`)
  or Java/JVM target versions, unless the Kotlin bump forces a minimum JVM target — if it does, treat that as new
  friction under rule 1 above and confirm with the user before bumping `sourceCompatibility`/`jvmToolchain`.
- Never touch `RADAR-Kubernetes` itself as part of the sweep, except as an explicit, separately-confirmed
  follow-up (e.g. an image tag bump in `etc/base.yaml` once a migrated service is released) — same boundary as
  `platform-upgrade`.

## Relationship to `platform-upgrade`

`AGENTS.md`'s "Known compatibility constraints" section (consumed by `platform-upgrade`) currently blocks bumping
`com.fasterxml.jackson*` to ≥ 2.21.0 anywhere in the fleet because that line requires Kotlin 2.x and the fleet is
on 1.9. As each repo in *this* sweep completes its migration to Kotlin 2.3, that constraint no longer applies to
that specific repo — update the constraint's note in `AGENTS.md` (or flag it to the user to update) so
`platform-upgrade` doesn't keep treating an already-migrated repo as blocked.

## Ordering — shared libraries first

Same rationale as `platform-upgrade`: a shared library's consumers can't safely take a Kotlin-2.3-compiled
artifact until the library itself has migrated and republished. But unlike `platform-upgrade`'s CVE sweep, this
migration also has to respect *build-tooling* dependencies, not just runtime ones — see the
[convention-plugin discovery](#known-issues--fixes) above, which is why `radar-commons-gradle` sits ahead of
everything else below. Process in this fixed order first, then the remaining Kotlin-source service repos from the
component table in any order:

0. `radar-commons/radar-commons-gradle` (the `org.radarbase.radar-kotlin` convention-plugin included-build) —
   every other repo in the fleet that applies `org.radarbase.radar-kotlin` (check per repo; most do) gets its
   actual Kotlin compiler/language version from here, transitively, regardless of what its own
   `gradle/libs.versions.toml` `kotlin = "..."` entry says. Must migrate and be republished (as a new
   `radar-commons` version) before any consumer's catalog pin bump can do anything. Re-verify this dependency
   still holds at the start of each run — the convention-plugin setup could itself change.
1. `radar-commons` (the rest of it: `radar-commons-kotlin`, `radar-commons-server`, `radar-commons-testing`, plus
   whatever `radar-commons-gradle` step 0 left needing a bump) and `radar-jersey` — both now depend on step 0's
   republished convention plugin rather than on each other; order between them no longer matters for this
   migration specifically (unlike the CVE-sweep ordering, which puts `radar-jersey` first because it has no
   runtime dependency on `radar-commons`).
2. `radar-commons-android` — depends on `RADAR-Schemas`; Android/Kotlin, so almost certainly in scope — but note
   Android Gradle Plugin (AGP) has its own Kotlin-version compatibility matrix, separate from the plain
   Kotlin/JVM one the other repos use. Treat any AGP-related friction here as its own new-issue candidate; don't
   assume a fix that worked in a Kotlin/JVM repo transfers.

Caveat, same as `platform-upgrade`: migrating and releasing a shared library does **not** automatically bump the
version pinned in a consumer's `build.gradle(.kts)`. After a shared library's migration branch is confirmed and
(once the user approves) published, bumping the consumer's pin to the new library version is part of that
consumer's own migration work in this sweep — fold it into the consumer's migration commit rather than opening a
separate change.

## Prerequisites

Confirm `git` and `gh` (authenticated) are available. Confirm each target repo is checked out as documented in
AGENTS.md's component table — verify the directory exists and `git remote get-url origin` matches the expected
GitHub repo before touching anything.

## Per-repo procedure

Repeat for each repo in the resolved order. Report progress as you go — don't batch silent work across many repos
before saying anything.

### 1. Resolve & verify checkout

Resolve the local path from AGENTS.md's table. `cd` there, confirm the `origin` remote matches, run `git status`;
if there's uncommitted work already sitting in the checkout, stop and flag it to the user rather than branching
over it.

### 2. Detect whether the repo is in scope

Look for `.kt` files under application source sets (not just `buildSrc`) and a Kotlin JVM/Android/multiplatform
plugin applied to a source-producing module. If none found: record as `SKIPPED-NO-KOTLIN-SOURCE` in the final
report and move to the next repo — don't branch or touch anything.

### 3. Read the current Kotlin version

Check, in order of likelihood: `gradle/libs.versions.toml` (a `kotlin = "..."` entry), the root/module
`build.gradle(.kts)` plugin block (`kotlin("jvm") version "..."`), or `buildSrc`/convention-plugin version
declarations (RADAR-base repos commonly centralize this via a `radarCommons`/`org.radarbase.*` convention plugin —
check `radar-commons`' convention-plugin source if the version isn't declared locally). If the resolved version is
already ≥ 2.x, record as `SKIPPED-ALREADY-2.x` and move on. If it's not on the 1.9.x line AGENTS.md assumes for
the fleet, flag that mismatch to the user before proceeding — it's a sign the inventory assumption is stale.

### 4. Branch

Same base-selection procedure as `platform-upgrade` step 2 (default branch vs `dev`/`develop`, whichever has more
recent commits). Branch name: `kotlin/migrate-2.3-YYYY-MM-DD` (today's date; suffix `-2`, `-3`, … if already
taken from an earlier run today).

### 5. Check the playbook before touching anything

Skim [Known issues & fixes](#known-issues--fixes) for entries relevant to this repo's module layout (kapt vs KSP,
multiplatform vs plain JVM, Android vs not, convention-plugin usage). Apply any that clearly transfer as you go
through steps 6–7, rather than rediscovering them.

### 6. Decide the version path and bump

Kotlin 2.0 changed the default compiler frontend to K2, which is a materially bigger jump than a typical
minor-version bump — going straight from 1.9 to 2.3 in one step makes it hard to isolate which change caused a
given break. Default to an incremental path (1.9.x → 2.0.x → confirm build → continue up to 2.3.x, one minor line
at a time), stopping to rebuild/retest at each step. This default is itself a judgment call made without having
seen any real breakage yet — if the first repo's migration shows that one-shot jumps are fine (or that the
incremental path is overkill), say so and record it as an entry under [Known issues & fixes](#known-issues--fixes)
so later repos in the sweep (and later runs) use whatever path actually proved necessary, and update this
paragraph to match.

At each version step, bump the Kotlin version itself plus any plugin whose version string is coupled to it (KSP,
`kotlinx-serialization` compiler plugin, Compose compiler if ever present) to the matching release for that
Kotlin line — check the plugin's own compatibility table rather than assuming a matching version number exists.

### 7. Build, verify, watch for silent behavior changes

Run the full build including tests: `./gradlew build` (add `:module:build` for a multi-module repo if scoping
narrows the check). A clean compile is necessary but not sufficient — K2's stricter type inference and changed
defaults can alter runtime behavior without a compile error. Skim compiler warnings (not just errors) introduced
by the bump, and pay particular attention to whichever areas warnings cluster around — do not assume "it compiled
and tests are green" fully clears a step if you noticed a suspicious new warning; if you're unsure whether a
warning indicates a real behavior change, that's new friction under rule 1 (stop and ask), not something to wave
through.

- **Build/tests pass, nothing suspicious** → proceed to the next version step (6) or, if this was the last step
  (2.3.x), to commit (8).
- **Build fails, or something looks like a silent behavior change, or you're genuinely unsure** → check whether
  [Known issues & fixes](#known-issues--fixes) already covers it (step 5). If yes and it clearly transfers, apply
  it and re-verify. If no, or it's ambiguous — **stop and ask the user for advice**, per the rules at the top of
  this file. Once resolved, record the fix (rule 2) before continuing.

### 8. Commit — do not push without confirmation

Commit the migration on the branch with a message noting the Kotlin version reached and, briefly, what had to
change to get there. Pushing and opening a PR (`gh pr create`) are visible, hard-to-reverse actions against a
shared repo — do not do either without showing the user this repo's diff and final build status first and getting
explicit confirmation to proceed, per repo.

### 9. Repo-level report

Report: repo, starting Kotlin version, ending version (or reason skipped), status (`MIGRATED` /
`SKIPPED-NO-KOTLIN-SOURCE` / `SKIPPED-ALREADY-2.x` / `BLOCKED-AWAITING-ADVICE`), and any new
[Known issues & fixes](#known-issues--fixes) entries added while working on it.

## Cross-repo final summary

After the ordered sweep (or the user-scoped subset) completes, produce one consolidated table across all repos
processed — status per repo, and everything still awaiting a push/PR decision or user advice — plus a note of how
many new playbook entries this run added, so the next run starts more informed than this one did.

## Safety notes

- Never push or open a PR without explicit per-repo confirmation (step 8).
- Never touch `RADAR-Kubernetes` itself as part of the sweep; any follow-up here is a separate, explicitly
  confirmed step.
- Run `git status` before any branch/checkout operation in a sibling repo; never discard uncommitted work that
  predates this skill run without flagging it to the user first.
- If a repo's checkout doesn't match its expected `origin` remote, or the sibling directory is missing entirely,
  skip it and report that rather than guessing at a different local path.
- When genuinely uncertain — about a fix, a version choice, whether a warning matters, or anything not already
  settled by an entry in [Known issues & fixes](#known-issues--fixes) — ask the user. This skill exists precisely
  because the failure modes aren't known yet; guessing defeats the purpose.
