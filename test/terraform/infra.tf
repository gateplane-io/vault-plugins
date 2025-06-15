module "infra" {
  source = "git@github.com:gateplane-io/terraform-gateplane-setup?ref=0.1.0"

  // To showcase the WebUI locally
  // Allows CORS and IFrames
  domains = ["*"]

  mock_plugin = {
    filename = "gateplane-mock"
    version  = var.plugin_version_mock
    sha256   = filesha256("${path.module}/../../dist/mock_linux_amd64_v1/gateplane-mock")
  }

  policy_gate_plugin = {
    filename = "gateplane-policy-gate"
    version  = var.plugin_version_policy_gate
    sha256   = filesha256("${path.module}/../../dist/policy-gate_linux_amd64_v1/gateplane-policy-gate")
  }

  plugin_directory = "/etc/vault/plugins"
}


// used by the tests
module "mock" {
  depends_on = [module.infra]
  source     = "git@github.com:gateplane-io/terraform-test-modules.git//gateplane-mock?ref=0.1.0"

  name            = "mock"
  path_prefix     = ""
  endpoint_prefix = ""
}

module "access" {
  depends_on = [module.infra]
  source     = "git@github.com:gateplane-io/terraform-gateplane-policy-gate.git?ref=0.1.0"

  name            = "pgate"
  path_prefix     = ""
  endpoint_prefix = ""

  protected_path_map = {
    "secret/data/*" = ["read"]
  }
}

module "tokens" {
  source = "git@github.com:gateplane-io/terraform-test-modules.git//tokens?ref=0.1.0"

  entity_groups = {
    # To showcase the WebUI
    "demo" = {
      "quantity" = 2,
      "policies" = [
        module.access.policy_names["requestor"],
        module.access.policy_names["approver"],
      ]
    },
    # Used by the tests
    "user" = {
      "quantity" = 3,
      "policies" = [
        module.access.policy_names["requestor"],
        module.mock.policy_names["requestor"],
      ]
    },
    "gtkpr" = {
      "quantity" = 3,
      "policies" = [
        module.access.policy_names["approver"],
        module.mock.policy_names["approver"],
      ]
    }
  }
}


output "token_map" {
  value = module.tokens.token_map
}

output "policy_map" {
  value = module.tokens.policy_map
}
