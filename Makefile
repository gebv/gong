all: run
run:
	GOPATH=${GOPATH}:${PWD} go run src/main.go -stderrthreshold=INFO
	