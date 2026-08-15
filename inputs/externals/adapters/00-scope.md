# IASI Copilot adapter — Scope

## Goal

Implement the first platform adapter for IASI:

```text
IASI instructions
      ↓
Copilot adapter
      ↓
GitHub Copilot custom instructions
```

This iteration implements the **Copilot adapter**, but adapts only the IASI artifact type `instructions`.

Do not adapt yet:

- commands;
- skills;
- MCP;
- custom agents;
- prompt files.

Those concepts will be adapted only after their IASI model exists.

## Architectural boundary

IASI remains the source of truth.

```text
.iasi/instructions/
        │
        ▼
iasi adapt copilot
        │
        ▼
.github/
```

The generated Copilot files are projections of IASI. They must not become a second source of methodology.

The adapter must read the installed `.iasi` found for the current location. It must not read the source repository and must not use the copy of the methodology embedded in the executable as the source for adaptation.

This distinction is intentional:

```text
iasi.exe embedded methodology
        ↓
iasi install --workspace
        ↓
workspace/.iasi/          ← installed methodology
        ↓
iasi adapt copilot
        ↓
project/.github/
```

A project may therefore use an IASI installation inherited from a parent workspace.

## Repository structure

Add the methodological adapter definition under:

```text
agentics/
└── adapters/
    ├── schema/
    │   └── adapter.md
    └── copilot/
        ├── README.md
        └── adapter.yml
```

Add the Go implementation under the existing Go source root:

```text
src/go/
└── internal/
    └── adapters/
        └── copilot/
```

The CLI command remains part of the existing `iasi` executable.

## Installation impact

`iasi install --workspace` currently installs the contents of `agentics`.

Extend the installation so that adapters are installed too:

```text
.iasi/
├── instructions/
├── commands/
├── skills/
├── mcp/
└── adapters/
    ├── schema/
    └── copilot/
```

Update `manifest.yml` so the installed set also records:

```yaml
installed:
  adapters: all
```

If `status` reports installed categories, add an `Adapters` entry.

## Out of scope

Any previous specification that listed `adapters` as out of scope referred only to the earlier IASI versioning milestone. Adapters are explicitly in scope for the current milestone.

Do not implement:

```text
iasi adapt codex
iasi adapt claude
iasi update
iasi resolve
iasi validate
```

Do not generate `AGENTS.md`.

Do not add IASI concepts merely because Copilot supports them.

The milestone is complete when installed IASI instructions can be projected safely and deterministically into Copilot custom-instruction files.
