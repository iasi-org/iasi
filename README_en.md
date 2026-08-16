[🇪🇸 Castellano](README.md) | **English**

# iasi

**Intelligent Systems Assisted Engineering Methodology**

`iasi` is the repository where the IASI methodology is defined, formalized and evolved.

The IASI ecosystem contains books, products, tools, documentation, infrastructure and engineering memory. This repository serves a different purpose: it contains the rules, structures and artifacts that describe **how we work with Intelligent Systems to perform engineering**.

It is not tied to a specific model, provider or tool. The methodology must be expressible independently and, when necessary, later adapted to environments such as Codex, Copilot, Claude or other systems.

---

# Purpose

The goal of `iasi` is to turn a way of working discovered through practice into an explicit, reproducible and challengeable methodology.

This repository formalizes aspects such as:

- how an agent should behave;
- which instructions it must follow;
- which capabilities can be reused;
- which actions can be exposed as commands;
- how external tools are integrated;
- how completed work is validated;
- which decisions require human intervention;
- and how the knowledge required for another human or agent to continue the work is preserved.

The methodology is not intended to remove human judgment or replace engineering.

It aims to create a framework in which humans and Intelligent Systems can collaborate in a controlled, traceable and reproducible way.

---

# Our approach

IASI does not start with a closed methodology and then attempt to fit problems into it.

We start from problems, observe how we solve them and formalize only what proves useful.

> **We do not adopt methodologies to solve problems. We start from problems to discover the methodology.**

Therefore, `iasi` must not become a collection of theoretical rules disconnected from practice.

Every element added to the methodology should address a real need, be explainable and remain open to review.

The methodology evolves through experience.

---

# IASI Structure

Canonical IASI artifacts live under `iasi/`. They define how AI agents
participate in the engineering process.

```text
iasi/
├── instructions/
├── commands/
├── skills/
├── mcp/
└── adapters/
```

Each part has a different purpose.

| Area | Purpose |
|------|---------|
| `instructions/` | Defines how an agent should behave and which rules it must follow when producing work. |
| `commands/` | Defines explicit actions that can be requested from agents. |
| `skills/` | Defines reusable capabilities that can be applied in different contexts. |
| `mcp/` | Defines the methodological integration with tools and services exposed through Model Context Protocol. |

These definitions are conceptual and should remain independent of each platform's concrete implementation.

---

# Instructions

`iasi/instructions/` contains persistent rules governing agent behavior.

Instructions may define, among other things:

- general behavior;
- human control;
- uncertainty handling;
- validation;
- tool use;
- rule precedence;
- code style;
- testing;
- document structure and writing style;
- source handling;
- diagram conventions.

Instructions are designed to be:

- **atomic**, focused on a single concern;
- **declarative**, without depending on a provider;
- **composable**, so they can be combined according to the task;
- **observable**, so compliance can be reviewed;
- **versionable**, because the methodology itself evolves.

Their common structure is defined in:

```text
iasi/instructions/schema/instructions.md
```

`iasi/commands/validate.md`, `iasi/commands/plan.md`,
`iasi/commands/execute.md`, `iasi/commands/verify.md`, and
`iasi/commands/archive.md` define the canonical agentic commands. Adapters may
project them to native mechanisms, but must not redefine them.

## CLI

```text
iasi install
iasi reinstall
iasi status
iasi version
iasi adapt copilot
```

A valid local installation is identified by `.iasi/manifest.yml`. Installed
layers in parent directories compose with local layers; `validation.json` and
`workflow.json` are local workflow state, not installed layers.

---

# Platform independence

IASI distinguishes between **the methodology** and **the concrete mechanism through which a tool consumes it**.

For example, the same instruction may eventually be expressed as:

- repository instructions;
- configuration files;
- prompts;
- skills;
- specialized agents;
- development environment rules;
- or any other mechanism supported by a platform.

The methodological source of truth remains in `iasi`.

Adaptations for specific tools should not determine the design of the methodology.

---

# Human and Intelligent System

The availability of an Intelligent System does not imply unlimited human availability.

IASI considers that the sustainable speed of an engineering system is not determined only by agent execution speed, but by the human capacity to:

- understand;
- decide;
- review;
- validate;
- and assimilate results.

Automation should reduce mechanical work, not turn every removed pause into another decision for the human.

The goal is not to maximize activity.

The goal is to improve our ability to perform good engineering.

---

# Validation

In IASI, executing a task is not the same as completing it.

The methodology distinguishes between:

```text
implemented
validated
accepted
```

A solution may be implemented without having been validated.

It may be technically validated without having been accepted by the human.

Whenever possible, validation criteria should rely on observable results, tests and acceptance criteria.

---

# Principles

- Engineering is the goal.
- Intelligent Systems are engineering assistants.
- Humans retain control over goals, priorities and acceptance.
- The methodology must remain independent of models and providers.
- Knowledge must be persisted explicitly.
- Instructions must be composable and verifiable.
- Automation should reduce complexity for the user, not transfer it.
- Implementation does not replace validation.
- Decisions must remain open to challenge.
- The methodology itself evolves.
- If implementation keeps getting cheaper, thinking becomes more valuable.

---

# Repository status

`iasi` is built incrementally.

We do not attempt to design the entire methodology in advance.

We first formalize the pieces that have already emerged as necessary in real projects. We then use, validate and modify them.

A structure is considered useful when it survives application across different cases without forcing artificial exceptions.

For that reason, the contents of this repository should be understood as a living methodology.

---

# Relationship with the IASI ecosystem

`iasi` defines the methodology.

The other repositories in the ecosystem provide the places where that methodology is used, tested, documented or turned into products.

The books explain the ideas.

The products materialize solutions.

Engineering memory preserves the path taken.

`iasi` attempts to turn what is learned along that path into an explicit and reusable way of working.

---

# An open project

IASI comes from practice and must remain open to challenge from practice.

The rules in this repository are not dogma. They are engineering decisions that should be justifiable, testable, refinable or replaceable when a better alternative exists.

**Every well-founded contribution is welcome.**

---

> *"The best methodology is not the one we follow, but the one we discover while solving the right problem."*

---

**IASI**

*Intelligent Systems Assisted Engineering*
