package service

import (
	"PriceGenerator/internal/generator"
	"PriceGenerator/internal/models"
	"PriceGenerator/internal/producer"
	"context"
	log "github.com/sirupsen/logrus"
	"time"
)

// PriceService structure for PriceService object
type PriceService struct {
	producer  *producer.PriceProducer
	generator *generator.PriceGenerator
}

func contains(s []models.Price, e models.Price) bool {
	for _, a := range s {
		if a.Name == e.Name {
			return true
		}
	}
	return false
}

// StartStream is starts infinity cycle to send butches of prices in redis stream
func (p *PriceService) StartStream(ctx context.Context) error {

	err := p.producer.ClearDatabase(ctx)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "main",
			"action":  "clear redis database",
		}).Errorf("unable to clear redis database %v", err.Error())
	}

	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-ctx.Done():
			return nil
		case t := <-ticker.C:
			var butchOfPrices []models.Price

			for i := 0; i < 20; i++ {
				newPrice := p.generator.Generate()
				if contains(butchOfPrices, *newPrice) {
					continue
				}
				butchOfPrices = append(butchOfPrices, *newPrice)
			}

			err := p.producer.SendBatchOfPrices(ctx, butchOfPrices)
			if err != nil {
				log.WithFields(log.Fields{
					"handler": "PriceService",
					"action":  "SendBatchOfPrices",
				}).Errorf("unable to send butch of prices %v", err.Error())
			}

			log.Infof("Tick at %v", t.UTC())
		}
	}
}

// NewPriceService creates new object of PriceService
func NewPriceService(producer *producer.PriceProducer, generator *generator.PriceGenerator) *PriceService {
	return &PriceService{producer: producer, generator: generator}
}
