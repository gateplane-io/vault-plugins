
build-plugin:
	goreleaser build --clean --single-target --snapshot


load-resources:
	cd test/terraform && terraform init
	cd test/terraform && terraform apply -auto-approve


export-resources: load-resources
	cd test/terraform && terraform output -json > ../terraform-output.json


unload-resources: load-resources
	cd test/terraform && terraform init
	cd test/terraform && terraform destroy -auto-approve


exec-vault:
	docker exec -ti ${VAULT_CONTAINER} sh


test-infra:
	pytest -v test/scenarios/


load-infra:
	docker compose -f test/compose.yaml up -d
	sleep 1


unload-infra:
	docker compose -f test/compose.yaml down
