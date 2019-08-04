package main

import (
	"context"
	"encoding/json"
	"fmt"
	pb "gocode/learn_shippy/consignment-service/proto/consignment"
	"io/ioutil"
	"log"
	"os"

	"google.golang.org/grpc"
)

const (
	address         = "localhost:50051"
	defaultFilename = "consignment.json"
)

func parseFile(file string) (*pb.Consignment, error) {
	var consignment *pb.Consignment
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(data, &consignment)
	return consignment, err
}

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewShippingServiceClient(conn)

	file := defaultFilename
	if len(os.Args) > 1 {
		file = os.Args[1]
	}

	consignment, err := parseFile(file)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(consignment)

	r, err := client.CreateConsignment(context.Background(), consignment)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Created at: %t", r.Created)

	getAll, err := client.GetConsignments(context.Background(), &pb.GetRequest{})
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range getAll.Consignments {
		log.Println(v)
	}
}
