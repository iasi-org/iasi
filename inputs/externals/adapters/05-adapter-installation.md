# Copilot adapter installation and lifecycle

## Purpose

Define how the GitHub Copilot adapter is distributed, installed and applied.

The Copilot adapter is **part of IASI**.

It is not installed as an independent product and does not have its own installer.

The lifecycle is:

```text
embed
  ↓
install
  ↓
adapt
```

Each stage has a different meaning.

---

# 1. Embed

The source adapter lives in the IASI repository:

```text
iasi/
└── agentics/
    └── adapters/
        ├── schema/
        │   └── adapter.md
        └── copilot/
            ├── README.md
            └── adapter.yml
```

The complete `agentics/` tree, including adapters, is included in the standalone IASI executable during the build process.

Conceptually:

```text
iasi repository
      │
      ▼
   build
      │
      ▼
iasi.exe
├── VERSION
└── embedded agentics/
    ├── instructions/
    └── adapters/
        └── copilot/
```

The adapter MUST travel with the same IASI version as the methodology it adapts.

There is no independent Copilot-adapter release mechanism in this iteration.

---

# 2. Install

Running:

```bash
iasi install --workspace
```

installs the complete available IASI distribution into:

```text
workspace/.iasi/
```

The installation MUST include adapters.

Expected structure:

```text
workspace/
└── .iasi/
    ├── manifest.yml
    ├── instructions/
    ├── commands/
    ├── skills/
    ├── mcp/
    └── adapters/
        ├── schema/
        │   └── adapter.md
        └── copilot/
            ├── README.md
            └── adapter.yml
```

Therefore, after `iasi install --workspace`, the Copilot adapter is **available**.

It has not yet modified any project.

No `.github/` files are generated during installation.

`install` installs IASI.

It does not apply a platform adapter.

---

# 3. Manifest

The workspace manifest MUST record adapters as part of the installed distribution.

Example:

```yaml
version: 0.1.0
profile: workspace

installed:
  instructions: all
  commands: all
  skills: all
  mcp: all
  adapters: all
```

The installed adapter belongs to the same IASI version declared by the manifest.

Do not introduce a separate installed version for the Copilot adapter.

---

# 4. Adapt

The Copilot adapter is applied explicitly with:

```bash
iasi adapt copilot
```

This command MUST NOT install or download the adapter.

It MUST use the adapter already present in the effective installed `.iasi`.

Conceptually:

```text
workspace/.iasi/
├── instructions/
└── adapters/
    └── copilot/
        │
        ▼
iasi adapt copilot
        │
        ▼
project/.github/
```

The adapter combines:

```text
installed IASI instructions
        +
installed Copilot adapter definition
```

to generate Copilot-native artifacts.

---

# 5. Workspace inheritance

A workspace installation may serve multiple projects.

Example:

```text
workspace/
├── .iasi/
│   ├── instructions/
│   └── adapters/
│       └── copilot/
│
├── project-a/
├── project-b/
└── project-c/
```

From:

```text
workspace/project-a/
```

running:

```bash
iasi adapt copilot
```

must search upward for `.iasi` and use:

```text
workspace/.iasi/
```

as its source.

The generated Copilot files belong to the current project:

```text
workspace/project-a/.github/
```

They MUST NOT be generated in:

```text
workspace/.github/
```

unless the workspace root itself is the current target.

This allows one IASI installation to support several independent projects without duplicating the methodology.

---

# 6. Availability and application

Keep these concepts distinct:

```text
available
active
adapted
```

## Available

The adapter exists in the installed `.iasi`.

Example:

```text
.iasi/adapters/copilot/
```

## Active

The installed IASI instructions applicable to the current context are available to be processed.

Full instruction resolution will be implemented later. For this iteration, use the installed active instructions according to the existing adapter specification.

## Adapted

The platform-specific representation has been generated in the target project.

Example:

```text
.github/copilot-instructions.md
.github/instructions/*.instructions.md
```

An available adapter is not automatically an adapted project.

---

# 7. No separate adapter installation command

Do not implement:

```bash
iasi install adapter copilot
```

Do not implement:

```bash
iasi adapter install copilot
```

Do not require users to copy adapter files manually.

For a workspace installation, every adapter included in that IASI distribution is installed automatically as part of:

```bash
iasi install --workspace
```

The explicit user action is only required when applying an adapter:

```bash
iasi adapt copilot
```

---

# 8. Source of truth

There are three different representations and each has a clear role:

```text
iasi repository
    │
    │ build
    ▼
iasi.exe
    │
    │ install
    ▼
workspace/.iasi/
    │
    │ adapt
    ▼
project/.github/
```

Their responsibilities are:

```text
iasi repository   → source distribution
iasi.exe          → self-contained packaged distribution
workspace/.iasi   → installed methodology
project/.github   → generated Copilot representation
```

For runtime adaptation, `.iasi` is the source of truth.

`iasi adapt copilot` MUST NOT bypass the installed `.iasi` and adapt directly from the methodology embedded in the executable.

This is important because the installed workspace is the methodological context being used by the project.

---

# 9. Regeneration

Generated Copilot files are derived artifacts.

They MAY be regenerated by running:

```bash
iasi adapt copilot
```

again.

Only files carrying the IASI generated-file ownership marker may be overwritten automatically.

Files created manually by the user MUST NOT be overwritten.

The rules for collision handling and regeneration are defined in the Copilot adapter CLI specification.

---

# 10. Future project installations

A future command:

```bash
iasi install --project
```

may install or activate only a subset of IASI.

Do not implement this behavior in the current milestone.

The Copilot adapter implementation SHOULD avoid assumptions that require every future installation to contain the complete IASI catalog.

For the current workspace implementation:

```bash
iasi install --workspace
```

installs all available adapters.

---

# Acceptance criteria

The implementation is correct when the following sequence works with a standalone executable:

```bash
iasi version
iasi install --workspace
```

and produces:

```text
.iasi/
├── instructions/
└── adapters/
    └── copilot/
```

Then, from a child project:

```bash
cd project-a
iasi adapt copilot
```

must use the parent workspace installation and generate:

```text
project-a/.github/
```

without:

- accessing the source IASI repository;
- downloading the Copilot adapter;
- installing another copy of IASI;
- generating Copilot files during `iasi install`;
- modifying another project in the workspace.

The intended lifecycle is:

```text
IASI distribution
      ↓
install once in workspace
      ↓
adapter available
      ↓
adapt explicitly per project
```
