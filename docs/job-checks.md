# Reviewed job check composition

`zeo-job-checks` publishes a complete shell file from two reviewed `zeo.job-checks/v1`
catalogs. It does not run checks, approve a job, or establish source trust.
Build it on an admitted Linux builder with Go 1.27.0; no local Mac build is
required. The caller supplies a baseline file and its independently reviewed
SHA256, a current catalog, and the directory containing its pinned inputs.

```sh
zeo-job-checks --baseline baseline.json --baseline-sha256 "$REVIEWED_SHA" \
  --catalog current.json --inputs "$INPUTS" --out "$NEW_ABSOLUTE_FRAGMENT" && /bin/sh "$NEW_ABSOLUTE_FRAGMENT"
```

Always preserve the composer's exit status. `--out` is mandatory and must name a
fresh absolute path. The CLI emits no executable fragment to stdout. It writes
and fsyncs a private temporary file beside the destination and publishes it with
a no-replace link only after full validation and a complete write. Refusal or a
short temporary-file write leaves the requested new fragment absent; an existing
output is never replaced. Use a fresh job path and `&&`; ignoring an error and
running an unrelated old file cannot establish this job's acceptance. The tool
does not execute the generated file. This publication is not job admission.

The catalog contains `schema`, `sources` (path/SHA256), `checks`
(ID/source/contains/argv), and optional `retired` and `amended` lists.
Each active obligation has a permanent ID. Dropping or renaming it requires a
reviewed retirement naming the old ID and a reason. Changing its invocation,
source path or assertion anchor requires an amendment binding the old and new
definition digests plus a reason. Updating source bytes without changing that
definition requires a new source pin. A retained ID alone does not preserve a
check whose invocation was replaced by `true`.

The repair rejects duplicate and case-folded object keys at every schema level
before typed JSON decoding. Input files are opened under the validated root;
relative escapes are refused and generated commands enter that same root.
Repeated identical invocations execute once while every obligation is recorded.
Literal arguments are shell-quoted; only a complete `${NAME}` argument expands
as a required environment variable. The composer does not print its value.

Source hashes and anchor presence cannot prove a check's semantics or prevent
concurrent administrative mutation after validation. Review owns executable
selection, bounded runners, job isolation and the complete assertion set.
The generated script is not a signature or release-trust record.

`make verify` includes race-enabled regression tests for lost scanner/consumer
checks, same-ID invocation changes, retirements/amendments, wrong source pins,
root binding and shell quoting, plus ambiguity/refusal/short-write publication
regressions. Those repair source and test bytes were exercised by the Linux gate
at `42ac5d163fdb2723a4dd402e93917541a770a48b`: Go 1.27.0 vet and race tests passed
across seven packages. Original full log SHA256:
`f326d2802dd7c3b84496365923aa0d3a5ceb09a91f984d3e9260d047d55a3946`.

This is source-test evidence for that exact commit, not native job admission,
installation or a whole-current-head approval. Later documentation-only changes
preserve executable/test bytes but still require their own review disposition;
no exact-head gate or approval is silently transferred.
