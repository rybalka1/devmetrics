

all: clean-startup server agent tests clean-endup

clean-startup:
	rm -rf cmd/server/server cmd/agent/agent

server:
	go build -buildvcs=false -o ./cmd/server/server ./cmd/server/main.go

agent:
	go build -buildvcs=false -o ./cmd/agent/agent ./cmd/agent/main.go

# Tests

tests: check1 check2 check3 check4 check5 check6 check7

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

check6:
	bash ./tests/check6.sh

check7:
	bash ./tests/check7.sh

clean-endup:
	rm -rf cmd/server/server cmd/agent/agent

# Build componrnts use flag -race
all-race: build-race-server build-race-agent

build-race-server:
	go1.21.8 build  -race -buildvcs=false -o ./cmd/server/server ./cmd/server/main.go

build-race-agent:
	go1.21.8 build  -race -buildvcs=false -o ./cmd/agent/agent ./cmd/agent/main.go
