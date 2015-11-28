default: rrproxy

.PHONY: default build clean update-godeps docker-build

rrproxy: server.go
	CGO_ENABLED=0 godep go build -v

update-godeps:
	rm -rf Godeps
	godep save

clean:
	rm -f rrproxy

docker-build-env:
	docker build -t dod-workshop-rrproxy .
