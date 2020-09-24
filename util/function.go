package util

import (
	"strconv"

	. "github.com/xwb1989/sqlparser"
	"proxy.go/conf"
)

var operator = []interface{}{"=", "<", ">", "<=", ">="}

func GetString(expr Expr) string {
	buf := NewTrackedBuffer(nil)
	expr.Format(buf)
	return buf.String()
}

func GetInt(expr Expr, defaultValue int) int {
	str := GetString(expr)
	istr, err := strconv.Atoi(str)
	if err != nil {
		return defaultValue
	}
	return istr
}

var config = conf.NewConfig()

func ParseWhere(expr Expr) []interface{} {
	if expr == nil {
		return nil
	}

	ce := expr.(*ComparisonExpr)
	rule := config.Rule.(*conf.RangeRule)
	column := rule.Column
	if GetString(ce.Left) == column {
		if Contains(operator, ce.Operator) {
			node := rule.GetNode(GetInt(ce.Right, 0), ce.Operator)
			return node
		}
	}
	return nil
}

func ParseMultiWhere(where *Where) []interface{} {
	if where == nil {
		return nil
	}

	exps := PaseWhereToSlice(where.Expr)

	ret := make([]interface{}, 0)
	for _, exp := range exps {
		parseNode := ParseWhere(exp)
		if parseNode == nil || len(parseNode) == 0 {
			continue
		}
		if len(ret) == 0 {
			ret = parseNode
		}
		ret = IntersectSlice(ret, parseNode)
	}
	return ret
}

func AliasedTableSQL(stmt *Select) []string {

	config := conf.NewConfig()
	node := ParseMultiWhere(stmt.Where)
	sqls := make([]string, 0)

	for _, tc := range stmt.From {
		tableName := tc.(*AliasedTableExpr).Expr.(TableName).Name.String()
		as := tc.(*AliasedTableExpr).As
		if mtables, ok := config.Models[tableName]; ok {
			for _, mtable := range mtables {
				if node != nil && len(node) > 0 && !Contains(node, mtable) {
					continue
				}
				sql := forSQL(stmt, mtable, as)
				sqls = append(sqls, sql)
			}
		}
	}
	return sqls
}

func forSQL(stmt *Select, mtable string, as TableIdent) string {
	newSQL := &Select{}
	newSQL.SelectExprs = stmt.SelectExprs
	newTe := &AliasedTableExpr{As: as, Expr: TableName{Name: NewTableIdent(mtable)}}
	newSQL.From = append(newSQL.From, newTe)
	newSQL.Where = stmt.Where
	newSQL.OrderBy = stmt.OrderBy
	newSQL.Limit = stmt.Limit
	buf := NewTrackedBuffer(nil)
	newSQL.Format(buf)
	return buf.String()
}

/*
func getNode(expr Expr, isLeft bool) string {
	if andExpr, ok := expr.(*AndExpr); ok {
		if isLeft {
			return getNode(andExpr.Left, isLeft)
		}
		return getNode(andExpr.Right, isLeft)
	} else if cExpr, ok := expr.(*ComparisonExpr); ok {
		return ParseWhere(cExpr)
	} else {
		return ""
	}
}
*/
func PaseWhereToSlice(expr Expr) []Expr {
	exprList := make([]Expr, 0)
	temp := expr
	for {
		if andExpr, ok := temp.(*AndExpr); ok {
			exprList = append(exprList, andExpr.Right)
			temp = andExpr.Left
		} else {
			exprList = append(exprList, temp)
			break
		}
	}
	return exprList
}
