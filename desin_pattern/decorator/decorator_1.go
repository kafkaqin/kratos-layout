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

	if _, exists := pt.records[method]; !exists {
		pt.records[method] = &PaymentRecord{
			Method: method,
		}
	}

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

// Payment 支付接口
type Payment interface {
	Pay(amount float64) (float64, string)
}

// CreditCardPayment 信用卡支付
type CreditCardPayment struct {
	tracker *PaymentTracker
}

func (c *CreditCardPayment) Pay(amount float64) (float64, string) {
	c.tracker.RecordPayment("Credit Card", amount, 0)
	return amount, fmt.Sprintf("Paid %.2f using Credit Card", amount)
}

// PaymentDecorator 装饰器基类
type PaymentDecorator struct {
	Payment Payment
	tracker *PaymentTracker
}

func (d *PaymentDecorator) Pay(amount float64) (float64, string) {
	return d.Payment.Pay(amount)
}

// CouponPaymentDecorator 优惠券装饰器
type CouponPaymentDecorator struct {
	*PaymentDecorator
	couponDiscount float64
}

func (c *CouponPaymentDecorator) Pay(amount float64) (float64, string) {
	// 应用优惠券
	discountedAmount := amount - c.couponDiscount
	if discountedAmount < 0 {
		discountedAmount = 0
	}

	// 记录优惠券支付信息
	c.tracker.RecordPayment("Coupon", amount-discountedAmount, c.couponDiscount)

	// 调用下一个支付装饰器
	remainingAmount, paymentDetail := c.Payment.Pay(discountedAmount)

	return remainingAmount, fmt.Sprintf("Applied coupon discount of %.2f. %s", c.couponDiscount, paymentDetail)
}

// PointsPaymentDecorator 积分支付装饰器
type PointsPaymentDecorator struct {
	*PaymentDecorator
	points float64
}

func (p *PointsPaymentDecorator) Pay(amount float64) (float64, string) {
	// 计算可使用的积分
	discount := (amount / 100) * 10
	if p.points >= discount {
		p.points -= discount
		remainingAmount := amount - discount

		// 记录积分支付信息
		p.tracker.RecordPayment("Points", amount-remainingAmount, discount)

		// 调用下一个支付装饰器
		finalRemainingAmount, paymentDetail := p.Payment.Pay(remainingAmount)

		return finalRemainingAmount, fmt.Sprintf("Used %.2f points. %s Remaining points: %.2f",
			discount, paymentDetail, p.points)
	}

	// 如果积分不足，直接调用下一个支付装饰器
	return p.Payment.Pay(amount)
}

// 创建支付装饰器构建器
type PaymentDecoratorBuilder struct {
	tracker *PaymentTracker
}

func NewPaymentDecoratorBuilder() *PaymentDecoratorBuilder {
	return &PaymentDecoratorBuilder{
		tracker: NewPaymentTracker(),
	}
}

func (b *PaymentDecoratorBuilder) BuildPaymentChain() Payment {
	// 创建基础信用卡支付
	cardPayment := &CreditCardPayment{tracker: b.tracker}

	// 添加优惠券装饰器
	couponPayment := &CouponPaymentDecorator{
		PaymentDecorator: &PaymentDecorator{
			Payment: cardPayment,
			tracker: b.tracker,
		},
		couponDiscount: 20, // 20元优惠券
	}

	// 添加积分支付装饰器
	pointsPayment := &PointsPaymentDecorator{
		PaymentDecorator: &PaymentDecorator{
			Payment: couponPayment,
			tracker: b.tracker,
		},
		points: 50, // 50积分
	}

	return pointsPayment
}

func main() {
	// 创建支付装饰器构建器
	builder := NewPaymentDecoratorBuilder()

	// 构建支付链
	payment := builder.BuildPaymentChain()

	// 模拟多次支付
	payments := []float64{100, 150, 200}
	for _, amount := range payments {
		fmt.Printf("\nProcessing payment of %.2f:\n", amount)
		remainingAmount, details := payment.Pay(amount)
		fmt.Println("Payment Details:", details)
		fmt.Printf("Remaining Amount: %.2f\n", remainingAmount)
	}

	// 打印支付汇总
	summaryRecords := builder.tracker.GetPaymentSummary()
	fmt.Println("\nPayment Summary:")
	for _, record := range summaryRecords {
		fmt.Printf("%s - Total Paid: %.2f, Total Discount: %.2f\n",
			record.Method, record.AmountPaid, record.TotalDiscount)
	}
}
