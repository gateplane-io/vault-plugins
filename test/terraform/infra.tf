module "infra" {
  source = "github.com/gateplane-io/terraform-gateplane-setup?ref=0.2.1"
  # source = "./../../../terraform-gateplane-setup"

  // To showcase the WebUI locally
  // Allows CORS and IFrames
  domains = ["*"]

  mock_plugin = {
    filename = "gateplane-mock"
    version  = var.plugin_test_version
    sha256   = filesha256("${path.module}/../../dist/mock_linux_amd64_v1/gateplane-mock")
  }

  okta_group_gate_plugin = {
    filename = "gateplane-okta-group-gate"
    version  = var.plugin_test_version
    sha256   = filesha256("${path.module}/../../dist/okta-group-gate_linux_amd64_v1/gateplane-okta-group-gate")
  }

  policy_gate_plugin = {
    filename = "gateplane-policy-gate"
    version  = var.plugin_test_version
    sha256   = filesha256("${path.module}/../../dist/policy-gate_linux_amd64_v1/gateplane-policy-gate")
  }

  plugin_directory = "/etc/vault/plugins"
}


// used by the tests
module "mock" {
  depends_on = [module.infra]
  source     = "github.com/gateplane-io/terraform-test-modules.git//gateplane-mock?ref=1.0.0"
  # source     = "./../../../terraform-test-modules/gateplane-mock"

  name            = "mock"
  path_prefix     = ""
  endpoint_prefix = ""

  lease_max_ttl = "3h"
}

module "access" {
  depends_on = [module.infra]
  source     = "github.com/gateplane-io/terraform-gateplane-policy-gate?ref=1.0.0"
  # source     = "./../../../terraform-gateplane-policy-gate"

  name            = "pgate"
  path_prefix     = ""
  endpoint_prefix = ""

  protected_path_map = {
    "secret/data/*" = ["read"]
  }
}

module "okta" {
  depends_on = [module.infra]
  source     = "github.com/gateplane-io/terraform-gateplane-okta-group-gate?ref=1.0.0"
  # source     = "../../../terraform-gateplane-okta-group-gate"

  name            = "oktagate"
  path_prefix     = ""
  endpoint_prefix = ""

  plugin_options = {
    "required_approvals" : 0
  }

  lease_ttl     = "2s"
  okta_group_id = var.okta_test_group_id

  okta_mount_accessor = vault_jwt_auth_backend.okta.accessor
  okta_api = {
    org_url   = local.okta_url
    api_token = var.okta_mount_api_token
  }
}

module "tokens" {
  source = "github.com/gateplane-io/terraform-test-modules.git//tokens?ref=1.0.0"
  # source = "./../../../terraform-test-modules/tokens"

  entity_groups = {
    # To showcase the WebUI
    "demo" = {
      "quantity" = 2,
      "policies" = [
        module.access.policy_names["requestor"],
        module.access.policy_names["approver"],
        module.okta.policy_names["requestor"],
        module.okta.policy_names["approver"],
        module.infra.ui_policy,
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
    },
    "okta" = {
      "quantity" = 1,
      "policies" = [
        module.okta.policy_names["requestor"],
        module.okta.policy_names["approver"],
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
