from scenarios import (
    VAULT_URLS,
    vault_api_request,
    get_token_for,
    configure_plugin,
    approval_scenario,
)


class TestMock:
    "Mock Plugin tests"

    def test_endpoint_access(self, setup_vault_resources):
        # logger.debug(setup_vault_resources)
        # logger.debug(VAULT_URLS)
        tf_output = setup_vault_resources  # just rename
        token = get_token_for(tf_output, gatekeeper=False)
        gtkpr_token = get_token_for(tf_output, gatekeeper=True)

        status, output = vault_api_request(VAULT_URLS["mock"]["request"], method="POST")
        assert 403 == status

        status, output = vault_api_request(
            VAULT_URLS["mock"]["request"], token=token, method="POST"
        )
        assert 200 == status

        status, output = vault_api_request(
            VAULT_URLS["mock"]["request"], token=token, method="GET"
        )
        assert 200 == status

        status, output = vault_api_request(
            VAULT_URLS["mock"]["request"], token=gtkpr_token, method="LIST"
        )
        assert 200 == status

    def test_e2e_scenario_simple(self, setup_vault_resources):
        tf_output = setup_vault_resources  # just rename
        user = get_token_for(tf_output, gatekeeper=False)
        gtkpr = get_token_for(tf_output, gatekeeper=True)

        approval_scenario("mock", user, [gtkpr])

    def test_e2e_scenario_multiple_gtkprs(self, setup_vault_resources):
        tf_output = setup_vault_resources  # just rename
        user = get_token_for(tf_output, gatekeeper=False)
        gtkprs = [get_token_for(tf_output, gatekeeper=True, index=i) for i in range(3)]

        approval_scenario("mock", user, gtkprs)

    def test_e2e_scenario_mandatory_reason(self, setup_vault_resources):
        tf_output = setup_vault_resources  # just rename
        user = get_token_for(tf_output, gatekeeper=False)

        configure_plugin("mock", {"require_reason": True})
        status, output = vault_api_request(
            VAULT_URLS["mock"]["request"], token=user, method="POST"
        )
        print(output)
        assert "required" in output["errors"][0].lower()
        assert 403 == status
        configure_plugin("mock", {"require_reason": False})

    def test_e2e_zero_approvals(self, setup_vault_resources):
        tf_output = setup_vault_resources  # just rename
        user = get_token_for(tf_output, gatekeeper=False)

        approval_scenario("mock", user, [])
