golangci-lint run; go fmt ./...

go test ./... -v;

go test -v -run "Config*" .
go test -v -run "CallRefreshToken" .

