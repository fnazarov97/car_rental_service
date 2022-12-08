package main

import (
	"car_rental/config"
	"car_rental/genprotos/brand"
	"car_rental/genprotos/car"
	bService "car_rental/services/barnd"
	cService "car_rental/services/car"
	"car_rental/storage"
	"car_rental/storage/postgres"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	conf := config.Load()
	AUTH := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		conf.PostgresHost,
		conf.PostgresPort,
		conf.PostgresUser,
		conf.PostgresPassword,
		conf.PostgresDatabase,
	)
	var inter storage.StorageI
	inter, err := postgres.InitDB(AUTH)
	if err != nil {
		panic(err)
	}

	log.Printf("\ngRPC server running port%s with tcp protocol!\n", conf.GRPCPort)

	listener, err := net.Listen("tcp", conf.GRPCPort)
	if err != nil {
		panic(err)
	}

	c := &bService.BrandService{
		Stg: inter,
	}
	s := grpc.NewServer()
	brand.RegisterBrandServiceServer(s, c)

	c1 := &cService.CarService{
		Stg: inter,
	}
	car.RegisterCarServiceServer(s, c1)
	reflection.Register(s)

	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
