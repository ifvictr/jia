MAIN=./cmd/jia.go

all: run

build:
	go build -o ./bin/jia ${MAIN}

clean:
	rm -rfv ./bin

run:
	go run -race ${MAIN}
