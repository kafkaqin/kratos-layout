package main

import (
	"fmt"
	"sync"
)

// PaymentRecord 记录每种支付方式的支付详情
type PaymentRecord struct {
	Method        string
	AmountPaid    float64
	TotalDiscount float64
}

// PaymentTracker 负责追踪支付记录
type PaymentTracker struct {
	records map[string]*PaymentRecord
	mu      sync.RWMutex
}

// NewPaymentTracker 创建新的支付追踪器
func NewPaymentTracker() *PaymentTracker {
	return &PaymentTracker{
		records: make(map[string]*PaymentRecord),
	}
}

// RecordPayment 记录支付信息
func (pt *PaymentTracker) RecordPayment(method string, amountPaid, discount float64) {
	pt.mu.Lock()
	defer pt.mu.Unlock()

	// 如果记录不存在，创建新记录
	if _, exists := pt.records[method]; !exists {
		pt.records[method] = &PaymentRecord{
			Method: method,
		}
	}

	// 更新记录
	record := pt.records[method]
	record.AmountPaid += amountPaid
	record.TotalDiscount += discount
}

// GetPaymentSummary 获取支付方式的汇总信息
func (pt *PaymentTracker) GetPaymentSummary() []PaymentRecord {
	pt.mu.RLock()
	defer pt.mu.RUnlock()

	summaryRecords := make([]PaymentRecord, 0, len(pt.records))
	for _, record := range pt.records {
		summaryRecords = append(summaryRecords, *record)
	}
	return summaryRecords
}

// PaymentStrategy 定义支付策略接口
type PaymentStrategy interface {
	Pay(amount float64) (float64, string, float64)
	GetName() string
}

// CreditCardStrategy 信用卡支付策略
type CreditCardStrategy struct{}

func (c *CreditCardStrategy) Pay(amount float64) (float64, string, float64) {
	return amount, fmt.Sprintf("Paid %.2f using Credit Card", amount), 0
}

func (c *CreditCardStrategy) GetName() string {
	return "Credit Card"
}

// CouponStrategy 优惠券支付策略
type CouponStrategy struct {
	couponDiscount float64
}

func (c *CouponStrategy) Pay(amount float64) (float64, string, float64) {
	discountedAmount := amount - c.couponDiscount
	if discountedAmount < 0 {
		discountedAmount = 0
	}
	return discountedAmount,
		fmt.Sprintf("Applied coupon discount of %.2f. Remaining amount: %.2f", c.couponDiscount, discountedAmount),
		c.couponDiscount
}

func (c *CouponStrategy) GetName() string {
	return "Coupon"
}

// PointsStrategy 积分支付策略
type PointsStrategy struct {
	points float64
}

func (p *PointsStrategy) Pay(amount float64) (float64, string, float64) {
	// 计算可使用的积分（每100元抵10元）
	discount := (amount / 100) * 10
	if p.points >= discount {
		p.points -= discount
		remainingAmount := amount - discount
		return remainingAmount,
			fmt.Sprintf("Used %.2f points. Remaining amount: %.2f. Remaining points: %.2f", discount, remainingAmount, p.points),
			discount
	}
	return amount, "Insufficient points to apply.", 0
}

func (p *PointsStrategy) GetName() string {
	return "Points"
}

// PaymentContext 支付上下文
type PaymentContext struct {
	strategies     []PaymentStrategy
	paymentTracker *PaymentTracker
}

// NewPaymentContext 创建新的支付上下文
func NewPaymentContext() *PaymentContext {
	return &PaymentContext{
		paymentTracker: NewPaymentTracker(),
	}
}

// AddStrategy 添加支付策略
func (pc *PaymentContext) AddStrategy(strategy PaymentStrategy) {
	pc.strategies = append(pc.strategies, strategy)
}

// Pay 执行支付流程
func (pc *PaymentContext) Pay(amount float64) (float64, string) {
	remainingAmount := amount
	var paymentDetails string
	var totalDiscount float64

	// 按优先级依次应用支付策略
	for _, strategy := range pc.strategies {
		var detail string
		var discount float64
		remainingAmount, detail, discount = strategy.Pay(remainingAmount)
		paymentDetails += detail + " "
		totalDiscount += discount

		// 记录支付信息
		pc.paymentTracker.RecordPayment(strategy.GetName(), amount-remainingAmount, discount)

		// 如果金额已经降为0，则停止继续支付
		if remainingAmount <= 0 {
			break
		}
	}

	return remainingAmount, paymentDetails
}

// PrintPaymentSummary 打印支付汇总信息
func (pc *PaymentContext) PrintPaymentSummary() {
	summaryRecords := pc.paymentTracker.GetPaymentSummary()
	fmt.Println("\nPayment Summary:")
	for _, record := range summaryRecords {
		fmt.Printf("%s - Total Paid: %.2f, Total Discount: %.2f\n",
			record.Method, record.AmountPaid, record.TotalDiscount)
	}
}

func main() {
	// 创建支付上下文
	paymentContext := NewPaymentContext()

	// 添加支付策略（按优先级）
	paymentContext.AddStrategy(&CouponStrategy{couponDiscount: 20}) // 优先使用优惠券
	paymentContext.AddStrategy(&PointsStrategy{points: 50})         // 其次使用积分
	paymentContext.AddStrategy(&CreditCardStrategy{})               // 最后使用信用卡

	// 模拟多次支付
	payments := []float64{100, 150, 200}
	for _, amount := range payments {
		fmt.Printf("\nProcessing payment of %.2f:\n", amount)
		remainingAmount, details := paymentContext.Pay(amount)
		fmt.Println("Payment Details:", details)
		fmt.Printf("Remaining Amount: %.2f\n", remainingAmount)
	}

	// 打印支付汇总信息
	paymentContext.PrintPaymentSummary()
}

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
