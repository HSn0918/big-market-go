package code

import "github.com/hsn0918/BigMarket/pkg/xcode"

var (
	StrategyArmoryEmpty = xcode.New(10001, "装载策略不存在")
	StrategyArmoryFail  = xcode.New(10002, "转载失败")
	RaffleFail          = xcode.New(10003, "抽奖失败")
)
