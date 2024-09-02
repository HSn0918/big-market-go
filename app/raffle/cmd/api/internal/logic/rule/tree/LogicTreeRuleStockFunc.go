package tree

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/hsn0918/BigMarket/common"

	"github.com/hsn0918/BigMarket/app/strategy/raffle/cmd/api/internal/svc"
)

func RuleStockFunc(ctx context.Context, svc *svc.ServiceContext, userId string, strategyId int64, awardId int, ruleValue string) (TreeActionEntity, error) {
	// 1.扣减库存
	status := subtractionAwardStock(ctx, svc, strategyId, awardId)
	// 2.判断库存是否充足
	if status {
		awardStockConsumeSendQueue(ctx, svc, strategyId, awardId)
		return TreeActionEntity{
			RuleLogicCheckTypeVO: TAKE_OVER,
			StrategyAwardVO: StrategyAwardVO{
				AwardId:        awardId,
				AwardRuleValue: ruleValue,
				End:            true,
			},
		}, nil
	}
	return TreeActionEntity{
		RuleLogicCheckTypeVO: ALLOW,
		StrategyAwardVO: StrategyAwardVO{
			AwardId:        101,
			AwardRuleValue: "1,100",
			End:            false,
		},
	}, nil

}
func subtractionAwardStock(ctx context.Context, svc *svc.ServiceContext, strategyId int64, awardId int) bool {
	cacheKey := fmt.Sprintf(common.StrategyAwardCountKey, strategyId, awardId)
	surplus, _ := svc.BizRedis.Decr(cacheKey)
	if surplus < 0 {
		// 库存小于0，恢复为0个
		svc.BizRedis.Set(cacheKey, "0")
		return false
	}
	// 1. 按照cacheKey decr 后的值，如 99、98、97 和 key 组成为库存锁的key进行使用。
	// 2. 加锁为了兜底，如果后续有恢复库存，手动处理等，也不会超卖。因为所有的可用库存key，都被加锁了。
	lockKey := fmt.Sprintf(common.StrategyAwardCountQueryKey+"#%d", strategyId, awardId, surplus)
	lock, _ := svc.BizRedis.SetnxEx(lockKey, "lock", 60*60*24)
	if !lock {
		logx.Info("策略奖品库存加锁失败：", lockKey)
		return false
	}
	return lock
}
func awardStockConsumeSendQueue(ctx context.Context, svc *svc.ServiceContext, strategyId int64, awardId int) (bool, error) {
	//userId := ctx.Value("userId").(int)
	cacheKey := fmt.Sprintf(common.StrategyAwardCountQueryList)
	data := fmt.Sprintf("strategyAward:%d,%d", strategyId, awardId)
	if _, err := svc.BizRedis.Rpush(cacheKey, data); err != nil {
		logx.Error("Rpush error：", err)
		return false, err
	}
	delayKey := fmt.Sprintf(common.StrategyAwardCountQueryListDelay)
	// 将相同的数据推送到延迟队列，使用 ZSet 模拟
	// 计算延迟时间（3秒后的 UNIX 时间戳）
	delay := time.Duration(3) * time.Second
	score := float64(time.Now().Add(delay).Unix())
	// 添加到 Redis 的有序集合
	if _, err := svc.BizRedis.Zadd(delayKey, int64(score), data); err != nil {
		logx.Error("Zadd error：", err)
		return false, err
	}
	return true, nil
}
