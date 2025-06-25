locals {
  # oidc/okta: mount path
  # /oidc/callback: endpoint
  vault_jwt_okta_path = "oidc/okta"
  redirect_path       = "ui/vault/auth/${local.vault_jwt_okta_path}/oidc/callback"

  # Vault CLI
  redirect_uris = ["http://localhost:8250/oidc/callback"]

  okta_url    = "https://${local.okta_domain}"
  okta_domain = "${var.okta_org}.okta.com"

  oidc_role_name = "role1"
}

# ============== OIDC Backend for Okta
resource "vault_jwt_auth_backend_role" "okta" {
  backend   = vault_jwt_auth_backend.okta.path
  role_name = local.oidc_role_name
  token_policies = [
    "default",
    module.okta.policy_names["requestor"],
    module.okta.policy_names["approver"],
  ]

  # The Alias Name will be the JWT sub,
  # which is the Okta UserID (00up4...)
  user_claim            = "sub"
  role_type             = "oidc"
  allowed_redirect_uris = local.redirect_uris

  oidc_scopes = ["profile"]
}

resource "vault_jwt_auth_backend" "okta" {
  description        = "Access Vault through Okta App in '${var.okta_org}'"
  path               = local.vault_jwt_okta_path
  type               = "oidc"
  oidc_discovery_url = local.okta_url
  oidc_client_id     = "<not_needed>"
  oidc_client_secret = "<it_will_never_authenticate>"

  default_role = local.oidc_role_name
  tune {
    listing_visibility = "unauth"
  }
}

# ============== Okta Entity Alias
/*
This entity alias would be implicitly created
(https://developer.hashicorp.com/vault/docs/concepts/identity#implicit-entities)
if it was authenticated through the Okta OIDC Backend

It serves as the target of the Okta Group Plugin Test
*/
resource "vault_identity_entity_alias" "okta" {
  name           = var.okta_test_user_id
  mount_accessor = vault_jwt_auth_backend.okta.accessor
  # Piggy-back Entity Alias to existing Vault Entity
  canonical_id = keys(module.tokens.token_map["okta"])[0]
}
