from scenarios import (
    VAULT_API,
    VAULT_TOKEN_ROOT,
    VAULT_URLS,
    vault_api_request,
    get_token_for,
    configure_plugin,
    approval_scenario,
)

import time


class TestMock:
    "Mock Plugin tests"

    def test_endpoint_access(self, setup_vault_resources):
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

        configure_plugin("mock", {"require_justification": True})
        status, output = vault_api_request(
            VAULT_URLS["mock"]["request"], token=user, method="POST"
        )
        print(output)
        assert "requires a justification" in output["errors"][0].lower()
        assert 403 == status
        configure_plugin("mock", {"require_justification": False})

    def test_e2e_zero_approvals(self, setup_vault_resources):
        tf_output = setup_vault_resources  # just rename
        user = get_token_for(tf_output, gatekeeper=False)

        approval_scenario("mock", user, [])

    def test_e2e_auto_expiration_status(self, setup_vault_resources):
        tf_output = setup_vault_resources  # just rename
        user = get_token_for(tf_output, gatekeeper=False)
        gtkpr = get_token_for(tf_output, gatekeeper=True)

        configure_plugin(
            "mock",
            {"lease": "1s"},
            url=VAULT_URLS["mock"]["config/lease"],
        )

        request = approval_scenario("mock", user, [gtkpr])

        print(request)
        assert "active" == request["request"]["status"]
        time.sleep(1.1)
        status, request_raw = vault_api_request(
            VAULT_URLS["mock"]["request"], token=user, method="GET"
        )
        assert "expired" == request_raw["data"]["status"]
        # Reset TTL
        configure_plugin(
            "mock",
            {"lease": "10m"},
            url=VAULT_URLS["mock"]["config/lease"],
        )

    def test_e2e_revocation(self, setup_vault_resources):
        tf_output = setup_vault_resources  # just rename
        user = get_token_for(tf_output, gatekeeper=False)
        gtkpr = get_token_for(tf_output, gatekeeper=True)

        configure_plugin(
            "mock",
            {"lease": "10m"},
            url=VAULT_URLS["mock"]["config/lease"],
        )

        request = approval_scenario("mock", user, [gtkpr])

        print(request)
        assert "active" == request["request"]["status"]
        status, _ = vault_api_request(
            f"{VAULT_API}/sys/leases/revoke",
            token=VAULT_TOKEN_ROOT,
            method="POST",
            data={"lease_id": request["claim"]["lease_id"]},
        )
        assert 204 == status
        status, request_raw = vault_api_request(
            VAULT_URLS["mock"]["request"], token=user, method="GET"
        )

        assert "revoked" == request_raw["data"]["status"]
        # Reset TTL
        configure_plugin(
            "mock",
            {"lease": "10m"},
            url=VAULT_URLS["mock"]["config/lease"],
        )

    def test_e2e_auto_abandoned_status(self, setup_vault_resources):
        tf_output = setup_vault_resources  # just rename
        user = get_token_for(tf_output, gatekeeper=False)

        configure_plugin(
            "mock",
            {"request_ttl": "1s"},
            url=VAULT_URLS["mock"]["config"],
        )
        status, request_raw = vault_api_request(
            VAULT_URLS["mock"]["request"], token=user, method="POST"
        )
        print(request_raw)
        assert "pending" == request_raw["data"]["status"]
        time.sleep(1.1)
        # No approval after the Request TTL has passed
        status, request_raw = vault_api_request(
            VAULT_URLS["mock"]["request"], token=user, method="GET"
        )
        assert "abandoned" == request_raw["data"]["status"]

        # Reset TTL
        configure_plugin(
            "mock",
            {"request_ttl": "1800s"},
            url=VAULT_URLS["mock"]["config"],
        )
