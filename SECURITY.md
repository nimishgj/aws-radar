# Security Policy

## Supported Versions

AWS Radar is under active development. Security fixes are applied to the
latest release and the `main` branch. We recommend always running the most
recent release.

| Version        | Supported          |
| -------------- | ------------------ |
| Latest release | :white_check_mark: |
| Older releases | :x:                |

## Reporting a Vulnerability

Please **do not report security vulnerabilities through public GitHub issues,
discussions, or pull requests.**

Instead, report them privately using GitHub's built-in
[private vulnerability reporting](https://github.com/nimishgj/aws-radar/security/advisories/new):

1. Go to the **Security** tab of this repository.
2. Click **Report a vulnerability**.
3. Provide as much detail as possible: affected version, a description of the
   issue, steps to reproduce, and potential impact.

You can expect an initial acknowledgement within a few business days. Once the
report is confirmed, we will work on a fix and coordinate a disclosure timeline
with you.

## Scope

AWS Radar only makes **read-only** AWS API calls to count resources and exposes
the results as Prometheus metrics. When reporting, issues that are especially
relevant include:

- Handling of AWS credentials or the metrics endpoint.
- Exposure of sensitive account information in metrics or logs.
- Vulnerabilities in the collection or HTTP-serving code paths.

Thank you for helping keep AWS Radar and its users safe.
