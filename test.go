package main

import (
	"fmt"
	"log"

	"github.com/xwb1989/sqlparser"
	"proxy.go/util"
)

func main() {
	sql := "select * from users as a where pname='avb' and id = 1000"
	stmt, err := sqlparser.Parse(sql)
	if err != nil {
		log.Fatal(err)
	}
	switch stmt := stmt.(type) {
	case *sqlparser.Select:
		/*
			for _, node := range stmt.From {
				getTable := node.(*sqlparser.AliasedTableExpr)
				fmt.Println("string", getTable.As.String())
				fmt.Println(getTable.Expr.(sqlparser.TableName).Name)
			}
		*/
		buf := sqlparser.NewTrackedBuffer(nil)
		stmt.SelectExprs.Format(buf)
		fmt.Println(buf.String())

		sqls := util.AliasedTableSQL(stmt)
		for _, sql := range sqls {
			fmt.Println(sql)
		}
	default:
		fmt.Println("default")
	}

}
