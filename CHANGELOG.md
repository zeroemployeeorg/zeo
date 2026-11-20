# Changelog

Notable changes are recorded here using Keep a Changelog categories. There
are no release tags in the verified history, so published changes are grouped
by month. Work on review branches is kept under Unreleased.

## Status — 2026-11-20

ZEO is an experimental Linux distribution project for persistent agent
organizations. Its goal is one installation entry point for compatible,
independently versioned platform components.

- **Activity:** Experimental and actively developed. The recorded main branch
  began on 2026-11-17; the latest substantive candidate source changes are
  dated 2026-11-19.
- **Release:** No releases yet; current tags and release metadata were checked.
- **Works today:** The recorded main branch contains the project scaffold.
  Candidate branches contain offline validation, staging and service-packaging
  tools with recorded source-test results; those are not a complete product.
- **Not available yet:** A qualified integrated installer, installed-system
  health check or demonstrated organization with a persistent agent session.
- **Check it yourself:** Candidate source branches provide `make verify` for
  source checks. The recorded initial main branch has no verification target.
  A source-test result does not establish a working installation.

## Unreleased

The following work is present on review branches and is not recorded as a
published release or completed installation.

### Added

- Distribution requirements, component inventory and interface documentation
  to explain what the platform needs and which integrations remain incomplete.
- Offline release-preflight validation for component versions, platform
  compatibility, application profiles and pinned artifacts.
- Local artifact staging and guest service-unit generation, allowing callers
  to prepare pinned session components without claiming to activate them.
- Candidate SSH target enrollment with explicit command registration,
  enrollment receipts and revocation controls.
- Candidate job-check composition that preserves approved checks and makes
  changes or retirements explicit instead of silently dropping checks.
- Candidate application service packaging and memory checks, with source
  tests kept separate from native deployment qualification.
- A single-entry-point installation specification covering task preservation,
  verified recovery and state-safe replacement.

### Fixed

- Candidate job composition now binds each invocation to its check definition,
  uses the intended working directory and publishes output atomically.

### Security

- Hardened candidate input decoding, authorization boundaries and revocation
  checks, with regression tests for refused inputs and interrupted operations.

## 2026-11

### Added

- Created the project scaffold on 2026-11-17 as the starting point for the
  distribution. The initial scaffold did not implement the integrated platform.
