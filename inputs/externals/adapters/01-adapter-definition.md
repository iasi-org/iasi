# Adapter definition

## Purpose

Define the minimum methodological representation of a platform adapter.

An adapter describes how IASI artifacts are projected to a specific platform. It contains platform mapping rules, not the IASI rules themselves.

## `agentics/adapters/schema/adapter.md`

Create this document with the following minimum model.

An adapter has:

```text
identity
platform
supported IASI artifact types
targets
mapping rules
generation constraints
```

The first implementation only needs to support `instructions`.

The schema must explicitly state:

- adapters MUST NOT duplicate IASI instruction content;
- adapters MUST NOT redefine IASI behavior;
- platform-specific paths and formats belong to the adapter;
- unsupported IASI artifact types remain unsupported rather than being approximated;
- generated artifacts are projections and are not a source of truth.

Do not design a generic plugin system yet.

## `agentics/adapters/copilot/README.md`

Create a short README explaining that this adapter projects IASI instructions into GitHub Copilot repository custom instructions.

Document these targets:

```text
.github/copilot-instructions.md
.github/instructions/*.instructions.md
```

Explain that custom agents, prompt files, skills and MCP are intentionally outside the current adapter scope.

## `agentics/adapters/copilot/adapter.yml`

Use this canonical declarative structure:

```yaml
schema_version: 1
id: copilot
platform: github-copilot

supports:
  instructions: true
  commands: false
  skills: false
  mcp: false
  agents: false

instructions:
  general:
    type: repository
    target: .github/copilot-instructions.md

  documentation:
    type: path
    target: .github/instructions/documentation.instructions.md
    applyTo: "**/*.md,**/*.qmd"

  code:
    type: path
    target: .github/instructions/code.instructions.md
    applyTo: "**/*.go,**/*.py,**/*.js,**/*.jsx,**/*.ts,**/*.tsx,**/*.java,**/*.kt,**/*.kts,**/*.c,**/*.h,**/*.cpp,**/*.hpp,**/*.cs,**/*.rs,**/*.rb,**/*.php,**/*.R,**/*.r,**/*.lua,**/*.sh,**/*.ps1,**/*.sql"

  diagrams:
    type: path
    target: .github/instructions/diagrams.instructions.md
    applyTo: "**/*.puml,**/*.plantuml,**/*.mmd,**/*.dot"

  knowledge:
    type: path
    target: .github/instructions/inputs.instructions.md
    applyTo: "**/inputs/**/*"
```

The field:

```yaml
version:
```

MUST NOT be used.

The distribution version of the adapter is always the IASI version defined by the repository-level:

```text
VERSION
```

`schema_version` is mandatory and has no default value.

If `schema_version` is missing or unsupported, the adapter is invalid and adaptation must stop during preflight.

The mapping belongs in `adapter.yml`, not scattered through Go code.

The Go implementation may validate supported values, but must obtain target paths and `applyTo` mappings from the adapter definition.

If YAML parsing is not already available, using a small maintained YAML dependency is acceptable. Do not implement a home-grown YAML parser.

## Minimum adapter validity

An adapter is valid only when:

- its directory exists under `.iasi/adapters/`;
- `adapter.yml` exists;
- the YAML parses successfully;
- `schema_version` exists and is supported;
- `id` exists;
- `platform` exists;
- `id` matches the adapter requested by the user;
- every declared supported artifact type contains the required mapping definition;
- every configured target is valid for the adapter.

For:

```bash
iasi adapt copilot
```

the descriptor MUST contain:

```yaml
id: copilot
```

An invalid adapter must fail during preflight and must not modify the target project.

## Versioning

The adapter version is part of the IASI distribution.

Do not create independent release mechanics for adapters.

IASI `VERSION` remains the distribution version and the unit that is installed and shipped.
