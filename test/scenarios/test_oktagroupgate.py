from scenarios import get_token_for, approval_scenario
import os
import requests
import time

OKTA_API_KEY = os.getenv("OKTA_API_KEY")
OKTA_ORG = os.getenv("OKTA_ORG")
OKTA_GROUP_ID = os.getenv("OKTA_GROUP_ID")

OKTA_URL = f"https://{OKTA_ORG}.okta.com"
OKTA_GROUP_API_ENDPOINT = f"{OKTA_URL}/api/v1/groups/{OKTA_GROUP_ID}/users"


def get_okta_group_members():
    resp = requests.get(
        OKTA_GROUP_API_ENDPOINT,
        # API Key Privs: Group Membership Admin
        headers={"Authorization": f"SSWS {OKTA_API_KEY}"},
    )
    try:
        return [entry["id"] for entry in resp.json()]
    except TypeError:
        print("Okta Group API Call failed!")
        return []


class TestOktaGroupGate:
    "Okta Group Gate Plugin tests"

    def test_e2e_okta_group_membership(self, setup_vault_resources):
        tf_output = setup_vault_resources  # just rename
        user = get_token_for(tf_output, type="okta")

        # Group is empty at first
        members = get_okta_group_members()
        assert len(members) == 0

        # No Gatekeeper Tokens: sets 'required_approvals' to 0
        data = approval_scenario("oktagate", user, [], authenticate_claim=True)

        # The user has to be there
        members = get_okta_group_members()
        assert data["okta_user_id"] in members

        # Wait for the TTL (which is set to 2 seconds)
        # The PeriodicFunc runs every minute by default:
        # https://pkg.go.dev/github.com/hashicorp/vault/sdk@v0.18.0/framework#Backend.PeriodicFunc
        # https://github.com/hashicorp/vault/blob/main/vault/core.go#L1256
        # So wait max 1 minute divided in 5s re-checks
        for t in range(1, 12):
            # The Group must be empty again
            members = get_okta_group_members()
            if len(members) == 0:
                continue
            time.sleep(5)

        # The Group must be empty again
        members = get_okta_group_members()
        assert len(members) == 0
