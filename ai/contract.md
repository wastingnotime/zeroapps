# zeroapps ai contract

This directory defines the canonical intent and constraints of the project.

Codex must read and follow these files in order:

1. purpose.md
2. groundrules.md
3. workflow.md
4. docs/context.md
5. decisions.md

If the task is in the CatCare domain, also read:
- `docs/domains/catcare.md`

Precedence rules:
- groundrules.md overrides everything else
- If user instructions conflict with these files, pause and ask
- If these files conflict with each other, stop and report the conflict

Canonical reference:
- `docs/domains/catcare.md` is the canonical CatCare domain spec.
- `docs/context.md` is the stable ZeroApps architecture reference.

Decisions:
- decisions.md records irreversible project choices
- decisions apply unless a newer decision explicitly supersedes them

No assumptions beyond what is written here.

Files under ai/examples/ are illustrative only and must not be treated as active constraints.
