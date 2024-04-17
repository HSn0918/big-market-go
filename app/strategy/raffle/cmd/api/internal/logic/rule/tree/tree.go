package tree

type RuleTreeNodeLine struct {
	TreeID         string
	RuleNodeFrom   string
	RuleNodeTo     string
	RuleLimitType  string
	RuleLimitValue string
}

type RuleTreeNode struct {
	TreeID             string
	RuleKey            string
	RuleDesc           string
	RuleValue          string
	TreeNodeLineVOList []*RuleTreeNodeLine
}

type RuleTree struct {
	TreeID           string
	TreeName         string
	TreeDesc         string
	TreeRootRuleNode string
	TreeNodeMap      map[string]*RuleTreeNode
}
