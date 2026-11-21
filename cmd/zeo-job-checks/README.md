# Job check composition

Candidate developer tool. The repaired source at
`42ac5d163fdb2723a4dd402e93917541a770a48b` passed the Linux `make verify` gate:
vet and race-enabled tests passed across all seven packages. The complete
original log SHA256 is
`f326d2802dd7c3b84496365923aa0d3a5ceb09a91f984d3e9260d047d55a3946`.
This is source qualification; target enrollment and an installed distribution
remain unqualified.

`zeo-job-checks --baseline approved.json --baseline-sha256 SHA256 --catalog owned.json --inputs DIR --out /absolute/new-checks.sh`
publishes a complete shell file only after validating every current source pin and check
anchor and comparing stable check IDs with the pinned last approved catalog.
The CLI requires a fresh absolute output path and emits no executable stdout.
A same-directory temporary file is fully written and fsynced before no-replace
publication. Refusal and short writes leave no new runnable fragment. Invoke the
file only after the composer succeeds (`&&`); never run an old file after an
error. Do not copy selected commands into a separate list. See
[usage and limits](../../docs/job-checks.md).

The trusted, reviewed catalog is the single list of required checks and their
ordered invocations. Several checks may share an invocation; it is emitted once
at the first occurrence, while every stable ID remains in the output. Removing
or renaming an approved ID fails unless `retired` names that missing ID with a
nonempty owner decision/reason. A retirement for an unknown or active ID fails.
Keeping an ID while changing its source, anchor or argv requires a reviewed
`amended` entry with `id`, `reason`, `previous_sha256` and `current_sha256`.
These digests bind the old/new check objects encoded by Go JSON in struct field
order (id, source, contains, argv). Ordinary source-file hash repins do not
change a check definition. Decisions remain trusted reviewed inputs, not approval
that the generator can grant itself.

Catalog shape:

```json
{
  "schema": "zeo.job-checks/v1",
  "sources": [{"path": "scanner.py", "sha256": "OWNER_FILE_SHA256"}],
  "checks": [{
    "id": "app.old-stamp-refusal",
    "source": "scanner.py",
    "contains": "def test_old_stamp_refused(",
    "argv": ["python3", "-m", "unittest", "scanner"]
  }],
  "retired": [],
  "amended": []
}
```

Use one ID per required test or assertion, including consumer assertions inside
a harness, and bind those assertions to the harness invocation. Paths are
confined to the input root. The emitted fragment first changes to that same
canonical input root, so relative commands do not resolve against the caller's
directory. Exact file hashes preserve the whole owner handoff;
anchors check that required method/assertion text survives an owner repin.
The generator does not interpret program semantics: an anchor in a comment or
an unreachable branch is not behavioral proof. Review and the real gate remain
required. It also cannot establish that the baseline was actually approved;
the reviewed baseline hash and catalog provenance belong in the job input record.

Arguments are shell-quoted literals, except an entire argument of the form
`${UPPERCASE_NAME}`, which expands a required shell variable as one quoted
argument. Executables must be literal. Commands are not executed by the tool.
The surrounding admitted runner still owns bounds, source checkout, environment,
logs and copy-back. Catalogs are trusted job definitions, not a sandbox for
untrusted commands or a guarantee against concurrent source mutation.

Regression source covers the historical dropped old-stamp and consumer checks,
renaming, missing/invalid retirement reasons, assertion loss on a legitimate
file repin, modified owner/baseline bytes, and complete invocation generation.
Run the repository `make verify` on the admitted Linux runner; do not execute
candidate code outside the admitted Linux runner.
