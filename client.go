package main

import (
	"context"
	"flag"

	"io"
	"log"

	"github.com/gopiesy/project-protos/policies"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	serverAddr = flag.String("server_addr", "localhost:9111", "The server address in the format of host:port")
	count      = flag.Int("count", 5, "Number of snapshots")
)

func runExecute(client policies.PolicyServiceClient) {
	ctx := context.Background()

	stream, err := client.StreamSnapshots(ctx)
	if err != nil {
		log.Fatalf("%v.StreamSnapshots(ctx) = %v, %v: ", client, stream, err)
	}

	for {
		if *count > 0 {
			snapshot, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					log.Println("EOF reached exiting Client")
					return
				}
				log.Fatalf("Err in Recv: %v", err)
			}

			log.Println("Snapshot received: ", snapshot.Name)
			*count -= 1

			// send response
			status := policies.SnapshotStatus{SnapshotName: snapshot.Name}
			if err := stream.Send(&status); err != nil {
				log.Fatalf("Err in Send: %v", err)
			}
			log.Println("Ack Sent")
		} else {
			log.Println("Bye Bye server")
			if e := stream.CloseSend(); e != nil {
				log.Fatalln(e)
			}
			return
		}
	}
}

func main() {
	flag.Parse()
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	opts = append(opts, grpc.WithBlock())
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := policies.NewPolicyServiceClient(conn)

	runExecute(client)
}
