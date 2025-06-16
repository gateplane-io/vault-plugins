from scenarios import get_token_for, configure_plugin, randomword, approval_scenario


class TestPolicyGate:
    "Policy Gate Plugin tests"

    def test_e2e_policies(self, setup_vault_resources):
        new_policies = [randomword() for i in range(3)]
        configure_plugin("pgate", {"policies": new_policies})

        tf_output = setup_vault_resources  # just rename
        user = get_token_for(tf_output, gatekeeper=False)
        gtkpr = get_token_for(tf_output, gatekeeper=True)

        auth = approval_scenario("pgate", user, [gtkpr])
        print(new_policies, auth)
        assert all([policy in auth["policies"] for policy in new_policies])
