# Run Code Review with Security

Step-by-step execution for the three-axis review.

## Preconditions

- There is a diff to review (either uncommitted changes or a branch against a base)
- The codebase has been built/compiled successfully
- Tests are passing (or at least not broken worse than before)

## 1. Invoke Standards review

Call the `/code-review` skill's Standards axis.

This checks:
- Code style and formatting
- Naming conventions
- Documentation completeness
- Test coverage
- Code smells and anti-patterns

Collect the findings.

## 2. Invoke Spec review

Call the `/code-review` skill's Spec axis.

This checks:
- Does the implementation match the ticket description?
- Are all acceptance criteria met?
- Are there any out-of-scope changes?

Collect the findings.

## 3. Invoke Security review

Call `/security-audit --diff --critical`.

This checks:
- Injection vulnerabilities (SQL, XSS, command injection)
- Authentication/authorization bypasses
- Hardcoded secrets
- Path traversal
- SSRF
- Insecure cryptography
- Missing security controls

Collect the findings (CRITICAL only in diff mode).

## 4. Combine results

Merge the three axes into one unified report.

### Structure:

```markdown
# Code Review Results

**Branch/Commit:** [branch name or commit SHA]  
**Files changed:** [count]  
**Lines changed:** +[additions] -[deletions]

---

## Standards

[Findings from axis 1]

**Summary:** [X issues] - [list severity breakdown]

---

## Spec

[Findings from axis 2]

**Summary:** [Y issues] - [what's missing vs. ticket requirements]

---

## Security

[Findings from axis 3]

**Summary:** [Z CRITICAL findings]

---

## Overall Assessment

**Ready to merge:** [Yes/No]

**Blockers (must fix before merge):**
- [List all CRITICAL security findings]
- [List any Standards violations that break builds]
- [List any missing acceptance criteria from Spec]

**Recommended (fix soon):**
- [List HIGH severity issues]
- [List code quality issues]

**Nice to have:**
- [List minor improvements]
```

## 5. Determine merge readiness

**Block merge if:**
- Any CRITICAL security findings
- Tests are failing
- Spec acceptance criteria are not met
- Build is broken

**Allow merge with follow-up if:**
- HIGH security findings (but create issues for them)
- Code quality issues (refactoring opportunities)
- Missing tests for non-critical paths

**Green light if:**
- No CRITICAL or HIGH findings
- All acceptance criteria met
- Tests passing
- Standards followed

## 6. Output recommendations

End with clear next steps:

```markdown
## Next Steps

1. [Fix SQL injection in tickets.go:313]
2. [Add test coverage for error paths]
3. [Update ticket to mark acceptance criteria complete]
4. [Ready to merge after fixes]
```

## 7. (Optional) Create follow-up issues

If there are HIGH findings that don't block merge but should be tracked:

- Create GitHub issues for each
- Label them appropriately: `security`, `tech-debt`, `quality`
- Reference the PR/branch in the issue
