#!/usr/bin/make

VAULT_CONTAINER="vault-inst-plugin-test"
TEST_VERSION="v0.0.0-dev"

.PHONY:build-plugin
build-plugin:
	RELEASE="" \
	VERSION_mock="${TEST_VERSION}" \
	VERSION_policy_gate="${TEST_VERSION}" \
	VERSION_okta_group_gate="${TEST_VERSION}" \
	goreleaser build --clean --single-target --snapshot

.PHONY:load-resources
load-resources:
	cd test/terraform && terraform init
	cd test/terraform && terraform apply -auto-approve

.PHONY:export-resources
export-resources: load-resources
	cd test/terraform && terraform output -json > ../terraform-output.json

.PHONY:unload-resources
unload-resources: load-resources
	cd test/terraform && terraform init
	cd test/terraform && terraform destroy -auto-approve

.PHONY:exec-vault
exec-vault:
	docker exec -ti ${VAULT_CONTAINER} sh

.PHONY:test-infra
test-infra:
	pytest -v test/scenarios/

.PHONY:load-infra
load-infra:
	docker compose -f test/compose.yaml up -d
	sleep 1

.PHONY:unload-infra
unload-infra:
	docker compose -f test/compose.yaml down
