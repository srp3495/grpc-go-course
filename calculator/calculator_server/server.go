package main

import (
	"context"
	"fmt"
	"grpc-g-course/greet/calculator/calculatorpb"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	calculatorpb.UnimplementedCalculatorServiceServer//why i am not sure
}

func(*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error){
     
	fmt.Printf("Got the request rpc : %v",req)
	first:=req.FirstNum
	second:=req.SecondNum

	 sum:=first+second

	 res:=&calculatorpb.SumResponse{
	 	Result: sum,
	 }
	 return res,nil
}

func(*server) PrimeDecompositionStream(req *calculatorpb.PrimeDecompositionStreamRequest,stream calculatorpb.CalculatorService_PrimeDecompositionStreamServer) error{
	fmt.Printf("Got the request rpc : %v",req)

	number:=req.Value
	divisor:= int64(2)

	for number > 1 {
		if number%divisor == 0 {
			stream.Send(&calculatorpb.PrimeDecompositionStreamResponse{
				Result: divisor,
			})
			number = number/divisor
		}else{
			divisor= divisor+1
			fmt.Printf("Divisor added : %v",divisor)
		}
	}
	return nil
}



func(*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error{
	sum:=float32(0)
	count:=float32(0)
	for{
	msg,err:=stream.Recv()
	if err==io.EOF{
		fmt.Printf("End of file reached sending response now %v",msg)

		return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
			Result: sum/count,
		})
	}
	if err!=nil{
		log.Fatalf("Some error occured while streaming : %v",err)
	}
    
	sum+=float32(msg.GetValue())
	count+=1
  }

}


func(*server) SendMaximum(stream calculatorpb.CalculatorService_SendMaximumServer) error{
	fmt.Println("Got the request from RPC")
	num:=int32(0)
	for{
		req,err:=stream.Recv()
		if err==io.EOF{
			fmt.Printf("\n End of file reached :%v",err)
			return nil
		}
		if err!=nil{
			log.Fatalf("\n Error receiving stream : %v",err)
		}
        if req.GetValue()>num{
           num=req.GetValue()
		}
       stream.Send(&calculatorpb.MaxResponse{
		   Result: num,
	   })
   }
}

func main() {

	fmt.Println("Calculator Server")
	lis, err := net.Listen("tcp", "0.0.0.0:50052")
	if err != nil {
		log.Fatal("Failed to listen", err)
	}
	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s,&server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
