# Security Policy

This file defines how the `/security-audit` wrapper operates. 

**Important:** All OWASP guidelines, vulnerability detection, and security standards are defined by the `code-security-skills@agent-security-playbook` plugin. This skill just wraps the plugin to:
1. Determine what to scan
2. Filter output by severity
3. Format results consistently

## Plugin skills available

The `code-security-skills@agent-security-playbook` plugin provides these skills:

- `code-security-skills:code-review-security` - General code security review (OWASP Top 10)
- `code-security-skills:api-security-review` - API-specific security (OWASP API Top 10)
- `code-security-skills:sca-audit` - Dependency/supply chain audit (CVE scanning)
- `code-security-skills:secrets-scan` - Detect hardcoded secrets
- `code-security-skills:mobile-code-review` - Mobile app security (OWASP MASVS)
- `code-security-skills:iac-security-review` - Infrastructure as Code security

## Severity filtering

The plugin assigns severity (CRITICAL, HIGH, MEDIUM, LOW). This wrapper filters them:

- **Full scan**: Report CRITICAL, HIGH, and MEDIUM
- **Diff scan (`--diff`)**: Report CRITICAL and HIGH only (faster, focus on blockers)
- **Critical-only (`--critical`)**: Report CRITICAL only

## Auto-issue creation

When `--create-issues` is passed:
- Create a GitHub issue for each **CRITICAL** finding
- Label: `security`, `needs-triage`
- Assign: none (let `/triage` handle assignment)
- Body: Use the template from `plays/issue-template.md`

## Default plugin skill

For general use (PR review, code audit), invoke:
- **`code-security-skills:code-review-security`**

This is the comprehensive code-level review covering OWASP Top 10.
