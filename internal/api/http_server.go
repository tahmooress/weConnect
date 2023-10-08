package api

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/tahmooress/weConnect-task/internal/api/handler"
	"github.com/tahmooress/weConnect-task/internal/service"
)

const (
	defaultReadTimeout  = 30 * time.Second
	defaultWriteTimeout = 120 * time.Second
)

var ErrEmptyPortOrIP = errors.New("ip or port is not set")

func NewHTTPServer(usecase service.UseCase) (
	io.Closer, error,
) {
	ip := os.Getenv("IP")
	port := os.Getenv("PORT")

	if ip == "" || port == "" {
		return nil, ErrEmptyPortOrIP
	}

	handler := handler.New(usecase)

	router := mux.NewRouter()

	router.HandleFunc("/statistics/{id}", handler.GetByID()).Methods(http.MethodGet)
	router.HandleFunc("/statistics", handler.GetAll()).Methods(http.MethodGet)
	router.HandleFunc("/statistics", handler.Create()).Methods(http.MethodPost)
	router.HandleFunc("/statistics/{id}", handler.Delete()).Methods(http.MethodDelete)

	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", "localhost", 6060),
		Handler:      router,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			fmt.Printf("handler: ListenAndServe() error: %s\n", err)
		}
		fmt.Println("server run on port: ", port, "and ip: ", ip)
	}()

	return srv, nil
}
