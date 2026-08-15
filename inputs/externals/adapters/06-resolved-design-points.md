# 06 — Resolved design points

> **Status: normative, subject to `08-final-clarifications.md`**
>
> This document records decisions made before the final clarification pass.
> If any wording here conflicts with `08-final-clarifications.md`, document 08 takes precedence.
> Known conflicts have been removed from this revision.

# 1. Adapter versioning

Adapters do not have an independent distribution version.

The distribution version is always the IASI version defined by the repository-level:

```text
VERSION
```

`adapter.yml` MUST NOT contain:

```yaml
version:
```

The adapter descriptor MUST use:

```yaml
schema_version: 1
id: copilot
platform: github-copilot
```

`schema_version` versions the descriptor format, not the adapter release.

`schema_version` is mandatory and has no default value.

# 2. Invalid or unknown adapters

`iasi adapt <name>` MUST perform complete preflight before modifying the target project.

The requested adapter is invalid when, among other validation failures:

- its directory does not exist;
- `adapter.yml` does not exist;
- the YAML is invalid;
- `schema_version` is missing or unsupported;
- `id` is missing;
- `platform` is missing;
- the descriptor `id` does not match the requested adapter;
- required mappings are missing;
- a configured target is invalid.

If the requested adapter does not exist or is unsupported, the command MUST fail with exit code `1`.

No target file may be modified.

# 3. Stale generated files

The Copilot adapter tracks its generated files in:

```text
.github/.iasi/copilot-manifest.yml
```

The manifest is an IASI-managed artifact.

A previously generated file may be removed as stale only when:

1. it is recorded in the previous valid Copilot manifest;
2. its normalized path remains inside `.github/`;
3. it exists;
4. it still contains the `IASI-GENERATED` ownership marker.

If ownership cannot be established, adaptation MUST fail rather than delete the file.

# 4. Meaning of `active`

In this iteration, instruction activation is defined only by:

```yaml
status: active
```

There are no project profiles, precedence resolution, task activation rules or dynamic resolution mechanisms yet.

Valid status values are:

```text
active
draft
deprecated
```

Their behavior is:

```text
active       → validate and adapt
draft        → validate, do not adapt
deprecated   → validate, do not adapt
```

Any missing or unknown status is invalid and MUST fail preflight.

# 5. Available adapters

For a workspace installation, all real adapters included in the current IASI distribution under:

```text
agentics/adapters/
```

are installed into:

```text
.iasi/adapters/
```

The directory:

```text
agentics/adapters/schema/
```

contains the adapter schema and is not itself an adapter.

Installed adapters are available, not automatically applied.

Only an explicit command such as:

```bash
iasi adapt copilot
```

materializes platform-specific output.

# 6. Authority of final clarification

The detailed rules for:

- unknown scopes;
- candidate discovery order;
- valid instruction statuses;
- first execution without a Copilot manifest;
- final execution order;

are defined normatively by:

```text
08-final-clarifications.md
```
