# Final clarifications for Copilot adapter implementation

> **Normative status: FINAL**
>
> This document is normative for the current Copilot adapter milestone. Where an earlier document conflicts with this one, this document takes precedence. Earlier documents should nevertheless be corrected when the conflict is known, so the specification remains internally consistent.

This document closes the remaining behavior required before implementation.

No additional behavior should be inferred outside the explicit specification.

---

# 1. Canonical `adapter.yml`

The canonical Copilot adapter descriptor is the one defined in:

```text
01-adapter-definition.md
```

It uses:

```yaml
schema_version: 1
id: copilot
platform: github-copilot
```

No other document or implementation may use the obsolete field:

```yaml
version:
```

`schema_version` is mandatory and has no default value.

Any earlier wording that described `schema_version` as optional is superseded by this document and MUST be removed or updated.

If an older example remains elsewhere, it must be replaced rather than interpreted as an alternative format.

---

# 2. Unknown instruction scopes

If an instruction is valid and has:

```yaml
status: active
```

but its `scope` is not mapped by the Copilot adapter, adaptation MUST fail during preflight.

Example:

```yaml
---
id: example.rule
status: active
scope: unknown-scope
---
```

must produce an error equivalent to:

```text
Unsupported instruction scope for Copilot adapter: unknown-scope
```

The command MUST:

- return exit code `1`;
- modify no project files;
- generate no partial output.

For this iteration, exit codes are:

```text
0 → success
1 → error
```

No richer exit-code taxonomy is required.

---

# 3. Exclusion order for `schema/` and README files

Non-instruction content MUST be excluded before instruction parsing, front-matter validation, ID collection or duplicate-ID detection.

The discovery order is:

```text
walk .iasi/instructions/
        ↓
exclude schema/**
        ↓
exclude README.md and README_*.md
        ↓
parse remaining Markdown candidates
        ↓
validate front matter
        ↓
validate IDs
        ↓
detect duplicate IDs
```

Therefore:

- files below `instructions/schema/` are not instruction candidates;
- `README.md` and localized README variants are not instruction candidates;
- excluded files do not participate in duplicate-ID validation.

---

# 4. Unknown instruction status

The only valid instruction status values in this iteration are:

```text
active
draft
deprecated
```

Any other value makes the instruction invalid.

Example:

```yaml
status: experimental
```

must fail during preflight.

A missing `status` field is also invalid.

The behavior of valid statuses is:

```text
active       → validate and adapt
draft        → validate but do not adapt
deprecated   → validate but do not adapt
```

`draft` and `deprecated` remain valid instructions. They are simply excluded from generated Copilot output.

Unknown statuses MUST NOT be silently treated as inactive.

---

# 5. First execution without `copilot-manifest.yml`

The absence of:

```text
.github/.iasi/copilot-manifest.yml
```

is valid and means:

```text
first adaptation
```

It is not an error.

During the first adaptation:

1. calculate all required targets;
2. inspect all existing target files;
3. if a target does not exist, it may be generated;
4. if a target exists and contains the IASI ownership marker, it may be regenerated;
5. if a target exists without the IASI ownership marker, adaptation MUST fail;
6. do not attempt stale-file deletion because no previous manifest exists;
7. after a successful commit, create the first `copilot-manifest.yml`.

The first manifest becomes the baseline for future stale-output detection.

From the second successful run onward:

```text
previous generated set
        -
current generated set
        =
stale candidates
```

Stale-file removal must follow the ownership and safety rules already defined in the specification.

---


# 6. Scope clarification from the previous versioning milestone

A previous IASI versioning specification listed `adapters` among the items outside the scope of that specific versioning iteration.

That statement applied only to the earlier versioning milestone.

For the current milestone:

```text
Copilot adapter
```

is explicitly in scope.

The current scope remains limited to adapting:

```text
instructions
```

Commands, skills, MCP, custom agents and prompt files remain outside the current adapter implementation.


# Final execution flow

The final implementation flow is:

```text
discover installed .iasi
        ↓
load and validate Copilot adapter
        ↓
discover instruction candidates
        ↓
exclude schema and README content
        ↓
validate front matter, status and IDs
        ↓
detect duplicate IDs
        ↓
select status: active
        ↓
validate mapped scopes
        ↓
load and validate prior Copilot manifest if present
        ↓
calculate targets and stale candidates
        ↓
preflight collisions and ownership
        ↓
stage complete output
        ↓
commit atomically
        ↓
write new copilot-manifest.yml
```

Any error before commit must leave the project unchanged.

Any ambiguity about ownership, target validity, instruction validity or adapter validity must stop adaptation rather than producing partial or approximate output.
