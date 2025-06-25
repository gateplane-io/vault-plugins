import logging
import requests

import random
import string

# Configure logging
logging.basicConfig(
    level=logging.DEBUG,  # Log at DEBUG level or higher
    format="%(asctime)s - %(levelname)s - %(message)s",
    handlers=[
        logging.StreamHandler(),  # Logs to console
        # logging.FileHandler('test_log.log')  # Logs to file
    ],
)
logger = logging.getLogger(__name__)

VAULT_ADDR = "http://127.0.0.1:8200"
VAULT_TOKEN_ROOT = "root"  # set on dev-server
VAULT_API = VAULT_ADDR + "/v1"
# VAULT_PLUGIN_BASE = VAULT_API+"/auth/pgate"
VAULT_URLS = {
    mount: {
        ep: VAULT_API + f"/auth/{mount}/{ep}"
        for ep in ["request", "approve", "claim", "config"]
    }
    for mount in ("mock", "pgate", "oktagate")
}

CLAIM_KEYS = {
    "mock": "data.claimed",
    "pgate": "auth.client_token",
    "oktagate": "data.okta_user_id",
}

TERRAFORM_OUTPUT_FILE = "test/terraform-output.json"
PLUGIN_CONFIG = {
    "ttl": 1800,
    "approval_ttl": 3600,
    "request_ttl": 3600,
    "required_approvals": 1,
    "require_reason": False,
}


def vault_api_request(url, data={}, token=None, method="GET"):
    # Set the headers for the request
    headers = {"Content-Type": "application/json"}
    if token:
        headers["Authorization"] = f"Bearer {token}"

    # Choose the request method (GET, POST, etc.)
    response = requests.request(method.upper(), url, json=data, headers=headers)

    return (
        response.status_code,
        response.json(),
    )  # Return the response as a JSON dictionary


def get_token_for(tf_output, gatekeeper=False, index=-1, type=None):
    role_ = "user"
    if type is not None:
        role_ = type
    if gatekeeper:
        role_ = "gtkpr"

    tokens_for_type = list(tf_output["token_map"]["value"][role_].items())

    if index == -1:
        id_to_token = random.choice(tokens_for_type)
    else:
        id_to_token = tokens_for_type[index]

    return id_to_token[1]


def configure_plugin(plugin, data, token=VAULT_TOKEN_ROOT):
    return vault_api_request(
        url=VAULT_URLS[plugin]["config"], data=data, token=token, method="POST"
    )


def randomword(length=8):
    letters = string.ascii_lowercase + "-_"
    return "".join(random.choice(letters) for i in range(length))


def approval_scenario(plugin, user_token, gtkpr_tokens, authenticate_claim=False):
    claim_key, claim_subkey = CLAIM_KEYS[plugin].split(".")

    configure_plugin(plugin, {"required_approvals": len(gtkpr_tokens)})

    status, output = vault_api_request(
        VAULT_URLS[plugin]["request"], token=user_token, method="POST"
    )
    assert 200 == status

    request_id = output["data"]["request_id"]

    for gtkpr_token in gtkpr_tokens:
        status, output = vault_api_request(
            VAULT_URLS[plugin]["request"], token=user_token, method="GET"
        )
        assert not output["data"]["grant_code"]  # It's empty

        status, output = vault_api_request(
            VAULT_URLS[plugin]["approve"],
            token=gtkpr_token,
            method="POST",
            data={"request_id": request_id},
        )
        assert 200 == status

    status, output = vault_api_request(
        VAULT_URLS[plugin]["request"], token=user_token, method="GET"
    )
    assert 200 == status
    assert output["data"]["grant_code"]

    status, output = vault_api_request(
        VAULT_URLS[plugin]["claim"],
        method="POST",
        token=user_token if authenticate_claim else None,
        data={"grant_code": output["data"]["grant_code"]},
    )
    assert 200 == status

    assert claim_subkey in output[claim_key]
    return output[claim_key]
