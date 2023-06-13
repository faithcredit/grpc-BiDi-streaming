package main

import (
	"context"
	"fmt"
	pb "grpc-go-course/greet/greetpb"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello, I'm a client.")
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}

	defer conn.Close()
	c := pb.NewGreetServiceClient(conn)

	// fmt.Println("Created client: %f", c)
	//doUnary(c)
	//doServerStreaming(c)
	// doClientStreaming(c)
	doBiDiStreaming(c)

}

func doUnary(c pb.GreetServiceClient) {
	fmt.Println("Starting to do a Unary RPC...")
	req := &pb.GreetRequest{

		Greeting: &pb.Greeting{
			FirstName: "Grant",
			LastName:  "John",
		},
	}
	res, err := c.Greet(context.Background(), req)

	if err != nil {
		log.Fatalf("Error while calling Greet RPC:%v", err)
	}

	log.Printf("Response from Greet: %v", res.Result)
}

func doServerStreaming(c pb.GreetServiceClient) {
	fmt.Println("Starting to do a Server Streaming RPC...")

	req := &pb.GreetManyTimesRequest{
		Greeting: &pb.Greeting{
			FirstName: "Grant",
			LastName:  "John",
		},
	}
	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling GreetManyTime RPC: %v", err)
	}

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			//We have reached the end of the stream
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stream: %v", err)
		}

		log.Printf("Response from GreetManyTimes: %v", msg.GetResult())
	}
}

func doClientStreaming(c pb.GreetServiceClient) {
	fmt.Println("Starting to do a Client Streaming RPC...")

	requests := []*pb.LongGreetRequest{
		&pb.LongGreetRequest{
			Greeting: &pb.Greeting{
				FirstName: "Grant",
			},
		},
		&pb.LongGreetRequest{
			Greeting: &pb.Greeting{
				FirstName: "John",
			},
		},
		&pb.LongGreetRequest{
			Greeting: &pb.Greeting{
				FirstName: "Docker",
			},
		},
		&pb.LongGreetRequest{
			Greeting: &pb.Greeting{
				FirstName: "David",
			},
		},
		&pb.LongGreetRequest{
			Greeting: &pb.Greeting{
				FirstName: "Wesely",
			},
		},
	}
	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Error while calling LongGreet: %v", err)
	}
	//we iterate over our slice and send each message individually
	for _, req := range requests {
		fmt.Printf("Sending req: %v\n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving reasponse from LongGreet: %v", err)
	}

	fmt.Printf("LongGreet Response: %v\n", res)
}
func doBiDiStreaming(c pb.GreetServiceClient) {
	fmt.Println("Starting to do a BiDi Streaming RPC...")

	requests := []*pb.GreetEveryoneRequest{
		{
			Greeting: &pb.Greeting{
				FirstName: "Grant",
			},
		},
		{
			Greeting: &pb.Greeting{
				FirstName: "John",
			},
		},
		{
			Greeting: &pb.Greeting{
				FirstName: "Docker",
			},
		},
		{
			Greeting: &pb.Greeting{
				FirstName: "David",
			},
		},
		{
			Greeting: &pb.Greeting{
				FirstName: "Wesely",
			},
		},
	}
	//We create a stream by invoking the client

	stream, err := c.GreetEveryone(context.Background())

	if err != nil {
		log.Fatalf("Error while creating stream: %v", err)
		return
	}

	// requests := []*pb.LongGreetRequest{
	// 	&pb.LongGreetRequest{
	// 		Greeting: &pb.Greeting{
	// 			FirstName: "Grant",
	// 		},
	// 	},
	// 	&pb.LongGreetRequest{
	// 		Greeting: &pb.Greeting{
	// 			FirstName: "John",
	// 		},
	// 	},
	// 	&pb.LongGreetRequest{
	// 		Greeting: &pb.Greeting{
	// 			FirstName: "Docker",
	// 		},
	// 	},
	// 	&pb.LongGreetRequest{
	// 		Greeting: &pb.Greeting{
	// 			FirstName: "David",
	// 		},
	// 	},
	// 	&pb.LongGreetRequest{
	// 		Greeting: &pb.Greeting{
	// 			FirstName: "Wesely",
	// 		},
	// 	},
	// }

	waitc := make(chan struct{})
	//We send a bunch of messages to the client (go routine)
	go func() {
		//function to send a bunch of messages
		for _, req := range requests {
			fmt.Printf("Sending message: %v\n", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()
	//We receive a bunch of messages from the client (go routine)
	go func() {
		//function to receive a bunch of messages
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving: %v", err)
				break
			}

			fmt.Printf("Received: %v", res.GetResult())
		}
		close(waitc)
	}()
	//Block until everything is done
	<-waitc
}
