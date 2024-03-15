

all: clean server

clean:
	rm -rf cmd/server/server cmd/agent/agent

server:
	go build -buildvcs=false -o ./cmd/server/server ./cmd/server/main.go

agent:
	go build -buildvcs=false -o ./cmd/agent/agent ./cmd/agent/main.go

# Tests

check1:
	metricstest -test.v -test.run=^TestIteration1$ -binary-path=cmd/server/server

check2:
	metricstest -test.v -test.run=^TestIteration2[AB]*$ -source-path=. -agent-binary-path=cmd/agent/agent
