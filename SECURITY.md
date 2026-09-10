# Security Policy

## Reporting Security Vulnerabilities

If you discover a security vulnerability in Ngumpul Host, please send an email to security@ngumpul.example (or open a private security advisory on GitHub). Please do not disclose vulnerabilities publicly until they have been addressed.

Please provide:
- A description of the vulnerability.
- Steps to reproduce or proof-of-concept.
- Potential impact and affected versions.

We appreciate your assistance in making Ngumpul Host safe for everyone.

---

## Architectural Security Principles

Ngumpul Host adheres to strict architectural boundaries:

1. **No Arbitrary Remote Execution:**
   Ngumpul Host intentionally does not provide arbitrary command execution, shell access, VPS provisioning, or automatic Docker execution from user requests. Hosting requests are metadata-only. An administrator handles deployment externally and records public endpoints.

2. **No User-Controlled Reverse Proxying:**
   The application does not allow users to specify arbitrary upstream proxies. Nginx routes only explicitly configured upstreams. Unmatched requests fallback to a designated destination (`http://1111:80`).

3. **Server-Side Authorization Boundary:**
   Every private API endpoint enforces object-level and role-based checks (`USER`, `ADMIN`). The frontend is never treated as the security boundary.

4. **Credential Isolation:**
   Infrastructure secrets (DNS tokens, SSH keys, server credentials) must never be stored inside Ngumpul Host database records or returned in public API responses.

5. **Safe Authentication:**
   Passwords are encrypted with Argon2id using unique cryptographic salts. Sessions are stored server-side in PostgreSQL and transmitted via `HttpOnly`, `SameSite=Lax` cookies.
