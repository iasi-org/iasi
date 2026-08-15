# 09 — Pre-implementation readiness

> **Status: blocking check before implementation**

The Copilot adapter specification is ready for implementation once the installed IASI instruction set is itself valid.

# Current blocking condition

The current methodology contains a duplicated instruction ID:

```text
documentation.style
```

Duplicate instruction IDs are invalid under the adapter contract.

Before implementation or acceptance testing proceeds, the two files declaring this ID MUST be identified and the methodology corrected so that each instruction has a unique ID.

The adapter MUST NOT invent which instruction should be renamed and MUST NOT resolve the collision implicitly.

The correction belongs to the IASI instruction source.

After correction, verify that a scan of instruction candidates contains exactly one declaration of:

```text
documentation.style
```

and no other duplicated IDs.

# Implementation readiness

Implementation may proceed when:

- the duplicate `documentation.style` ID is resolved in the methodology source;
- the canonical `adapter.yml` uses `schema_version` and no independent `version`;
- documents 06 and 07 are consistent with document 08;
- `.github/.iasi/` behavior is implemented according to the CLI specification.

Document 08 remains the final normative clarification for the current milestone.
