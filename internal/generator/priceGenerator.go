package generator

import (
	"PriceGenerator/internal/models"
	"fmt"
	"math/rand"
	"time"
)

var names = []string{"Apple", "Amazon", "Google", "Effective Soft", "GitHub", "Electronic Arts", "Xiaomi", "Tesla", "Disney", "Sony"}

const (
	maxPrice = 10000
	minPrice = 1000
)

// PriceGenerator is structure for PriceGenerator objects
type PriceGenerator struct {
}

// Generate is generate and return new price object with random values
func (p *PriceGenerator) Generate() *models.Price {
	price := models.Price{
		Id:   fmt.Sprintf("%d", time.Now().Unix()),
		Name: names[rand.Intn(10-1)+1],
		Bid:  uint64(rand.Intn(maxPrice-minPrice) + maxPrice),
		Ask:  uint64(rand.Intn(maxPrice-minPrice) + maxPrice),
	}
	return &price
}

// NewPriceGenerator create new PriceGenerator object
func NewPriceGenerator() *PriceGenerator {
	return &PriceGenerator{}
}
