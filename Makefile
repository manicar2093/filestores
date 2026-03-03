test:
	@ go run github.com/onsi/ginkgo/v2/ginkgo -v --label-filter="!integration" ./...

test-integration:
	@ bash setup-localstack.sh start
	@ go run github.com/onsi/ginkgo/v2/ginkgo -v --label-filter="integration" ./...
	@ bash setup-localstack.sh stop

test-all: test test-integration
