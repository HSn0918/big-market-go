// Code generated by goctl. DO NOT EDIT.
package types

type RaffleAwardListRequest struct {
	StrategyId int64 `json:"strategy_id"` // 策略ID，用于查询对应的抽奖奖品列表
}

type RaffleAwardListResponse struct {
	AwardId       int    `json:"award_id"`       // 奖品ID
	AwardTitle    string `json:"award_title"`    // 奖品标题
	AwardSubtitle int    `json:"award_subtitle"` // 奖品副标题
	Sort          int    `json:"sort"`           // 奖品排序索引
}

type RaffleRequest struct {
	StrategyId int64 `json:"strategy_id"` // 策略ID，用于执行抽奖操作
}

type RaffleResponse struct {
	AwardId    int `json:"award_id"`    // 获得的奖品ID
	AwardIndex int `json:"award_index"` // 奖品在列表中的索引
}

type StrategyArmoryRequest struct {
	StrategyId int64 `json:"strategy_id"` // 策略ID，用于查询对应的策略库存
}

type StrategyArmoryResponse struct {
	IsSuccess bool `json:"is_success"` // 表明请求是否成功
}