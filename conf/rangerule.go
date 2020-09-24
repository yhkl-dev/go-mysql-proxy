package conf

type Rule interface {
}

type Range struct {
	Max  int
	Min  int
	Node string
}

type RangeRule struct {
	Column string
	Ranges []*Range
}

func NewRangeRule(column string) *RangeRule {
	return &RangeRule{Column: column, Ranges: make([]*Range, 0)}
}

func (this *RangeRule) AddRange(max int, min int, node string) {
	this.Ranges = append(this.Ranges, &Range{Max: max, Min: min, Node: node})
}

func (this *RangeRule) GetNode(value int) string {
	for _, node := range this.Ranges {
		if value >= node.Min && value <= node.Max {
			return node.Node
		}
	}
	return ""
}
