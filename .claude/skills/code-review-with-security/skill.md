---
name: code-review-with-security
description: Run code review (Standards + Spec) plus security audit (three-axis review)
tags: [review, security, quality]
---

# Code Review with Security

Runs the standard `/code-review` (Standards + Spec axes), then adds a **Security axis** by invoking `/security-audit --diff --critical`.

This is a **wrapper** around the Matt Pocock `/code-review` skill - it doesn't modify the original, just extends it with security.

## Usage

```bash
/code-review-with-security
```

Use this exactly where you'd normally use `/code-review`:
- After `/implement` finishes building a feature
- Standalone to review a branch or PR
- Before merging to main

## What it does

1. **Invoke `/code-review`** (unchanged, from the Matt Pocock skills)
   - Standards axis: Does the code follow repo coding standards?
   - Spec axis: Does the code match what the ticket asked for?

2. **Invoke `/security-audit --diff --critical`** (which calls the `code-security-skills@agent-security-playbook` plugin)
   - Security axis: Are there CRITICAL security vulnerabilities per OWASP guidelines?

3. **Combine results** into a unified three-axis review

**Note:** The OWASP security guidelines and vulnerability checks come from the installed `code-security-skills@agent-security-playbook` plugin, not defined here.

## Output format

```markdown
# Code Review Results

## Standards
[Output from /code-review Standards axis]

## Spec
[Output from /code-review Spec axis]

## Security
[Output from /security-audit]

---

## Summary

- **Standards:** [X issues found]
- **Spec:** [Y issues found]
- **Security:** [Z CRITICAL findings]

**Ready to merge:** [Yes/No]  
**Blockers:** [List of must-fix items before merge]
```

## When to use vs. `/code-review`

- **Use this** when reviewing code that touches:
  - Authentication/authorization
  - Data handling (user input, file uploads, database queries)
  - API endpoints
  - Cryptography or password handling
  - External integrations

- **Use plain `/code-review`** when reviewing:
  - Documentation changes
  - UI-only changes (no backend logic)
  - Minor refactors with no security surface

## Execution

See `plays/run-review-with-security.md` for the step-by-step procedure.

## Configuration

To make `/implement` use this by default in your project, add to `.claude/settings.json`:

```json
{
  "skills": {
    "implement": {
      "review_skill": "code-review-with-security"
    }
  }
}
```

(Note: If your harness doesn't support this config yet, invoke this skill manually after `/implement` finishes)
