all: run
build:
	DICO_TEMPLATES=./templates/* dico ./src *.go
	gofmt -w ./src
	GOPATH=${GOPATH}:${PWD} go build -o bin/gong src/main.go
run: build
	GOPATH=${GOPATH}:${PWD} go run src/main.go -stderrthreshold=INFO
test: build
	GOPATH=${GOPATH}:${PWD} go test ./src/store/... -stderrthreshold=INFO
	
test_travis: build
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
		github.com/GeertJohan/go.rice
	
setup:
	ln -s ../settings-app/index.html src/server/web-static/index.html
	ln -s ../settings-app/dist/js src/server/web-static/js