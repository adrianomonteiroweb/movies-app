package main

import (
	"context"
	"log"
	"math/rand"
	"net"
	"packagemovies/pb"
	"strconv"

	"google.golang.org/grpc"
)

const (
	PORT = ":50051"
)

var movies []*pb.MovieInfo

type movieServer struct {
	pb.UnimplementedMovieServer
}

func main() {
	listener, err := net.Listen("tcp", PORT)

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	server := grpc.NewServer()

	pb.RegisterMovieServer(server, &movieServer{})

	log.Printf("Server listening at %v", listener.Addr())

	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to server: %v", err)
	}

}

func initMovies() {
	movie1 := &pb.MovieInfo{
		ID: "1",
		Isbn: "0593310438",
		Title: "The Batman",
		Director: &pb.Director{
			Firstname: "Matt",
			Lastname: "Reeves",
		}}
	
	movie2 := &pb.MovieInfo{
		ID: "2",
		Isbn: "3430220302",
		Title: "Doctor Strange in the Multiverse of Madness",
		Director: &pb.Director{
			Firstname: "Sam",
			Lastname: "Raimi",
		}}

	movies = append(movies, movie1)
	movies = append(movies, movie2)
}

func (s *movieServer) GetMovies(in *pb.Empty, stream pb.Movie_GetMoviesServer) error {
	log.Printf("Received: %v", in)

	for _, movie := range movies {
		if err := stream.Send(movie); err != nil {
			return err
		}
	}

	return nil
}

func (s *movieServer) GetMovie(ctx context.Context, in *pb.ID) (*pb.MovieInfo, error) {
	log.Printf("Received: %v", in)

	res := &pb.MovieInfo{}

	for _, movie := range movies {
		if movie.GetID() == in.GetValue() {
			res = movie
			break
		}
	}

	return res, nil
}

func (s *movieServer) CreateMovie(ctx context.Context, in *pb.MovieInfo) (*pb.ID, error) {
	log.Printf("Received: %v", in)

	res := pb.ID{}
	res.Value = strconv.Itoa(rand.Intn(100000000))
	in.ID = res.GetValue()

	movies = append(movies, in)

	return &res, nil
}

func (s *movieServer) UpdateMovie(ctx context.Context, in *pb.MovieInfo) (*pb.Status, error) {
	log.Printf("Received: %v", in)

	res := pb.Status{}

	for index, movie := range movies {
		if movie.GetID() == in.GetID() {
			movies = append(movies[:index], movies[index+1:]...)

			in.ID = movie.GetID()

			movies = append(movies, in)
			res.Value = 1

			break
		}
	}

	return &res, nil
}

func (s *movieServer) DeleteMovie(ctx context.Context, in *pb.ID) (*pb.Status, error) {
	log.Printf("Received: %v", in)

	res := pb.Status{}

	for index, movie := range movies {
		if movie.GetID() == in.GetValue() {
			movies = append(movies[:index], movies[index+1:]...)
			
			res.Value = 1

			break
		}
	}

	return &res, nil
}