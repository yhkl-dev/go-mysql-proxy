package util

import (
	"fmt"
	"log"

	"github.com/siddontang/go-mysql/client"
	"github.com/siddontang/go-mysql/mysql"
	"github.com/siddontang/go-mysql/server"
)

type MyHandler struct {
	server.EmptyHandler
	conn *client.Conn
}

func (h MyHandler) HandleQuery(query string) (*mysql.Result, error) {
	fmt.Println("query: ", query)
	result, err := h.conn.Execute(query)
	if err != nil {
		return nil, fmt.Errorf("error ", err)
	}
	return result, nil
}

func NewMyHandler() MyHandler {
	conn, err := client.Connect("47.94.221.199:3306", "yhkl", "123456", "dbdms")
	if err != nil {
		log.Fatal(err)
	}
	return MyHandler{conn: conn}
}
