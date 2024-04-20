package code

import "github.com/hsn0918/BigMarket/pkg/xcode"

var (
	StrategyArmoryEmpty = xcode.New(10001, "装载策略不存在")
	StrategyArmoryFail  = xcode.New(10002, "转载失败")
	RaffleFail          = xcode.New(10003, "抽奖失败")
	RaffleAwardEmpty    = xcode.New(10004, "抽奖记录为空")
	RaffleAwardNotFound = xcode.New(10005, "抽奖记录不存在")
)
