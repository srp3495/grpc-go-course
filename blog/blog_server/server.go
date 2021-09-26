package main

import (
	"context"
	"fmt"
	"grpc-g-course/greet/blog/blogpb"
	"log"
	"net"
	"os"
	"os/signal"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

var collection *mongo.Collection

type server struct {
	blogpb.UnimplementedBlogServiceServer
}

type blogItem struct {
	Id       primitive.ObjectID `bson:"_id,omitempty"` //mapping with db field
	AuthorId string             `bson:"author_id"`
	Content  string             `bson:"content"`
	Title    string             `bson:"title"`
}

func (*server) CreateBlog(ctx context.Context, req *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error) {
	fmt.Println("UPDATE blog request received")
	blog := req.GetBlog()

	data := blogItem{
		AuthorId: blog.GetAuthorId(),
		Title:    blog.GetTitle(),
		Content:  blog.GetContent(),
	}

	res, err := collection.InsertOne(context.Background(), data)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprint("Internal error : %v", err))
	}
	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(codes.Internal, fmt.Sprint("Cannot convert to OID : %v", err))
	}
	return &blogpb.CreateBlogResponse{
		Blog: &blogpb.Blog{Id: oid.Hex(), AuthorId: blog.GetAuthorId(), Title: blog.GetTitle(), Content: blog.GetContent()},
	}, nil
}

func(*server) Readblog(ctx context.Context, in *blogpb.ReadBlogRequest) (*blogpb.ReadBlogResponse, error){
fmt.Println("Read blog Request")
	id:=in.GetId()
	oid,err:=primitive.ObjectIDFromHex(id)
	if err!=nil{
		return nil,status.Errorf(codes.InvalidArgument,fmt.Sprint("Error while converting OID : %v",err))
	}
	data:=&blogItem{}

	filter:=bson.M{"_id":oid}

	res:=collection.FindOne(context.Background(),filter)
	if err=res.Decode(data);err!=nil{
		return nil,status.Errorf(codes.NotFound,fmt.Sprint("Id not found in database: %v",err))
	}
	return &blogpb.ReadBlogResponse{
		Blog: &blogpb.Blog{
			Id: data.Id.Hex(),
			AuthorId: data.AuthorId,
			Title: data.Title,
			Content: data.Content,

		},
	},nil

}


func(*server) UpdateBlog(ctx context.Context, in *blogpb.UpdateBlogRequest) (*blogpb.UpdateBlogResponse, error){
	fmt.Println("UPDATE blog Server started")
	blog:=in.GetBlog()
	oid,err:=primitive.ObjectIDFromHex(blog.GetId())
	if err!=nil{
		return nil,status.Errorf(codes.InvalidArgument,fmt.Sprintf("Error while converting OID : %v",err))
	}
	data:=&blogItem{}
	filter:=bson.M{"_id":oid}

	res:=collection.FindOne(context.Background(),filter)
	if err=res.Decode(data);err!=nil{
		return nil,status.Errorf(codes.NotFound,fmt.Sprintf("Id not found in database: %v",err))
	}

	data.AuthorId=blog.GetAuthorId()
	data.Title=blog.GetTitle()
	data.Content=blog.GetContent()

	_,err3:=collection.ReplaceOne(context.Background(),filter,data)
      if err3!=nil{
		  return nil,status.Errorf(codes.Internal,fmt.Sprintf("Error in updating mongodb %v",err3))
	  }

	return &blogpb.UpdateBlogResponse{
		Blog: &blogpb.Blog{
			Id: data.Id.Hex(),
			AuthorId: data.AuthorId,
			Title: data.Title,
			Content: data.Content,

		},
	},nil

}

func(*server) DeleteBlog(ctx context.Context, in *blogpb.DeleteBlogRequest) (*blogpb.DeleteBlogResponse, error){

	fmt.Println("Delete blog Request")
	oid,err:=primitive.ObjectIDFromHex(in.GetId())
	filter:=bson.M{"_id":oid}
	if err!=nil{
		return nil,status.Errorf(codes.InvalidArgument,fmt.Sprint("Error while converting OID : %v",err))
	}
	res,err1:=collection.DeleteOne(context.Background(),filter)
	if err1!=nil{
		return nil,status.Errorf(codes.Internal,fmt.Sprintf("Error while deleting %v",err))
	}
    
	if res.DeletedCount==0{
		return nil,status.Errorf(codes.Internal,fmt.Sprintf("Cannot find blog in MongoDb %v",err))
	}
	return &blogpb.DeleteBlogResponse{
		Id: in.GetId(),
	},nil

}

func (*server) ListBlog(_ *blogpb.ListBlogRequest, stream blogpb.BlogService_ListBlogServer) error {
	fmt.Println("List blog request")

	cur, err := collection.Find(context.Background(), primitive.D{{}})
	if err != nil {
		return status.Errorf(
			codes.Internal,
			fmt.Sprintf("Unknown internal error: %v", err),
		)
	}
	defer cur.Close(context.Background()) // Should handle err
	for cur.Next(context.Background()) {
		data := &blogItem{}
		err := cur.Decode(data)
		if err != nil {
			return status.Errorf(
				codes.Internal,
				fmt.Sprintf("Error while decoding data from MongoDB: %v", err),
			)

		}
		stream.Send(&blogpb.ListBlogResponse{Blog: dataToBlogPb(data)}) // Should handle err
	}
	if err := cur.Err(); err != nil {
		return status.Errorf(
			codes.Internal,
			fmt.Sprintf("Unknown internal error: %v", err),
		)
	}
	return nil
}
func dataToBlogPb(data *blogItem) *blogpb.Blog {
	return &blogpb.Blog{
		Id:       data.Id.Hex(),
		AuthorId: data.AuthorId,
		Title:    data.Title,
		Content:  data.Content,
	}
}


func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile) //will log the line where error happened if go code fails
	fmt.Println("Recieved create blog request")

	//connect to mongodb
	fmt.Println("Connecting to mongo db")
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Printf("Error connecting to mongodb : %v", err)
	}
	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatalf("Error in connecting to mongo context : %v", err)
	}

	collection = client.Database("mydb").Collection("blog")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatal("Failed to listen", err)
	}

	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)
	blogpb.RegisterBlogServiceServer(s, &server{})

	reflection.Register(s)

	go func() {
		fmt.Println("Starting server")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}

	}()

	//wait for ctrl+c to exit
	ch := make(chan os.Signal, 1)   //making a channel to receive message
	signal.Notify(ch, os.Interrupt) //os.Interrupt is nothing but ctrl+c

	//block until a signal is received
	<-ch

	fmt.Println("Stopping server")
	s.Stop()
	fmt.Println("closing the listener")
	lis.Close()
	fmt.Println("closing mongodb connection")
	client.Disconnect(context.TODO())

	fmt.Println("End of program")

	//just to end program in a better manner

}
