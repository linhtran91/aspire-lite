package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"aspire-lite/config"
	"aspire-lite/internals/handlers"
	"aspire-lite/internals/must"
	"aspire-lite/internals/repositories"

	"github.com/gorilla/mux"
)

func main() {
	cfgPath, err := config.ParseFlags()
	if err != nil {
		log.Fatal(err)
	}
	cfg, err := config.NewConfig(cfgPath)
	if err != nil {
		log.Fatal(err)
	}

	db, err := must.Connect("postgres", cfg.BuildDSN())
	if err != nil {
		panic(err)
	}
	repaymentRepository := repositories.NewRepayment(db)
	loanRepository := repositories.NewLoan(db)

	repaymentHandler := handlers.NewRepayment(repaymentRepository)
	loanHandler := handlers.NewLoan(loanRepository)
	router := mux.NewRouter()
	router.HandleFunc("/api/customers/{customer_id}/loans", loanHandler.CreateLoan).Methods("POST")
	router.HandleFunc("/api/customers/{customer_id}/loans", loanHandler.List).Methods("GET")
	router.HandleFunc("/api/loans/{loan_id}/approve", loanHandler.ApproveLoan).Methods("PUT")
	router.HandleFunc("/api/repayments/{repayment_id}", repaymentHandler.SubmitRepay).Methods("PUT")

	run(cfg, router)
}

func run(cfg *config.Config, router http.Handler) {
	ch := make(chan os.Signal, 1)
	server := &http.Server{
		Addr:    cfg.Server.Host + ":" + cfg.Server.Port,
		Handler: router,
	}
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	log.Printf("Server is starting on %s\n", server.Addr)

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("Server failed to start due to err: %v", err)
		}
	}()

	interrupt := <-ch
	log.Printf("Server is shutting down due to %+v\n", interrupt)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server was unable to gracefully shutdown due to err: %+v", err)
	}
}
