

all: clean-startup server agent tests clean-endup

clean-startup:
	rm -rf cmd/server/server cmd/agent/agent

server:
	go build -buildvcs=false -o ./cmd/server/server ./cmd/server/main.go

agent:
	go build -buildvcs=false -o ./cmd/agent/agent ./cmd/agent/main.go

# Tests

tests: check1 check2 check3 check4 check5

check1:
	bash ./tests/check1.sh

check2:
	bash ./tests/check2.sh

check3:
	bash ./tests/check3.sh

check4:
	bash ./tests/check4.sh

check5:
	bash ./tests/check5.sh

clean-endup:
	rm -rf cmd/server/server cmd/agent/agent
