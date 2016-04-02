all: run
run:
	GOPATH=${GOPATH}:${PWD} go run src/main.go -stderrthreshold=INFO
setup:
	ln -s ../settings-app/index.html src/server/web-static/index.html
	ln -s ../settings-app/dist/js src/server/web-static/js