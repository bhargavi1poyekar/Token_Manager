package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"google.golang.org/grpc" // for rpc calls

	pb "token_manage/proto" // import the generated protobuf code
)

func main() {
	// Parse command line flags
	serverAddr := flag.String("host", "localhost", "The server address")
	serverPort := flag.Int("port", 50051, "The server port")

	// flags for functions with default value false
	createFlag := flag.Bool("create", false, "Create a new token")
	writeFlag := flag.Bool("write", false, "Write to a token")
	readFlag := flag.Bool("read", false, "Read from a token")
	dropFlag := flag.Bool("drop", false, "Drop a token")

	// default values for all properties as emty string or 0
	id := flag.String("id", "", "The ID of the token to operate on")
	name := flag.String("name", "", "The name of the token")
	low := flag.Uint64("low", 0, "The lower bound of the token's domain")
	mid := flag.Uint64("mid", 0, "The midpoint of the token's domain")
	high := flag.Uint64("high", 0, "The upper bound of the token's domain")

	// execute parsing
	flag.Parse()

	// Set up connection to server
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", *serverAddr, *serverPort), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}

	// close the connection before exiting function
	defer conn.Close()

	// Create client
	client := pb.NewTokenManagerClient(conn)

	// Create context for rpc call
	ctx := context.Background()

	// Call the required function and then print the response received.
	switch {
	// create token
	case *createFlag:
		res, err := client.CreateToken(ctx, &pb.CreateTokenRequest{Id: *id})
		if err != nil {
			log.Fatalf("Failed to create token: %v", err)
		}
		log.Println(res.Success)

	//write token
	case *writeFlag:
		res, err := client.WriteToken(ctx, &pb.WriteTokenRequest{
			Id:   *id,
			Name: *name,
			Low:  *low,
			Mid:  *mid,
			High: *high,
		})
		if err != nil {
			log.Fatalf("Failed to write to token: %v", err)
		}
		log.Printf("Partial value: %d", res.Partial)

	//read token
	case *readFlag:
		res, err := client.ReadToken(ctx, &pb.ReadTokenRequest{Id: *id})
		if err != nil {
			log.Fatalf("Failed to read from token: %v", err)
		}
		log.Printf("Final value: %d", res.Final)

	// drop/delete token
	case *dropFlag:
		res, err := client.DropToken(ctx, &pb.DropTokenRequest{Id: *id})
		if err != nil {
			log.Fatalf("Failed to drop token: %v", err)
		}
		log.Println(res.Success)

	// if method not mentioned
	default:
		log.Fatalf("Please specify a valid operation (-create, -write, -read, -drop)")
	}
}
