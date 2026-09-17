---
name: security-audit
description: Run OWASP security analysis on code and generate prioritized findings
tags: [security, review]
---

# Security Audit

Runs the OWASP security skills suite against the codebase (or a diff) and produces prioritized findings.

## Usage

```bash
/security-audit              # Full codebase scan
/security-audit --diff       # Only scan the current diff (faster)
/security-audit --critical   # Only report CRITICAL findings
```

## What it does

This is a **thin wrapper** around the `code-security-skills@agent-security-playbook` plugin. It:

1. **Determines scope**: Full codebase or just the diff (`--diff` flag)
2. **Invokes the plugin**: Calls `code-security-skills:code-review-security` (or other specialized skills)
3. **Filters output**: By severity (CRITICAL/HIGH/MEDIUM) based on flags
4. **Formats results**: Consistent structure with file:line references
5. **Optional: Creates issues**: GitHub issues for CRITICAL findings (`--create-issues`)

**The plugin does all the security work** - OWASP guidelines, vulnerability detection, CWE mapping, fix suggestions. This skill just orchestrates it.

## Integration points

- **After `/implement`**: Run this before committing to catch issues early
- **In `/code-review-with-security`**: The wrapper invokes this automatically as a third axis
- **Standalone**: Run anytime you want to audit the security posture

## Execution

Follow the procedure in `plays/run-audit.md`.

## Customization

Edit `plays/security-policy.md` to configure:
- Severity thresholds
- Which file patterns to always scan (auth, crypto, file I/O)
- Auto-issue-creation rules
