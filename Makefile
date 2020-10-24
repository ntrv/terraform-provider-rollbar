TEST?=$$(go list ./... | grep -v vendor | grep -v babbel-rollbar-client | grep -v babbel-terraform-provider)
HOSTNAME=github.com
NAMESPACE=rollbar
NAME=rollbar
BINARY=terraform-provider-${NAME}
VERSION=0.2.0
OS_ARCH=linux_amd64

#default: install
default: dev

dev: install _dev_cleanup _dev_init _dev_apply _dev_log
plan: install _dev_cleanup _dev_init _dev_plan

dev_auto_apply: install _dev_cleanup _dev_init _dev_apply_auto _dev_log

dev_no_debug: install _dev_cleanup _dev_init _dev_apply_nodebug


_dev_cleanup:
	# Cleanup last run
	rm -vrf example/.terraform /tmp/terraform-provider-rollbar.log
_dev_init:
	# Initialize terraform
	(cd example && terraform init)
_dev_apply:
	# Test the provider
	(cd example && TERRAFORM_PROVIDER_ROLLBAR_DEBUG=1 terraform apply) || true
_dev_apply_nodebug:
	# Test the provider
	(cd example && terraform apply) || true
_dev_apply_auto:
	# Test the provider
	(cd example && TERRAFORM_PROVIDER_ROLLBAR_DEBUG=1 terraform apply --auto-approve) || true
_dev_log:
	# Print the debug log
	cat /tmp/terraform-provider-rollbar.log

_dev_plan:
	# Test the provider
	(cd example && terraform plan)


build:
	go build -o ${BINARY}

release:
	GOOS=darwin GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_darwin_amd64
	GOOS=freebsd GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_freebsd_386
	GOOS=freebsd GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_freebsd_amd64
	GOOS=freebsd GOARCH=arm go build -o ./bin/${BINARY}_${VERSION}_freebsd_arm
	GOOS=linux GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_linux_386
	GOOS=linux GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_linux_amd64
	GOOS=linux GOARCH=arm go build -o ./bin/${BINARY}_${VERSION}_linux_arm
	GOOS=openbsd GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_openbsd_386
	GOOS=openbsd GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_openbsd_amd64
	GOOS=solaris GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_solaris_amd64
	GOOS=windows GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_windows_386
	GOOS=windows GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_windows_amd64

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

test: 
	go test -covermode=atomic -coverprofile=coverage.out $(TEST) || exit 1
	@#echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4                    

testacc: 
	TF_ACC=1 go test -covermode=atomic -coverprofile=coverage.out $(TEST) -v $(TESTARGS) -timeout 120m   

slscan:
	./.slscan.sh