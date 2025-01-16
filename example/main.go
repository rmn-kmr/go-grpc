package main

import (
	"context"
	"log"
	"time"

	pb "github.com/rmnkmr/lsp/proto"

	"google.golang.org/grpc"
)

func main() {
	// Set up a connection to the server
	conn, err := grpc.Dial("localhost:8085", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Create a client using the connection
	c := pb.NewExternalAPIClient(conn)

	// Prepare a context with a timeout
	ctx := context.Background()
	ctx = context.WithValue(ctx, "X-Request-Id", "123")
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	// Call the SayHello method
	r, err := c.LoanProviders(ctx, &pb.LoanProvidersRequest{})
	if err != nil {
		log.Fatalf("could not get response: %v", err)
	}

	log.Printf("Reply: %s", r)
}
