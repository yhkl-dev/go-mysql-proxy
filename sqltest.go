package main

import (
	"fmt"
	"log"

	"github.com/xwb1989/sqlparser"
	"proxy.go/util"
)

func main() {
	sql := `
		select user_id, user_name from users u where id < 800 and id > 50 and name = 'abc' and age = 19 order by u.id desc limit 0, 10
	`

	stmt, err := sqlparser.Parse(sql)
	if err != nil {
		log.Fatal(err)
	}

	switch stmt := stmt.(type) {
	case *sqlparser.Select:
		sqls := util.AliasedTableSQL(stmt)
		for _, sql := range sqls {
			fmt.Println(sql)
		}

	default:
		fmt.Println("default")
	}
}
