# Roadmap

## Overview

At Ratify, our mission is to safeguard the container supply chain by ratifying trustworthy and compliant artifacts. We achieve this through a robust and pluggable verification engine that includes built-in verifiers. These verifiers can be customized to validate supply chain metadata associated with artifacts, covering essential aspects such as signatures and attestations (including vulnerability reports, SBOM, provenance data, and VEX documents). As the landscape of supply chain security evolves, we actively develop new verifiers, which can be seamlessly integrated into our verification engine. Additionally, if you have a specific use case, you can create your own verifier following our comprehensive guidance. Each verifier will generate detailed verfication reports, which can be consumed by various policy controllers to enforce policies.

Ratify serves several major scenarios. It seamlessly integrates with OPA Gatekeeper as the Kubernetes policy controller. By doing so, it safeguards your Kubernetes cluster against untrustworthy and non-compliant container images. Ratify functions as an external data provider for Gatekeeper, supplying artifact verification results that align with defined policies. Furthermore, Ratify can be employed at container runtime via the plugin interface, enabling policy evaluation based on its artifact verification outcomes. This enables more granular artifact verification at the Kubernetes node level. Last but not least, Ratify can be integrated into your CI/CD pipeline to ensure the trustworthiness and compliance of container images before using them. 

This document presents the roadmap of Ratify that translates our strategy into practical steps.

## Milestones

The Ratify roadmap is divided into milestones, each with a set of features (high level) and timeline. The milestones marked as `Tentative` are subject to change based on the project’s priorities and the community’s feedback. We will prioritize releases for security or urgent fixes, so the roadmap may be adjusted and new features may be postponed to the next milestone. Any dates and features listed below in a given milestone are subject to change. See the [GitHub milestones](https://github.com/deislabs/ratify/milestones?state=open) for the most up-to-date issues and their status. We are targeting to release a new Ratify version every 3 or 4 months.

### v1.0

**Status**: Completed

**Released date**: Sep 27, 2023

**Release link**: [v1.0.0 Release Notes](https://github.com/deislabs/ratify/releases/tag/v1.0.0)

**Major features**

- Ratify as an external Data Provider for Gatekeeper
- Plugin framework for extensibility
- Policies for Notary Project signatures verification at admission control in kubernetes
- Policies for Cosign keyless verification at admission control in kubernetes
- High Availability support in Kubernetes (Experimental)

### v1.1

**Status**: Completed

**Release date**: Dec 12, 2023

**Release link**: [v1.1.0 Release Notes](https://github.com/deislabs/ratify/releases/tag/v1.1.0)

**Major features**

- Policies for assessing vulnerability reports at admission control in kubernetes
- Policies for assessing software license at admission control in kubernetes
- New diagnostic logs

### v1.2

**Status**: In progress

**Target date**: Apr 30, 2024

**major features**

- Kubernetes multi-tenancy support
- OCI v1.1 compliance
- Cosign signatures verification using keys in AKV
- Error logs improvements

### v1.3

**Status**: Not started

**Target date**: Jun 30, 2024

**Major features**

- Cosign keyless verification using OIDC settings
- Notary Project signature verification with Time-stamping support
- Signing Certificate/key rotation support

### v1.4

**Status**: Tentative

**Target date**: Sep 30, 2024

**Major features**

- Attestations support
- Use Ratify at container runtime (Preview)

### v2.0

Status: Tentative

Target date: TBD

**Major features**

- Use Ratify in CI/CD pipelines (Preview)
- Support CEL as additional policy language