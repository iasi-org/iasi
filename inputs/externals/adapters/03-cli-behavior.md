# `iasi adapt copilot`

## Command

Add:

```bash
iasi adapt copilot
```

No additional flags are required in this iteration.

## Source lookup

Use the same upward `.iasi` lookup already used by `iasi status`.

For example:

```text
C:/workspace/.iasi/
C:/workspace/project-a/
```

Running:

```bash
cd C:/workspace/project-a
iasi adapt copilot
```

must use:

```text
C:/workspace/.iasi/
```

as the installed IASI source.

## Target

For this iteration, the current working directory is the target project root.

Generated files therefore go below:

```text
<cwd>/.github/
```

The Copilot adapter also owns an internal metadata area:

```text
<cwd>/.github/.iasi/
```

This directory is reserved for IASI metadata associated with generated Copilot artifacts.

For this iteration it contains:

```text
.github/.iasi/copilot-manifest.yml
```

`.github/.iasi/` is not a user-facing Copilot instruction target. It exists only to support safe regeneration, ownership tracking and stale-output management.

The same `.github/` containment rules apply to this metadata directory. The adapter MUST NOT write adapter metadata outside the target project's `.github/` tree.

Do not attempt Git repository discovery yet.

Do not write generated Copilot files into the directory that contains `.iasi` unless that directory is also the current target.

This allows one workspace installation to serve multiple projects.

## Preflight

Before writing any file:

1. resolve the `.iasi` source;
2. load the Copilot adapter definition;
3. discover active instructions;
4. calculate every target file;
5. inspect all existing target files.

Only after preflight succeeds may files be changed.

This prevents partial generation.

## Existing files

Generated artifacts must be safe to regenerate.

Rules:

### Target does not exist

Create it.

### Target exists and contains the IASI ownership marker

Replace it with newly generated content.

### Target exists and does not contain the IASI ownership marker

Fail the command.

Do not modify that file.

Do not modify any other generated target during the failed run.

Do not implement `--force` yet.


## Copilot manifest behavior

Before writing generated outputs, inspect:

```text
.github/.iasi/copilot-manifest.yml
```

If it does not exist, treat the operation as the first adaptation. This is valid.

If it exists, validate it before using it for ownership or stale-file decisions.

A malformed, unsupported or unsafe manifest is a preflight error.

The new manifest is written only after a successful commit of the generated Copilot files.

## Directories

Create only the directories required by the actual operation.

Possible managed directories in this iteration are:

```text
.github/
.github/instructions/
.github/.iasi/
```

`.github/.iasi/` is required when writing `copilot-manifest.yml`.

## Success output

Keep console output concise.

Example:

```text
IASI Copilot adapter

Source : C:/workspace/.iasi
Target : C:/workspace/project-a

Generated:
  .github/copilot-instructions.md
  .github/instructions/documentation.instructions.md
  .github/instructions/code.instructions.md
  .github/instructions/diagrams.instructions.md
  .github/instructions/inputs.instructions.md
```

Only list files actually generated.

## Failure examples

No installed IASI:

```text
IASI is not installed for this location.
```

Missing Copilot adapter:

```text
Copilot adapter is not available in this IASI installation.
```

Human-owned target collision:

```text
Cannot generate Copilot instructions because this file already exists and is not managed by IASI:
.github/copilot-instructions.md
```

Errors go to stderr and return a non-zero exit code.

## Runtime source rule

`iasi adapt copilot` must operate from installed `.iasi` content.

It must still work when the source `iasi` repository is completely unavailable.

It must also reflect deliberate changes made to the installed `.iasi` rather than silently falling back to the methodology embedded in the executable.
