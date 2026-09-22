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

### Known issues & fixes

*No entries yet — this section is populated as issues are discovered during migration runs. Do not delete this
placeholder line until the first real entry is added; its absence is how you'll know at a glance whether any
migration has actually been attempted yet.*

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
artifact until the library itself has migrated and republished. Process in this fixed order first, then the
remaining Kotlin-source service repos from the component table in any order:

1. `radar-jersey` — independent.
2. `radar-commons` — depends on `RADAR-Schemas`, but `RADAR-Schemas` is Avro schema definitions with no Kotlin
   source of its own; confirm that with the detection step before assuming it needs migrating at all.
3. `radar-commons-android` — depends on `RADAR-Schemas`; Android/Kotlin, so almost certainly in scope — but note
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
