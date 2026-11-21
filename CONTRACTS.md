# Entrypoints and integration contracts

**Updated:** 2026-11-20 · **Status:** Maintainer source interface index; installed distribution and end-user entrypoint not qualified.

This repository assembles the product. Its current executables are assembly helpers. The Runtime-owned action dispatcher and maintained application/property recipes are separate capabilities; this index does not replace their catalog, create a new runner or establish their installed coverage. A component name or version pin does not authorize its use.

For product requirements see [REQUIREMENTS.md](REQUIREMENTS.md); source and component evidence is summarized in [INVENTORY.md](INVENTORY.md). Normal installation is intended to consume a complete compatible release of prebuilt artifacts. The developer commands below are not the one-click installation experience.

## Current source interfaces

The implementation files are authoritative. This document points to them rather than maintaining a second executable schema. The repair source at `42ac5d163fdb2723a4dd402e93917541a770a48b` passed Go 1.27.0 vet and race tests across seven packages. The original full gate-log SHA256 is `f326d2802dd7c3b84496365923aa0d3a5ceb09a91f984d3e9260d047d55a3946`. This includes the input, authorization, revocation and output-publication repairs; it supersedes the older source-only evidence, not native acceptance. Documentation-only successors preserve those executable/test bytes but need an explicit current-head review and gate-carry disposition. Documentation and source presence do not qualify installation.

| Entrypoint | Authoritative implementation | Input and output | Scope of successful exit |
|---|---|---|---|
| `zeo-release` | [CLI](cmd/zeo-release/main.go), [manifest types/validation](internal/release/manifest.go) | `zeo.release/v1` plus explicit Linux target/profile; optional local artifact directory; JSON `zeo.release-plan/v1` to stdout or atomic file | Metadata/local-byte checks completed; output always retains `ready_for_install: false`. |
| `zeo-service` | [CLI](cmd/zeo-service/main.go), [unit renderer](internal/install/session.go) | Organization/seat, allocated identity names, Pomar digest and explicit resource caps; systemd unit text on stdout | Candidate unit rendered; no account allocation, artifact verification, installation or activation. |
| `zeo-stage` | [CLI](cmd/zeo-stage/main.go), [artifact staging](internal/install/artifact.go) | Absolute local Pomar payload and existing trusted artifact root, digest and size; relative staged path on stdout | Exact copied bytes published or matching existing bytes verified; no release authentication, executable qualification or service activation. |
| `zeo-target-bootstrap` | [CLI](cmd/zeo-target-bootstrap/main_linux.go), [enrollment contract](docs/ssh-target-bootstrap.md) | Reviewed public packet and byte pins; plan/apply/revoke with exact retained receipts | Narrow enrollment boundary only; repaired source tests passed at the exact commit above, while native authentication, installed revocation and containment remain unqualified. |
| `zeo-job-checks` | [CLI](cmd/zeo-job-checks/main.go), [usage](docs/job-checks.md) | Reviewed baseline digest/current catalog/pinned input root and mandatory fresh absolute `--out` | Complete check script published without replacement; no check execution, job approval or installation. |
| `make verify` | [Makefile](Makefile) | Maintainer checkout and admitted Linux Go/race toolchain | Formatting prerequisite, vet and race tests; not installed health or product acceptance. |
| One-entry installation and `make doctor` | Not implemented | Planned complete release/profile and actual installed instance | No success, green diagnostic or installation behavior is available from the present source. |

CLI parse/usage errors generally use exit 2; validation or I/O failures use exit 1; successful helper work uses exit 0. The job-check CLI supports successful `--help`; bootstrap currently treats its flag help request as usage exit2. Missing or invalid semantic inputs can be validation errors rather than usage errors. An exit code must be interpreted within the narrow helper scope above.

The exact source usages are:

```text
zeo-release --manifest FILE --arch arm64|amd64 --linux-id ID --linux-version VERSION --profile ID --profile-version VERSION [--artifacts DIR] [--out FILE]
zeo-service --organization ID --seat ID --user USER --group GROUP --pomar-sha256 HEX --memory-bytes N --tasks N --cpu-percent N
zeo-stage --pomar FILE --sha256 HEX --size BYTES --artifact-root DIR
zeo-target-bootstrap --packet FILE --packet-sha256 HEX --binary-sha256 HEX [--apply|--revoke]
zeo-job-checks --baseline FILE --baseline-sha256 HEX --catalog FILE --inputs DIR --out ABSOLUTE_NEW_FILE
```

These examples describe the interfaces after tools are built in the maintainer development environment. They do not imply a published tool bundle, qualified component payload or supported end-user installation.

## Release metadata boundary

The manifest pins release version, platform, independently versioned environment, component source/artifact identities, contracts, state schemas, runtime dependencies, upgrade paths, agent tools and provenance. Exactly one component owns each required capability; one component may own several. Required capabilities are runtime, coordination, messaging, agenda, controlling, identity and application. They are responsibilities, not seven required daemons.

`requires` describes integration dependencies and may be cyclic. `install_after` is separate and must form an acyclic ordering. This does not invent a bootstrap protocol for a mutually dependent service group. Owners must supply the actual compatible bootstrap/configuration interfaces before installation can use that order.

Field names and their Go JSON tags in the linked implementation define the current wire shape. Decoding rejects duplicate, unknown and noncanonical keys. Version identifiers are opaque: compatibility is explicitly listed rather than inferred from ordering version strings. The typed-metadata digest is not a raw-document signature. HTTPS syntax and a supplied SHA256 do not establish the signer's identity, public availability, executable target or runtime compatibility.

Runtime owns release/signature authority and trust bootstrap. Distribution consumes its qualified contract. The current helpers must not be used as an independent signer or as a way to turn hash-only artifact selection into authenticated release trust.

## Component boundary map

The table states what installation needs, not that those interfaces have been adopted or that the components do not exist elsewhere. A source contract, passing component tests, an installed tuple and successful integration are distinct evidence levels.

| Boundary / responsible capability owner | Required receiving contract | Current distribution qualification |
|---|---|---|
| Runtime release authority | Signed release/profile/artifact closure, trusted verifier/root reference, platform/contract binding and supported switch/recovery | Qualified Linux trust and complete service lifecycle receiving remain open. |
| Action dispatcher and property/application catalog | Maintained action/recipe versions, required inputs, producer contracts, policy/refusals and released invocation | Current owned dispatcher/catalog snapshot and installed coverage requested; no global completeness claim. |
| Runtime organization/placement | Authoritative schema/revision, admitted creation caller, stable identity and desired environment/session binding | Proposal and adoption pending; a binding command or chart observation is not a creation grant. |
| Pomar environments/session transport | Observed directory/heartbeat, exact environment/session generation, scoped attach authority, continuity and lifecycle | Local session CLI source is known; qualified remote browser/agent and installation receiving remain open. |
| App and Connect | Actual authenticated frontend/backend, operator delegation, accepted transport/grants and installed readiness | Component web evidence does not qualify the real seat browser-terminal path. |
| Messaging | Fresh-instance bootstrap, service principal, delivery/acknowledgement semantics and owner-alert route | A queued message is not delivered-user proof; platform bootstrap and failure-alert acceptance remain open. |
| Agenda and scheduling | Content/authority registration, real timer/caller/sender, state compatibility and owned supervision | Preserve the existing scheduler and jobs. Structural render or an effect-guarded command is not adopted firing. |
| Coordination, Controlling and Recorder | Actual payload, identity, protected state, authenticated health, Linux lifecycle and recovery | Component/runtime owner installation and Linux gaps remain in the inventory. |
| Application profile | Complete required payload list, source/tool contracts, runtime/assets, config/state, ingress/delegation and recovery | Missing required applications must keep the complete profile incomplete. No source checkout or developer credential inheritance. |
| Installed doctor | Actual required probes bound to release/profile/target/instance, truthful unknown/refusal and actionable owner | Not implemented. A manifest plan, imported PASS document, public HTTP 200 or successful source test cannot mint GREEN. |

A complete profile must include all of its required components and app workloads. Optional roles must be declared explicitly; missing required roles are not silently skipped. Actor-held shells, duplicate schedules, shared organizational authority and an unverified artifact path do not close a contract gap.

## Packaging and acceptance

The guest session layout and keeper-restart limitations are described in [session-service.md](docs/session-service.md). Runtime must supply qualified actual identities; the renderer only validates names and requested caps. The artifact root and executable custody must be verified before activation. Never install a unit merely because rendering succeeded.

The integrated milestone still requires a clean pinned installation, independently supervised services, organization/seat creation, real App terminal and resize, same-session reconnect across App/host-agent restarts, one bounded model task, and stopped-VM failure visibility with delivered owner alert within sixty seconds. The keeper's own restart loses live PTYs. A generic missed-job watchdog is not automatically the VM-failure detector.

Contract changes require the supplying and consuming owners to bind exact source/schema/artifact compatibility and run their relevant failure and integration checks. Distribution records the resulting evidence before consuming the release. This index grants no network, credential, signing, migration, publication or deployment authority.
