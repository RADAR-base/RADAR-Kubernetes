# E2E tests: Behave & Playwright

This folder contains end-to-end tests for RADAR-Kubernetes using:
- Behave (BDD)
- Playwright (browser automation)

## Prerequisites
- Python 3.10+ installed
- Node not required for tests, but Playwright needs its browsers installed once

## Setup

Create and activate a virtual environment. 
For macOS/Linux:
```
python3 -m venv test/features/venv
source test/features/venv/bin/activate
```

Install Python dependencies:
```
pip install -r test/features/requirements.txt
```

Install Playwright browsers (one time per machine/venv):
```
python -m playwright install
```

## Configuration

Optional runtime configuration via Behave userdata:
- Timeout in seconds for UI waits (default 20):
  - behave -D timeout_s=30
- Dev mode (default false):
  - behave -D dev_mode=true

Note: The default Playwright launch mode is headless as defined in `test/features/environment.py`. 
If you need headed runs, with a display in a browser, set the dev_mode to `dev_mode=true`.

## Running tests
From repository root (this directory):
- Run all tests:
  `behave -f pretty -f progress3 -D dev_mode=true`
- Run a specific feature file:
  `behave test/features/user_authentication.feature`
- Run by scenario name (exact match):
  `behave -n "User changes password"`
- Run scenarios that match a name pattern (regex):
  `behave -i "signs in"`
- Run by tag (if you add tags like @playwright):
  `behave --tags=@playwright`

### UI scenarios in the default run

`management_portal_login.feature` signs in through the Management Portal's own login page (internal IDP and
authserver, the default in this release). It runs as part of a plain `behave` run, so CI needs the Playwright
browser installed (`python -m playwright install chromium`).

### Scenarios that need Ory Kratos/Hydra (`@ory`)

The UI login scenarios in `user_authentication.feature` are tagged `@ory`. They test the login flow of the
Ory Kratos/Hydra/self-enrolment setup, which this release does not deploy yet. `behave.ini` excludes them
by default (`default_tags = -@ory`), so a plain `behave` run (as in CI) skips them.
To check an Ory-enabled deployment, run only those scenarios:
  `behave --tags=@ory`

Screenshots are saved automatically by page helpers to:
- test/features/reports/screenshots/

## Playwright Codegen (to prototype selectors)
Generate steps and inspect selectors by recording your interactions:
```
playwright codegen <address>
```

Example:
```
playwright codegen http://localhost/managementportal
```

Codegen opens a browser window and generates Playwright code you can adapt into Page Objects under `test/features/pages/`.
