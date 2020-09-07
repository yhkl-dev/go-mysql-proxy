package util

import (
	. "github.com/xwb1989/sqlparser"
	"proxy.go/conf"
)

func AliasedTableSQL(selectFields SelectExprs, from TableExprs) []string {

	config := conf.NewConfig()
	sqls := make([]string, 0)

	for _, tc := range from {
		tableName := tc.(*AliasedTableExpr).Expr.(TableName).Name.String()
		as := tc.(*AliasedTableExpr).As
		if mtables, ok := config.Models[tableName]; ok {
			for _, mtable := range mtables {

				newSQL := &Select{}
				newSQL.SelectExprs = selectFields
				newTe := &AliasedTableExpr{As: as, Expr: TableName{Name: NewTableIdent(mtable)}}
				newSQL.From = append(newSQL.From, newTe)
				buf := NewTrackedBuffer(nil)
				newSQL.Format(buf)
				sqls = append(sqls, buf.String())
			}
		}
	}
	return sqls
}
