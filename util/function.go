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

func ParseWhere(where *Where) string {
	if where == nil {
		return ""
	}
	ce := where.Expr.(*ComparisonExpr)

	rule := config.Rule.(*conf.RangeRule)
	column := rule.Column
	if GetString(ce.Left) == column {
		node := rule.GetNode(GetInt(ce.Right, 0))

		return node
	}
	return ""
}

func AliasedTableSQL(stmt *Select) []string {

	config := conf.NewConfig()
	sqls := make([]string, 0)

	for _, tc := range stmt.From {
		tableName := tc.(*AliasedTableExpr).Expr.(TableName).Name.String()
		as := tc.(*AliasedTableExpr).As
		node := ParseWhere(stmt.Where)
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
	ParseWhere(stmt.Where)
	newSQL.Where = stmt.Where
	newSQL.OrderBy = stmt.OrderBy
	newSQL.Limit = stmt.Limit
	buf := NewTrackedBuffer(nil)
	newSQL.Format(buf)
	return buf.String()
}
