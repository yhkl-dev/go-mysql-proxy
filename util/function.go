package util

import (
	"strconv"

	. "github.com/xwb1989/sqlparser"
	"proxy.go/conf"
)

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

func ParseWhere(expr Expr) string {
	if expr == nil {
		return ""
	}

	ce := expr.(*ComparisonExpr)

	rule := config.Rule.(*conf.RangeRule)
	column := rule.Column
	if GetString(ce.Left) == column {
		node := rule.GetNode(GetInt(ce.Right, 0))

		return node
	}
	return ""
}

func ParseMultiWhere(where *Where) string {
	if where == nil {
		return ""
	}

	retl := getNode(where.Expr, true)
	if retl == "" {
		return getNode(where.Expr, false)
	}
	return retl
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
				if node != "" && node != mtable {
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
