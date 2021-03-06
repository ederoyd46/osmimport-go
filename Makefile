BASE_DIR=`pwd`

default: run

update_deps:
	go get -u github.com/ederoyd46/osm/osmformat
	go get -u github.com/ederoyd46/osm/fileformat
	go get -u github.com/golang/lint/golint
	go get -u golang.org/x/tools/cmd/gotype
	go get -u github.com/nsf/gocode
	go get -u golang.org/x/tools/cmd/goimports
	go get -u github.com/rogpeppe/godef
	go get -u golang.org/x/tools/cmd/oracle
	go get -u golang.org/x/tools/cmd/gorename
	go get -u github.com/kisielk/errcheck
	go get -u github.com/jstemmer/gotags
	go get -u github.com/golang/protobuf/proto
	go get -u github.com/golang/protobuf/protoc-gen-go
	go get -u github.com/k0kubun/pp
	go get -u golang.org/x/tools/cmd/godoc
	go get -u github.com/motemen/gore
	go get -u github.com/dancannon/gorethink
	#Worker Pool Framework
	go get -u github.com/jeffail/tunny
	
run:
	@find . -name '*.go' -print0 | xargs -0 go run

build:
	@go build

install:
	@go install

start_db_matts_mac:
	rethinkdb -d /Volumes/External/db/rethinkdb &> /tmp/rethinkdb.log &

stop_db_matts_mac:
	@pgrep rethinkdb | xargs kill

open_db_matts_mac:
	@open -a firefox http://localhost:8080
