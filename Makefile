test:
	go test -v `go list ./... | egrep -v /vendor/`

build: 
	cd cmd/ && CGO_ENABLED=0 GOOS=linux go build -o cmd .

build-docker:
	docker build -t mangirdas/image-inspector .

run-cmd: build
	./cmd/cmd	

run: build build-docker
	docker run -ti --rm --privileged -p 8000:8000 -v /var/run/docker.sock:/var/run/docker.sock mangirdas/image-inspector
