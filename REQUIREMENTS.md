This Go application should produce a binary which can be leveraged to authenticate to Salesforce.com, and return a JSON with an access token, refresh token, and instance URL. It should prompt for client ID and client secret, and listen on localhost for the callback from OAuth2. The CLI should be well-structured, maintainable, and have the right ability to accept flags and arguments. Automated tests should be included to ensure the application is robust and reliable. The application should be built via github actions automatic releases which are stored as artifacts for download. A comprehensive README.md should be created to cover this application.

Requirements:

- REQ-1: prompts for client ID and client secret
- REQ-2: listens on localhost for the callback from OAuth2
- REQ-3: returns a JSON with an access token, refresh token, and instance URL
- REQ-4: leverage a stable CLI framework to ensure that the CLI is well-structured, maintainable, and has the right ability to accept flags and arguments
- REQ-5: include automated tests to ensure the application is robust and reliable
- REQ-6: all code must pass linting validation (golangci-lint) with zero errors before considering any development work complete
- REQ-7: allow specifying of a custom salesforce domain (i.e. company.my.salesforce.com) to use for login instead of login.salesforce.com

Build Automation:

- BUILD-1: include automation to build via github actions automatic releases which are stored as artifacts for download
- BUILD-2: include a Makefile to build the application locally
- BUILD-3: release this as a zipped binary for Windows, Mac, and Linux.

Quality Assurance:

- QA-1: all code changes must pass golangci-lint validation with zero errors
- QA-2: linting must be integrated into CI/CD pipeline to prevent merging of non-compliant code
- QA-3: development workflow must include linting as a mandatory step before considering any work complete

Docs:

- DOCS-1: Create a comprehensive README.md to cover this application.
