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

// StartStream is starts infinity cycle to send butches of prices in redis stream
func (p *PriceService) StartStream(ctx context.Context) error {

	err := p.producer.ClearDatabase(ctx)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "main",
			"action":  "clear redis database",
		}).Errorf("unable to clear redis database %v", err.Error())
	}

	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ctx.Done():
			return nil
		case t := <-ticker.C:
			var butchOfPrices = make([]models.Price, 10)

			for i := 0; i < 10; i++ {
				newPrice := p.generator.Generate()
				butchOfPrices[i] = *newPrice
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
