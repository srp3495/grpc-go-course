#!/bin/bash

protoc greet/greetpb/greet.proto --go_out=plugins=grpc:.  //old version before module


protoc greet/greetpb/greet.proto --go-grpc_out=. 
protoc greet/greetpb/greet.proto --go_out=. 


protoc calculator/calculatorpb/calculator.proto --go-grpc_out=. 
protoc calculator/calculatorpb/calculator.proto --go_out=. 

go run greet/greet_client/client.go 
go run greet/greet_server/server.go
go run calculator/calculator_client/client.go
go run calculator/calculator_server/server.go 