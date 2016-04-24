all: run
build:
	DICO_TEMPLATES=./templates/* dico ./src *.go
	gofmt -w ./src
	GOPATH=${GOPATH}:${PWD} go build -o bin/gong src/main.go
run: build
	GOPATH=${GOPATH}:${PWD} go run src/main.go -stderrthreshold=INFO -v=2
run1: build
	GOPATH=${GOPATH}:${PWD} go run src/main.go -stderrthreshold=ERROR -v=0
test: build
	GOPATH=${GOPATH}:${PWD} go test ./src/store/...
	
build_travis:
	DICO_TEMPLATES=./templates/* ./vendor/bin/dico ./src *.go
	gofmt -w ./src
	GOPATH=${PWD}:${PWD}/vendor go build -o bin/gong src/main.go
test_travis: build_travis
	GOPATH=${PWD}:${PWD}/vendor go test ./src/store/... -stderrthreshold=INFO
	
vendor_clean:
	rm -Rf ./vendor
	mkdir ./vendor
	rm -Rf ./bin
	mkdir ./bin
	# mkdir ./logs
	rm -Rf ./pkg
	
vendor_get: vendor_clean
	GOPATH=${PWD}/vendor go get -u -v \
		github.com/shurcooL/github_flavored_markdown \
		github.com/BurntSushi/toml \
		github.com/golang/glog \
		github.com/labstack/echo \
		github.com/blevesearch/bleve \
		github.com/boltdb/bolt \
		github.com/satori/go.uuid \
		github.com/GeertJohan/go.rice \
		github.com/gebv/dico \
		gopkg.in/go-playground/validator.v8
 	
setup:
	ln -s ../settings-app/index.html src/server/web-static/index.html
	ln -s ../settings-app/dist/js src/server/web-static/js