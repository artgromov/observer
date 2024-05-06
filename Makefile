cmd/server/server:
	go build -o $@ cmd/server/*.go

cmd/agent/agent:
	go build -o $@ cmd/agent/*.go

build: cmd/server/server cmd/agent/agent

run-server:
	go run cmd/server/main.go

run-agent:
	go run cmd/agent/main.go

.PHONY: test-my
test-my:
	go test -v ./...

.PHONY: test-iter1
test-iter1: build
	./metricstest -test.v -test.run=^TestIteration1$$ \
		-binary-path=cmd/server/server

.PHONY: test-iter2
test-iter2: build
	./metricstest -test.v -test.run=^TestIteration2[AB]*$$ \
		-source-path=. \
		-agent-binary-path=cmd/agent/agent

.PHONY: test
test: test-my test-iter1 test-iter2


.PHONY: clean
clean:
	rm cmd/server/server || true
	rm cmd/agent/agent || true
