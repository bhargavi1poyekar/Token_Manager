# CMSC 621. Advanced Operating System.

## Bhargavi Poyekar (CH33454)

## Project 2: Client-Server Token Manager

<br>

## Problem Statement: 

* You are to implement a client-server application in Go for the management of tokens.
* Token ADT: 
A token is an abstract data type, with the following properties: id, name, domain, and state. Tokens
are uniquely identified by their id, which is a string. The name of a token is another string. The
domain of a token consists of three uint64 integers: low, mid, and high. The state of a token
consists of two unit64 integers: a partial value and a final value, which is defined at the integer x
in the range [low, mid) and [mid, high), respectively, that minimizes h(name, x) for a hash function
h.
* Tokens support the following methods:

    * create(id): create a token with the given id. Reset the token’s state to “undefined/null”. Return a success or fail response.
    * drop(id): to destroy/delete the token with the given id
    * write(id, name, low, high, mid):
        1. set the properties name, low, mid, and high for the token with the given id. Assume uint64 integers low <= mid < high.
        2. compute the partial value of the token as argmin_x H(name, x) for x in [low, mid), and reset the final value of the token.
        3. return the partial value on success or fail response
    * read(id):
        1. find argmin_x H(name, x) for x in [mid, high).
        2. set the token’s final value to the minimum of the value in step#1 and its partial value
        3. return the token’s final value on success or fail response
* Implement a client-server solution for managing tokens. Your server should maintain an initially
empty (non-persistent) collection of tokens. Clients issue RPC calls to the server to execute
create, drop, read-write methods on tokens. The server executes such RPCs and returns an
appropriate response to each call. Client-server communication is assumed synchronous. 
* Your solution should in developed in Go and utilize the gRPC and Google Protocol Buffer frameworks. 
* Your server, upon completing an RPC call by a client for a token, it should “dump” on stderr or
stdout all the information it has for that token, followed by a list of the ids of all of its tokens. The server upon starting, it gets the port number to listen to from the command line arguments/flags.
* The server terminates upon receiving a CTRL-C (SIGINT) signal. Your server should allow for
maximum concurrency among non-conflicting RPC calls (eg any pair of concurrent RPC calls on
the same token conflict unless both are read); conflicting RPC calls should be executed in a serial
manner.
* Your client executes a single RPC call (to a server on the specified host and port) upon fetching
the command line arguments/flags
* Your client should print on stdout the response it received from the RPC call to the server, and
then terminate.

## Program Description:

* The hierarchy of my go setup is:

    * /home/bhargavi/go

        * /bin 
        * /pkg
        * /src
            * /token_manage (my project)
                * proto
                    * tokens.proto
                    * tokens_grpc.pb.go (server stub)
                    * tokens.pb.go (client stub)
                * server.go
                * client.go
                * go.mod (module file for dependencies)

* The bin and pkg consists of the external packages required for my go code which are imported in client and server.

* Proto files: 
    * proto folder is present to store the protobuf files, which is required for communication using protocol buffers. They are used to marshal and unmarshal the arguments of the server and clients. The tokens.proto is written by me and it is compiled using the protoc command to create 2 new pb files that consists of the methods for doing the necessary serializations.
    * server and client code imports the protobuf files to communicate with each other.
    * In the proto file, the numbers 1,2,3 given to the properties or arguments are the order of the parameters. The message token becomes a type that is sent as request and response between client and server.
    * All the methods that are called from client to server are mentioned in the proto file.
* Server code:
    * It has the token struct that consists of all the properties
    * The TokenManagerServer has map of tokens and mutex for the shared resources.
    * The Hash function is given in the problem statement that uses SHA-256
    * The main function gets the command-line argument and parse it.
    * It gets the port number and makes it to listen.
    * grpc server is created and the TokenManager Service is registered there.
    * Server then starts listening on the port.
    * All the calls use mutex for proper sharing of resources.
    * The respective methods perform their functions and give out proper responses. 
    * Before returning, dumpToken is called to print the details of tokens.
    * The write function finds the value from [low,mid) that has the minimum Hash value and assigns it to Partial State.
    * The read function finds the value from [mid, high) that has the minimum Hash value and then finds min(minHash, hash(partial)) and assigns it to Final State.
* Client code:
    * The client code gets the command-line arguments.
    * It then calls the required functions using rpc and then prints the response.


## Commands Executed:

* Install golang

        sudo apt-get install golang-go
        mkdir ~/go

* Add these lines at end of .bashrc file to set the paths 

        export GOPATH=~/go
        export PATH=$PATH:$GOPATH/bin

    Then do 

        source ~/.bashrc

* Create project directory:

        mkdir token_manage
        cd token_manage
     Initialize the module token_manage
     
        go mod init token_manage
    
    Fetches and installs the dependencies from go.mod

        go get -v


* Install required packages for grpc and protocol buffer:

        sudo apt-get install -y protobuf-compiler
        go get -u google.golang.org/grpc
        go get -u github.com/golang/protobuf/protoc-gen-go
        go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

        
* Create server.go and client.go files in token_manage
* Create proto files:

        mkdir proto
        cd proto

* Create tokens.proto file in this and run:

        protoc -I=$GOPATH/src --go_out=. --go-grpc_out=. $GOPATH/src/token_manage/proto/tokens.proto

    This will create 2 pb.go files

* Now cd .. to go back to token_manage directory and run client and server codes:

        go run server.go -port 50050
    
    On other terminals run client code.

        go run client.go -create -id 1234 -host localhost -port 50050

        go run client.go -write -id 1234 -name abc -low 0 -mid 10 -high 100 -host localhost -port 50050

        go run client.go -read -id 1234 -host localhost -port 50050

        go run client.go -drop 1234 -host localhost -port 50051

## Output:

* The output is stored in script.txt files for both server and client

* The script started using script server.txt and script client.txt

* You can see the content of both the script files using cat server.txt and cat client.txt in any linux terminal.

## Conclusion:

1. I understood how to use grpc for the rpc calls.
2. I understood how proto files are used for the marshalling and unmarshalling of the parameters.
3. I got to know how to implement the mutex in go lang. 
4. I also learnt how to use flags for command line arguments.
5. I became more comfortable with golang and its syntax and module structure. 

## References:

1. https://grpc.io/docs/languages/go/quickstart/
2. https://protobuf.dev/getting-started/gotutorial/
3. https://grpc.io/docs/languages/go/basics/
4. https://pkg.go.dev/sync
5. https://gobyexample.com/command-line-flags
6. https://go.dev/tour/basics/7









 
    



 