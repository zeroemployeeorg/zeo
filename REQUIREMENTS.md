# Product requirements

**Updated:** 2026-11-18 · **Status:** Implementation in progress; no qualified distribution release.

This register describes public product behavior in neutral terms. The original 312 identifiers remain stable; section 20 adds common environment-profile requirements, section 21 platform operation, section 22 the one-seat installation milestone and section 23 complete application profiles and installed-system diagnostics. IDs are never reused or renumbered. Superseded behavior must retain its ID and a replacement reference.

Each row separates implementation from verification. `Not started`, `In progress`, `Implemented`, `Verified` and `Blocked` have distinct meanings; only repeatable behavioral evidence qualifies a product requirement as Verified. Unknown limitations and dependencies remain unresolved. Distribution maintainers are accountable for the whole product; component maintainers retain their component responsibilities.

Section 20 clarifies the seat-environment policy used by sections 11–12. A seat instance does not, by itself, prove VM isolation: capacity, custody and actual boundaries must be measured and reported.

## 0. Distribution governance

| Work ID | Product requirement | Owner | State | Implementation revision | Verification evidence | Known limitations | Remaining dependencies |
|---|---|---|---|---|---|---|---|
| ZEO-00-01 | Maintain dedicated ownership of the integrated ZEO distribution. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-00-02 | Use this repository as the public home for the distribution. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-00-03 | Deliver released software that can be installed and operated without assembling its infrastructure by hand. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-00-04 | Maintain an authoritative delivery charter and preserve its requirements in the product register. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-00-05 | Keep stable work identifiers and separate implementation progress from verification results. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |

## 1. Research organizations

| Work ID | Product requirement | Owner | State | Implementation revision | Verification evidence | Known limitations | Remaining dependencies |
|---|---|---|---|---|---|---|---|
| ZEO-01-01 | Provide persistent agent-organization infrastructure suitable for ML research groups and small research laboratories. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-01-02 | Enable users to install the platform, create an organization, connect models and perform useful work with budgets, permissions and human oversight. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-01-03 | Support multiple independently configured organizations through the same product. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-01-04 | Support an organization that develops ZEO itself. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-01-05 | Support an independent personal research organization. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-01-06 | Support a separate organization for a configured external workload. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-01-07 | Support organizations created by unrelated adopters. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-01-08 | Allow temporary organizations for controlled experiments. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-01-09 | Require no access to the maintainers' private services or machine setup when using a release. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-01-10 | Support research both with the platform and into its organizational behavior. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-01-11 | Enable reproducible publication implementations, experiments and ablations with inspectable outputs. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-01-12 | Enable comparisons of organizational policy, communication, planning, memory and resource allocation. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-01-13 | Treat useful coordinated execution as a product capability beyond conversational interaction. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |

## 2. Product constraints

| Work ID | Product requirement | Owner | State | Implementation revision | Verification evidence | Known limitations | Remaining dependencies |
|---|---|---|---|---|---|---|---|
| ZEO-02-01 | Retain the established product and platform constraints throughout delivery. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-02-02 | Assign distribution maintenance to a dedicated role. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-02-03 | Keep a clearly identified public distribution repository. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-02-04 | Publish product repositories in the project's public organization after publication review. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-02-05 | Use Linux as the execution baseline. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-02-06 | Provide a single starting point for installation and system configuration. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-02-07 | Support the normal installation without requiring Docker. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-02-08 | Permit an Apple Silicon host adapter to provision Linux guests using the native virtualization framework. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-02-09 | Extend Pomar's existing environment facilities when they meet the product requirements. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-02-10 | Coordinate component versions, releases and lifecycle operations as one product. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-02-11 | Enforce boundaries between organizations and between seats. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-02-12 | Back completion claims with repeatable verification. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-02-13 | Use the same Linux application baseline on physical machines, hosted servers and supported guests. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-02-14 | Justify any language replacement with an observed engineering need. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-02-15 | Retain useful repository boundaries while integrating the product; repository consolidation is optional. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |

## 3. Distribution responsibilities

| Work ID | Product requirement | Owner | State | Implementation revision | Verification evidence | Known limitations | Remaining dependencies |
|---|---|---|---|---|---|---|---|
| ZEO-03-01 | Make distribution maintainers accountable for the complete release and adoption experience. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-03-02 | Assign ownership of each whole-product delivery responsibility. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-03-03 | Maintain the public distribution repository. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-03-04 | Own and publish the supported platform contract. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-03-05 | Own release assembly and compatible component manifests. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-03-06 | Own installation and onboarding across components. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-03-07 | Resolve integration behavior across component boundaries. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-03-08 | Own product-level release acceptance. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-03-09 | Own coordinated upgrade, backup, restore and removal workflows. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-03-10 | Maintain public user and contributor guidance. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-03-11 | Publish release artifacts and notes. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-03-12 | Coordinate changes needed from component maintainers to deliver usable releases. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-03-13 | Retain component maintainers' responsibility for their implementations. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-03-14 | Keep integration defects open until the complete product works, regardless of isolated component test results. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-03-15 | Give each blocking dependency an owner and work item, with an explicit escalation path for ownership or contract conflicts. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-03-16 | Keep the distribution maintenance role focused on product delivery. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |

## 4. Component inventory

| Work ID | Product requirement | Owner | State | Implementation revision | Verification evidence | Known limitations | Remaining dependencies |
|---|---|---|---|---|---|---|---|
| ZEO-04-01 | Inspect source, installed services and deployment configuration before asserting current capabilities. | Distribution maintainers | In progress | Candidate inventory changes; release qualification pending | Source observations summarized in INVENTORY.md; deployment evidence pending | Inventory incomplete; integration behavior unqualified | Remaining component and installation evidence |
| ZEO-04-02 | Maintain an inventory spanning every required product component. | Distribution maintainers | In progress | Candidate inventory changes; release qualification pending | Source observations summarized in INVENTORY.md; deployment evidence pending | Inventory incomplete; integration behavior unqualified | Remaining component and installation evidence |
| ZEO-04-03 | Include Runtime in the component inventory. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-04-04 | Include organizational coordination in the inventory. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-04-05 | Include Messenger in the inventory. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-04-06 | Include Agenda in the inventory. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-04-07 | Include Controlling in the inventory. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-04-08 | Include Connect in the inventory. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-04-09 | Inspect both the user application and its backend. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-04-10 | Inspect Pomar's environment implementation. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-04-11 | Inventory identity, connection and credential integrations. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-04-12 | Inventory recording, evidence storage and artifact handling. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-04-13 | Inspect shared libraries and interface contracts. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-04-14 | Add required components discovered during inspection. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-04-15 | Use a consistent evidence record for every component. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-04-16 | Identify the inspected repository and immutable source revision. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-04-17 | Identify the responsible component maintainer. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-04-18 | Describe responsibility demonstrated by implementation. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-04-19 | Record languages, build tooling and runtime requirements. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-04-20 | Record OS, architecture and host-integration assumptions. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-04-21 | Describe the actual process or library execution model. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-04-22 | Identify durable state, its ownership and migration behavior. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-04-23 | Record APIs, events, contract versions and authentication boundaries. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-04-24 | Record installation and supervision mechanisms. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-04-25 | List product and external-service dependencies. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-04-26 | Track publication obstacles, including private dependencies and unresolved rights. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-04-27 | Track functionality gaps and behavior that still lacks release evidence. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-04-28 | Cite the observations or checks supporting inventory findings. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-04-29 | Identify machine-, account- and organization-specific constants that obstruct independent installation. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-04-30 | Label unresolved facts explicitly and substantiate capabilities from implementation rather than plans or names. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |

## 5. Linux support

| Work ID | Product requirement | Owner | State | Implementation revision | Verification evidence | Known limitations | Remaining dependencies |
|---|---|---|---|---|---|---|---|
| ZEO-05-01 | Publish a supported Linux configuration with an explicit compatibility boundary. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-05-02 | Name the initial Linux distribution and supported version. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-05-03 | Name architectures qualified for the first release. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-05-04 | Measure minimum and recommended resource needs. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-05-05 | Specify required OS services and facilities. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-05-06 | Document necessary network connectivity. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-05-07 | Specify storage capacity and persistence requirements. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-05-08 | Identify installation steps requiring elevated privileges. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-05-09 | Document supported model and agent authentication flows. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-05-10 | Design for ARM64 and AMD64 while claiming support only after architecture-specific release tests. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-05-11 | Keep Linux application execution independent of macOS credential, service and filesystem facilities. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-05-12 | Isolate host-specific virtualization code from Linux product logic. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-05-13 | Preserve needed host-specific workflows through a documented migration or optional adapter. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |

## 6. Component boundaries

| Work ID | Product requirement | Owner | State | Implementation revision | Verification evidence | Known limitations | Remaining dependencies |
|---|---|---|---|---|---|---|---|
| ZEO-06-01 | Inspect existing contracts before assigning final implementation boundaries. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-06-02 | Maintain distinct responsibilities for distribution, execution, environments, OS supervision and organizational coordination. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-06-03 | Have the distribution assemble compatible software and manage its release lifecycle. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-06-04 | Have Runtime own agent execution and organizational runtime behavior. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-06-05 | Have Pomar supply execution environments through supported backends. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-06-06 | Use Linux service management for OS service supervision. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-06-07 | Keep organizational work coordination in the coordination component. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-06-08 | Reuse schedulers, registries, credential facilities and environment managers unless a documented gap requires an alternative. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-06-09 | Choose process boundaries from lifecycle, isolation, resource and failure needs rather than repository count. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-06-10 | Version inter-component contracts explicitly. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |

## 7. Release assembly

| Work ID | Product requirement | Owner | State | Implementation revision | Verification evidence | Known limitations | Remaining dependencies |
|---|---|---|---|---|---|---|---|
| ZEO-07-01 | Pin one compatible component set for each distribution release. | Distribution maintainers | In progress | Preflight source a68ef3ad in PR #1; compatible release set unqualified | Linux make verify passed at ddd4666; whole-release acceptance pending | Metadata checks do not qualify a release or installation | Published compatible artifacts, trust verification and integration acceptance |
| ZEO-07-02 | Make the release manifest sufficient to identify and reproduce the assembly. | Distribution maintainers | In progress | Preflight source a68ef3ad in PR #1; compatible release set unqualified | Linux make verify passed at ddd4666; whole-release acceptance pending | Metadata checks do not qualify a release or installation | Published compatible artifacts, trust verification and integration acceptance |
| ZEO-07-03 | Record the distribution's release version. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-07-04 | Pin component versions and source revisions. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-07-05 | Identify the target OS and architecture. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-07-06 | Record artifact integrity digests. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-07-07 | Declare runtime dependencies. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-07-08 | Record interface contract versions. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-07-09 | Describe schema versions and migration compatibility. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-07-10 | Retain build provenance. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-07-11 | List supported predecessor releases for upgrade. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-07-12 | Pin relevant agent tools and their compatibility constraints. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-07-13 | Install prebuilt release artifacts rather than requiring users to compile the component repositories. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-07-14 | Keep distribution implementation, assembly instructions, integration tests, examples and documentation in this repository. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-07-15 | Consume component releases without maintaining manually copied source forks inside the distribution. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |

## 8. Integrated installation

| Work ID | Product requirement | Owner | State | Implementation revision | Verification evidence | Known limitations | Remaining dependencies |
|---|---|---|---|---|---|---|---|
| ZEO-08-01 | Offer one documented installation entry point on a fresh supported Linux host. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-08-02 | Make that entry point carry out the full installation sequence. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-08-03 | Check host compatibility and available resources. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-08-04 | Explain elevated changes before applying them. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-08-05 | Fetch the selected release and verify its integrity. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-08-06 | Install prerequisites through controlled and documented procedures. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-08-07 | Create service identities and protect persistent state directories. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-08-08 | Configure communication between required components. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-08-09 | Install and activate the required OS service definitions. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-08-10 | Run health and integration verification. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-08-11 | Provide onboarding access, including a usable remote-access route. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-08-12 | Declare installation successful only when the required product is usable. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-08-13 | Recover safely from interrupted or repeated installation without duplicating organizations or damaging state. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-08-14 | Configure the required application and organizational services as an integrated installation. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-08-15 | Allow necessary credentials and user choices without requiring manual infrastructure assembly. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-08-16 | Support unattended-server environments without a desktop session. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |

## 9. Onboarding

| Work ID | Product requirement | Owner | State | Implementation revision | Verification evidence | Known limitations | Remaining dependencies |
|---|---|---|---|---|---|---|---|
| ZEO-09-01 | Guide users from first installation to a useful initial project. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-09-02 | Guide creation of an independent organization. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-09-03 | Configure a supported model or coding-agent connection. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-09-04 | Establish spending and execution-concurrency limits. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-09-05 | Select or create the working directory. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-09-06 | Start a bounded first project. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-09-07 | Show progress and completed results. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-09-08 | Allow users to answer requests for input. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-09-09 | Avoid reliance on maintainers' accounts, subscriptions, credentials or private services. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-09-10 | Qualify API-key and subscription authentication separately and document only tested modes. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-09-11 | Complete a small coordinated execution workflow with a verifiable deliverable, beyond service-health checks. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |

## 10. Organization independence

| Work ID | Product requirement | Owner | State | Implementation revision | Verification evidence | Known limitations | Remaining dependencies |
|---|---|---|---|---|---|---|---|
| ZEO-10-01 | Carry explicit organization identity throughout product interfaces and state. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-10-02 | Partition each organization's authority and recovery scope. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-10-03 | Separate organizational coordinator identity and continuity. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-10-04 | Separate agent identities by organization. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-10-05 | Keep organization configuration independent. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-10-06 | Keep organization workspaces independent. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-10-07 | Keep message storage independent. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-10-08 | Separate project and commitment records. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-10-09 | Scope credentials and grants to the owning organization. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-10-10 | Separate budgets and usage accounting. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-10-11 | Separate memory and evidence. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-10-12 | Separate durable application state. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-10-13 | Define organization-specific backup and restore boundaries. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-10-14 | Permit shared executable files while keeping authority and organizational data separate. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-10-15 | Operate, inspect, back up and remove each organization independently. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-10-16 | Document and enforce whether multiple distribution versions can coexist on one host. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-10-17 | Demonstrate concurrent organizations and reject unauthorized cross-organization access. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |

## 11. Seat isolation

| Work ID | Product requirement | Owner | State | Implementation revision | Verification evidence | Known limitations | Remaining dependencies |
|---|---|---|---|---|---|---|---|
| ZEO-11-01 | Represent seat identity and isolation requirements explicitly rather than inferring them from seat counts. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-11-02 | Support enforced identities, workspace permissions, process controls and resource limits within Linux; apply the environment-profile policy in section 20. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-11-03 | Run agents without unrestricted host-administration rights by default. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-11-04 | Specify the complete seat isolation and lifecycle contract. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-11-05 | Define allocation of seat identities. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-11-06 | Define workspace ownership. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-11-07 | Define access to credentials. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-11-08 | Define access to shared resources. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-11-09 | Define process and session ownership. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-11-10 | Enforce CPU, memory and process-count limits. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-11-11 | Define seat provisioning, suspension, recovery and deletion. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-11-12 | Provide a controlled route for privileged operations. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-11-13 | Explain the security limits of shared-kernel execution. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-11-14 | Select stronger VM isolation where required; measure allocations and follow section 20 without silently weakening requested boundaries. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |

## 12. Environment backends

| Work ID | Product requirement | Owner | State | Implementation revision | Verification evidence | Known limitations | Remaining dependencies |
|---|---|---|---|---|---|---|---|
| ZEO-12-01 | Inspect existing Pomar capabilities before choosing backend changes. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-12-02 | Qualify the required environment delivery routes. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-12-03 | Support an existing Linux server as an execution target. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-12-04 | Support an ARM64 Linux guest on an Apple Silicon host. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-12-05 | Provide disposable environments for experiments. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-12-06 | Verify host-adapter behavior as part of Linux guest qualification. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-12-07 | Qualify the ARM64 guest userspace and artifacts. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-12-08 | Verify guest creation and provisioning. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-12-09 | Verify CPU and memory allocation. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-12-10 | Verify connectivity and access. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-12-11 | Verify durable virtual disks. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-12-12 | Verify guest start, stop and restart. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-12-13 | Verify behavior after the physical host reboots. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-12-14 | Provide useful guest logs and failure reporting. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-12-15 | Verify backup and recovery of the guest's required state. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-12-16 | Measure resource contention under representative load. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-12-17 | Require demonstrated capability before depending on nested virtualization or GPU access. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-12-18 | Evaluate management and worker topology through resource measurements. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-12-19 | Pilot coordinator and agent workloads in host-provided Linux while retaining recovery onto a separate supported server. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-12-20 | Expand a live migration only after a representative pilot passes. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |

## 13. Lifecycle and recovery

| Work ID | Product requirement | Owner | State | Implementation revision | Verification evidence | Known limitations | Remaining dependencies |
|---|---|---|---|---|---|---|---|
| ZEO-13-01 | Provide coordinated product lifecycle operations. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-13-02 | Report component and organization health with actionable failures. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-13-03 | Start and stop services in dependency order while preserving state. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-13-04 | Recover required services after reboot and accurately identify interrupted work. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-13-05 | Check compatibility and preserve data through coordinated upgrades and migrations. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-13-06 | Keep failed upgrades diagnosable and recoverable. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-13-07 | Capture a consistent and documented backup set. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-13-08 | Restore supported state onto a fresh supported host. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-13-09 | Collect useful diagnostics without exposing secrets. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-13-10 | Remove installed services cleanly and state explicitly what data remains. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-13-11 | Define backup-based recovery for irreversible migrations instead of promising unsupported rollback. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-13-12 | Establish consistency across components before backing up live databases and files. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-13-13 | Qualify cross-architecture recovery explicitly, including limits on process and agent-session portability. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |

## 14. Research workflows

| Work ID | Product requirement | Owner | State | Implementation revision | Verification evidence | Known limitations | Remaining dependencies |
|---|---|---|---|---|---|---|---|
| ZEO-14-01 | Deliver a bounded end-to-end research workflow and its implementation plan. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-14-02 | Represent enough experimental context to interpret and reproduce results. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-14-03 | Record the research question and hypothesis. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-14-04 | Record protocol and baseline choices. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-14-05 | Identify workloads and datasets. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-14-06 | Pin source and execution environments. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-14-07 | Record model and agent settings. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-14-08 | Record applicable random seeds. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-14-09 | Record budgets and termination conditions. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-14-10 | Track run identities and retries. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-14-11 | Record evaluation methods and metrics. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-14-12 | Preserve artifacts and provenance. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-14-13 | Record failures and human interventions. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-14-14 | Link conclusions to supporting evidence. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-14-15 | Provide a Python-oriented research integration while reusing established experiment and training tools. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-14-16 | Version experimental dimensions independently. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-14-17 | Identify the distribution release separately. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-14-18 | Identify organization configuration separately. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-14-19 | Identify experimental treatment separately. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-14-20 | Identify workload and evaluator revisions separately. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-14-21 | Compare policies on the same software release without introducing unrelated software changes. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-14-22 | Provide a small ML experiment or reproduction and an organizational-policy comparison, with honest replay and hosted-model limitations. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-14-23 | Prevent experiments from accidentally acquiring production credentials or external write authority. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |

## 15. Public source distribution

| Work ID | Product requirement | Owner | State | Implementation revision | Verification evidence | Known limitations | Remaining dependencies |
|---|---|---|---|---|---|---|---|
| ZEO-15-01 | Map every product repository to its intended public destination. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-15-02 | Include auxiliary build, runtime, test and explanatory dependencies in publication scope. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-15-03 | Execute authorized publication and transfer steps, recording concrete permission or ownership blockers. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-15-04 | Complete publication checks before exposing additional source or history. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-15-05 | Review tracked files and history for secrets and private operational information. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-15-06 | Resolve publication blockers and rotate credentials if exposure requires it. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-15-07 | Eliminate dependencies that require inaccessible private repositories. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-15-08 | Update package identifiers, module references, documentation and release automation for public distribution. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-15-09 | Record rights and licence status for included source and assets. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-15-10 | Resolve missing open-source licence decisions. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-15-11 | Prove that public users can build and operate the system without private repository access. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-15-12 | Exclude organizational databases, confidential records, messages, credentials and live configuration from publication. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-15-13 | Coordinate any necessary shared-history remediation and its effect on existing clones. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-15-14 | Distinguish released, experimental and planned features in public product material. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |

## 16. External contribution

| Work ID | Product requirement | Owner | State | Implementation revision | Verification evidence | Known limitations | Remaining dependencies |
|---|---|---|---|---|---|---|---|
| ZEO-16-01 | Make external adoption and contribution practical. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-16-02 | Explain the product clearly. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-16-03 | Provide a verified quickstart. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-16-04 | Explain the product architecture. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-16-05 | Publish qualified platforms and measured resource needs. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-16-06 | Provide contributor setup independent of the maintainers' organization. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-16-07 | Document component ownership and issue routing. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-16-08 | Document contract changes and compatibility handling. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-16-09 | Publish release and support expectations. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-16-10 | Provide a security-reporting channel. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-16-11 | State known product limitations. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-16-12 | Publish examples and reproducible acceptance workflows. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-16-13 | Allow a contributor to create a development organization without private state or credentials. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |

## 17. Release acceptance

| Work ID | Product requirement | Owner | State | Implementation revision | Verification evidence | Known limitations | Remaining dependencies |
|---|---|---|---|---|---|---|---|
| ZEO-17-01 | Require product-level acceptance evidence before accepting a release candidate. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-17-02 | Demonstrate complete installation on a fresh qualified Linux host. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-17-03 | Have a configured organization finish a project with independently checkable results. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-17-04 | Demonstrate installation and useful operation without private services or repositories. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-17-05 | Demonstrate correct service and interrupted-work recovery after reboot. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-17-06 | Inject component failures and verify visible failures, recovery and duplicate-work prevention. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-17-07 | Run two organizations together and reject cross-organization state or credential access. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-17-08 | Demonstrate that worker identities cannot read another seat's protected resources. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-17-09 | Upgrade a supported predecessor release while preserving its state. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-17-10 | Restore a consistent backup onto a different clean host. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-17-11 | Qualify the release inside the supported Apple-hosted Linux guest. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-17-12 | Measure registered identities, idle sessions and active execution separately during load tests. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-17-13 | Export an experiment result package with traceable provenance. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-17-14 | Verify removal against the documented data-retention behavior. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-17-15 | Reproduce the documented build and checks from an external-style public checkout. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-17-16 | Use real model execution in a bounded acceptance workflow; keep regression fixtures distinct from live-operation evidence. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-17-17 | Publish actual workload measurements rather than inferring active capacity from identity counts. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |

## 18. Delivery milestones

| Work ID | Product requirement | Owner | State | Implementation revision | Verification evidence | Known limitations | Remaining dependencies |
|---|---|---|---|---|---|---|---|
| ZEO-18-01 | Deliver the product through observable milestones. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-18-02 | Complete the inventory and ownership milestone. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-18-03 | Produce source-backed component/dependency findings, platform gaps, publication blockers and assigned work. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-18-04 | Complete the first integrated Linux-installation milestone. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-18-05 | Install a release candidate on a clean Linux host and finish a bounded workflow. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-18-06 | Complete the lifecycle and isolation milestone. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-18-07 | Demonstrate independent organizations, seat boundaries, reboot recovery, upgrade, backup, restore and uninstall. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-18-08 | Complete the supported guest-platform pilot milestone. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-18-09 | Run the same release through Pomar's host adapter and migrate representative workloads before expanding. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-18-10 | Complete the research-workflow milestone. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-18-11 | Deliver experiment and publication-reproduction workflows with traceable outputs. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-18-12 | Complete the public-release milestone. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-18-13 | Finish public source publication, installation tests, documentation and release artifacts. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-18-14 | Investigate publication blockers alongside implementation instead of postponing discovery until release. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |

## 19. Progress reporting

| Work ID | Product requirement | Owner | State | Implementation revision | Verification evidence | Known limitations | Remaining dependencies |
|---|---|---|---|---|---|---|---|
| ZEO-19-01 | Track delivery facts for every requirement. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-19-02 | Name the accountable owner. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-19-03 | Assign a permanent work identifier. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-19-04 | Record the current delivery state. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-19-05 | Record the implementing source revision. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-19-06 | Link verification evidence. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-19-07 | State known limitations. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-19-08 | Identify unresolved dependencies. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-19-09 | Use defined progress states consistently. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-19-10 | Use Not started, In progress, Implemented, Verified and Blocked as distinct states. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-19-11 | Distinguish code implementation, verification, merge, deployment and successful organization operation. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-19-12 | Publish a release report with enough detail to audit its claims. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-19-13 | Identify exact product and component revisions. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-19-14 | Identify installable artifacts and public download locations. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-19-15 | List supported configurations. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-19-16 | Report acceptance outcomes. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-19-17 | Report measured resource usage. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-19-18 | Provide upgrade and recovery evidence. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-19-19 | Report publication and repository-migration progress. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-19-20 | Disclose outstanding limitations. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-19-21 | Name every unmet requirement and explain why it remains open. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-19-22 | Start delivery reporting with inspected components, assigned work and blockers, and continue into implementation. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-19-23 | Enable an external user to install, configure, operate and recover an independent research organization without personal assistance. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |

## 20. Common Linux environment profiles

| Work ID | Product requirement | Owner | State | Implementation revision | Verification evidence | Known limitations | Remaining dependencies |
|---|---|---|---|---|---|---|---|
| ZEO-20-01 | Install the ZEO distribution into a versioned Linux environment provided by Pomar. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-20-02 | Keep environment profiles, product releases, organization settings and experiment treatments independently versioned. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-20-03 | Separate Linux-target creation or attachment from reusable Linux provisioning. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-20-04 | Target a Pomar instance per seat on Apple Silicon, with separate host-user custody and measured guest capacity. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-20-05 | Support ARM64 and AMD64 hosted Linux targets without requiring nested virtualization. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-20-06 | Keep Docker optional for the standard environment path. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-20-07 | Describe Linux baseline, layout, tools, dependencies, services and resource limits in the environment profile. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-20-08 | Store credential references in reusable profiles while delivering secret material through controlled mechanisms. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-20-09 | Report the actual isolation boundary and reject unsupported placements without silent weakening. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-20-10 | Retain Runtime's placement authority and Pomar's environment authority without duplicating schedulers or seat registries. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-20-11 | Expose repeatable inspect, attach/create, provision, drift, start/stop, health, logs, backup/restore, reset and removal operations. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-20-12 | Run the same environment conformance suite on Apple-hosted ARM64 guests and hosted ARM64 and AMD64 Linux targets. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-20-13 | Run the same bounded released-product workflow on each qualified backend. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-20-14 | Separate reconstructable environments, durable organizational state, work products, disposable data and secrets. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-20-15 | Recover across architectures using appropriate rebuilt artifacts and compatible state rather than assuming VM-image portability. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-20-16 | Qualify agent-session portability and provide a work-preserving handover when resumption is unsupported. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-20-17 | Fence live migration so only one copy of an organization has execution authority. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-20-18 | Qualify a small isolated pilot before scaling the seat population. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-20-19 | Select architecture-specific artifacts and verify their versions and digests automatically. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-20-20 | Check requested workload capabilities against the target and return useful placement refusals. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-20-21 | Record architecture and hardware provenance without promising identical performance or ML results across hosts. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-20-22 | Inspect and reuse existing released environment features before introducing replacement implementations. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-20-23 | Consume qualified published environment releases and require adoption evidence before moving gates or production deployment. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-20-24 | Measure host capacity before choosing per-seat guest allocations or any grouped topology; disclose the resulting isolation. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |


## 21. Actor-independent unified platform

| Work ID | Product requirement | Owner | State | Implementation revision | Verification evidence | Known limitations | Remaining dependencies |
|---|---|---|---|---|---|---|---|
| ZEO-21-01 | Assemble the platform machinery into one versioned Linux release and one installation flow. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-21-02 | Run platform services under systemd with explicit health checks, restart policy and state custody. | Distribution maintainers | In progress | Session packaging4218cb5 and artifact stagingddd4666 in PR #1 | Linux make verify passed at ddd4666; real systemd/service activation pending | Renderer only; no platform installation or service qualification | Owner artifacts, allocated identities and adopted environment profile |
| ZEO-21-03 | Keep platform service availability independent of builder sessions, panes and consumer actors. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-21-04 | Run consumer seats on pinned released artifacts outside builder working trees. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-21-05 | Support platform builders and release consumers on separate physical hosts. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-21-06 | Retain Runtime authority for desired organizations, seats and placement. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-21-07 | Use Pomar for observed environments, heartbeats and reconciliation; report unexpected state without silently changing placement. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-21-08 | Expose the seat directory and authenticated real terminal attachment through App. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-21-09 | Execute seats inside Pomar environments with explicit supported isolation and identity mappings. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-21-10 | Keep observation, App access and owner notification outside the pilot VM whose failure is tested. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-21-11 | Reuse owner-adopted bootstrap, configuration, credential and lifecycle interfaces in distribution installation. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-21-12 | Record each service artifact, runtime location, supervisor, state, dependencies, actor dependence, Linux gap and responsible owner; assign unknowns to owners. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |

## 22. Single-host installation with one working seat

| Work ID | Product requirement | Owner | State | Implementation revision | Verification evidence | Known limitations | Remaining dependencies |
|---|---|---|---|---|---|---|---|
| ZEO-22-01 | Install an exact pinned component set on a fresh supported Linux host through one entry point. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-22-02 | Demonstrate that installed systemd platform services continue when the builder session ends. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-22-03 | Create one organization and seat through the intended Runtime and Pomar ownership boundary. | Distribution maintainers | Blocked | Component integration contract/adoption pending | No integrated demonstration held | No adopted Linux creation/placement or browser-terminal boundary held by distribution | Runtime, environment, application and delegation owners; interface/work dates requested |
| ZEO-22-04 | List the seat in App and use a real terminal whose resize reaches the running guest session. | Distribution maintainers | Blocked | Component integration contract/adoption pending | No integrated demonstration held | No adopted Linux creation/placement or browser-terminal boundary held by distribution | Runtime, environment, application and delegation owners; interface/work dates requested |
| ZEO-22-05 | Detach and reconnect to the same immutable session, preserving it across App and host-agent restarts. | Distribution maintainers | Blocked | Component integration contract/adoption pending | No integrated demonstration held | No adopted Linux creation/placement or browser-terminal boundary held by distribution | Runtime, environment, application and delegation owners; interface/work dates requested |
| ZEO-22-06 | Complete a bounded task with a real model and independently inspect its deliverable and budget evidence. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-22-07 | Stop the pilot seat VM and demonstrate App failure visibility and actual owner alert delivery within sixty seconds. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |
| ZEO-22-08 | Bind every pilot demonstration to the installed release, component and environment versions and original receipts; keep broader acceptance obligations open. | Distribution maintainers | Not started | — | — | UNKNOWN | UNKNOWN |

## 23. Whole-stack application profiles and installed-system doctor

| Work ID | Product requirement | Owner | State | Implementation revision | Verification evidence | Known limitations | Remaining dependencies |
|---|---|---|---|---|---|---|---|
| ZEO-23-01 | Install the complete declared platform, seat-tooling and application profile through one supported entry point. | Distribution maintainers | Not started | — | — | Installer and installed-system doctor are not implemented | Adopted owner release/profile/lifecycle interfaces |
| ZEO-23-02 | Enumerate required application payloads and dependencies explicitly; a missing required app must prevent successful installation. | Distribution maintainers | Not started | — | — | Installer and installed-system doctor are not implemented | Adopted owner release/profile/lifecycle interfaces |
| ZEO-23-03 | Consume runtime-owned release identities, signatures and trust bootstrap without creating a separate distribution signing authority. | Distribution maintainers | Not started | — | — | Installer and installed-system doctor are not implemented | Adopted owner release/profile/lifecycle interfaces |
| ZEO-23-04 | Derive staged installation and bootstrap order from adopted component contracts, including mutually dependent core services. | Distribution maintainers | Not started | — | — | Installer and installed-system doctor are not implemented | Adopted owner release/profile/lifecycle interfaces |
| ZEO-23-05 | Finish installation by running make doctor from the installed release bundle without requiring component source checkouts or a compiler. | Distribution maintainers | Not started | — | — | Installer and installed-system doctor are not implemented | Adopted owner release/profile/lifecycle interfaces |
| ZEO-23-06 | Report doctor GREEN only when every required installed-profile probe passes; fail closed on missing, unknown or unimplemented capabilities. | Distribution maintainers | Not started | — | — | Installer and installed-system doctor are not implemented | Adopted owner release/profile/lifecycle interfaces |
| ZEO-23-07 | Preserve application drafts, history, state and credential references through idempotent installation, interruption and supported lifecycle operations. | Distribution maintainers | Not started | — | — | Installer and installed-system doctor are not implemented | Adopted owner release/profile/lifecycle interfaces |
| ZEO-23-08 | Keep external-service credentials with the owning integration service and grant applications only their approved scoped capabilities. | Distribution maintainers | Not started | — | — | Installer and installed-system doctor are not implemented | Adopted owner release/profile/lifecycle interfaces |
| ZEO-23-09 | Apply each profile’s declared ingress and publication policy using the owning service interfaces; do not expose private applications by default. | Distribution maintainers | Not started | — | — | Installer and installed-system doctor are not implemented | Adopted owner release/profile/lifecycle interfaces |
| ZEO-23-10 | Distinguish independent-adopter and authorized operator profiles explicitly without silently calling a reduced catalog the complete requested stack. | Distribution maintainers | Not started | — | — | Installer and installed-system doctor are not implemented | Adopted owner release/profile/lifecycle interfaces |
| ZEO-23-11 | Reuse versioned shared artifact and interface schemas across component process boundaries rather than maintaining divergent copies. | Distribution maintainers | Not started | — | — | Installer and installed-system doctor are not implemented | Adopted owner release/profile/lifecycle interfaces |
| ZEO-23-12 | Bind installed-system diagnostic results to release, profile, target, instance and actual probes, separately from source tests and full release acceptance. | Distribution maintainers | Not started | — | — | Installer and installed-system doctor are not implemented | Adopted owner release/profile/lifecycle interfaces |
