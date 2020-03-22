package main

import (
	ct "context"
	"fmt"
	"io"
	"net/http"
	"time"

	c "github.com/Fathi122/workerqueue/conf"
	pc "github.com/Fathi122/workerqueue/workerproto"

	log "github.com/sirupsen/logrus"

	gpool "github.com/processout/grpc-go-pool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	certPath = "../certs/ca.crt"
)

var address string = c.GetConfig().Parameters.Grpc.Host + c.GetConfig().Parameters.Grpc.Port

// holds the grpc Pool state
var grpcPool *gpool.Pool

// Factory function for creating grpc connnection
func createConnection() (*grpc.ClientConn, error) {
	// Set up a connection to the server.
	certFile := certPath
	creds, sslErr := credentials.NewClientTLSFromFile(certFile, "")
	if sslErr != nil {
		log.Errorln("Error while loading CA trust certificate : ", sslErr)
		return nil, sslErr
	}
	log.Debugln("createConnection : ", address)
	opts := grpc.WithTransportCredentials(creds)
	conn, sslErr := grpc.Dial(address, opts)
	if sslErr != nil || conn == nil {
		log.Errorln("Dial Failed : ", sslErr)
		return nil, sslErr
	}
	return conn, nil
}

// writeData for writting data via gRPC
func writeData(dataToWrite string) []string {
	responses := make([]string, 0)

	// init context and get connection from the pool
	ctx := ct.Background()
	conn, err := grpcPool.Get(ctx)
	if conn != nil && err == nil {
		defer func() {
			log.Debugln("Closing grpc connection ", conn)
			conn.Close()
		}()
		cli := pc.NewWorkerServerClient(conn.ClientConn)
		if cli != nil {
			req := &pc.WriteRequest{
				DataTowrite: dataToWrite,
			}
			// call api
			stream, err := cli.WriteData(ctx, req)
			if err != nil || stream == nil {
				log.Errorln("error while calling WriteData RPC : ", err)
				return nil
			}
			for {
				res, err := stream.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					log.Errorln("Error while reading : ", err)
					return nil
				}
				responses = append(responses, res.GetDatawritten())
			}
			return responses
		}
	}

	log.Errorln("Error in writeData while get connection from Pool : ", err)
	return nil
}

// getData for getting data via gRPC
func getData(keyName string) string {

	// init context and get connection from the pool
	ctx := ct.Background()
	conn, err := grpcPool.Get(ctx)
	if conn != nil && err == nil {
		defer func() {
			log.Debugln("Closing grpc connection ", conn)
			conn.Close()
		}()
		cli := pc.NewWorkerServerClient(conn.ClientConn)
		if cli != nil {
			req := &pc.GetRequest{
				Key: keyName,
			}
			// call api
			data, err := cli.GetData(ctx, req)
			if err != nil {
				log.Errorln("error while calling GetData RPC : ", err)
				return ""
			}
			return data.GetDataResponse()
		}
	}

	log.Errorln("Error in getData while get connection from Pool : ", err)
	return ""
}

// main
func main() {
	// init the grpc pool
	var grpcErr error

	log.SetLevel(log.DebugLevel)
	log.Debugln("------ Starting Client ------")

	grpcPool, grpcErr = gpool.New(createConnection, 0, 5, time.Second)
	if grpcErr != nil {
		grpcPool = nil
		log.Fatalf("Connection Failed")
		return
	}
	// Handler for writting/reading to datastore
	datastoreHandler := func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case "POST":
			dataToWrite, ok := req.URL.Query()["data"]
			if !ok || len(dataToWrite[0]) < 1 {
				log.Errorln("Url Param 'data' is missing")
				fmt.Fprintf(w, "Url Param 'data' is missing")
				return
			}
			responses := writeData(dataToWrite[0])
			if responses != nil && len(responses) > 0 {
				fmt.Fprintf(w, responses[0])
			} else {
				fmt.Fprintf(w, "Got Error")
			}
		case "GET":
			keyName, ok := req.URL.Query()["key"]
			if !ok || len(keyName[0]) < 1 {
				log.Debugln("Url Param 'Key' is missing")
				fmt.Fprintf(w, "Url Param 'Key' is missing")
				return
			}
			data := getData(keyName[0])
			fmt.Fprintf(w, "Data returned : "+data)
		default:
			fmt.Fprintf(w, "Unsupported Method")
		}
	}
	// default handler
	defaultHandler := func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Welcome to page access datastore with /datastore path")
	}

	http.HandleFunc("/datastore", datastoreHandler)
	http.HandleFunc("/", defaultHandler)

	log.Fatalln(http.ListenAndServe(":8080", nil))
}
