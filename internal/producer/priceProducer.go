package producer

import (
	"PriceGenerator/internal/config"
	"PriceGenerator/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
)

// PriceProducer is structure for PriceProducer object
type PriceProducer struct {
	client *redis.Client
}

// ClearDatabase send clear command to redis database
func (p *PriceProducer) ClearDatabase(ctx context.Context) error {
	return p.client.Do(ctx, "DEL", "Prices").Err()
}

// SendBatchOfPrices send butch of prices on redis database
func (p *PriceProducer) SendBatchOfPrices(ctx context.Context, newPrices []models.Price) error {

	var err error

	convertedPrice, err := json.Marshal(newPrices)
	if err != nil {
		return err
	}

	err = p.client.XAdd(ctx, &redis.XAddArgs{
		Stream:     "Prices",
		NoMkStream: false,
		MaxLen:     0,
		MinID:      "",
		Approx:     false,
		Limit:      0,
		ID:         "",
		Values:     map[string]interface{}{"butch": convertedPrice},
	}).Err()
	if err != nil {
		return err
	}

	for _, price := range newPrices {
		log.WithFields(log.Fields{
			"Price name ": price.Name,
			"Price bid ":  price.Bid,
			"Price ask ":  price.Ask,
		}).Infof("New price sended [%v]", price.Id)
	}

	return err
}

// getRedis returns new redis client connection
func getRedis(cfg *config.Config) *redis.Client {

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddres,
		Password: cfg.RedisPassword,
		Username: cfg.RedisUserName,
		DB:       0,
	})

	res, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Errorf("unable ti ping redis connection %v", err.Error())
	}

	fmt.Println(res)

	return client
}

// NewPriceProducer creates new PriceProducer object
func NewPriceProducer(cfg *config.Config) *PriceProducer {
	client := getRedis(cfg)
	return &PriceProducer{client: client}
}
