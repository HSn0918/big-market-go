package tree

import (
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"github.com/hsn0918/BigMarket/app/strategy/raffle/cmd/api/internal/model"

	"github.com/hsn0918/BigMarket/app/strategy/raffle/cmd/api/internal/svc"
)

// DefaultTreeFactor 决策树工厂
type DefaultTreeFactor struct {
	logicTreeGroup map[string]Logic
	tree           *RuleTree
	ctx            context.Context
	svcCtx         *svc.ServiceContext
}

// DecisionTreeEngine 决策树引擎
type DecisionTreeEngine struct {
	logicTreeGroup map[string]Logic
	tree           *RuleTree
}

func NewDefaultTreeFactor(ctx context.Context, svcCtx *svc.ServiceContext) DefaultTreeFactor {
	return DefaultTreeFactor{
		ctx:            ctx,
		svcCtx:         svcCtx,
		tree:           new(RuleTree),
		logicTreeGroup: NewLogicTreeGroup(),
	}
}
func NewLogicTreeGroup() map[string]Logic {
	mp := make(map[string]Logic)
	mp[RULE_LOCK.Code()] = RuleLockFunc
	mp[RULE_STOCK.Code()] = RuleStockFunc
	mp[RULE_LUCK_AWARD.Code()] = RuleLuckAwardFunc
	return mp
}
func (d DefaultTreeFactor) OpenLogicTree(strategyAward *model.StrategyAward) *DecisionTreeEngine {
	// 1.获取树
	tree, err := d.svcCtx.RuleTreeModel.FindOneByTreeId(d.ctx, strategyAward.RuleModels.String)
	if err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			return nil
		}
		logx.Error("FindOneByTreeId", err)
		return nil
	}

	// 2.获取树节点建联系
	lines, err := d.svcCtx.RuleTreeNodeLineModel.QueryRuleTreeNodeLinesByTreeId(d.ctx, strategyAward.RuleModels.String)
	treeNodeLineMap := make(map[string][]*RuleTreeNodeLine)
	for _, line := range lines {
		vo := &RuleTreeNodeLine{
			TreeID:         line.TreeId,
			RuleNodeFrom:   line.RuleNodeFrom,
			RuleNodeTo:     line.RuleNodeTo,
			RuleLimitType:  line.RuleLimitType,
			RuleLimitValue: line.RuleLimitValue,
		}
		treeNodeLineMap[line.RuleNodeFrom] = append(treeNodeLineMap[line.RuleNodeFrom], vo)
	}
	// 3.获取树节点
	nodes, err := d.svcCtx.RuleTreeNodeModel.QueryRuleTreeNodesByTreeId(d.ctx, strategyAward.RuleModels.String)
	if err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			return nil
		}
		logx.Error("QueryRuleTreeNodesByTreeId", err)
		return nil
	}
	// 创建节点映射
	treeNodeMap := make(map[string]*RuleTreeNode)
	for _, node := range nodes {
		vo := &RuleTreeNode{
			TreeID:             node.TreeId,                   // 填充TreeID
			RuleKey:            node.RuleKey,                  // 填充RuleKey
			RuleDesc:           node.RuleDesc,                 // 填充RuleDesc
			RuleValue:          node.RuleValue.String,         // 填充RuleValue
			TreeNodeLineVOList: treeNodeLineMap[node.RuleKey], // 可以根据需要填充或处理
		}
		// 将节点按TreeID分类并加入到映射中
		treeNodeMap[node.RuleKey] = vo
	}
	// 4.创建树
	ruleTree := &RuleTree{
		TreeID:           tree.TreeId,
		TreeName:         tree.TreeName,
		TreeDesc:         tree.TreeDesc.String,
		TreeRootRuleNode: tree.TreeNodeRuleKey,
		TreeNodeMap:      treeNodeMap,
	}
	return &DecisionTreeEngine{
		logicTreeGroup: d.logicTreeGroup,
		tree:           ruleTree,
	}
}
func (d DecisionTreeEngine) decisionLogic(matterValue string, nodeLine *RuleTreeNodeLine) bool {
	switch nodeLine.RuleLimitType {
	case EQUAL:
		return matterValue == nodeLine.RuleLimitValue
	case GT:
	case LT:
	case GE:
	case LE:
	default:
		return false
	}
	return false
}
func (d DecisionTreeEngine) Process(ctx context.Context, svcCtx *svc.ServiceContext, userId string, strategyId int64, awardId int) (StrategyAwardVO, error) {
	strategyAwardVO := StrategyAwardVO{}
	// 1.获取基本信息
	nextNode := d.tree.TreeRootRuleNode
	treeNodeMap := d.tree.TreeNodeMap
	// 2.根节点
	ruleTreeNode := treeNodeMap[nextNode]
	// 3.判断是否满足条件
	for ruleTreeNode != nil {
		// 3.1 获取决策节点
		logicTreeNode := d.logicTreeGroup[ruleTreeNode.RuleKey]
		ruleValue := ruleTreeNode.RuleValue
		// 3.2 决策节点计算
		logicEntity, err := logicTreeNode(ctx, svcCtx, userId, strategyId, awardId, ruleValue)
		if err != nil {
			return StrategyAwardVO{
				AwardId:        100,
				AwardRuleValue: ruleValue,
				End:            true,
			}, err
		}

		ruleLogicCheckType := logicEntity.RuleLogicCheckTypeVO
		strategyAwardVO = logicEntity.StrategyAwardVO
		// 3.3 获取下个节点
		nextNode = d.nextNode(ruleLogicCheckType.Code(), ruleTreeNode.TreeNodeLineVOList)
		ruleTreeNode = treeNodeMap[nextNode]
	}

	return StrategyAwardVO{
		AwardId:        strategyAwardVO.AwardId,
		AwardRuleValue: strategyAwardVO.AwardRuleValue,
		End:            true,
	}, nil
}
func (d DecisionTreeEngine) nextNode(matterValue string, treeNodeLineVOList []*RuleTreeNodeLine) string {
	if treeNodeLineVOList == nil || len(treeNodeLineVOList) == 0 {
		return ""
	}
	for _, line := range treeNodeLineVOList {
		if d.decisionLogic(matterValue, line) {
			return line.RuleNodeTo
		}
	}
	return ""

}
