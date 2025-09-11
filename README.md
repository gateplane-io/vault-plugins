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
  - [🚧 Contributing](#-contributing)
  - [⚖️ License](#-license)

<!--TOC-->


---


## 1. 💡 Overview

### ❓ What is GatePlane?
GatePlane is a project *made by Security Professionals for Security Professionals* to make time-based, conditional access management approachable and transparent.

Gateplane implements Privileged Access Management (PAM) and Just-In-Time (JIT) access, helping tech groups and companies to give, revoke and monitor permissions across the whole organization. The ultimate goal is to solve problems people in tech organizations stumble upon frequently, such as:
* **Developer**: "*I need to debug production but I don't have access*"
* **IT**: “*I get around 60 requests for sensitive access per day, and I have to manually set them up and tear them down*”
* **Security Officer**: “*I have no idea who has access to what after the last incident and the auditor comes next week*”
* **Security Engineer**: “*Setting up this PAM solution will take forever and will change everything we know about our infrastructure*”

It does so by using the Open-Source and battle-tested tools [Vault](https://developer.hashicorp.com/vault/docs) or [OpenBao](https://openbao.org/) as its underlying infrastructure,
drawing from their *Authentication*, *Authorization*, *Integrations* and *Structured Logging* principles, avoiding to re-invent the wheel.

*GatePlane Community* comes **free of charge** and available for everyone to audit and contribute to it under the [*Elastic License v2*](./LICENSE).

##### More about GatePlane can be found in its [Website]() (*under contruction* 🚧)

### 🎯 Key Features

* **Approval-Based Access**: Access is granted only if someone else approves of it
* **Access that expires**: Every granted access has an expiration date, no exceptions
* **Vast support of Integrations**: Supports Kubernetes, AWS/GCP/Azure, SSH, Databases and the whole suite of Vault/OpenBao Secret Engines, as well as JIT Okta Group Membership management
* [**WebUI to Request, Approve and Claim access**](https://app.gateplane.io) : All available accesses are listed and can be managed through a Web Application tied with your Vault/OpenBao instance
* **Notifications**: When a user requests, approves or claims access, a Slack/Discord/MSTeams/*you-name-it* message is sent
* 🚧 **Metrics**: Measurements on the friction relief, numbers of elevated accesses needed and statistics on mean-time-to-claim

### ❔ Why GatePlane
*GatePlane is a community-first project*.

It makes *Just-In-Time* (JIT) *Privileged Access Management* (PAM) accessible to everyone through Open-Source tools and auditable code.

Our rationale is that trust on the market of PAM solutions is never really gained if auditability is not part of the equation. Also buzzword marketing, and vague promises that often lack technical ground are not in our line of work.

Our mission is enabling everyone to provide these guarantees to their organization without closed-source software, unclear security requirements, opaque components in their threat model and tying their sensitive access to SaaS systems they do not control. All this with direct and honest communication.

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
The *Policy Gate* plugin can utilize all Vault/OpenBao *Secret Engines* (see for: [Vault](https://developer.hashicorp.com/vault/docs/secrets)/[OpenBao](https://openbao.org/docs/secrets/)), providing a conditional access layer to all of them.

This effectively means that AWS, Kubernetes, SSH or Database access configured through Vault/OpenBao can now
be provided through a request/approval flow and also be expiring through configured TTLs.

#### Okta Group Gate
The *Okta Group Gate* recreates the request/approval flow to gain expiring access to Okta Groups,
for users that access Vault/OpenBao through Okta.

The plugin reaches the Okta API through API credentials stored in Vault/OpenBao and
adds the user requesting access to a configured Okta Group for a configured time duration.

##### 🚧 API documentation for each of the plugins is provided [here](./docs/).

##### Documentation on setting up the plugins is provided [below](#-how-to-install).

### 📋 Features in Detail

The *always-free* **Community Edition**, includes all Vault/OpenBao plugins,
that achieve Just-In-Time, Conditional Privileged Access.

These components are *self-hostable*, with *source always available* in this repository,
and it is allowed to modify to suit your needs, under the terms of [*Elastic License v2*](./LICENSE).

Additionally, the [GatePlane WebUI under `app.gateplane.io` domain](https://app.gateplane.io)
can be used by setting up your Vault / OpenBao instance, using the instructions provided [below](#-how-to-install).

The **Team** and **Enterprise** packages are only used through the WebUI, and are tied to services
hosted by GatePlane, such as Notifications and Metrics.

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
* **Support**: Get support on designing your infrastructure access management with GatePlane

* **Zero-Knowledge**: GatePlane infrastructure does only get non-sensitive metadata to provide the above features.

**Your organization's access control *NEVER* depends on or is shared with GatePlane infrastructure**

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
    auth \  # All GatePlane plugins are of type "auth"
    gateplane-policy-gate
```

####  Registering to Plugin Catalog - With GatePlane Terraform modules

GatePlane provides helper Terraform modules, that can be used in Infrastructure-as-Code environments.

This helps keeping version handling at check, while also avoiding manual tinkering with high privilege tokens (like ones allowing plugin registration).

In this case, the [GatePlane Setup Terraform](https://github.com/gateplane-io/terraform-gateplane-setup) module can be used.

```hcl
module "gateplane_setup" {
  source = "github.com/gateplane-io/terraform-gateplane-setup?ref=0.2.0"

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

  plugin_directory = "/etc/vault/plugins"   // The value provided in the 'plugin_directory' configuration key
}
```

##### This module can also be used to set Vault/OpenBao instance's CORS headers to [`app.gateplane.io`](app.gateplane.io), through the `domains` parameter.

#### Enabling a Gate - Manually

In this example, the Policy Gate plugin will be used to protect a Vault/OpenBao path,
such as `aws/prod/object-writer`, which can be an AWS Secrets Engine (see for [Vault](https://developer.hashicorp.com/vault/docs/secrets/aws)/[OpenBao](https://github.com/openbao/openbao-plugins)), providing AWS Credentials of an IAM User that can do `s3:PutObject` actions to critical S3 buckets (e.g: the company's website).

In that case, a Vault/OpenBao policy must exist (e.g: `aws-prod-object-writer`) that allows access to this path, as follows:

`aws-prod-object-writer.hcl`
```hcl
path "aws/prod/object-writer" {
    capabilities = ["read"]
}
```

To create expiring Vault/OpenBao tokens of this policy, based on approvals, a *Gate* must be created using the Policy Gate plugin:
```bash
vault enable gateplane-policy-gate -path auth/gateplane/aws-prod-object-writer
```
##### Note: The `auth/gateplane/aws-prod-object-writer` path is used for clarity. Any `auth/`-prefixed path can be used.

Then, configuring this Gate to grant access to the `aws-prod-object-writer` policy requires accessing the `/config` endpoint:
```bash
vault write auth/gateplane/aws-prod-object-writer/config \
    policies=aws-prod-object-writer \  # multiple policies can be protected at once - separated by comma
    required_approvals=1 \             # additional options can be provided
    require_reason=true
```

With that, the `/request`, `/approve` and `/claim` endpoints of `auth/gateplane/aws-prod-object-writer` will be usable as in the [Usage Example](#usage-example).

#### Enabling a Gate - With GatePlane Terraform modules

The [Policy Gate Terraform module](https://github.com/gateplane-io/terraform-gateplane-policy-gate) simplifies the above task, also creating helper policies that allow access to `/request`, `/approve` and `/claim` endpoints, ready to be assigned to Vault/OpenBao Entities.

```hcl
module "gateplane_aws-prod-object-writer" {
  depends_on = [module.gateplane_setup] // the module registering the plugins
  source     = "github.com/gateplane-io/terraform-gateplane-policy-gate?ref=0.1.0"

  name            = "aws-prod-object-writer"    // The name to be used in the endpoint and policies
  path_prefix     = "gateplane"                 // The path prefix
  endpoint_prefix = ""                          // A prefix for the endpoint

  // The Vault/OpenBao path to protect can be used directly
  protected_path_map = {
    "auth/gateplane/aws-prod-object-writer" = ["read"]
  }

  // The configuration provided to /config
  plugin_options = {
    "required_approvals" : 1,
    "require_reason": true,
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

##### This module also adds capabilities to the `requestor` and `approver` policies so they can be used with [GatePlane WebUI](app.gateplane.io), through the `enable_ui` parameter.

### ▶️ Usage Example

#### Requesting Access
Vault/OpenBao Entities writing to the `auth/gateplane/aws-prod-object-writer/request` will create an access request:
```bash
$ VAULT_TOKEN="<requestor-token>" \
    vault write auth/gateplane/aws-prod-object-writer/request \
        reason="I want to get in"  # Reason is configured as mandatory for this gate
Key           Value
---           -----
exp           2025-07-01T10:46:52.418873353Z
iat           2025-07-01T09:46:52.418873353Z
overwrite     false     # Whether a request by this entity has been already created
reason        I want to get in
request_id    5ec53023-d998-6b3d-f58f-49976f3b1af7 # The Entity ID of the Requestor
status        pending   # status can be: pending / approved / active / expired
```

By design, each Requestor can have exactly one request against a Gate.

#### Approving Access
Then the Approver can approve using the RequestID:
```bash
$ VAULT_TOKEN="<approver-token>" \
    vault write auth/gateplane/aws-prod-object-writer/approve \
        request_id="5ec53023-d998-6b3d-f58f-49976f3b1af7"
Key           Value
---           -----
access_approved true
exp             2025-07-01T10:20:00.326612898Z
iat             2025-07-01T09:50:00.326612898Z
message         access approved
approval_id     5ec53023-d998-6b3d-f58f-49976f3b1af7:dbd64311-28e8-7e28-b1ac-1e5c9aa490e7:+
```

##### The Approver gets to know the RequestID either by an out-of-band communication, a LIST to the `/request` endpoint or the Notification Feature

#### Claiming Access

The Requestor is notified for the approval by polling the `/request` endpoint:
```bash
$ VAULT_TOKEN="<requestor-token>" \
    vault read auth/gateplane/aws-prod-object-writer/request
Key           Value
---           -----
exp           2025-07-01T10:46:52.418873353Z
grant_code    b4608697-73d5-447d-84e2-e244c78b3165  # Generated once the request is approved
iat           2025-07-01T09:46:52.418873353Z
reason        I want to get in
request_id    5ec53023-d998-6b3d-f58f-49976f3b1af7
status        approved # The state changes to approved
```

Using the `grant_code` against the `/claim` endpoint finally grants the Vault/OpenBao token:
```bash
$ VAULT_TOKEN="" \  # see the note
    vault write auth/gateplane/aws-prod-object-writer/claim \
     grant_code="b4608697-73d5-447d-84e2-e244c78b3165"
Key               Value
---               -----
token             s.raPGTZdARXdY0KvHcWSpp5wWZIHNT
token_renewable   false
# Some fields are omitted
policies          ["aws-prod-object-writer"]
```

##### Endpoints that issue Vault/OpenBao tokens (like the Policy Gate's `/claim` endpoint) reject authenticated requests: https://github.com/hashicorp/vault/issues/6074.

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


### 🚧 Contributing


<!-- ### Community -->

<!--- ### Security Notes

This is a security-first project. If you believe you have found a security issue in OpenBao, please responsibly disclose by contacting the maintainers:
* a[at]b.com
 -->

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
