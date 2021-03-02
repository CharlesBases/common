module charlesbases/http

go 1.15

replace github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.5

replace github.com/CharlesBases/common => ../

require (
	github.com/CharlesBases/common v0.0.0-00010101000000-000000000000
	github.com/cihub/seelog v0.0.0-20170130134532-f561c5e57575
	github.com/gorilla/mux v1.8.0
	github.com/urfave/negroni v1.0.0
)
