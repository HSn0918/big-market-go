package chain

import (
	"context"
	"fmt"
	"sort"
	"strconv"

	"github.com/hsn0918/BigMarket/common"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/hsn0918/BigMarket/app/strategy/raffle/cmd/api/internal/svc"
)

func WeightFunc(ctx context.Context, svcCtx *svc.ServiceContext, strategyId int64) (StrategyAwardVO, error) {
	// 0.获取rule_weight
	cacheStrategyRuleWeightKey := fmt.Sprintf(common.StrategyRuleWeightKey, strategyId)
	ruleWeightMap, err := svcCtx.BizRedis.HgetallCtx(ctx, cacheStrategyRuleWeightKey)
	if err != nil {
		logx.Error("HgetallCtx error:", err)
		return StrategyAwardVO{
			AwardId:    100,
			LogicModel: RULE_WEIGHT.Code(),
			End:        false,
		}, nil
	}
	WeightValue := make([]int, 0)
	for k := range ruleWeightMap {
		intValue, err := strconv.Atoi(k)
		if err != nil {
			logx.Error("strconv.Atoi error:", err)
			return StrategyAwardVO{
				AwardId:    100,
				LogicModel: RULE_WEIGHT.Code(),
				End:        false,
			}, nil
		}
		WeightValue = append(WeightValue, intValue)
	}
	sort.Slice(WeightValue, func(i, j int) bool {
		return WeightValue[i] < WeightValue[j]
	})
	// 1.获取用户使用抽奖积分
	usedPoints := ctx.Value("usedPoints").(string)
	usedPointsInt, err := strconv.Atoi(usedPoints)

	// 2.找到符合的rule_weight
	// 例如用户消费积分是4002，WeightValue：【4000，5000，6000】
	// 找到4000
	// 如果用户积分小于最小的4000
	// 直接返回
	selectedWeight := 0
	for _, weight := range WeightValue {
		if weight <= usedPointsInt {
			selectedWeight = weight
		} else {
			break
		}
	}
	if selectedWeight < WeightValue[0] {
		// 用户积分小于最小的时候
		// 直接返回，因为chain最后都为defaultFunc
		return StrategyAwardVO{
			AwardId:    101,
			LogicModel: RULE_WEIGHT.Code(),
			End:        false,
		}, nil
	}
	cacheStrategyRateRangeRuleWeightKey := fmt.Sprintf(common.StrategyRateRangeRuleWeightKey, strategyId, strconv.Itoa(selectedWeight))
	weightRangMap, err := svcCtx.BizRedis.HgetallCtx(ctx, cacheStrategyRateRangeRuleWeightKey)
	var awardId int
	for _, value := range weightRangMap {
		awardId, err = strconv.Atoi(value)
		if err != nil {
			logx.Error("strconv.Atoi error:", err)
			return StrategyAwardVO{}, err
		}
		break
	}
	if err != nil {
		logx.Error("HgetallCtx error:", err)
		return StrategyAwardVO{}, err
	}
	return StrategyAwardVO{
		AwardId:    awardId,
		End:        true,
		LogicModel: RULE_WEIGHT.Code(),
	}, nil
}
