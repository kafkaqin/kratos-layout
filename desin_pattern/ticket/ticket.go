package lottery

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

// 你现在是彩票专家，请帮我分析下以下几种彩票的数据结构如何存储比较好扩展，比如说用到那种存储服务器，如何存储，并提供对应的golang代码例子，要用到设计模式，包括存储用户下注的数据，方便查看，存储数据源，出奖结果，要求易于扩展: 1.双色球 2.排列五 3.排列三 4.大乐透 5.任选九 6.胜负彩 7.竞彩篮球 8.竞彩足球 9.北京单场 10.七乐彩 11.快乐8 12.福彩3D
// 彩票类型枚举

// 彩票类型枚举
type LotteryType int

const (
	DoubleBall LotteryType = iota
	ArrangeV5
	ArrangeV3
	SuperLotto
	SelectNine
	FootballLottery
	BasketballLottery
	SingleMatch
	SevenHappy
	Happy8
	Welfare3D
)

// 彩票注单结构
type LotteryTicket struct {
	ID          string       `json:"id" bson:"_id"`
	UserID      string       `json:"user_id" bson:"user_id"`
	LotteryType LotteryType  `json:"lottery_type" bson:"lottery_type"`
	Numbers     [][]int      `json:"numbers" bson:"numbers"`
	BetAmount   float64      `json:"bet_amount" bson:"bet_amount"`
	BetTime     time.Time    `json:"bet_time" bson:"bet_time"`
	Multiple    int          `json:"multiple" bson:"multiple"`
	PlayType    string       `json:"play_type" bson:"play_type"`
	Status      TicketStatus `json:"status" bson:"status"`
}

// 开奖结果结构
type LotteryDrawResult struct {
	ID             string      `json:"id" bson:"_id"`
	LotteryType    LotteryType `json:"lottery_type" bson:"lottery_type"`
	DrawDate       time.Time   `json:"draw_date" bson:"draw_date"`
	WinningNumbers []int       `json:"winning_numbers" bson:"winning_numbers"`
	Jackpot        float64     `json:"jackpot" bson:"jackpot"`
	WinningDetails []WinDetail `json:"winning_details" bson:"winning_details"`
}

type WinDetail struct {
	PrizeLevel  string  `json:"prize_level" bson:"prize_level"`
	WinnerCount int     `json:"winner_count" bson:"winner_count"`
	PrizeAmount float64 `json:"prize_amount" bson:"prize_amount"`
}

type TicketStatus int

const (
	Pending TicketStatus = iota
	Winning
	Lost
	Claimed
)

// 存储策略接口
type StorageStrategy interface {
	SaveTicket(ticket *LotteryTicket) error
	SaveDrawResult(result *LotteryDrawResult) error
	FindTicketsByUser(userID string) ([]LotteryTicket, error)
	FindDrawResultByType(lotteryType LotteryType) ([]LotteryDrawResult, error)
}

// MongoDB存储策略
type MongoDBStorageStrategy struct {
	client     *mongo.Client
	database   *mongo.Database
	ticketColl *mongo.Collection
	resultColl *mongo.Collection
}

func NewMongoDBStorageStrategy(uri string) (*MongoDBStorageStrategy, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	database := client.Database("lottery_db")
	return &MongoDBStorageStrategy{
		client:     client,
		database:   database,
		ticketColl: database.Collection("tickets"),
		resultColl: database.Collection("draw_results"),
	}, nil
}

func (m *MongoDBStorageStrategy) SaveTicket(ticket *LotteryTicket) error {
	_, err := m.ticketColl.InsertOne(context.Background(), ticket)
	return err
}

func (m *MongoDBStorageStrategy) SaveDrawResult(result *LotteryDrawResult) error {
	_, err := m.resultColl.InsertOne(context.Background(), result)
	return err
}

func (m *MongoDBStorageStrategy) FindTicketsByUser(userID string) ([]LotteryTicket, error) {
	var tickets []LotteryTicket
	cursor, err := m.ticketColl.Find(context.Background(), map[string]string{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	err = cursor.All(context.Background(), &tickets)
	return tickets, err
}

func (m *MongoDBStorageStrategy) FindDrawResultByType(lotteryType LotteryType) ([]LotteryDrawResult, error) {
	var results []LotteryDrawResult
	cursor, err := m.resultColl.Find(context.Background(), map[string]LotteryType{"lottery_type": lotteryType})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	err = cursor.All(context.Background(), &results)
	return results, err
}

// Redis缓存策略
type RedisCacheStrategy struct {
	client *redis.Client
}

func NewRedisCacheStrategy(addr string) *RedisCacheStrategy {
	return &RedisCacheStrategy{
		client: redis.NewClient(&redis.Options{
			Addr: addr,
		}),
	}
}

func (r *RedisCacheStrategy) CacheTicket(ticket *LotteryTicket) error {
	data, err := json.Marshal(ticket)
	if err != nil {
		return err
	}
	return r.client.Set(context.Background(), fmt.Sprintf("ticket:%s", ticket.ID), data, 24*time.Hour).Err()
}

// 数据仓库
type LotteryDataWarehouse struct {
	storageStrategy StorageStrategy
	cacheStrategy   *RedisCacheStrategy
	pgxPool         *pgxpool.Pool
}

func NewLotteryDataWarehouse(
	storageStrategy StorageStrategy,
	cacheStrategy *RedisCacheStrategy,
	pgxPool *pgxpool.Pool,
) *LotteryDataWarehouse {
	return &LotteryDataWarehouse{
		storageStrategy: storageStrategy,
		cacheStrategy:   cacheStrategy,
		pgxPool:         pgxPool,
	}
}

// 复杂的查询和统计方法
func (ldw *LotteryDataWarehouse) GetUserLotteryStats(userID string) (map[LotteryType]int, error) {
	// 实现复杂的数据统计逻辑
	return nil, nil
}

// 工厂方法：创建不同类型的彩票
type LotteryFactory struct {
	warehouse *LotteryDataWarehouse
}

func (lf *LotteryFactory) CreateLotteryTicket(
	lotteryType LotteryType,
	userID string,
	numbers [][]int,
) *LotteryTicket {
	return &LotteryTicket{
		ID:          fmt.Sprintf("%d_%s", time.Now().UnixNano(), userID),
		UserID:      userID,
		LotteryType: lotteryType,
		Numbers:     numbers,
		BetTime:     time.Now(),
		Status:      Pending,
	}
}

func main() {
	// 示例使用
	mongoStrategy, _ := NewMongoDBStorageStrategy("mongodb://localhost:27017")
	redisCache := NewRedisCacheStrategy("localhost:6379")
	pgxConfig, _ := pgxpool.ParseConfig("postgres://user:pass@localhost:5432/lottery_db")
	pgxPool, _ := pgxpool.ConnectConfig(context.Background(), pgxConfig)

	dataWarehouse := NewLotteryDataWarehouse(mongoStrategy, redisCache, pgxPool)
	lotteryFactory := &LotteryFactory{warehouse: dataWarehouse}

	// 创建双色球彩票
	ticket := lotteryFactory.CreateLotteryTicket(
		DoubleBall,
		"user123",
		[][]int{{1, 2, 3, 4, 5, 6}, {7}},
	)

	// 保存彩票
	mongoStrategy.SaveTicket(ticket)
	redisCache.CacheTicket(ticket)
}

type LotteryType string

const (
	DoubleBall        LotteryType = "DOUBLE_BALL"
	ArrangeV5         LotteryType = "ARRANGE_V5"
	ArrangeV3         LotteryType = "ARRANGE_V3"
	SuperLotto        LotteryType = "SUPER_LOTTO"
	SelectNine        LotteryType = "SELECT_NINE"
	FootballLottery   LotteryType = "FOOTBALL_LOTTERY"
	BasketballLottery LotteryType = "BASKETBALL_LOTTERY"
	SingleMatch       LotteryType = "SINGLE_MATCH"
	SevenHappy        LotteryType = "SEVEN_HAPPY"
	Happy8            LotteryType = "HAPPY_8"
	Welfare3D         LotteryType = "WELFARE_3D"
)

// 投注类型
type BetType string

const (
	DirectBet   BetType = "DIRECT"
	GroupBet    BetType = "GROUP"
	CombineBet  BetType = "COMBINE"
	SingleMatch BetType = "SINGLE_MATCH"
)

// 奖项等级
type PrizeLevel string

const (
	FirstPrize  PrizeLevel = "FIRST"
	SecondPrize PrizeLevel = "SECOND"
	ThirdPrize  PrizeLevel = "THIRD"
	// 可根据不同彩票类型扩展
)

// 彩票投注记录
type LotteryTicket struct {
	ID          uuid.UUID    `bson:"_id"`
	UserID      uuid.UUID    `bson:"user_id"`
	LotteryType LotteryType  `bson:"lottery_type"`
	BetType     BetType      `bson:"bet_type"`
	Numbers     [][]int      `bson:"numbers"`
	BetAmount   float64      `bson:"bet_amount"`
	Multiple    int          `bson:"multiple"`
	IssueNumber string       `bson:"issue_number"`
	BetTime     time.Time    `bson:"bet_time"`
	Status      TicketStatus `bson:"status"`
}

// 开奖结果
type DrawResult struct {
	ID             uuid.UUID   `bson:"_id"`
	LotteryType    LotteryType `bson:"lottery_type"`
	IssueNumber    string      `bson:"issue_number"`
	DrawTime       time.Time   `bson:"draw_time"`
	WinningNumbers []int       `bson:"winning_numbers"`
	Jackpot        float64     `bson:"jackpot"`
	Prizes         []PrizeInfo `bson:"prizes"`
}

// 奖项信息
type PrizeInfo struct {
	Level       PrizeLevel `bson:"level"`
	WinnerCount int        `bson:"winner_count"`
	PrizeAmount float64    `bson:"prize_amount"`
}

// 彩票存储接口
type LotteryRepository interface {
	SaveTicket(ticket *LotteryTicket) error
	SaveDrawResult(result *DrawResult) error
	GetTicketsByUser(userID uuid.UUID) ([]LotteryTicket, error)
	GetDrawResultByIssue(lotteryType LotteryType, issueNumber string) (*DrawResult, error)
}

// MongoDB实现
type MongoLotteryRepository struct {
	client           *mongo.Client
	database         *mongo.Database
	ticketCollection *mongo.Collection
	resultCollection *mongo.Collection
	logger           *zap.Logger
}

// 创建MongoDB仓库
func NewMongoLotteryRepository(uri string, logger *zap.Logger) (*MongoLotteryRepository, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	database := client.Database("lottery_system")

	return &MongoLotteryRepository{
		client:           client,
		database:         database,
		ticketCollection: database.Collection("lottery_tickets"),
		resultCollection: database.Collection("lottery_results"),
		logger:           logger,
	}, nil
}

func (r *MongoLotteryRepository) SaveTicket(ticket *LotteryTicket) error {
	_, err := r.ticketCollection.InsertOne(context.Background(), ticket)
	if err != nil {
		r.logger.Error("Failed to save ticket", zap.Error(err))
	}
	return err
}

func (r *MongoLotteryRepository) SaveDrawResult(result *DrawResult) error {
	_, err := r.resultCollection.InsertOne(context.Background(), result)
	if err != nil {
		r.logger.Error("Failed to save draw result", zap.Error(err))
	}
	return err
}

func (r *MongoLotteryRepository) GetTicketsByUser(userID uuid.UUID) ([]LotteryTicket, error) {
	var tickets []LotteryTicket
	filter := bson.M{"user_id": userID}

	cursor, err := r.ticketCollection.Find(context.Background(), filter)
	if err != nil {
		r.logger.Error("Failed to find tickets", zap.Error(err))
		return nil, err
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &tickets); err != nil {
		r.logger.Error("Failed to decode tickets", zap.Error(err))
		return nil, err
	}

	return tickets, nil
}

func (r *MongoLotteryRepository) GetDrawResultByIssue(lotteryType LotteryType, issueNumber string) (*DrawResult, error) {
	var result DrawResult
	filter := bson.M{
		"lottery_type": lotteryType,
		"issue_number": issueNumber,
	}

	err := r.resultCollection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		r.logger.Error("Failed to find draw result", zap.Error(err))
		return nil, err
	}

	return &result, nil
}

// 彩票服务
type LotteryService struct {
	repository LotteryRepository
	logger     *zap.Logger
}

func NewLotteryService(repository LotteryRepository, logger *zap.Logger) *LotteryService {
	return &LotteryService{
		repository: repository,
		logger:     logger,
	}
}

func (s *LotteryService) PlaceBet(userID uuid.UUID, lotteryType LotteryType, numbers [][]int, betAmount float64) (*LotteryTicket, error) {
	ticket := &LotteryTicket{
		ID:          uuid.New(),
		UserID:      userID,
		LotteryType: lotteryType,
		Numbers:     numbers,
		BetAmount:   betAmount,
		BetTime:     time.Now(),
		IssueNumber: fmt.Sprintf("%d", time.Now().Unix()), // 简单的期号生成
		Status:      Pending,
	}

	err := s.repository.SaveTicket(ticket)
	if err != nil {
		s.logger.Error("Failed to place bet", zap.Error(err))
		return nil, err
	}

	return ticket, nil
}

func (s *LotteryService) RecordDrawResult(lotteryType LotteryType, winningNumbers []int, prizes []PrizeInfo) (*DrawResult, error) {
	result := &DrawResult{
		ID:             uuid.New(),
		LotteryType:    lotteryType,
		DrawTime:       time.Now(),
		WinningNumbers: winningNumbers,
		IssueNumber:    fmt.Sprintf("%d", time.Now().Unix()),
		Prizes:         prizes,
	}

	err := s.repository.SaveDrawResult(result)
	if err != nil {
		s.logger.Error("Failed to record draw result", zap.Error(err))
		return nil, err
	}

	return result, nil
}

func main() {
	// 初始化日志
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// 创建MongoDB仓库
	repository, err := NewMongoLotteryRepository("mongodb://localhost:27017", logger)
	if err != nil {
		logger.Fatal("Failed to create repository", zap.Error(err))
	}

	// 创建彩票服务
	lotteryService := NewLotteryService(repository, logger)

	// 用户下注示例
	userID := uuid.New()
	ticket, err := lotteryService.PlaceBet(
		userID,
		DoubleBall,
		[][]int{{1, 2, 3, 4, 5, 6}, {7}},
		100.00,
	)

	// 记录开奖结果示例
	_, err = lotteryService.RecordDrawResult(
		DoubleBall,
		[]int{1, 2, 3, 4, 5, 6},
		[]PrizeInfo{
			{
				Level:       FirstPrize,
				WinnerCount: 1,
				PrizeAmount: 5000000,
			},
		},
	)
}
