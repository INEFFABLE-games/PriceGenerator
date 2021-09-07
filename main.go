package main

import (
	"PriceGenerator/generator"
	log "github.com/sirupsen/logrus"
	"time"
)

func main(){

	priceGenerator := generator.NewPriceGenerator()

	for  {
		newPrice := priceGenerator.Generate()

		log.Println(newPrice)

		time.Sleep(time.Second*1)
	}

}
