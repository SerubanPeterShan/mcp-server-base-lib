# Security Policy

## Supported Versions

We currently support the following versions with security updates:

| Version | Supported          |
| ------- | ------------------ |
| 1.x.x   | :white_check_mark: |

## Reporting a Vulnerability

We take the security of the MCP Server Library seriously. If you believe you have found a security vulnerability, please follow these steps:

1. **DO NOT** disclose the vulnerability publicly
2. **DO NOT** open a public GitHub issue
3. Email your findings to security@serubanpetershan.com with the following information:
   - Description of the vulnerability
   - Steps to reproduce
   - Potential impact
   - Any suggested fixes (if available)

## Security Best Practices

When using this library, please follow these security best practices:

1. **WebSocket Security**:
   - Always implement proper origin checking in production
   - Use WSS (WebSocket Secure) in production environments
   - Implement proper authentication mechanisms

2. **Message Handling**:
   - Validate all incoming messages
   - Implement rate limiting
   - Sanitize message payloads

3. **Server Configuration**:
   - Use secure logging practices
   - Implement proper error handling
   - Configure appropriate timeouts

4. **Dependencies**:
   - Regularly update dependencies
   - Monitor for known vulnerabilities
   - Use dependency scanning tools

## Security Updates

Security updates will be released as patch versions (e.g., 1.0.1, 1.0.2) and will be clearly marked in the release notes. We recommend always using the latest patch version of the library.

## Security Contact

For security-related questions or concerns, please contact:
- Email: security@serubanpetershan.com
- PGP Key: [Your PGP Key if available]

## Responsible Disclosure

We follow a responsible disclosure policy:
1. We will acknowledge receipt of your vulnerability report
2. We will investigate and verify the vulnerability
3. We will work on a fix
4. We will notify you when the fix is ready
5. We will credit you in the security advisory (unless you prefer to remain anonymous)

## Security Features

The MCP Server Library includes several built-in security features:
- Concurrent client handling with proper synchronization
- Structured logging for security auditing
- Health check endpoint for monitoring
- Configurable message handlers for custom security implementations 