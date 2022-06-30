// Package main implements a server for Greeter service.
package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/mattchw/go-onboard/grpc/pb/book"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type BookServiceServer struct {
	book.UnimplementedBookServiceServer
}

type UnimplementedBookServiceServer struct {
}

// Client instance
var DB *mongo.Client = ConnectDB()
var bookCollection *mongo.Collection = GetCollection(DB, "books")

func ConnectDB() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB")
	return client
}

// getting database collections
func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("go-onboard").Collection(collectionName)
	return collection
}

type BookItem struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Title       string             `bson:"title"`
	Description string             `bson:"description"`
}

func (s *BookServiceServer) ReadBook(ctx context.Context, req *book.ReadBookReq) (*book.ReadBookRes, error) {
	// convert string id (from proto) to mongoDB ObjectId
	oid, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Could not convert to ObjectId: %v", err))
	}
	result := bookCollection.FindOne(ctx, bson.M{"_id": oid})
	// Create an empty BookItem to write our decode result to
	data := BookItem{}
	// decode and write to data
	if err := result.Decode(&data); err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Could not find book with Object Id %s: %v", req.GetId(), err))
	}
	// Cast to ReadBookRes type
	response := &book.ReadBookRes{
		Book: &book.Book{
			Id:          oid.Hex(),
			Title:       data.Title,
			Description: data.Description,
		},
	}
	return response, nil
}

func (s *BookServiceServer) CreateBook(ctx context.Context, req *book.CreateBookReq) (*book.CreateBookRes, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	b := req.GetBook()
	payload := BookItem{
		// ID:    Empty, so it gets omitted and MongoDB generates a unique Object ID upon insertion.
		Title:       b.GetTitle(),
		Description: b.GetDescription(),
	}

	result, err := bookCollection.InsertOne(ctx, payload)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v", err),
		)
	}

	oid := result.InsertedID.(primitive.ObjectID)
	b.Id = oid.Hex()

	return &book.CreateBookRes{Book: b}, nil
}

func (s *BookServiceServer) UpdateBook(ctx context.Context, req *book.UpdateBookReq) (*book.UpdateBookRes, error) {
	b := req.GetBook()

	// Convert the Id string to a MongoDB ObjectId
	oid, err := primitive.ObjectIDFromHex(b.GetId())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Could not convert the supplied blog id to a MongoDB ObjectId: %v", err),
		)
	}

	// Convert the data to be updated into an unordered Bson document
	update := bson.M{
		"title":       b.GetTitle(),
		"description": b.GetDescription(),
	}

	// Convert the oid into an unordered bson document to search by id
	filter := bson.M{"_id": oid}

	// Result is the BSON encoded result
	// To return the updated document instead of original we have to add options.
	result := bookCollection.FindOneAndUpdate(ctx, filter, bson.M{"$set": update}, options.FindOneAndUpdate().SetReturnDocument(1))

	// Decode result and write it to 'decoded'
	decoded := BookItem{}
	err = result.Decode(&decoded)
	if err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Could not find blog with supplied ID: %v", err),
		)
	}
	return &book.UpdateBookRes{
		Book: &book.Book{
			Id:          decoded.ID.Hex(),
			Title:       decoded.Title,
			Description: decoded.Description,
		},
	}, nil
}

func (s *BookServiceServer) DeleteBook(ctx context.Context, req *book.DeleteBookReq) (*book.DeleteBookRes, error) {
	oid, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Could not convert to ObjectId: %v", err))
	}
	// DeleteOne returns DeleteResult which is a struct containing the amount of deleted docs (in this case only 1 always)
	// So we return a boolean instead
	_, err = bookCollection.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Could not find/delete blog with id %s: %v", req.GetId(), err))
	}
	return &book.DeleteBookRes{
		Success: true,
	}, nil
}

func (s *BookServiceServer) ListBooks(req *book.ListBooksReq, stream book.BookService_ListBooksServer) error {
	// Initiate a BookItem type to write decoded data to
	data := &BookItem{}
	// collection.Find returns a cursor for our (empty) query
	results, err := bookCollection.Find(context.Background(), bson.M{})
	if err != nil {
		return status.Errorf(codes.Internal, fmt.Sprintf("Unknown internal error: %v", err))
	}
	// An expression with defer will be called at the end of the function
	defer results.Close(context.Background())
	// cursor.Next() returns a boolean, if false there are no more items and loop will break
	for results.Next(context.Background()) {
		// Decode the data at the current pointer and write it to data
		if err = results.Decode(data); err != nil {
			return status.Errorf(codes.Internal, fmt.Sprintf("Unknown cursor error: %v", err))
		} else {
			stream.Send(&book.ListBooksRes{
				Book: &book.Book{
					Id:          data.ID.Hex(),
					Title:       data.Title,
					Description: data.Description,
				},
			})
		}
	}

	return nil
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	fmt.Println("Starting server on port :50051...")

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	srv := &BookServiceServer{}
	book.RegisterBookServiceServer(s, srv)
	reflection.Register(s)
	log.Printf("server listening at %v", lis.Addr())

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()
	fmt.Println("Server successfully started on port :50051")

	// Right way to stop the server using a SHUTDOWN HOOK
	// Create a channel to receive OS signals
	c := make(chan os.Signal)

	// Relay os.Interrupt to our channel (os.Interrupt = CTRL+C)
	// Ignore other incoming signals
	signal.Notify(c, os.Interrupt)

	<-c

	// After receiving CTRL+C Properly stop the server
	fmt.Println("\nStopping the server...")
	s.Stop()
	lis.Close()
	fmt.Println("Done.")
}
