package factory

import "fmt"

type Bet interface {
	PlaceBet(amount float64) string
}

// 竞彩足球
type FootballBet struct{}

func (f *FootballBet) PlaceBet(amount float64) string {
	return fmt.Sprintf("Football bet placed: $%.2f", amount)
}

// 竞彩篮球
type BasketballBet struct{}

func (b *BasketballBet) PlaceBet(amount float64) string {
	return fmt.Sprintf("Basketball bet placed: $%.2f", amount)
}

// 工厂类
type BetFactory struct{}

func (f *BetFactory) CreateBet(betType string) Bet {
	switch betType {
	case "football":
		return &FootballBet{}
	case "basketball":
		return &BasketballBet{}
	default:
		return nil
	}
}

// 使用示例
func main() {
	factory := &BetFactory{}

	footballBet := factory.CreateBet("football")
	fmt.Println(footballBet.PlaceBet(100))

	basketballBet := factory.CreateBet("basketball")
	fmt.Println(basketballBet.PlaceBet(150))
}
