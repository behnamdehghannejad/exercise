package main

import "fmt"

func CalculateProductFinalPrice() {
	products := []Product{
		{Name: "Laptop", BasePrice: 999.99},
		{Name: "Mouse", BasePrice: 29.99},
		{Name: "Keyboard", BasePrice: 79.99},
	}

	noDiscount := NoDiscount{}
	percentDiscount := PercentageDiscount{Percent: 10}
	fixedDiscount := FixedDiscount{Discount: 50}

	fmt.Println("=== No CalculateProductFinalPrice ===")
	total := CalculateTotal(products, noDiscount)
	fmt.Printf("Total: %.2f\n", total)

	fmt.Println("")
	fmt.Println("=== 10% CalculateProductFinalPrice ===")
	total = CalculateTotal(products, percentDiscount)
	fmt.Printf("Total: %.2f\n", total)

	fmt.Println("")
	fmt.Println("=== $50 Fixed CalculateProductFinalPrice ===")
	total = CalculateTotal(products, fixedDiscount)
	fmt.Printf("Total: %.2f\n", total)
}

type Discounter interface {
	ApplyDiscount(price float64) float64
}

type PercentageDiscount struct {
	Percent float64
}

type FixedDiscount struct {
	Discount float64
}

type NoDiscount struct {
}

func (p PercentageDiscount) ApplyDiscount(price float64) float64 {
	return price - price*(1-p.Percent/100)
}

func (f FixedDiscount) ApplyDiscount(price float64) float64 {
	return price - f.Discount
}

func (NoDiscount) ApplyDiscount(price float64) float64 {
	return price
}

type Product struct {
	Name      string
	BasePrice float64
}

func (p Product) FinalPrice(dis Discounter) float64 {
	finalPrice := dis.ApplyDiscount(p.BasePrice)
	if finalPrice < p.BasePrice {
		finalPrice = p.BasePrice
	}
	fmt.Printf("CalculateProductFinalPrice: %.2f\n", finalPrice)
	return finalPrice
}

func CalculateTotal(products []Product, dis Discounter) float64 {
	total := 0.0
	for _, p := range products {
		total += p.FinalPrice(dis)
	}
	return total
}
