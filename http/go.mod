module charlesbases/http

go 1.16

replace github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.5

replace github.com/CharlesBases/common => ../

require (
	github.com/CharlesBases/common v0.0.0-00010101000000-000000000000
	github.com/golang/protobuf v1.4.3
	github.com/google/uuid v1.2.0
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v0.0.0-20170926233335-4201258b820c
	github.com/urfave/negroni v1.0.0
	google.golang.org/protobuf v1.24.0
)
