# Go versions to use when running containerized tests
goversion = 1.17 1.16

.PHONY: clean help test

help:
	@echo "available targets: clean, test"

clean:
	rm -f coverage.html coverage.out coverage.txt

test:
	@set -e; for GOVERSION in $(goversion); do \
		echo "🧪 🧪 🧪 Testing on Go $${GOVERSION}"; \
		docker build -t testbashertest:$${GOVERSION} --build-arg GOVERSION=$${GOVERSION} -f test/Dockerfile .;  \
		docker run -it --rm --name testbashertest_$${GOVERSION} testbashertest:$${GOVERSION}; \
	done; \
	echo "🎉 🎉 🎉 All tests passed"
	