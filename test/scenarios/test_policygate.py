from scenarios import (
    get_token_for,
    configure_plugin,
    randomword,
    approval_scenario,
    vault_api_request,
    VAULT_URLS,
    VAULT_API,
)


class TestPolicyGate:
    "Policy Gate Plugin tests"

    def test_e2e_policies(self, setup_vault_resources):
        new_policies = [randomword() for i in range(3)]
        configure_plugin(
            "pgate",
            {"policies": new_policies},
            url=VAULT_URLS["pgate"]["config/access"],
        )

        tf_output = setup_vault_resources  # just rename
        user = get_token_for(tf_output, gatekeeper=False)
        gtkpr = get_token_for(tf_output, gatekeeper=True)

        data = approval_scenario("pgate", user, [gtkpr])
        status, entity = vault_api_request(
            f"{VAULT_API}/auth/token/lookup-self",
            method="GET",
            token=user,
        )
        print(new_policies, data, entity)
        assert 200 == status
        identity_policies = entity["data"]["identity_policies"]
        assert all([policy in identity_policies for policy in new_policies])
