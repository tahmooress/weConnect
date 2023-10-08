package main

import (
	"context"
	"fmt"
	"log"

	"github.com/tahmooress/weConnect-task/cmd"
	"github.com/tahmooress/weConnect-task/internal/api"
	"github.com/tahmooress/weConnect-task/internal/repository/mongodb"
	"github.com/tahmooress/weConnect-task/internal/service"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	repo, err := mongodb.New(ctx)
	if err != nil {
		cancel()
		log.Fatal(err)
	}

	usecase := service.New(repo)

	fmt.Println("here")

	server, err := api.NewHTTPServer(usecase)
	if err != nil {
		cancel()
		repo.Close()
		log.Fatal(err)
	}

	fmt.Println("here")

	cmd.Shutdown(ctx, cancel, repo, server)
}
