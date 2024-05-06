cmd/server/server:
	go build -o $@ cmd/server/*.go

cmd/agent/agent:
	go build -o $@ cmd/agent/*.go

build: cmd/server/server cmd/agent/agent

run-server:
	go run cmd/server/main.go

test-iter1: build
	./metricstest -test.v -test.run=^TestIteration1$$ -binary-path=cmd/server/server

.PHONY: clean
clean:
	rm cmd/server/server || true
	rm cmd/agent/agent || true
