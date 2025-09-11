# Security Policy

## Reporting a vulnerability

**GatePlane is a security-first project**

If you believe you have found a vulnerability, please disclose it responsibly,
either by using [Github Responsible Disclosure program for `gateplane-io/vault-plugins`](https://github.com/gateplane-io/vault-plugins/security), or directly through email at [`maintainers@gateplane.io`](mailto:maintainers@gateplane.io).

### What has to be provided

In order for us to be able to recreate and assess the vulnerability, please include any of the following:

* Proof of Concept (PoC) CLI commands or API calls
* The version of Vault/OpenBao and the plugin's version (as reported by `vault auth list`)
* All logs and outputs involved

### In casse GatePlane is not the vulnerable component

As GatePlane plugins use Vault/OpenBao, it is possible the that vulnerable behaviour is part of their functionality, and not GatePlane plugins.

In that case you will be notified with evidence on that matter,
and Vault and/or OpenBao projects will be notified by us, according to Vault or OpenBao Security Programs (for [Vault](https://github.com/hashicorp/vault?tab=security-ov-file) / [OpenBao](https://github.com/openbao/openbao?tab=security-ov-file)), mentioning you as the author of the vulnerability.
