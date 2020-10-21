tests-unit:
	go test -coverprofile=coverage.out ./...

code-coverage:
	go tool cover -func=coverage.out

composed-container:
	go test -count=1 -v  ./container -run="TestComposition/Print a composed interface"
