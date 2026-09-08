Last updated: 2026-09-08T14:09:59Z (UTC)

# Producer Validation Scripts

These scripts belong to the producer. The maintainer entry points are
`docs/repo/runbooks/change-operating-kit.md` and `docs/repo/runbooks/release-operating-kit.md`.
Use their authority, source and tooling-readiness boundaries before invocation.

| Script | Role and caller | Inputs, behavior and evidence | Tests |
| --- | --- | --- | --- |
| `validate-guidance-candidate.py` | Read-only eligibility primitive; release preflight | Git, Python 3.11+, PyYAML; exact accepted broad-source `--baseline-ref`, optional `--candidate-ref` (HEAD). Compares committed trees only; concise success or nonzero blocker, optional `--json`. Mechanical eligibility does not establish owner acceptance or semantic policy review. | `tests/test_guidance_candidate.py` |
| `validate-release-identity.py` | Read-only identity primitive; feedback and candidates | Python 3.11+, PyYAML; repository root by default, optional `--root`. Checks authoritative versions, declaration checksums and mirrors; nonzero on mismatch. Go manifest tests separately prove compiled graph identity. | `tests/test_release_identity.py`, `internal/manifest/manifest_test.go` |
| `verify-guidance-lifecycle.py` | Isolated native smoke workflow; guidance candidate | Python 3.11+, PyYAML, native shell/installer; `--fresh-target`, explicit published `--upgrade-version`, staged `--catalog`, new `--work-dir`. Compares every managed source/target, verifies/downloads matching old CLI through native installer, then init/dry-run/failure/apply/check and preservation in isolated directories. Fails nonzero on any native error or mismatch; prints counts and tested versions. No live consumer or publication action. | `tests/test_guidance_lifecycle.py`, both hosted guidance lanes |
| `build-release-assets.py` | Existing deterministic pack builder; candidate | Go and native build tools; `--output-dir`, optional version/platform/base URL. Builds twice, verifies pack identity and writes candidate packs/catalog. Does not run broad source acceptance. | `tests/test_release_assets.py`, native candidate lanes |

The scripts reuse existing stdout/exit conventions. No authority registry, receipt parser or
persistent evidence service is introduced. Lifecycle smoke only writes its explicitly new test
workspace and the existing fresh test target is read-only. Its public download is limited to the
explicit published version through existing installers. Preserve its output in the workflow run;
record acceptance and applicability in the plan execution log.
