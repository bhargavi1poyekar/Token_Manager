package main

import (
	"context" // for grpc
	"flag"    // for command-line input
	"fmt"     // for printing
	"log"     // to log the errors
	"math"    // to get the min values
	"net"     // for server to listen
	"sync"    // for mutex

	"google.golang.org/grpc" // for rpc calls

	"crypto/sha256"   // required for Hash
	"encoding/binary" // output of hash

	pb "token_manage/proto" // import the generated protobuf code
)

// Token struct with the required properties
type Token struct {
	ID     string
	Name   string
	Domain struct {
		Low  uint64
		Mid  uint64
		High uint64
	}
	State struct {
		Partial uint64
		Final   uint64
	}
}

// TokenManagerServer is a structure given as an input to all the function
// This structure has the mutex which maintaint the management of shared resources.
type TokenManagerServer struct {
	pb.UnimplementedTokenManagerServer
	tokens map[string]*Token // map of tokens to store all tokens
	mu     sync.Mutex        // Sharing of resources
}

// Hash concatenates a message and a nonce and generates a hash value.
func Hash(name string, nonce uint64) uint64 {
	hasher := sha256.New()
	hasher.Write([]byte(fmt.Sprintf("%s %d", name, nonce)))
	return binary.BigEndian.Uint64(hasher.Sum(nil))
}

// NewTokenManager creates a new instance of TokenManager
// required for grpc registering
func NewTokenManager() *TokenManagerServer {
	return &TokenManagerServer{
		tokens: make(map[string]*Token),
	}
}

func main() {
	var port = flag.Int("port", 50051, "The server port")
	// flag input get Int pointer, with flag name port, default value 50051 and description

	flag.Parse() // execute the parsing

	port_no := fmt.Sprintf(":%d", *port) // get the value of port and convert to string
	// : => the colon means on all available network interfaces

	lis, err := net.Listen("tcp", port_no) // Make the listening port
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()        // create instance of grpc server
	token_m := NewTokenManager() // new instance of TokenManagerServer

	pb.RegisterTokenManagerServer(s, token_m) // Register the Service with grpc server

	log.Printf("Starting server on port %d", *port) // Print after starting the server

	if err := s.Serve(lis); err != nil { // call server on server
		log.Fatalf("Failed to serve: %v", err)
	}

}

// Create creates a new token with the given ID
func (tm *TokenManagerServer) CreateToken(ctx context.Context, req *pb.CreateTokenRequest) (*pb.CreateTokenResponse, error) {
	tm.mu.Lock() // start the critical code
	defer tm.mu.Unlock()
	// defer tells to unlock the resource before leaving the function,
	//even if it leaves with an error

	id := req.GetId()               // get the requested ID
	if _, ok := tm.tokens[id]; ok { // if already present
		return nil, fmt.Errorf("token with id %s already exists", id)
	}

	//create new token
	tm.tokens[id] = &Token{
		ID: id,
	}

	tm.DumpToken(id) // print the token details at server side

	return &pb.CreateTokenResponse{ //return success message
		Success: true,
	}, nil

}

// Drop deletes the token with the given ID
func (tm *TokenManagerServer) DropToken(ctx context.Context, req *pb.DropTokenRequest) (*pb.DropTokenResponse, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	id := req.GetId() // get the id
	if _, ok := tm.tokens[id]; !ok {
		return nil, fmt.Errorf("token with id %s does not exist", id)
	}

	// delete the token from token structure

	delete(tm.tokens, id)

	// print the token details
	tm.DumpToken(id)

	// return success message
	return &pb.DropTokenResponse{
		Success: true,
	}, nil
}

// Write sets the properties of the token with the given ID
func (tm *TokenManagerServer) WriteToken(ctx context.Context, req *pb.WriteTokenRequest) (*pb.WriteTokenResponse, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	id := req.GetId()      //get id
	t, ok := tm.tokens[id] // get token with the id
	if !ok {
		return nil, fmt.Errorf("token with id %s does not exist", id)
	}

	// Assign the Name, Low, Mid, High to the token
	t.Name = req.GetName()
	t.Domain.Low = req.GetLow()
	t.Domain.Mid = req.GetMid()
	t.Domain.High = req.GetHigh()

	// To find the min Hash Value
	var min uint64 = t.Domain.Low       // min index as low
	var minHash uint64 = math.MaxUint64 // Max possible integer

	for i := t.Domain.Low; i < t.Domain.Mid; i++ { // for range [low,mid)
		h := Hash(t.Name, i) // get hash
		if h < minHash {     // if hash is less then minHash
			min = i     // argmin
			minHash = h // new min hash
		}
	}

	t.State.Partial = min // set the Partial Value
	t.State.Final = 0     // set the Final Value to 0

	tm.DumpToken(id) // print token details

	return &pb.WriteTokenResponse{ // return response
		Partial: t.State.Partial,
	}, nil
}

// Read calculates the final state of token
func (tm *TokenManagerServer) ReadToken(ctx context.Context, req *pb.ReadTokenRequest) (*pb.ReadTokenResponse, error) {

	tm.mu.Lock()
	defer tm.mu.Unlock()

	// Find the token by id.
	id := req.GetId()
	token, exists := tm.tokens[id] // find token from id

	if !exists {
		return nil, fmt.Errorf("token with id %s does not exist", req.GetId())
	}

	// to get min hash value
	var minValue uint64 = math.MaxUint64
	var min uint64 = 0

	// fing min Hash and min arg from [mid, high)
	for x := token.Domain.Mid; x < token.Domain.High; x++ {
		hash := Hash(token.Name, x)
		if hash < minValue {
			minValue = hash
			min = x
		}

	}

	// hash value of partial value
	hashpartial := Hash(token.Name, token.State.Partial)

	// if min hash from [mid,high] is less than min hash from [low, mid)
	if minValue < hashpartial {
		token.State.Final = min // final=argmin of H() for [mid,high]
	} else {
		token.State.Final = token.State.Partial // final=partial
	}

	tm.DumpToken(id) // print token details

	//return final value as response
	return &pb.ReadTokenResponse{Final: token.State.Final}, nil
}

// Dump all the token info after rpc calls on stdout
func (tm *TokenManagerServer) DumpToken(id string) {

	t, ok := tm.tokens[id] //get token from id

	if !ok {
		fmt.Errorf("token with id %s does not exist", id)
	} else { // Print the details of that token
		fmt.Println("\n RPC ended for: \n Token Id: ", t.ID)
		fmt.Println("Token Name: ", t.Name)
		fmt.Println("Token Domain Low: ", t.Domain.Low)
		fmt.Println("Token Domain Mid: ", t.Domain.Mid)
		fmt.Println("Token Domain High: ", t.Domain.High)
		fmt.Println("Token State Partial: ", t.State.Partial)
		fmt.Println("Token State Final: ", t.State.Final)

	}

	// Print the id's of all tokens
	fmt.Println("\nAll Token id's:")
	for id := range tm.tokens {
		fmt.Println(id)
	}

	// return back to original function
	return
}
