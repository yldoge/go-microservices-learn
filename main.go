package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"yldoge.com/learn-microservices/handlers"
)

func main() {
	l := log.New(os.Stdout, "product-api | ", log.LstdFlags)

	ph := handlers.NewProducts(l)

	sm := mux.NewRouter()

	getRoute := sm.Methods(http.MethodGet).Subrouter()
	getRoute.HandleFunc("/", ph.GetProducts)

	putRoute := sm.Methods(http.MethodPut).Subrouter()
	putRoute.HandleFunc("/{id:[0-9]+}", ph.UpdateProduct)
	putRoute.Use(ph.MiddlewareProductValidation)

	postRoute := sm.Methods(http.MethodPost).Subrouter()
	postRoute.HandleFunc("/", ph.AddProduct)
	postRoute.Use(ph.MiddlewareProductValidation)

	s := &http.Server{
		Addr:         ":9090",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	l.Println("Received terminate, gracefully shutdown...", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)

}
