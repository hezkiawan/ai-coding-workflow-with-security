# Security Issue Template

Use this template when auto-creating GitHub issues for CRITICAL findings.

---

## 🔴 [CRITICAL] [Vulnerability Name]

**CWE:** [CWE-ID]  
**OWASP Category:** [OWASP Top 10 reference]  
**Discovered by:** Automated security audit  
**Date:** [YYYY-MM-DD]

### Location

**File:** `path/to/file.go`  
**Line:** [line number]

### Vulnerability Description

[Brief explanation of what the vulnerability is]

### Attack Scenario

[Concrete example of how an attacker could exploit this vulnerability]

Example payload:
```
[Example malicious input or request]
```

### Impact

- **Confidentiality:** [High/Medium/Low]
- **Integrity:** [High/Medium/Low]
- **Availability:** [High/Medium/Low]

### Affected Code

```go
// Current vulnerable code
[code snippet showing the vulnerability]
```

### Recommended Fix

```go
// Secure implementation
[code snippet showing the fix]
```

### References

- [CWE link]
- [OWASP Cheat Sheet link]
- [Any relevant documentation]

### Acceptance Criteria

- [ ] Vulnerability is fixed with recommended approach
- [ ] Unit test added to prevent regression
- [ ] Security audit re-run shows finding resolved
- [ ] Code review completed

---

**Priority:** P0 (Critical - Fix immediately)  
**Estimated effort:** [TBD after investigation]
