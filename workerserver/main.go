package main

import (
	"context"
	"math/rand"
	"net"
	"time"

	c "github.com/Fathi122/workerqueue/conf"
	pb "github.com/Fathi122/workerqueue/workerproto"

	log "github.com/sirupsen/logrus"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/coreos/etcd/clientv3"
	"go.etcd.io/etcd/etcdserver/api/v3rpc/rpctypes"
	"gopkg.in/natefinch/lumberjack.v2"
)

// server is used to implement WorkerServer services.
type (
	server struct{}
)

var etcdClient *clientv3.Client

// GetData get data from store
func (s *server) GetData(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	resp, err := etcdClient.Get(ctx, req.GetKey())
	if err != nil {
		log.Errorln("------ Get Key : ", req.GetKey(), "response Error : ", err, " ------")
		return &pb.GetResponse{DataResponse: "No Key"}, nil
	}
	if len(resp.Kvs) == 0 {
		log.Errorln("------ Get Key : ", req.GetKey(), " response no Key found ", "resp.Count ", resp.Count, " ------")
		return &pb.GetResponse{DataResponse: "No Key"}, nil
	}

	log.Debugln("------ Get response : ", string(resp.Kvs[0].Value), " ------")
	return &pb.GetResponse{DataResponse: string(resp.Kvs[0].Value)}, nil
}

// WriteData for writting data to store
func (s *server) WriteData(msg *pb.WriteRequest, stream pb.WorkerServer_WriteDataServer) error {
	// generate a key
	key := getRandomString(64)

	// set context with timeout and put key
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	resp, err := etcdClient.Put(ctx, key, msg.DataTowrite)
	if err != nil {
		switch err {
		case context.Canceled:
			log.Errorln("ctx is canceled by another routine : ", err)
		case context.DeadlineExceeded:
			log.Errorln("ctx is attached with a deadline is exceeded : ", err)
		case rpctypes.ErrEmptyKey:
			log.Errorln("client-side error : ", err)
		default:
			log.Errorln("bad cluster endpoints, which are not etcd servers : ", err)
		}
	} else {
		log.Debugln("------ Write response : ------", resp)
		// send
		stream.Send(&pb.WriteResponse{
			Datawritten: " Data : " + msg.DataTowrite + " written with" + " Key : " + key,
		})
	}

	cancel()
	return err
}

// Generate a random string of A-Z or a-z 0-9 chars with len = l
func getRandomString(len int) string {
	// set seed.
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	// generates random integer between min and max
	randInt := func(min, max int) int {
		return min + r.Intn(max-min)
	}
	// 48-57 ->  0 - 9
	// 97-122 -> a - z
	// 65-90 -> A - Z
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		flip := r.Intn(3)
		switch flip {
		case 0:
			bytes[i] = byte(randInt(48, 57))
		case 1:
			bytes[i] = byte(randInt(97, 122))
		case 2:
			bytes[i] = byte(randInt(65, 90))
		}
	}
	return string(bytes)
}

// init etcd
func initEtcd(serverConfig *c.Config) *clientv3.Client {
	// init etcd client
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{serverConfig.Parameters.Etcd.Host},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatalf("initEtcd Failed client : ", err)
	}

	log.Debugln("Etcd started Host : " + serverConfig.Parameters.Etcd.Host)
	return cli
}

func main() {
	// get Config
	serverConfig := c.GetConfig()
	log.SetLevel(log.DebugLevel)
	/*f, err := os.OpenFile("/tmp/log/workserver.log", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic("Failed to create trace file !")
	}
	defer f.Close()
	*/

	// set log to file with logrotate from lumberjack
	log.SetOutput(&lumberjack.Logger{
		Filename:   "/tmp/log/workserver.log",
		MaxSize:    1, // megabytes
		MaxBackups: 3,
		MaxAge:     1,    //days
		Compress:   true, // disabled by default
	})

	// set log json format
	log.SetFormatter(&log.JSONFormatter{})

	// init etcd client
	etcdClient = initEtcd(serverConfig)
	if etcdClient == nil {
		log.Fatalf("Failed to initialise Etcd")
		return
	}

	lis, err := net.Listen("tcp", serverConfig.Parameters.Grpc.Port)
	if err != nil {
		log.Fatalf("failed to listen : ", err)
		return
	}

	opts := []grpc.ServerOption{}
	certFile := "../certs/server.crt"
	keyFile := "../certs/server.pem"

	creds, sslErr := credentials.NewServerTLSFromFile(certFile, keyFile)
	if sslErr != nil {
		log.Fatalf("Failed loading certificates : ", sslErr)
		return
	}

	log.Debugln("Got Grpc Certs")
	opts = append(opts, grpc.Creds(creds))
	s := grpc.NewServer(opts...)

	log.Debugln("Starting Grpc Server")
	pb.RegisterWorkerServerServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve : ", err)
	}
}
