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
        ep: VAULT_API + f"/{mount}/{ep}"
        for ep in [
            "request",
            "approve",
            "claim",  # usage
            "config",
            "config/lease",
            "config/access",  # configuration
        ]
    }
    for mount in ("mock", "pgate", "oktagate")
}

CLAIM_KEYS = {
    "mock": "data.requestor_id",
    "pgate": "data.new_policies",
    "oktagate": "data.requestor_id",
}

TERRAFORM_OUTPUT_FILE = "test/terraform-output.json"
PLUGIN_CONFIG = {
    "ttl": 600,
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


def configure_plugin(plugin, data, token=VAULT_TOKEN_ROOT, url=None):
    if url is None:
        url = VAULT_URLS[plugin]["config"]
    return vault_api_request(url=url, data=data, token=token, method="POST")


def randomword(length=8):
    letters = string.ascii_lowercase + "-_"
    return "".join(random.choice(letters) for i in range(length))


def approval_scenario(plugin, user_token, gtkpr_tokens):
    claim_key, claim_subkey = CLAIM_KEYS[plugin].split(".")

    configure_plugin(plugin, {"required_approvals": len(gtkpr_tokens)})

    status, output = vault_api_request(
        VAULT_URLS[plugin]["request"], token=user_token, method="POST"
    )
    assert 200 == status

    requestor_id = output["data"]["requestor_id"]

    for gtkpr_token in gtkpr_tokens:
        status, output = vault_api_request(
            VAULT_URLS[plugin]["request"], token=user_token, method="GET"
        )
        assert output["data"]["status"] == "pending"

        status, output = vault_api_request(
            VAULT_URLS[plugin]["approve"],
            token=gtkpr_token,
            method="POST",
            data={"requestor_id": requestor_id},
        )
        assert 200 == status

    status, output = vault_api_request(
        VAULT_URLS[plugin]["request"], token=user_token, method="GET"
    )
    assert 200 == status
    assert output["data"]["status"] == "approved"

    status, claim_output = vault_api_request(
        VAULT_URLS[plugin]["claim"],
        method="POST",
        token=user_token,
    )
    print(claim_output)
    assert 200 == status

    status, output = vault_api_request(
        VAULT_URLS[plugin]["request"], token=user_token, method="GET"
    )
    assert 200 == status
    assert output["data"]["status"] == "active"

    assert claim_subkey in claim_output[claim_key]
    return {**claim_output["data"], **output[claim_key]}
