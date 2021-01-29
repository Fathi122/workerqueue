# simple rest to grpc test

## generate grpc certs

```
./gengrpcert.sh
```
## re-generate grpc protobuf

```
protoc -I workerproto/ workerproto/workerpro.proto --go_out=plugins=grpc:./workerproto
```

## start client

- set TLS certs for docker-compose
```
export COMPOSE_TLS_VERSION=TLSv1_2
```

- start client
```
./runclient.sh start
```

- stop client
```
./runclient.sh stop
```

## build and run server

- set TLS certs for docker-compose
```
export COMPOSE_TLS_VERSION=TLSv1_2
```

- start server
```
./runserver.sh start
```

- stop server
```
./runserver.sh stop
```

## test commands

- write test data to etcd

```
curl -XPOST "http://localhost:8080/datastore?data='test%20data%20to%20write'"
```

```
curl -XPOST "http://192.168.99.100:8080/datastore?data='test%20data%20to%20write'"
```

- read test data with key returned

```
curl -XGET "http://localhost:8080/datastore?key=6yw6y86ohg84y0bdnldwqnfkmn81n26lu6b2i1wiy58t404txk516nfdcyhmn7fu"
```

```
curl -XGET "http://192.168.99.100:8080/datastore?key=6yw6y86ohg84y0bdnldwqnfkmn81n26lu6b2i1wiy58t404txk516nfdcyhmn7fu"
```