package generator

import (
	"PriceGenerator/models"
	"math/rand"
	"time"
)

var names = []string{"Apple","Amazon","Google","Effective Soft","GitHub","Electronic Arts","Xiaomi","Tesla","Disney","Sony"}

const (
	maxPrice = 10000
	minPrice = 1000
)

type PriceGenerator struct {
}

func (p *PriceGenerator) Generate()*models.Price{
	price := models.Price{
		Id:   time.Stamp,
		Name: names[rand.Intn(10-1) + 1],
		Bid:  uint64(rand.Intn(maxPrice-minPrice) + maxPrice),
		Ask:  uint64(rand.Intn(maxPrice-minPrice) + maxPrice),
	}
	return &price
}

func NewPriceGenerator() *PriceGenerator{
	return &PriceGenerator{}
}