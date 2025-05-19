package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"google.golang.org/grpc/reflection"
	"github.com/captainhbb/movieexample-rating/internal/controller"
	"github.com/captainhbb/movieexample-rating/internal/ingester/kafka"
	"github.com/captainhbb/movieexample-rating/internal/repository/mysql"

	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"
	"github.com/captainhbb/movieexample-protoapis/gen"
	"github.com/captainhbb/movieexample-discovery/pkg/discovery"
	"github.com/captainhbb/movieexample-discovery/pkg/discovery/consul"

	grpchandler "github.com/captainhbb/movieexample-rating/internal/handler/grpc"
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
		panic(err)
	}

	ingester, err := kafka.NewIngester("localhost", "rating", "ratings")
	if err != nil {
		log.Fatalf("failed to initialize ingester: %v", err)
	}

	ctrl := controller.New(repo, ingester)
	err = ctrl.StartIngestion(ctx)
	if err != nil {
		log.Fatalf("failed to start ingestion: %v", err)
	}

	h := grpchandler.New(ctrl)

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", cfg.APIConfig.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer()
	reflection.Register(srv)
	gen.RegisterRatingServiceServer(srv, h)
	if err := srv.Serve(lis); err != nil {
		panic(err)
	}
}
