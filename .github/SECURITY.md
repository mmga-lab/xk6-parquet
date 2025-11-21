# Security Policy

## Supported Versions

We release patches for security vulnerabilities. The following versions are currently supported:

| Version | Supported          |
| ------- | ------------------ |
| latest  | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

We take the security of xk6-parquet seriously. If you believe you have found a security vulnerability, please report it to us as described below.

### Please Do Not

- **Do not** open a public issue on GitHub
- **Do not** share the vulnerability publicly until we've had a chance to address it

### How to Report

Please report security vulnerabilities by emailing the maintainers directly or using GitHub's private vulnerability reporting feature:

1. **GitHub Private Vulnerability Reporting** (Recommended)
   - Go to the [Security tab](https://github.com/mmga-lab/xk6-parquet/security)
   - Click "Report a vulnerability"
   - Fill out the form with details about the vulnerability

2. **Email**
   - Send an email to the maintainers listed in the CODEOWNERS file
   - Include as much information as possible about the vulnerability

### What to Include

Please include the following information in your report:

- Type of vulnerability (e.g., buffer overflow, SQL injection, cross-site scripting, etc.)
- Full paths of source file(s) related to the vulnerability
- Location of the affected source code (tag/branch/commit or direct URL)
- Any special configuration required to reproduce the issue
- Step-by-step instructions to reproduce the issue
- Proof-of-concept or exploit code (if possible)
- Impact of the issue, including how an attacker might exploit it

### What to Expect

- We will acknowledge receipt of your vulnerability report within 48 hours
- We will provide a more detailed response within 7 days indicating the next steps
- We will keep you informed of the progress towards a fix and announcement
- We may ask for additional information or guidance

### Disclosure Policy

- We will coordinate with you on the timing of the public disclosure
- We prefer to fully fix vulnerabilities before publicly disclosing them
- We will credit you in the security advisory unless you prefer to remain anonymous

## Security Best Practices

When using xk6-parquet:

1. **Keep Dependencies Updated**
   - Regularly update to the latest version of xk6-parquet
   - Keep k6 and Go updated to the latest stable versions
   - Monitor security advisories for dependencies

2. **Parquet File Validation**
   - Validate Parquet files from untrusted sources before processing
   - Be cautious with user-supplied file paths
   - Consider file size limits to prevent resource exhaustion

3. **Resource Limits**
   - Set appropriate memory limits when reading large Parquet files
   - Use chunked reading for large datasets
   - Monitor resource usage in production

4. **Input Validation**
   - Validate all user inputs in your k6 scripts
   - Sanitize file paths to prevent directory traversal attacks
   - Use absolute paths when possible

## Known Security Considerations

### File Path Handling

xk6-parquet accepts file paths as input. Ensure that:
- File paths are validated and sanitized
- Users cannot supply arbitrary paths that could access sensitive files
- Use absolute paths or restrict access to specific directories

### Memory Usage

Large Parquet files can consume significant memory. Consider:
- Using chunked reading for large files
- Setting appropriate memory limits
- Monitoring memory usage in production environments

### Dependency Vulnerabilities

We regularly monitor and update dependencies. Known security issues in dependencies are tracked and addressed promptly.

## Security Updates

Security updates will be released as:
- Patch releases for supported versions
- Security advisories on GitHub
- Announcements in the README and release notes

## Contact

For any questions about security, please contact the maintainers through the channels mentioned above.
