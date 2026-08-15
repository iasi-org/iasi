# 07 — Initial specification closure

> **Status: normative, subject to `08-final-clarifications.md`**
>
> This document defines the safety and validation contract established before the final clarification pass.
> `08-final-clarifications.md` is the final authority where later clarification was required.

# 1. Canonical adapter descriptor

The Copilot adapter MUST use:

```yaml
schema_version: 1
id: copilot
platform: github-copilot
```

`schema_version` is mandatory.

The adapter MUST NOT declare an independent:

```yaml
version:
```

IASI `VERSION` remains the only distribution version.

# 2. Minimum adapter validity

An adapter is valid only when:

- its directory exists under `.iasi/adapters/`;
- `adapter.yml` exists;
- YAML parsing succeeds;
- `schema_version` is present and supported;
- `id` is present;
- `platform` is present;
- `id` matches the adapter requested by the CLI;
- mappings required for supported artifact types are valid;
- every target is valid and contained inside `.github/`.

An invalid adapter MUST fail during preflight and MUST NOT modify the project.

# 3. Copilot manifest

The adapter-managed manifest lives at:

```text
.github/.iasi/copilot-manifest.yml
```

Initial format:

```yaml
schema_version: 1
adapter: copilot
iasi_version: 0.1.0

generated:
  - .github/copilot-instructions.md
  - .github/instructions/code.instructions.md
```

The manifest MUST carry an IASI ownership marker.

An existing manifest MUST be validated before use:

- YAML parses successfully;
- `schema_version` is supported;
- `adapter` is `copilot`;
- all generated paths normalize inside `.github/`.

An invalid existing manifest MUST fail preflight.

The absence of the manifest is valid and represents the first adaptation, as defined in `08-final-clarifications.md`.

# 4. Atomicity

Adaptation is logically:

```text
preflight
   ↓
staging
   ↓
commit
```

Preflight MUST calculate and validate the complete operation without modifying the project.

Staging MUST generate complete new content in temporary storage.

Commit may begin only after preflight and staging succeed.

If commit fails, the implementation MUST attempt to restore the prior state of every affected file.

The definitive Copilot manifest is written only after the generated outputs have been committed successfully.

# 5. Target boundary

All Copilot adapter outputs MUST remain inside:

```text
<target-project>/.github/
```

This includes adapter metadata stored under:

```text
.github/.iasi/
```

Paths MUST be normalized before validation.

Targets escaping `.github/` through absolute paths, `..`, normalization or any equivalent mechanism are invalid.

# 6. Instruction candidates

Before validating instructions, discovery MUST exclude:

```text
instructions/schema/**
README.md
README_*.md
```

The remaining Markdown candidates MUST contain valid IASI front matter with, at minimum:

```yaml
id:
status:
scope:
```

Only these candidates participate in instruction validation and duplicate-ID detection.

# 7. Instruction IDs

Instruction IDs MUST be unique within the effective installed `.iasi`.

If two candidate instructions declare the same ID, adaptation MUST fail during preflight.

The error SHOULD identify:

- the duplicated ID;
- every file declaring it.

Duplicate IDs MUST NOT be resolved by file order, path, modification date or any implicit precedence rule.

# 8. Active instructions and scopes

Valid status values are:

```text
active
draft
deprecated
```

Only `active` instructions are generated.

`draft` and `deprecated` instructions remain valid but are not projected.

An active instruction whose scope is not mapped by the Copilot adapter is a preflight error, not a warning or skip.

# 9. Stale outputs

A stale file may be deleted only when:

- it was recorded by the prior valid Copilot manifest;
- it is inside `.github/`;
- it still contains the IASI ownership marker.

Any ownership ambiguity MUST fail the operation.

# 10. Safety rule

When the adapter cannot determine unambiguously:

- what to generate;
- what belongs to IASI;
- what may be replaced;
- what may be deleted;
- where it may write;

it MUST fail rather than approximate.

# 11. Final clarification

The final rules for discovery order, unknown scopes, unknown statuses and first-run behavior are defined by:

```text
08-final-clarifications.md
```
