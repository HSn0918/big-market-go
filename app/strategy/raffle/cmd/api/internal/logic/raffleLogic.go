package logic

import (
	"context"
	"fmt"
	"math/rand/v2"
	"strconv"

	"github.com/hsn0918/BigMarket/app/strategy/raffle/cmd/api/internal/svc"
	"github.com/hsn0918/BigMarket/app/strategy/raffle/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RaffleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRaffleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RaffleLogic {
	return &RaffleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RaffleLogic) Raffle(req *types.RaffleRequest) (resp *types.RaffleResponse, err error) {
	var awardId int
	awardId, err = l.getRandomAwardId(req.StrategyId)
	// todo:
	resp = &types.RaffleResponse{
		AwardId: awardId,
	}
	return
}
func (l *RaffleLogic) getRandomAwardId(StrategyId int64) (awardId int, err error) {
	// 1.从redis中取RateRange
	cacheRateRange := fmt.Sprintf(cacheStrategyRateRange, StrategyId)
	rateRangeStr, err := l.svcCtx.BizRedis.Get(cacheRateRange)
	if err != nil {
		return -1, err
	}
	rateRange, err := strconv.Atoi(rateRangeStr)
	if err != nil {
		return -1, err
	}
	randInt := rand.IntN(rateRange)
	// 2.从redis中取AwardId
	cacheAwardId := fmt.Sprintf(cacheStrategyRateRangeKey, StrategyId, randInt)
	awardIdStr, err := l.svcCtx.BizRedis.Get(cacheAwardId)
	if err != nil {
		return -1, err
	}
	awardId, err = strconv.Atoi(awardIdStr)
	if err != nil {
		return -1, err
	}
	return
}
