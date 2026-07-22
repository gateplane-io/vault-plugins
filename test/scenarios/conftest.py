import json

import pytest

from scenarios import (
    TERRAFORM_OUTPUT_FILE,
    revoke_plugin_claim_leases,
)


PLUGIN_MOUNTS = ("mock", "pgate", "oktagate")


def revoke_all_plugin_claim_leases():
    failures = []
    for plugin in PLUGIN_MOUNTS:
        status, output = revoke_plugin_claim_leases(plugin)
        if status not in (200, 204):
            failures.append(f"{plugin}: HTTP {status}: {output}")

    if failures:
        pytest.fail(
            "Could not clean up plugin claim leases: " + "; ".join(failures),
            pytrace=False,
        )


@pytest.fixture(autouse=True)
def isolate_plugin_claim_leases():
    """Prevent shared test identities from retaining active claims between tests."""
    revoke_all_plugin_claim_leases()
    yield
    revoke_all_plugin_claim_leases()


@pytest.fixture(scope="class")
def setup_vault_resources():
    # The file must exist, so
    # 'make export-resources' must have been run
    with open(TERRAFORM_OUTPUT_FILE) as tf_output_file:
        tf_output = json.loads(tf_output_file.read())
        yield tf_output
