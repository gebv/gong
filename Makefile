all: run
build:
	DICO_TEMPLATES=./templates/* dico ./src *.go
	gofmt -w ./src
run: build
	GOPATH=${GOPATH}:${PWD} go run src/main.go -stderrthreshold=INFO
test: build
	GOPATH=${GOPATH}:${PWD} go test ./src/store/... -stderrthreshold=INFO
setup:
	ln -s ../settings-app/index.html src/server/web-static/index.html
	ln -s ../settings-app/dist/js src/server/web-static/js