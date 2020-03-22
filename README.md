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

```
./runclient.sh
```

## build and run server

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

- read test data with key returned

```
curl -XGET "http://localhost:8080/datastore?key=6yw6y86ohg84y0bdnldwqnfkmn81n26lu6b2i1wiy58t404txk516nfdcyhmn7fu"
```