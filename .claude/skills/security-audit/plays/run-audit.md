# Run Security Audit

This is the execution play. Follow these steps in order.

## 1. Parse arguments

Check for flags:
- `--diff`: Scan only changed files (faster, less context)
- `--critical`: Report CRITICAL findings only
- `--create-issues`: Auto-create GitHub issues for CRITICAL findings

## 2. Determine scope

### If `--diff` flag is present:
1. Get the list of changed files: `git diff --name-only HEAD`
2. Use this as the target scope

### Otherwise (full scan):
1. Scan the entire codebase (root directory)
2. Let the plugin handle exclusions (it already knows to skip `node_modules/`, `.git/`, etc.)

## 3. Detect which plugin skill to invoke

**The `code-security-skills@agent-security-playbook` plugin has multiple specialized skills. Choose the right one based on what you're reviewing:**

| If reviewing... | Invoke this skill | It will check... |
|-----------------|-------------------|------------------|
| **General codebase/PR** | `code-security-skills:code-review-security` | All OWASP Top 10 vulnerabilities at code level |
| **API endpoints** | `code-security-skills:api-security-review` | OWASP API Security Top 10 |
| **Dependencies** | `code-security-skills:sca-audit` | Known CVEs in dependencies |
| **Mobile code (iOS/Android)** | `code-security-skills:mobile-code-review` | OWASP MASVS mobile vulnerabilities |
| **Infrastructure (Terraform/K8s)** | `code-security-skills:iac-security-review` | IaC misconfigurations |
| **Looking for exposed secrets** | `code-security-skills:secrets-scan` | Hardcoded credentials, API keys |

**For most use cases (PR review, general audit), use `code-security-skills:code-review-security`** - it's the comprehensive code-level review.

## 4. Invoke the plugin skill

Call the chosen skill with the scope as an argument:

```bash
# Full scan
Skill: code-security-skills:code-review-security
Args: ./

# Diff scan (pass the directory containing changed files)
Skill: code-security-skills:code-review-security
Args: ./backend
```

**The plugin handles everything:**
- Language/framework detection
- OWASP guideline application
- Vulnerability detection
- Severity assignment
- Fix recommendations

Your job is just to invoke it and process the output.

## 5. Parse and filter the plugin output

The plugin returns findings with severity levels. Filter them based on flags:

- **If `--critical` flag**: Keep only CRITICAL findings
- **If `--diff` flag**: Keep CRITICAL and HIGH (diff scans should be fast, focus on blockers)
- **Otherwise (full scan)**: Keep CRITICAL, HIGH, and MEDIUM

## 6. Format output

Take the plugin's findings and format them consistently:

```markdown
# Security Audit Results

**Scope:** [full codebase | current diff]  
**Files scanned:** [count]  
**CRITICAL:** [count] findings  
**HIGH:** [count] findings  
**MEDIUM:** [count] findings (if applicable)

---

## Critical Findings

### 1. [Vulnerability Name] (CRITICAL)

**File:** `path/to/file.go:123`  
**CWE:** CWE-89 (SQL Injection)  
**OWASP:** A03:2021 - Injection

**Description:**  
[Brief explanation of the vulnerability]

**Attack Scenario:**  
[How an attacker could exploit this]

**Fix:**
```go
// Vulnerable code
query := fmt.Sprintf("SELECT * FROM users WHERE id = %s", userInput)

// Fixed code
query := "SELECT * FROM users WHERE id = ?"
db.Query(query, userInput)
```

**Confidence:** High

---

[Repeat for each finding]

## Summary

Total findings: [count]  
- **Must fix now:** [CRITICAL count]  
- **Should fix soon:** [HIGH count]  
- **Consider fixing:** [MEDIUM count]
```

## 7. (Optional) Create issues

If `--create-issues` was passed:

For each **CRITICAL** finding:
1. Check if an issue already exists (search by title pattern)
2. If not, create a GitHub issue using `gh issue create`:
   - Title: `[Security] [CRITICAL] [Vulnerability Name]`
   - Body: Use the template from `plays/issue-template.md`
   - Labels: `security`, `needs-triage`
3. Print the issue URL

## 8. Exit code

Return exit code based on findings:
- 0: No CRITICAL or HIGH findings
- 1: CRITICAL findings present
- 2: Only HIGH findings present
