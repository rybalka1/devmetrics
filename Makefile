

all: clean server agent tests

clean:
	rm -rf cmd/server/server cmd/agent/agent

server:
	go build -buildvcs=false -o ./cmd/server/server ./cmd/server/main.go

agent:
	go build -buildvcs=false -o ./cmd/agent/agent ./cmd/agent/main.go

# Tests

tests: check1 check2 check3

check1:
	bash ./tests/check1.sh

check2:
	bash ./tests/check2.sh

check3:
	bash ./tests/check3.sh
