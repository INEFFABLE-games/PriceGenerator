package main

import (
	"PriceGenerator/internal/config"
	"PriceGenerator/internal/generator"
	"PriceGenerator/internal/producer"
	"PriceGenerator/internal/service"
	"context"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
)

func main() {
	cfg := config.NewConfig()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	priceProducer := producer.NewPriceProducer(cfg)
	priceGenerator := generator.NewPriceGenerator()

	priceService := service.NewPriceService(priceProducer, priceGenerator)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		err := priceService.StartStream(ctx)
		if err != nil {
			log.WithFields(log.Fields{
				"handler": "main",
				"action":  "start stream",
			}).Errorf("unable to start stream %v", err.Error())
		}
	}()

	<-c
	cancel()
}
