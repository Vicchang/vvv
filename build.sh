[ -z "$GOPATH" ] && export GOPATH=$HOME/go

# api
go build -o "server" "api/cmd/main.go"
# lb
go build -o "loadbalancer" "lb/cmd/main.go"