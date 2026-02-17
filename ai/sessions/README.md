# Session Records

Store one file per Codex session:

- Path pattern: `ai/sessions/YYYY-MM-DD-<session_id>.md`
- Use `unknown` until the session ID is available.
- Keep records append-only (do not overwrite prior sessions).

Current session rule:
- Treat the newest file with `status` not equal to `done` as current.

