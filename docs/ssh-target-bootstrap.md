# SSH target enrollment candidate

`cmd/zeo-target-bootstrap` implements the initial Linux target enrollment boundary
for an existing Connect broker. **It is candidate Go code with recorded Linux
source-test coverage; target qualification remains pending.** It has not run on
a target and is not a whole-platform installer or a working broker. Compilation
and tests belong on an admitted Linux gate.

The initial administrator invokes an independently qualified, pinned executable
once with a reviewed PUBLIC enrollment packet. Packet and binary SHA256 checks
prove byte integrity; they do not authenticate a release, authorize a principal
or replace Runtime release trust. Connect retains private-key custody, grants,
approval, transport and execution logging. No private key belongs in the packet,
source tree, arguments or output.

The closed `zeo.ssh-target-bootstrap/v1` packet has these required fields:

| Field | Meaning |
|---|---|
| `schema` | Exactly `zeo.ssh-target-bootstrap/v1` |
| `target_hostname` | Exact short target name |
| `broker_source_ip` | One actual source IP as seen by the target, without a network range or IPv6 zone identifier |
| `broker_source_network` | `tailnet` or `public`; no inferred NAT/interface identity |
| `login_user` | Exactly `zeoconnect-broker-run`, also returned in the enrollment result |
| `public_key` | Public Ed25519 key, without key options or comment |
| `command_id` | One registered no-argument command ID |
| `privilege` | `read`, `write` or `root`, matching the broker registry |
| `timeout_seconds` | 1–300 seconds, aligned with the broker registration |
| `output_max_bytes` | 1–65,536 bytes across stdout and stderr |
| `script_base64` | Reviewed shell adapter, at most 64 KiB decoded |
| `script_sha256` | Digest of the adapter bytes |

Without `--apply` the tool validates the packet and reports `PLAN_ONLY`. Applying
also requires the actual executable digest and a conflict-free exact Linux
root target. Existing identities, groups, paths, unsupported custody, wrong host
and inconsistent prior receipts refuse. Durable intent precedes configuration
changes. Partial outcomes retain evidence and refuse blind replay. Completed
identical enrollment is read back without creating another account.

The dedicated account has a locked password and a root-owned home containing
only public authorization data. Its shell is `/bin/sh` because SSH invokes forced
commands through that shell; `nologin` would prevent dispatch. The root-owned
home and key cannot be modified by the login user. The key restricts source
address, command, terminal allocation and forwarding. Existing admin access and
global SSH configuration are unchanged. Actual PAM/SSH compatibility still
requires target acceptance.

The forced command is the installed Go executable's `dispatch` mode. It accepts
only the packet's exact ID in `SSH_ORIGINAL_COMMAND`, with no arguments or shell
syntax. Read/write classes execute as the dedicated unprivileged user. Only a
root-class packet installs the exact `zeo-target-bootstrap run` sudo rule; `run`
refuses other classes or identities. The packet's login user and privilege must
match Connect's registered host and command. An old root-only metadata reader
cannot be relabeled as an unprivileged read.

The runner verifies root-owned registered bytes and executes an anonymous sealed
Linux memory snapshot, with closed input and fixed environment. It never reopens
an already-checked source path to execute different bytes. An exclusive shared
lock prevents concurrent sessions. Stdout/stderr share the packet's output cap;
exceeding it fails and terminates the owned process group. Deadline, disconnect
signals and lingering background descendants also stop that group. Unknown or
failed outcomes must not be replayed. Broker transport bounds still apply;
truncation or nonzero exit cannot become complete evidence.

For root-class enrollment the sudoers rule is staged and parsed in a dot-named
file on the **destination filesystem**, before identity creation. It is published
without replacement. Separate `/var` and `/etc` mounts therefore do not create
an EXDEV failure after user allocation. The public key is installed last.

`--revoke` first verifies the original packet-derived current postimage. It
retains authorization as `authorized_keys.revoked`, requests that the owned
account be locked/expired, and retains account, script and receipt evidence.
For a root-class enrollment it also retains the exact sudo policy under a
`.revoked` name and removes the active NOPASSWD rule. Both dispatcher modes
require the exact active public authorization before admitting execution; the
root helper additionally verifies its exact policy. The private journal is not
made readable to the restricted account. An interrupted revoke after key removal
therefore blocks new execution even if account expiry or policy cleanup fails.
Changed objects refuse instead of overwriting another owner's work. Its result
explicitly does **not** claim that already authenticated sessions are drained.
Broker stop/drain and uncertainty reconciliation remain separate obligations.

This initial candidate enrolls one command. It supplies no arbitrary remote
shell, script upload, registry updater or complete release installation path.
Additional deployment admission must use the settled Connect/Runtime contract;
a metadata enrollment is not a working go-live distribution. Actual key custody,
host-key trust, source IP, broker installation, grants, native authentication,
state compatibility and real command acceptance remain qualification inputs.

`make test-bootstrap` runs the package's race-enabled tests; `make verify` includes
it through the normal Go test target. State-transition fixtures use temporary
files and model NSS/account/policy tools; they never run real useradd, sudo or SSH.
Other fixtures exercise native Linux sealed snapshots, output/deadline process
control and locking. They do not establish target account/PAM/SSH compatibility,
live enrollment, an installed broker or full interrupted-operation recovery.

The Linux source gate at `61c93d04e74371e8ca1714d0344b23c9207f8184`
passed vet and race-enabled tests in all six packages with Go 1.27.0. Its original
log SHA256 is `2cf9014f08f0c951ad64d76c4964161dcc92bdefa72472391ed93c54d7071530`.
Subsequent clean-source cross-compilation produced Linux AMD64 and ARM64
executables; neither was executed. Native account/PAM/SSHD acceptance remains
pending. Owned process-group cleanup does not establish containment of
`setsid` or double-fork descendants, and resource containment is not qualified.
No production root enrollment is cleared by these source results.

The subsequent repair source rejects IPv6 zones before key-line construction,
removes owned elevation during revocation, and checks active authorization in the
local dispatcher. These repairs and regression fixtures passed the Linux source
gate at `42ac5d163fdb2723a4dd402e93917541a770a48b`: Go 1.27.0 vet and race tests
across seven packages. Original full log SHA256:
`f326d2802dd7c3b84496365923aa0d3a5ceb09a91f984d3e9260d047d55a3946`.

That result covers the repair source, not native account/SSH acceptance,
installed revocation or full containment. Later documentation-only successors
change no executable or test bytes and do not automatically acquire a new gate
or whole-head approval. Existing JSON intent/completion journals still need
migration to the required SQLite-only service-state policy before a compliant
distribution can be offered.
