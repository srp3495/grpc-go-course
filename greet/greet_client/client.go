package main

import (
	"context"
	"fmt"
	"grpc-g-course/greet/greet/greetpb"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
)

func main(){
	fmt.Println("Hi from client")
	cc,err:=grpc.Dial("localhost:50051",grpc.WithInsecure())
	if err!=nil{
		log.Fatalf("failed to connect: %v",err)
	}
    defer cc.Close()  //to defer the execution to end of program
	
	c:=greetpb.NewGreetServiceClient(cc)
	//fmt.Println("created client : %f",c)
	//doUnary(c)
	//doServerStreaming(c)
	doClientStreaming(c)

}
func doUnary(c greetpb.GreetServiceClient){
	fmt.Println("starting to do unary RPC")
	req:=&greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{FirstName: "Sreeroop", LastName: "Shiv"},
	}

	res,err:= c.Greet(context.Background(), req)
		if err!=nil{
			log.Fatalf("Error while calling Greet RPC : %v",err)
		}
		log.Printf("Response from Greet: %v",res.Result )

	}

func doServerStreaming(c greetpb.GreetServiceClient){
	fmt.Println("starting to do server streaming RPC")
	req:=&greetpb.GreeetManyRequest{
		Greeting: &greetpb.Greeting{FirstName: "Sreeroop Shiv",
	LastName: "UK",},
	}

	resStream,err:=c.GreetManyTimes(context.Background(),req)
	if err!=nil{
		log.Fatalf("Error while calling GreetMany RPC: %v",err)
	}
	for { //to accept as many message 
		msg,err:=resStream.Recv()

	if err==io.EOF{
		//we have reached end of stream
		break
	}
	if err!=nil{
		log.Fatalf("error while reading stream : %v",err)
	}
	log.Printf("Streamed message is : %v",msg.Result)
}
}	

func doClientStreaming(c greetpb.GreetServiceClient){
	fmt.Println("starting to do client streaming RPC")

	request:=[]*greetpb.LongGreetRequest{
       &greetpb.LongGreetRequest{
		   Greeting: &greetpb.Greeting{
			   FirstName: "Sreeroop",
		   },
	   },
	   &greetpb.LongGreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Sreeshma",
		},
	},
	&greetpb.LongGreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Sivadasan",
		},
	},
	&greetpb.LongGreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Sudarsini",
		},
	},
	}

	    stream,err:=c.LongGreet(context.Background())
		if err!=nil{
			log.Fatalf("Error streaming client : %v",err)
		}

		for _,req :=range request{
			fmt.Printf("\nSending request : %v",req.Greeting.FirstName)
            stream.Send(req)
            time.Sleep(1000*time.Millisecond)
		}
		res,err:= stream.CloseAndRecv()
		if err!=nil{
			log.Fatalf("Error receiving the result from server : %v",err)
		}
        
		fmt.Printf("\n Recieved result is : %v",res)

	}
