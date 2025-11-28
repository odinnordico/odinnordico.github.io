# Security Policy

## Supported Versions

We release patches for security vulnerabilities. The following versions are currently being supported with security updates:

| Version | Supported          |
| ------- | ------------------ |
| v0.1.x   | :white_check_mark: |
| 0.1.x   | :x: |

## Reporting a Vulnerability

We take the security of odinnordico.github.io seriously. If you believe you have found a security vulnerability, please report it to us as described below.

### How to Report

**Please do not report security vulnerabilities through public GitHub issues.**

Instead, please report them via email to:

- 📧 **Email**: [odin.nordico90@gmail.com](mailto:odin.nordico90@gmail.com)
- **Subject**: `[SECURITY] Brief description of the issue`

Please include the following information in your report:

- Type of issue (e.g., buffer overflow, SQL injection, cross-site scripting, etc.)
- Full paths of source file(s) related to the manifestation of the issue
- The location of the affected source code (tag/branch/commit or direct URL)
- Any special configuration required to reproduce the issue
- Step-by-step instructions to reproduce the issue
- Proof-of-concept or exploit code (if possible)
- Impact of the issue, including how an attacker might exploit it

### What to Expect

After submitting a vulnerability report, you should receive:

1. **Acknowledgment**: Within 48 hours of submission
2. **Initial Assessment**: Within 5 business days, we'll provide an initial assessment
3. **Updates**: Regular updates on our progress as we investigate and address the issue
4. **Resolution**: A timeline for when you can expect a fix to be released

### Disclosure Policy

- We ask that you give us reasonable time to investigate and fix the issue before public disclosure
- We will keep you informed of our progress
- Once the issue is resolved, we will publicly acknowledge your responsible disclosure (unless you prefer to remain anonymous)

## Security Best Practices

When using this project, we recommend following these security best practices:

### For Users

1. **Keep Dependencies Updated**: Regularly update Go and project dependencies
   ```bash
   go get -u ./...
   go mod tidy
   ```

2. **Validate Data Files**: Ensure YAML data files come from trusted sources
   - Be cautious when accepting YAML files from external sources
   - Review YAML content before processing

3. **Environment Security**: When running the development server:
   - Only bind to `localhost` unless specifically needed
   - Use `--host 0.0.0.0` only in trusted network environments
   - Don't expose the development server to the public internet

4. **Asset Security**: 
   - Verify all assets (images, logos) are from trusted sources
   - Scan uploaded or user-provided files for malware

### For Contributors

1. **Input Validation**: Always validate and sanitize user inputs
2. **Dependency Management**: 
   - Review dependencies before adding them
   - Keep dependencies up to date
   - Use Go's vulnerability scanner: `go run golang.org/x/vuln/cmd/govulncheck@latest ./...`

3. **Code Review**: 
   - All code changes should be reviewed by at least one other maintainer
   - Security-sensitive changes require additional scrutiny

4. **Testing**: 
   - Write tests for security-critical functionality
   - Include negative test cases for validation logic

## Known Security Considerations

### YAML Processing
- This project processes YAML files which could potentially contain malicious content
- Always review YAML files from untrusted sources before processing
- The `yaml.v3` library is used with default safe parsing

### File System Access
- The application reads from and writes to the file system
- Ensure appropriate file permissions on data directories
- Be cautious with the `--data-dir` and `--output-dir` flags

### Development Server
- The development server (`serve` command) is intended for local development only
- Do not use the development server in production environments
- Production deployments should use proper web servers (nginx, Apache, etc.)

### Template Execution
- Templates are executed with access to resume data
- Custom templates should be reviewed for potential injection vulnerabilities
- Only use templates from trusted sources

## Security Updates

Security updates will be announced through:

1. **GitHub Security Advisories**: [Security Advisories](https://github.com/odinnordico/odinnordico.github.io/security/advisories)
2. **Release Notes**: Check our [Releases](https://github.com/odinnordico/odinnordico.github.io/releases) page
3. **README Updates**: Critical security information will also be added to the README

## Vulnerability Scanning

We use the following tools to scan for vulnerabilities:

- **Go Vulnerability Database**: `govulncheck`
- **Dependency Scanning**: GitHub Dependabot
- **Code Analysis**: GitHub CodeQL (if configured)

To check for vulnerabilities locally:

```bash
# Install govulncheck
go install golang.org/x/vuln/cmd/govulncheck@latest

# Run vulnerability check
govulncheck ./...
```

## Contact

For security-related questions or concerns that are not vulnerabilities, you can:

- Open a [GitHub Discussion](https://github.com/odinnordico/odinnordico.github.io/discussions)
- Email: [odin.nordico90@gmail.com](mailto:odin.nordico90@gmail.com)

---

**Thank you for helping keep odinnordico.github.io and its users safe!**
