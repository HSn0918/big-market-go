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
		logicTreeGroup: make(map[string]Logic),
	}
}
func NewLogicTreeGroup() map[string]Logic {
	mp := make(map[string]Logic)
	mp["default"] = nil
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
func (d DecisionTreeEngine) Process(userId string, strategyId int64, awardId int) (StrategyAwardVO, error) {
	return StrategyAwardVO{
		AwardId:        0,
		AwardRuleValue: "",
		End:            false,
	}, nil
}
