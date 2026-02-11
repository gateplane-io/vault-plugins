# GatePlane Vault/OpenBao Plugins
<!-- Badges -->
[![License: ElasticV2](https://img.shields.io/badge/ElasticV2-green?label=license&cacheSeconds=3600&link=https%3A%2F%2Fwww.elastic.co%2Flicensing%2Felastic-license)](https://www.elastic.co/licensing/elastic-license)
[![Test Plugins](https://github.com/gateplane-io/vault-plugins/actions/workflows/test.yaml/badge.svg)](https://github.com/gateplane-io/vault-plugins/actions/workflows/test.yaml)
[![GoReport Card](https://img.shields.io/badge/go%20report-A+-brightgreen.svg?style=flat)](https://goreportcard.com/report/github.com/gateplane-io/vault-plugins)
[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/gateplane-io/vault-plugins/badge)](https://scorecard.dev/viewer/?uri=github.com/gateplane-io/vault-plugins)

<h3 align="center">
An <i>Approval-Based</i>,<br/>
<i>Just-In-Time</i> (JIT),<br/>
<i>Privileged Access Management</i> (PAM)<br/>
layer for <i>Vault/OpenBao</i>
</h3>

<img src="https://raw.githubusercontent.com/gateplane-io/gateplane-io.github.io/main/assets/logo-trans-light.svg#gh-dark-mode-only" alt="logo-darkmode" width="45%" align="right">
<img src="https://raw.githubusercontent.com/gateplane-io/gateplane-io.github.io/main/assets/logo-trans-dark.svg#gh-light-mode-only" alt="logo-lightmode" width="45%" align="right">

<!-- Handled by md_toc hook -->
<!--TOC-->

- [1. 💡 Overview](#1--overview)
  - [❓ What is GatePlane?](#-what-is-gateplane)
  - [🎯 Key Features](#-key-features)
  - [❔ Why GatePlane](#-why-gateplane)
  - [🔍 This Repository](#-this-repository)
- [2. ✨ Features](#2--features)
  - [🧩 Plugins in Detail](#-plugins-in-detail)
    - [Policy Gate](#policy-gate)
    - [Okta Group Gate](#okta-group-gate)
  - [📋 Features in Detail](#-features-in-detail)
    - [Community Edition](#community-edition)
    - [Team and Enterprise Edition](#team-and-enterprise-edition)
- [3. 💻 Installation / Getting Started](#3--installation--getting-started)
  - [🧰 Prerequisites](#-prerequisites)
  - [🚀 How to Install](#-how-to-install)
    - [Downloading and Verifying the Plugins](#downloading-and-verifying-the-plugins)
    - [Registering to Plugin Catalog - Manually](#registering-to-plugin-catalog---manually)
    - [Registering to Plugin Catalog - With GatePlane Terraform modules](#registering-to-plugin-catalog---with-gateplane-terraform-modules)
    - [Enabling a Gate - Manually](#enabling-a-gate---manually)
    - [Enabling a Gate - With GatePlane Terraform modules](#enabling-a-gate---with-gateplane-terraform-modules)
  - [▶️ Usage Example](#-usage-example)
    - [Requesting Access](#requesting-access)
    - [Approving Access](#approving-access)
    - [Claiming Access](#claiming-access)
  - [🛠️ How to Build and Test](#-how-to-build-and-test)
    - [Building](#building)
    - [Testing](#testing)
- [4. 💬 Contact](#4--contact)
  - [Community](#community)
  - [⚖️ License](#-license)

<!--TOC-->


---


## 1. 💡 Overview

### ❓ What is GatePlane?
GatePlane is a project that makes time-based, conditional access management approachable and transparent, *made by Security Professionals for Security Professionals* .

Gateplane implements *Privileged Access Management* (PAM) and *Just-In-Time* (JIT) access, helping tech groups and companies to give, revoke and monitor permissions across the whole organization.

The ultimate goal is to solve problems people in tech organizations stumble upon frequently, such as:
* **Developer**: "*I need to debug production but I don't have access*"
* **IT**: “*I get around 60 requests for sensitive access per day, and I have to manually set them up and tear them down*”
* **Security Officer**: “*I have no idea who has access to what after the last incident and the auditor comes next week*”
* **Security Engineer**: “*Setting up this PAM solution will take forever and will change everything we know about our infrastructure*”

It does so by using the Open-Source and battle-tested tools [Vault](https://developer.hashicorp.com/vault/docs) or [OpenBao](https://openbao.org/) as its underlying infrastructure,
drawing from their *Authentication*, *Authorization*, *Integrations* and *Auditing* principles, avoiding to re-invent the wheel.

*GatePlane Community* comes **free of charge** and available for everyone to *use*, *audit* and *contribute* to, under the [*Elastic License v2*](./LICENSE).

##### More about GatePlane can be found in its [Website](https://gateplane.io)

### 🎯 Key Features

* **Approval-Based Access**: Access is granted only if someone else approves of it
* **Access that expires**: Every granted access has an expiration date, no exceptions
* ***Vast* support of Integrations**: Supports *Kubernetes*, *AWS/GCP/Azure*, *SSH*, *Databases* (from *PostgreSQL* to *Elasticsearch*) and the [whole suite of *Vault/OpenBao Secrets Engines*](https://developer.hashicorp.com/vault/integrations?components=database%2Csecrets-engine), as well as Just-In-Time *Okta Group Membership Management* drawing from the [150+ Okta Integrations](https://www.okta.com/integrations/)
* [**WebUI to Request, Approve and Claim access**](https://app.gateplane.io) : All available accesses are listed and can be managed through a Web Application tied with your Vault/OpenBao instance
* **Notifications**: When a user requests, approves or claims access, a Slack/Discord/MSTeams/*you-name-it* message is sent
* **Metrics**: Measurements on the friction relief, numbers of elevated accesses needed and statistics on mean-time-to-claim

### ❔ Why GatePlane

It makes *Just-In-Time* (JIT) *Privileged Access Management* (PAM) accessible to everyone through Open-Source tools and auditable code.

Its *Zero-Knowledge* and *Self-Hostable* natute ensures that access is neither controlled, nor proxied by GatePlane, **keeping your organization's accesses and secrets strictly inside your infrastructure**.

Our rationale is that trust on the market of PAM solutions is never really gained if transparency is not part of the equation. Also buzzword marketing, and vague promises that often lack technical ground are not in our line of work.

Our mission is enabling everyone to provide these guarantees to their organization without closed-source software, unclear security requirements, opaque components in their threat model and tying their sensitive access to SaaS vendor-locks they do not control. All this with direct and honest communication.

### 🔍 This Repository
This repository contains a series of *Vault/OpenBao Plugins* (see for: [Vault](https://developer.hashicorp.com/vault/docs/plugins)/[OpenBao](https://openbao.org/docs/plugins/)) that enable Vault/OpenBao to act as an Access Control Plane, providing *conditional access* to resources.

The plugins currently included:
* **Policy Gate**: Controls Access to Vault/OpenBao and its Integrations
* **Okta Group Gate**: Controls Access to Okta Groups

*... more plugins will come in the future ...*

## 2. ✨ Features

### 🧩 Plugins in Detail

Each plugin is considered a *Gate*, that needs conditions in order to provide time-limited access to the requesting users. *Conditions* can be approvals by a number of different users and providing a reason for requesting access. Time-To-Live (TTL) is enforced depending on the type of access provided.

#### Policy Gate
The *Policy Gate* plugin can utilize all Vault/OpenBao *Secrets Engines* (see for: [Vault](https://developer.hashicorp.com/vault/docs/secrets)/[OpenBao](https://openbao.org/docs/secrets/)), providing a conditional access layer to all of them.

With that, *all AWS, Kubernetes, SSH or Database access configured through Vault/OpenBao can now
be provided through a request/approval flow* and also be expiring through configured TTLs.

#### Okta Group Gate
The *Okta Group Gate* leverages the request/approval flow to gain expiring access to Okta Groups,
for users that access Vault/OpenBao through Okta.

The plugin reaches the Okta API with API credentials stored in Vault/OpenBao and
adds the user requesting access to a configured Okta Group for a set time duration.

##### API documentation for each of the plugins is provided [here](https://docs.gateplane.io).

##### Documentation on setting up the plugins is provided [below](#-how-to-install).

### 📋 Features in Detail

The *always-free* **Community Edition**, includes all Vault/OpenBao plugins,
that achieve Just-In-Time, Conditional Privileged Access.

These components are *self-hostable*, with *source always available* in this repository,
and it is allowed to modify to suit your needs, under the terms of [*Elastic License v2*](./LICENSE).

Additionally, the [GatePlane WebUI under `app.gateplane.io` domain](https://app.gateplane.io)
can be used by setting up your Vault / OpenBao instance, using the instructions provided [below](#-how-to-install).

The **Team** and **Enterprise** packages provide services hosted by GatePlane,
such as *Notifications* and *Metrics*, custom domains, WebUI integrations for automated Cloud Service Provider access,
Crypto Wallets and other External Services.

#### Community Edition

* **Approval-Based Access**: Accesses are only granted if selected users approve of it
* **Configurable number of approvers**: Sensitive accesses can be protected by more than 1 approvers
* **Just-In-Time Access**: All accesses granted expire through a configured TTL
* **Reason for access**: Mandated or optional reason for access is embedded in the access request and audit trail
* **GatePlane WebUI under `app.gateplane.io`**: An aggregated view, where one can create, approve and claim access requests

#### Team and Enterprise Edition

* **Notification Service**: Get notified for each *access request*, *approval*, *claim* or *revokation* on your organization's messaging app
* **Metrics Service**: Identify friction points, most used accesses and critical activity windows
* **Dedicated WebUI domain**: Access a pre-configured GatePlane WebUI under a custom domain (e.g: `<myorg>.app.gateplane.io`), allowing for security configurations (e.g: mTLS), Vault/OpenBao login integrations (e.g: Okta login, Userpass, etc) and Access Claim Integrations.
* **WebUI integrations**: Claiming conditional access to specific use-cases, such as Crypto Wallets
* **Support**: Get support on designing your infrastructure access management with GatePlane

* **Zero-Knowledge**: GatePlane infrastructure does only get non-sensitive metadata to provide the above features for *Team and Eneterprise* plans.

**Your organization's access control ***NEVER*** depends on or is shared with GatePlane infrastructure**

## 3. 💻 Installation / Getting Started

### 🧰 Prerequisites

A self-hosted Vault (community or enterprise) or an OpenBao instance is needed to set up the GatePlane Plugins. Install options are available in each respective documentation (see for [Vault](https://developer.hashicorp.com/vault/docs/get-vault#install-options)/[OpenBao](https://openbao.org/docs/install/)) for host and containerized deployments, both for single node and high-available setups.

Additionally, a *plugin directory* needs to be set under Vault/OpenBao configuration's `plugin_directory` directive (see for: [Vault](https://developer.hashicorp.com/vault/docs/configuration#plugin_directory)/[OpenBao](https://openbao.org/docs/configuration/)).

### 🚀 How to Install

#### Downloading and Verifying the Plugins

To register a plugin, it needs to be compiled and located in the Vault/OpenBao plugin directory.

The latest binaries can be downloaded from the [Github Releases Page](https://github.com/gateplane-io/vault-plugins/releases). As SHA256 checksums are needed by Vault/OpenBao to register plugins, the `checksums.txt` file can also be downloaded for quick reference. Additionally, the plugin version SemVer should be noted, as always provided in the Release description (e.g: `v0.1.0-base0.1.0`).

Finally, the builds can be verified using the [`slsa-verifier`](https://github.com/slsa-framework/slsa-verifier) project's `verify-artifact` command.

#### Registering to Plugin Catalog - Manually

As Vault/OpenBao documentation points out (see for [Vault](https://developer.hashicorp.com/vault/docs/plugins/plugin-management)/[OpenBao](https://openbao.org/docs/plugins/plugin-management/)), a `vault register` command should be issued to allow plugins to be *enabled*.

##### Note: `vault` will be used in the examples as the CLI tool, which is interchangeable with `bao`. If using OpenBao, an alias can be set to use copy-paste from this document, like: `alias vault='bao'`

Registering the `gateplane-policy-gate` plugin.
```bash
vault plugin register -sha256=<SHA256 found in the 'checksums.txt'> \
    -version="<SemVer found in the Github Release>" \
    secret \  # All GatePlane plugins are of type "secret"
    gateplane-policy-gate
```

####  Registering to Plugin Catalog - With GatePlane Terraform modules

GatePlane provides helper Terraform modules, that can be used in Infrastructure-as-Code environments.

This helps keeping version handling at check, while also avoiding manual tinkering with high privilege tokens (like ones allowing plugin registration).

In this case, the [GatePlane Setup Terraform](https://github.com/gateplane-io/terraform-gateplane-setup) module can be used.

```hcl
module "gateplane_setup" {
  source = "github.com/gateplane-io/terraform-gateplane-setup?ref=0.4.0"

  policy_gate_plugin = {
    filename = "gateplane-policy-gate"  // The name of the binary for Policy Gate
    version  = "v0.1.0-base0.1.0"       // The version provided in Github Release Page
    sha256   = "01ba4..."               // The SHA256 checksum found in the 'checksums.txt'
  }

  okta_group_gate_plugin = {
    filename = "gateplane-okta-group-gate"  // The name of the binary for Okta Group Gate
    version  = "v0.1.0-base0.1.0"           // The version provided in Github Release Page
    sha256   = "4355a..."                   // The SHA256 checksum found in the 'checksums.txt'
  }

  create_ui_policy = true

  // Allows Cross-Origin Resource Sharing (CORS)
  // to WebUIs to GatePlane WebUI
  url_origins = [
    "https://<your-instance>",
    "https://app.gateplane.io"
  ]
}
```

##### This module can also be used to set Vault/OpenBao instance's CORS headers to [`app.gateplane.io`](app.gateplane.io), through the `url_origins` parameter.

#### Enabling a Gate - Manually

In this example, the Policy Gate plugin will be used to protect a Vault/OpenBao path,
such as `aws/prod/object-writer`, which can be an AWS Secrets Engine (see for [Vault](https://developer.hashicorp.com/vault/docs/secrets/aws)/[OpenBao](https://github.com/openbao/openbao-plugins)), providing AWS Credentials of an IAM User that can do `s3:PutObject` actions to critical S3 buckets (e.g: the company's website).

In that case, a Vault/OpenBao policy must exist (e.g: `aws-prod-object-writer`) that allows access to this path, as follows:

`aws-prod-object-writer.hcl`
```hcl
path "aws/prod/creds/object-writer" {
    capabilities = ["read"]
}
```

To create expiring Vault/OpenBao tokens of this policy, based on approvals, a *Gate* must be created using the Policy Gate plugin:
```bash
vault enable gateplane-policy-gate -path gateplane/aws-prod-object-writer
```
##### Note: The `gateplane/aws-prod-object-writer` path is used for clarity - any path can be used.

Configuring the Gate manually:

```bash
vault write gateplane/aws-prod-object-writer/config \
    required_approvals=1 \  # additional options can be provided
    require_reason=true
```

```bash
vault write gateplane/aws-prod-object-writer/config/access \
    policies=aws-prod-object-writer
```

```bash
vault write gateplane/aws-prod-object-writer/config/api/vault \
  approle_id="..."      # Set an approle that can manipulate Entities
  approle_secret="..."
```

With that, the `/request`, `/approve` and `/claim` endpoints of `gateplane/aws-prod-object-writer` will be usable as in the [Usage Example](#usage-example).

#### Enabling a Gate - With GatePlane Terraform modules

The [Policy Gate Terraform module](https://github.com/gateplane-io/terraform-gateplane-policy-gate) simplifies the above task, also creating the AppRole needed and helper Policies that allow access to `/request`, `/approve` and `/claim` endpoints, ready to be assigned to Vault/OpenBao Entities.

```hcl
module "gateplane_aws-prod-object-writer" {
  depends_on = [module.gateplane_setup] // the module registering the plugins
  source     = "github.com/gateplane-io/terraform-gateplane-policy-gate?ref=1.0.0"

  name            = "aws-prod-object-writer"    // The name to be used in the endpoint and policies
  path_prefix     = "gateplane"                 // The path prefix
  endpoint_prefix = ""                          // A prefix for the endpoint

  // The Vault/OpenBao path to protect can be used directly
  // circumventing the need to create the policy manually.
  protected_path_map = {
    "aws/prod/creds/object-writer" = ["read"]
  }

  // The configuration provided to /config
  plugin_options = {
    "required_approvals" : 1,
    "require_justification": true,
  }
}

output "policies" {
    description = "These policies can be used to access the created Gate"
    value = [
        # Grants access to 'claim' and 'create' access requests
        module.gateplane_aws-prod-object-writer.policy_names["requestor"],
        # Grants access to 'list' and 'approve' access requests
        module.gateplane_aws-prod-object-writer.policy_names["approver"],
    ]
}
```

### ▶️ Usage Example

#### Requesting Access
Vault/OpenBao Entities writing to the `gateplane/aws-prod-object-writer/request` will create an access request:
```bash
$ VAULT_TOKEN="<requestor-token>" \
    vault write gateplane/aws-prod-object-writer/request \
        reason="I want to get in"  # Reason is configured as mandatory for this gate
Key                   Value
---                   -----
claim_ttl             30m
deleted_after         1760859773
exp                   1760776973
iat                   1760773373
justification         I want to get in
num_of_approvals      0
overwrite             false # Whether a request by this entity has been already created
requestor_id          c542f5ab-1e4b-2479-f0a6-ef8b32a3c39e
required_approvals    1
status                pending # status can be: pending / approved / active / expired
```

By design, each Requestor can have exactly one request against a Gate.

#### Approving Access
Then the Approver can approve using the RequestID:
```bash
$ VAULT_TOKEN="<approver-token>" \
    vault write gateplane/aws-prod-object-writer/approve \
        request_id="5ec53023-d998-6b3d-f58f-49976f3b1af7"
Key       Value
---       -----
status    pending
```

##### The Approver gets to know the RequestID either by an out-of-band communication, a LIST to the `/request` endpoint or the *Notification Feature*

#### Claiming Access
```bash
$ VAULT_TOKEN="<requestor-token>" \
    vault write -force gateplane/aws-prod-object-writer/claim
Key                  Value
---                  -----
lease_id             gateplane/aws-prod-object-writer/claim/h3hAUgVBoWMn6uc3vQ6CgEdp
lease_duration       30m
lease_renewable      false
new_policies         [aws-prod-object-writer]
previous_policies    [gateplane-aws-prod-object-writer-requestor]
requestor_id         c542f5ab-1e4b-2479-f0a6-ef8b32a3c39e
```

The requestor's Entity Policies now include `aws-prod-object-writer` until the lease is active (while it is not expired or revoked). The requestor finally can use `vault read aws/prod/creds/object-writer` to issue personalized, temporary AWS credentials.

### 🛠️ How to Build and Test

#### Building

Building the plugin binaries requires Git, Golang 1.24, GoReleaser and GNU make.

With the above requirements, building can be as easy as:
```bash
git clone https://github.com/gateplane-io/vault-plugins
cd vault-plugins
make build-plugin
```

##### Note: the plugins are built under `dist/` and their version is set to `v0.0.0-dev`. This version string must be provided to Vault/OpenBao to register the plugins to the Plugin Catalog.

#### Testing

Testing the plugins additionally requires Python 3.12, `pip`, Terraform/OpenTofu and Docker Compose.

With the above requirements, testing goes like:
```bash
pip install -r test/requirements-dev.txt
make load-infra
make export-resources
make test-infra
```

## 4. 💬 Contact

You can always reach the Dev Team through [Email](mailto:maintainers@gateplane.io).

### Community

We truly believe in the power of the Community, and we appreciate every new member!

You can join us here:
* [Slack](https://join.slack.com/t/gateplane-community/shared_invite/zt-3erzr2612-7Lhsx~cwpQ3kUvqcClIdiQ)
* [Github Discussions](https://github.com/orgs/gateplane-io/discussions)

### ⚖️ License
This project is licensed under the [Elastic License v2](https://www.elastic.co/licensing/elastic-license).

This means:

- ✅ You can use, fork, and modify it for **yourself** or **within your company**.
- ✅ You can submit Pull Requests and redistribute modified versions (with the license attached).
- ❌ You may **not** sell it, offer it as a paid product, or use it in a hosted service (e.g: SaaS).
- ❌ You may **not** re-license it under a different license.

In short: You can use and extend the code freely, privately or inside your business - just don’t build a business around it without our permission.
[This FAQ by Elastic](https://www.elastic.co/licensing/elastic-license/faq) greatly summarizes things.

See the [`./LICENSES/Elastic-2.0.txt`](./LICENSES/Elastic-2.0.txt) file for full details.
