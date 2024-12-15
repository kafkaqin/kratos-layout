package main

import "fmt"

// 通用接口
type BetStrategy interface {
	CalculatePayout(betAmount float64, odds float64) float64
}

// 竞彩足球策略
type FootballBet struct{}

func (f *FootballBet) CalculatePayout(betAmount, odds float64) float64 {
	return betAmount * odds
}

// 竞彩篮球策略
type BasketballBet struct{}

func (b *BasketballBet) CalculatePayout(betAmount, odds float64) float64 {
	return betAmount * odds * 0.95 // 示例：篮球有抽成
}

// 上下文类
type BetContext struct {
	strategy BetStrategy
}

func (c *BetContext) SetStrategy(strategy BetStrategy) {
	c.strategy = strategy
}

func (c *BetContext) Calculate(betAmount, odds float64) float64 {
	return c.strategy.CalculatePayout(betAmount, odds)
}

// 使用示例
func main() {
	context := &BetContext{}

	// 竞彩足球
	context.SetStrategy(&FootballBet{})
	fmt.Println("Football Payout:", context.Calculate(100, 2.0))

	// 竞彩篮球
	context.SetStrategy(&BasketballBet{})
	fmt.Println("Basketball Payout:", context.Calculate(100, 2.0))
}
