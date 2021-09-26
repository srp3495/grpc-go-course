package main

import (
	"context"
	"fmt"
	"grpc-g-course/greet/blog/blogpb"
	"io"
	"log"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Blog client welcomes you")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer cc.Close() //to defer the execution to end of program

	c := blogpb.NewBlogServiceClient(cc)

	fmt.Println("created blog")

	blog := &blogpb.Blog{
		AuthorId: "Sreeroop",
		Title:    "Crypto",
		Content:  "yet to be figured out",
	}
	res, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatalf("create blog failed: %v", err)
	}
	fmt.Printf("Blog has been created : %v", res)
	blogId := res.GetBlog().GetId()

	//read blog

	fmt.Println("Reading blog")
	_, err2 := c.Readblog(context.Background(), &blogpb.ReadBlogRequest{
		Id: "xxxx",
	})
	if err2 != nil {
		fmt.Printf("Error while reading blog, not found : %v \n", err2)
	}

	readBlogReq := &blogpb.ReadBlogRequest{
		Id: blogId,
	}

	readRes, read_Err := c.Readblog(context.Background(), readBlogReq)
	if read_Err != nil {
		fmt.Printf("Error while reading blog, not found : %v \n", err2)
	}
	fmt.Printf("Read the blog %v \n", readRes)

	//update blog
	newBlog := &blogpb.Blog{
		Id:       blogId,
		AuthorId: "Sreeroop Shiv UK",
		Title:    "Stock",
		Content:  "its been figured out",
	}

	updateRes, updateErr := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{
		Blog: newBlog,
	})
	if updateErr != nil {
		fmt.Printf("Error while updating blog, not found : %v \n", err2)
	}
	fmt.Printf("Read the updated blog %v \n", updateRes)

	//delete blog
	_, delErr := c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{
		Id: blogId,
	})
	if delErr != nil {
		fmt.Printf("Error while deleting blog:, %v", delErr)
	}
	fmt.Printf("Deleted Id is : %v \n", blogId)

	//listing blogs

	stream, err := c.ListBlog(context.Background(), &blogpb.ListBlogRequest{})
	if err != nil {
		log.Fatalf("error while calling ListBlog RPC: %v", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Something happened: %v", err)
		}
		fmt.Println(res.GetBlog())
	}

}
