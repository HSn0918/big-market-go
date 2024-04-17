package logic

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"strconv"

	"github.com/hsn0918/BigMarket/app/strategy/raffle/cmd/api/code"
	"github.com/hsn0918/BigMarket/pkg/xcode"

	"github.com/hsn0918/BigMarket/common"

	"github.com/ahmetb/go-linq/v3"
	"github.com/hsn0918/BigMarket/pkg/util"

	"github.com/hsn0918/BigMarket/app/strategy/raffle/cmd/api/internal/model"
	"github.com/hsn0918/BigMarket/app/strategy/raffle/cmd/api/internal/svc"
	"github.com/hsn0918/BigMarket/app/strategy/raffle/cmd/api/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
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
		if err != nil {
			logx.Error("cacheStrategyAwardCount error:", err)
		}
	}
	// 3.1 默认装配配置【全量抽奖概率】
	l.assembleLotteryStrategy(req.StrategyId, StrategyAwardList)
	resp = &types.StrategyArmoryResponse{IsSuccess: true}
	// 3.2 装配权重策略配置
	strategy, err := l.svcCtx.StrategyModel.QueryStrategy(l.ctx, req.StrategyId)
	if err != nil {
		// 3.2.1 "rule_weight"规则不存在，则直接返回
		if errors.Is(err, model.ErrNotFound) {
			return resp, nil
		}
		resp = &types.StrategyArmoryResponse{IsSuccess: false}
		logx.Error("QueryStrategy error:", err)
		return resp, code.StrategyArmoryEmpty
	}
	// 3.2.2 "rule_weight"规则存在，则需要进行权重策略配置

	strategyRule, err := l.svcCtx.StrategyRuleModel.QueryStrategyRule(l.ctx, req.StrategyId, strategy.GetRuleWeight())
	// 业务异常，策略规则中 rule_weight 权重规则已适用但未配置
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return resp, nil
		}
		logx.Error("QueryStrategyRule error:", err)
		return resp, code.StrategyArmoryFail
	}
	// 3.2.3 装配策略奖品概率查找表
	RuleWeightValueMap := strategyRule.GetRuleWeightValues()

	for k, v := range RuleWeightValueMap {
		l.storeStrategyRuleWeight(req.StrategyId, k)
		ruleWeightValues := v
		StrategyAwardListFilter := cloneAndFilterStrategyAward(StrategyAwardList, ruleWeightValues)
		_, err = l.assembleLotteryStrategyByRuleWeight(req.StrategyId, k, StrategyAwardListFilter)
		if err != nil {
			resp = &types.StrategyArmoryResponse{IsSuccess: false}
			logx.Error("assembleLotteryStrategyByRuleWeight error:", err)
			return resp, xcode.ServerErr
		}
	}

	return resp, nil
}
func cloneAndFilterStrategyAward(list []*model.StrategyAward, filter []int) (listCloneAndFilter []*model.StrategyAward) {
	listCloneAndFilter = make([]*model.StrategyAward, 0)
	for _, v := range list {
		if linq.From(filter).Contains(int(v.AwardId)) {
			listCloneAndFilter = append(listCloneAndFilter, v)
		}
	}
	return listCloneAndFilter
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
	cacheKey := fmt.Sprintf(common.StrategyAwardCountKey, StrategyId, AwardId)
	err = l.svcCtx.BizRedis.SetCtx(l.ctx, cacheKey, strconv.Itoa(AwardCount))
	if err != nil {
		return err
	}
	return nil
}

//	func (l *StrategyArmoryLogic) storeStrategyAwardSearchRateTable(id int64, size int, tables map[int]int64) (err error) {
//		cacheKey := fmt.Sprintf(StrategyRateRange, id)
//		l.svcCtx.BizRedis.SetCtx(l.ctx, cacheKey, strconv.Itoa(size))
//		for i, v := range tables {
//			err = l.svcCtx.BizRedis.SetCtx(l.ctx, fmt.Sprintf(StrategyRateRangeKey, id, i), strconv.Itoa(int(v)))
//			if err != nil {
//				return err
//			}
//		}
//		return
//	}
func (l *StrategyArmoryLogic) storeStrategyAwardSearchRateTable(id int64, size int, tables map[int]int64) (err error) {
	// 主键使用给定的常量格式化
	cacheKey := fmt.Sprintf(common.StrategyRateRange, id)

	// 存储表的大小信息
	err = l.svcCtx.BizRedis.SetCtx(l.ctx, fmt.Sprintf(common.StrategyRateRangeSize, id), strconv.Itoa(size))
	if err != nil {
		return err
	}

	// 遍历表格并存储每个条目
	for i, v := range tables {
		err = l.svcCtx.BizRedis.HsetCtx(l.ctx, cacheKey, strconv.Itoa(i), strconv.FormatInt(v, 10))
		if err != nil {
			return err
		}
	}
	return nil
}
func (l *StrategyArmoryLogic) assembleLotteryStrategyByRuleWeight(id int64, m string, list []*model.StrategyAward) (IsSuccess bool, err error) {
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
	err = l.storeStrategyAwardSearchRateTableByRuleWeight(id, m, shuffleStrategyAwardSearchRateTable)
	if err != nil {
		return false, err
	}
	return
}
func (l *StrategyArmoryLogic) storeStrategyAwardSearchRateTableByRuleWeight(id int64, m string, tables map[int]int64) (err error) {
	for i, v := range tables {
		if v == 0 {
			continue
		}
		err = l.svcCtx.BizRedis.HsetCtx(l.ctx,
			fmt.Sprintf(common.StrategyRateRangeRuleWeightKey, id, m),
			strconv.Itoa(i),
			strconv.FormatInt(v, 10))
		if err != nil {
			return err
		}
	}
	return
}
func (l *StrategyArmoryLogic) storeStrategyRuleWeight(id int64, rule string) (err error) {
	cacheStrategyRuleWeightKey := fmt.Sprintf(common.StrategyRuleWeightKey, id)
	err = l.svcCtx.BizRedis.HsetCtx(l.ctx, cacheStrategyRuleWeightKey, rule, rule)
	if err != nil {
		logx.Error("storeStrategyRuleWeight error:", err)
		return err
	}
	return nil
}
