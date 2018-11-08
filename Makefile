version := $(shell date +'%Y%m%d%H%M%S') 

clean:
	kubectl delete -f deploy/crds/influxdata_v1alpha1_influxdb_cr.yaml

.PHONY: build
build:
	operator-sdk build aaltameemi/influxdb-backup-operator:v$(version)
	docker push aaltameemi/influxdb-backup-operator:v$(version)
	@sed -E -i 's/(.*?)influxdb-backup-operator:v.*?/\1influxdb-backup-operator:v$(version)/g' deploy/crds/influxdata_v1alpha1_influxdb_cr.yaml
	@echo "Version should be $(version)"
	@cat deploy/crds/influxdata_v1alpha1_influxdb_cr.yaml | grep aal

.PHONY: deploy
deploy:
	kubectl create -f deploy/crds/influxdata_v1alpha1_influxdb_cr.yaml 

.PHONY: test
test:
	kubectl create -f deploy/crds/influxdata_v1alpha1_backup_cr.yaml		
