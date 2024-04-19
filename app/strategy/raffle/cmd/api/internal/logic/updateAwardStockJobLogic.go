package logic

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hsn0918/BigMarket/common"

	"github.com/robfig/cron/v3"

	"github.com/hsn0918/BigMarket/app/strategy/raffle/cmd/api/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type CronJob struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	Job    *cron.Cron
}

func NewCronJob(ctx context.Context, svcCtx *svc.ServiceContext, cron *cron.Cron) *CronJob {
	c := &CronJob{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
		Job:    cron,
	}
	c.addTasks()
	return c
}

func (c *CronJob) addTasks() {
	_, _ = c.Job.AddFunc("0/5 * * * * ?", func() {
		cacheKey := fmt.Sprintf(common.StrategyAwardCountQueryList)
		strategyString, err := c.svcCtx.BizRedis.Lpop(cacheKey)
		if err == nil {
			strategyAward := strings.Split(strings.Split(strategyString, common.COLON)[1], common.SPLIT)
			strategyId, _ := strconv.Atoi(strategyAward[0])
			awardId, _ := strconv.Atoi(strategyAward[1])
			fmt.Printf("strategyId: %d, awardId: %d\n", strategyId, awardId)
			_ = c.svcCtx.StrategyAwardModel.UpdateAwardStock(int64(strategyId), awardId)
		}
	})
}
