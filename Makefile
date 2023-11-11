run:
	go run cmd/main.go

test:
	go test  ./... -tags all -cover -coverprofile=cover.out.tmp -covermode=atomic -coverpkg ./...

test-cov:
	make test
	cat cover.out.tmp | grep -v "mocks" > cover.out
	go tool cover -func cover.out
	go tool cover -html=cover.out

test-unit:
	go test  ./... -tags unit -cover -coverprofile=cover.out.tmp -covermode=atomic -coverpkg ./...
 
test-int:
	go test  ./... -tags integration -cover -coverprofile=cover.out.tmp -covermode=atomic -coverpkg ./...

