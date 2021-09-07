package main

import (
	"PriceGenerator/internal/config"
	"PriceGenerator/internal/generator"
	"PriceGenerator/internal/models"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"time"
)

func getRedis(cfg *config.Config) *redis.Client {

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddres,
		Password: cfg.RedisPassword,
		Username: cfg.RedisUserName,
		DB:       0,
	})

	res, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res)

	return client
}

func main() {

	cfg := config.NewConfig()

	redisClient := getRedis(cfg)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	err := redisClient.Do(ctx, "DEL", "Prices").Err()
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "main",
			"action":  "clear redis db",
		}).Errorf("unable to clear redis db %v", err)
	}

	priceGenerator := generator.NewPriceGenerator()

	//---------------cycle
	ctx = context.Background()

	for {
		batchOfPrices := make([]models.Price, 10)
		for i := 0; i < 10; i++ {

			newPrice := priceGenerator.Generate()

			batchOfPrices[i] = *newPrice

			err = redisClient.Do(ctx, "XADD", "Prices", "*", "Id", newPrice.Id, "Name", newPrice.Name, "Bid", newPrice.Bid, "Ask", newPrice.Ask).Err()
			if err != nil {
				log.WithFields(log.Fields{
					"handler": "main",
					"action":  "Do",
				}).Errorf("unablde to send %v", err)
			}

			fmt.Printf("New price: %v\n", newPrice)
		}

		time.Sleep(time.Second * 1)
	}

}
