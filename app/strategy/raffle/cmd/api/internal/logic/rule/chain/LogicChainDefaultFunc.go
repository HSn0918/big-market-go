package chain

import (
	"context"
	"fmt"
	"math/rand/v2"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/hsn0918/BigMarket/common"

	"github.com/hsn0918/BigMarket/app/strategy/raffle/cmd/api/internal/svc"
)

func DefaultFunc(ctx context.Context, svc *svc.ServiceContext, strategyId int64) (StrategyAwardVO, error) {
	// 1.redis中取rateRange
	cacheRateRange := fmt.Sprintf(common.StrategyRateRangeSize, strategyId)
	rateRangeStr, err := svc.BizRedis.GetCtx(ctx, cacheRateRange)
	if err != nil {
		logx.Error("redis get rateRange error:", err)
		return StrategyAwardVO{}, err
	}
	rateRange, err := strconv.Atoi(rateRangeStr)
	if err != nil {
		logx.Error("strconv.Atoi error:", err)
		return StrategyAwardVO{}, err
	}
	// 2.redis中取awardId
	randInt := rand.IntN(rateRange)
	cacheStrategy := fmt.Sprintf(common.StrategyRateRange, strategyId)
	awardIdStr, err := svc.BizRedis.HgetCtx(ctx, cacheStrategy, strconv.Itoa(randInt))
	if err != nil {
		logx.Error("redis get awardId error:", err)
		return StrategyAwardVO{}, err
	}
	awardId, err := strconv.Atoi(awardIdStr)
	if err != nil {
		logx.Error("strconv.Atoi error:", err)
		return StrategyAwardVO{}, err
	}
	// 3.返回
	return StrategyAwardVO{
		AwardId:    awardId,
		End:        true,
		LogicModel: RULE_DEFAULT.Code(),
	}, nil
}
