package logic

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"

	"github.com/hsn0918/BigMarket/pkg/util"

	"github.com/ahmetb/go-linq/v3"
	"github.com/hsn0918/BigMarket/app/strategy/raffle/cmd/api/internal/model"
	"github.com/hsn0918/BigMarket/app/strategy/raffle/cmd/api/internal/svc"
	"github.com/hsn0918/BigMarket/app/strategy/raffle/cmd/api/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

const (
	cacheStrategyAwardCountKey = "strategy#award#%d#count#%d"
	cacheStrategyRateRangeKey  = "big#market#strategy#rate#range#key#%d#%d"
	cacheStrategyRateRange     = "big#market#strategy#rate#range#%d"
)

type StrategyArmoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStrategyArmoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StrategyArmoryLogic {
	return &StrategyArmoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StrategyArmoryLogic) StrategyArmory(req *types.StrategyArmoryRequest) (resp *types.StrategyArmoryResponse, err error) {
	// 1.查询策略配置
	var StrategyAwardList []*model.StrategyAward
	StrategyAwardList, err = l.svcCtx.StrategyAwardModel.QueryStrategyAwardList(l.ctx, req.StrategyId)
	// 2 缓存奖品库存【用于decr扣减库存使用】
	for _, v := range StrategyAwardList {
		err = l.cacheStrategyAwardCount(req.StrategyId, int(v.AwardId), int(v.AwardCount))
	}
	// 3.1 默认装配配置【全量抽奖概率】
	l.assembleLotteryStrategy(req.StrategyId, StrategyAwardList)
	// 3.2 权重策略配置 - 适用于 rule_weight 权重规则配置【4000:102,103,104,105 5000:102,103,104,105,106,107 6000:102,103,104,105,106,107,108,109】
	// todo:3.2
	// 业务异常，策略规则中 rule_weight 权重规则已适用但未配置
	resp = &types.StrategyArmoryResponse{IsSuccess: true}
	return
}
func (l *StrategyArmoryLogic) assembleLotteryStrategy(id int64, list []*model.StrategyAward) (IsSuccess bool, err error) {
	// 1. 获得最小概率值
	var minAward float64
	minAward = linq.From(list).
		SelectT(func(a *model.StrategyAward) float64 { return a.AwardRate }).
		Min().(float64)

	// 2. 循环计算找到概率范围值
	rateRange := util.Convert(minAward)
	// 3. 生成策略奖品概率查找表「这里指需要在list集合中，存放上对应的奖品占位即可，占位越多等于概率越高」
	strategyAwardSearchRateTables := make([]int64, int(rateRange))
	currentIndex := 0
	for _, v := range list {
		awardRate := int(v.AwardRate * rateRange) // 将概率转化为表中的索引范围
		for i := 0; i < awardRate && currentIndex < int(rateRange); i++ {
			strategyAwardSearchRateTables[currentIndex] = v.AwardId
			currentIndex++
		}
	}
	// 4. 对存储的奖品进行乱序操作
	rand.Shuffle(len(strategyAwardSearchRateTables), func(i, j int) {
		strategyAwardSearchRateTables[i], strategyAwardSearchRateTables[j] = strategyAwardSearchRateTables[j], strategyAwardSearchRateTables[i]
	})
	// 5. 生成出Map集合，key值，对应的就是后续的概率值。通过概率来获得对应的奖品ID
	shuffleStrategyAwardSearchRateTable := make(map[int]int64)
	for i, id := range strategyAwardSearchRateTables {
		shuffleStrategyAwardSearchRateTable[i] = id
	}
	// 6. 存放到 Redis
	err = l.storeStrategyAwardSearchRateTable(id, len(shuffleStrategyAwardSearchRateTable), shuffleStrategyAwardSearchRateTable)
	if err != nil {
		return false, err
	}
	return
}
func (l *StrategyArmoryLogic) cacheStrategyAwardCount(StrategyId int64, AwardId int, AwardCount int) (err error) {
	cacheKey := fmt.Sprintf(cacheStrategyAwardCountKey, StrategyId, AwardId)
	err = l.svcCtx.BizRedis.SetCtx(l.ctx, cacheKey, strconv.Itoa(AwardCount))
	if err != nil {
		return err
	}
	return nil
}
func (l *StrategyArmoryLogic) storeStrategyAwardSearchRateTable(id int64, size int, tables map[int]int64) (err error) {
	cacheKey := fmt.Sprintf(cacheStrategyRateRange, id)
	l.svcCtx.BizRedis.SetCtx(l.ctx, cacheKey, strconv.Itoa(size))
	for i, v := range tables {
		err = l.svcCtx.BizRedis.SetCtx(l.ctx, fmt.Sprintf(cacheStrategyRateRangeKey, id, i), strconv.Itoa(int(v)))
		if err != nil {
			return err
		}
	}
	return
}
