package main

import (
	"currency-checker/internal/client"
	"currency-checker/internal/config"
	"currency-checker/internal/handlers"
	"currency-checker/internal/logger"
	"currency-checker/internal/service"
	"currency-checker/internal/storage"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	cfg := config.GetConfig()
	fmt.Print(cfg)
	log := logger.New()
	strg := storage.Open(cfg.StoragePath)
	coinBaseClient := client.NewCoinBaseClient()
	smtpClient := client.NewSmtpClient(cfg)
	emailService := service.NewEmailService(strg, coinBaseClient, smtpClient)
	handler := handlers.NewHandler(emailService, log)
	router := handler.NewGorillaMux(mux.NewRouter())

	srv := http.Server{
		Addr:    cfg.Listen.BindIP + ":" + cfg.Listen.Port,
		Handler: router,
	}
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
		panic(err.Error())
	}
}
