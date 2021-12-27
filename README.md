# micro
```
.
├── README.md
├── gitconfig.sh
├── mkdir.bat
├── ql-docker
│   ├── docker-compose.yml
│   ├── etcd
│   │   ├── data
│   │   └── etcd-v3.5.1-linux-amd64.tar.gz
│   ├── go
│   │   ├── Dockerfile
│   │   └── protoc-3.17.3-linux-x86_64.zip
│   ├── mongo
│   │   └── Dockerfile
│   └── mysql
│       ├── Dockerfile
│       ├── conf
│       │   └── my.cnf
│       ├── data
│       └── mysql-files
├── ql-gateway
│   ├── cmd
│   │   └── service.go
│   ├── command
│   │   ├── app.go
│   │   ├── logger.go
│   │   ├── root.go
│   │   └── version.go
│   ├── conf
│   │   └── conf.go
│   ├── configs
│   │   └── config.toml
│   ├── go.mod
│   ├── go.sum
│   ├── internal
│   │   ├── app
│   │   │   ├── app.go
│   │   │   └── base.go
│   │   ├── component
│   │   │   ├── db
│   │   │   │   ├── cpnt.go
│   │   │   │   └── option.go
│   │   │   ├── grpccli
│   │   │   │   ├── client.go
│   │   │   │   └── etcd.go
│   │   │   ├── log
│   │   │   │   └── log.go
│   │   │   ├── mongo
│   │   │   └── redis
│   │   │       └── cpnt.go
│   │   ├── constant
│   │   │   └── constant.go
│   │   ├── encoding
│   │   │   └── http.go
│   │   ├── errutil
│   │   │   └── error.go
│   │   ├── middleware
│   │   │   ├── auth.go
│   │   │   └── jwt.go
│   │   ├── model
│   │   │   ├── auth
│   │   │   │   └── auth.go
│   │   │   └── model.go
│   │   ├── protocol
│   │   │   ├── admin.go
│   │   │   ├── common.go
│   │   │   └── order.go
│   │   └── service
│   │       ├── auth
│   │       │   ├── handler.go
│   │       │   └── server.go
│   │       ├── order
│   │       │   ├── handler.go
│   │       │   ├── rpc.go
│   │       │   └── server.go
│   │       └── user
│   │           ├── handler.go
│   │           ├── rpc.go
│   │           └── service.go
│   ├── log
│   │   ├── access.log
│   │   └── run.log
│   ├── router
│   │   ├── api.go
│   │   ├── auth.go
│   │   ├── handler.go
│   │   ├── mp.go
│   │   └── order.go
│   └── utils
│       └── utils.go
├── ql-mp
│   ├── cmd
│   │   └── api.go
│   ├── command
│   │   ├── app.go
│   │   ├── logger.go
│   │   ├── root.go
│   │   └── version.go
│   ├── conf
│   │   └── conf.go
│   ├── configs
│   │   └── config.toml
│   ├── go.mod
│   ├── go.sum
│   ├── internal
│   │   ├── app
│   │   │   ├── app.go
│   │   │   └── base.go
│   │   ├── component
│   │   │   ├── db
│   │   │   │   ├── cpnt.go
│   │   │   │   └── option.go
│   │   │   ├── hold
│   │   │   │   └── hold.go
│   │   │   ├── log
│   │   │   │   └── log.go
│   │   │   ├── mongo
│   │   │   └── redis
│   │   │       └── cpnt.go
│   │   ├── constant
│   │   │   └── constant.go
│   │   ├── encoding
│   │   │   └── http.go
│   │   ├── endpoint
│   │   │   └── user.go
│   │   ├── errutil
│   │   │   └── error.go
│   │   ├── middleware
│   │   │   └── jwt.go
│   │   ├── model
│   │   │   ├── model.go
│   │   │   └── user.go
│   │   ├── protocol
│   │   │   └── protocol.go
│   │   ├── service
│   │   │   └── user
│   │   │       ├── handler.go
│   │   │       └── service.go
│   │   └── transport
│   │       ├── auth.go
│   │       └── user.go
│   ├── log
│   │   ├── access.log
│   │   └── run.log
│   ├── router
│   │   └── router.go
│   └── utils
│       ├── hashid.go
│       ├── utils.go
│       └── wxbizdatacrypt.go
├── ql-order
│   ├── cmd
│   │   └── service.go
│   ├── command
│   │   ├── app.go
│   │   ├── root.go
│   │   └── version.go
│   ├── conf
│   │   └── conf.go
│   ├── configs
│   │   └── config.toml
│   ├── go.mod
│   ├── go.sum
│   ├── internal
│   │   ├── app
│   │   │   ├── app.go
│   │   │   └── base.go
│   │   ├── component
│   │   │   ├── db
│   │   │   │   ├── cpnt.go
│   │   │   │   └── option.go
│   │   │   ├── log
│   │   │   │   └── log.go
│   │   │   ├── mongo
│   │   │   └── redis
│   │   │       └── cpnt.go
│   │   ├── constant
│   │   │   └── constant.go
│   │   ├── encoding
│   │   │   └── grpc.go
│   │   ├── endpoint
│   │   │   └── order
│   │   │       └── order.go
│   │   ├── errutil
│   │   │   └── error.go
│   │   ├── model
│   │   │   └── model.go
│   │   └── service
│   │       └── order
│   │           ├── handler.go
│   │           └── server.go
│   ├── log
│   │   └── run.log
│   ├── transport
│   │   └── order
│   │       ├── auth.go
│   │       └── order.go
│   └── utils
│       └── utils.go
└── ql-protos
    ├── go.mod
    ├── go.sum
    ├── grpcui.txt
    ├── mk.bat
    ├── mk.sh
    ├── mp
    │   ├── user.pb.go
    │   └── user.proto
    └── order
        ├── order.pb.go
        └── order.proto
```

```
         .__            .__                     
  ______ |  |     _____ |__| ___________  ____  
 / ____/ |  |    /     \|  |/ ___\_  __ \/  _ \ 
< <_|  | |  |__ |  Y Y  \  \  \___|  | \(  <_> )
 \__   | |____/ |__|_|  /__|\___  >__|   \____/ 
    |__|              \/        \/              
```