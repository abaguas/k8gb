.PHONY: up-local
up-local:
	kubectl create ns test-gslb
	kubectl apply -f ./deploy/crds/ohmyglb.absa.oss_gslbs_crd.yaml
	kubectl apply -f deploy/crds/ohmyglb.absa.oss_v1beta1_gslb_cr.yaml
	operator-sdk up local --namespace=test-gslb

.PHONY: test
test:
	go test -v ./pkg/controller/gslb/

.PHONY: e2e-test
e2e-test:
	operator-sdk test local ./pkg/test/