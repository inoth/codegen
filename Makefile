
#all: build down setup clean log

start: clean build package
start-arm: clean build-arm package-arm

init:
	go mod tidy -go=1.16 && go mod tidy -go=1.17

version:
	echo v0.0.1

log:
	docker logs -f codegenfornet

clean:
	- rm -vrf ./release/*

build:
	- go mod tidy
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./release/codegenfornet main.go
	cp Dockerfile release/
	- mkdir release/config
	cp config/*.yaml release/config/

build-arm:
	- go mod tidy
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ./release/codegenfornet main.go
	cp Dockerfile-arm release/
	- mkdir release/config
	cp config/*.yaml release/config/.

package:
	#sudo docker build -t codegenfornet:latest .
	sudo docker buildx build -t codegenfornet:latest --platform linux/amd64 .

package-arm:
	export DOCKER_CLI_EXPERIMENTAL=enabled
	sudo docker buildx build -t codegenfornet:latest -f Dockerfile-arm --platform linux/arm64 .

setdev:
	rsync -avz ./release/ root@localhost:/home/admin/codegenfornet/.
