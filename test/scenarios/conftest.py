import pytest
import json

from scenarios import TERRAFORM_OUTPUT_FILE


@pytest.fixture(scope="class")
def setup_vault_resources():
    # The file must exist, so
    # 'make export-resources' must have been run
    with open(TERRAFORM_OUTPUT_FILE) as tf_output_file:
        tf_output = json.loads(tf_output_file.read())
        yield tf_output
