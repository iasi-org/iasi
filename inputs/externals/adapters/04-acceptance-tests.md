# Acceptance tests

## Existing commands

All existing tests for:

```text
iasi install --workspace
iasi status
iasi version
```

must continue to pass.

## Installation

Verify that a workspace installation contains:

```text
.iasi/adapters/schema/adapter.md
.iasi/adapters/copilot/README.md
.iasi/adapters/copilot/adapter.yml
```

Verify that `manifest.yml` records adapters as installed.

## Source resolution

Create a temporary structure:

```text
workspace/
├── .iasi/
└── project/
```

Run the adapter with `project/` as current directory.

Verify that it uses `workspace/.iasi`.

## Repository-wide instructions

Given active `general` instructions, verify creation of:

```text
project/.github/copilot-instructions.md
```

Verify that:

- it contains the IASI ownership marker;
- it contains the installed IASI version;
- it contains active general instructions;
- YAML front matter from the source instructions is not copied;
- instructions are ordered deterministically.

## Path-specific instructions

Verify generation of at least:

```text
.github/instructions/documentation.instructions.md
.github/instructions/code.instructions.md
.github/instructions/diagrams.instructions.md
.github/instructions/inputs.instructions.md
```

when their scopes contain active instructions.

Verify that each begins with the exact `applyTo` value defined in `adapter.yml`.

## Status filtering

Create representative instructions with:

```yaml
status: active
status: draft
status: deprecated
```

Verify that only `active` instructions are projected.

Verify that:

```text
README.md
schema/
```

are not projected.

## Unknown scopes

Create a valid active instruction whose scope is not mapped by the Copilot adapter.

Run:

```bash
iasi adapt copilot
```

Verify that:

- the command fails with exit code `1`;
- the unsupported scope is reported;
- the instruction is not inserted into another scope;
- no generated target is created, modified or deleted.

This verifies that unknown active scopes are preflight errors rather than silently skipped content.

## Empty scope

If a mapped scope has no active instructions, its target file must not be created.

## Human-owned collision

Pre-create:

```text
.github/copilot-instructions.md
```

without the IASI ownership marker.

Run:

```bash
iasi adapt copilot
```

Verify that:

- the command fails;
- the existing file is unchanged;
- no other target file is modified or created.

This test verifies preflight atomicity.

## Regeneration

Generate the Copilot files once.

Modify an active instruction inside the installed `.iasi`.

Run:

```bash
iasi adapt copilot
```

again.

Verify that IASI-owned files are regenerated and contain the modified installed instruction.

This proves that adaptation uses installed `.iasi`, not embedded source data.

## Idempotency

Run the adapter twice without changing inputs.

Verify that generated files are byte-for-byte identical.

## Independence from source repository

The acceptance test must work using:

```text
iasi.exe
workspace/.iasi/
project/
```

without access to the source `iasi` repository.

## Build

The milestone is complete when:

```bash
go test ./...
```

passes and the existing build process still produces the standalone executable:

```bash
go build -o iasi.exe ./cmd/iasi
```

A copied standalone executable must be able to adapt Copilot from an installed `.iasi`.
