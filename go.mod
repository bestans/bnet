module bnet

go 1.13

require (
	github.com/golang/protobuf v1.4.0
	google.golang.org/protobuf v1.4.0
)

replace (
	github.com/golang/protobuf => ../github.com/golang/protobuf-master-2020-3-12
	github.com/google/go-cmp => ../github.com/google/go-cmp
	golang.org/x/xerrors => ../golang.org/x/xerrors
	google.golang.org/protobuf => ../google.golang.org/protobuf-go-1.20.1
)
