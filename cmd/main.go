package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/captainhbb/movieexample-protoapis/gen"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"

	"github.com/captainhbb/movieexample-discovery/pkg/discovery/consul"
	"github.com/captainhbb/movieexample-discovery/pkg/discovery"
	"github.com/captainhbb/movieexample-rating/internal/controller/rating"
	grpchandler "github.com/captainhbb/movieexample-rating/internal/handler/grpc"
	"github.com/captainhbb/movieexample-rating/internal/repository/mysql"
)

const serviceName = "rating"

func main() {

	f, err := os.Open("base.yaml")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var cfg serviceConfig
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		panic(err)
	}

	registry, err := consul.NewRegistry("localhost:8500")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	hostPort := fmt.Sprintf("localhost:%v", cfg.APIConfig.Port)
	instanceID := discovery.GenerateInstanceID(serviceName)

	err = registry.Register(ctx, instanceID, serviceName, hostPort)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
				log.Println("Failed to report healthy state: " + err.Error())
			}

			time.Sleep(1 * time.Second)
		}
	}()
	defer registry.Deregister(ctx, instanceID, serviceName)

	log.Printf("Starting the rating service, listening on %v\n", cfg.APIConfig.Port)
	repo, err := mysql.New()
	if err != nil {
		panic("can't initialize mysql db object")
	}
	ctrl := rating.New(repo)
	h := grpchandler.New(ctrl)

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", cfg.APIConfig.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer()
	gen.RegisterRatingServiceServer(srv, h)
	srv.Serve(lis)
}
