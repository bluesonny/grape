package models

import (
	"fmt"
	"log"
)

type Grape struct {
	ArticleId   int    `db:"article_id"`
	Link        string `db:"link"`
	Title       string `db:"title"`
	Abstract    string `db:"abstract"`
	ArticleTime string `db:"article_time"`
	Site        string `db:"site"`
}

func (gp *Grape) Insert(list []*Grape) (err error) {
	//defer DB.Close()
	if len(list) == 0 {
		log.Printf("数据库逻辑没有执行-----888----，空切片%v\n", list)
		return
	}

	for i := len(list)/2 - 1; i >= 0; i-- {
		opp := len(list) - 1 - i
		list[i], list[opp] = list[opp], list[i]
	}
	log.Println("...写数据库...")
	statement := "insert into grape (article_id, link, title, abstract, article_time, site) values (:article_id, :link, :title, :abstract, :article_time, :site)"
	_, err = DB.NamedExec(statement, list)
	if err != nil {
		log.Printf("数据库失败-----999----%v\n", err)
		return
	}
	return
}

func (gp *Grape) Get(where string) bool {

	var id int
	sql := fmt.Sprintf("select id from grape where %s", where)
	log.Printf("sql语句...%v", sql)
	err := DB.Get(&id, sql)
	if err != nil {
		log.Printf("返回结果什么错误...%v", err)
	}
	log.Printf("记录id=%d", id)
	if id > 0 {
		return true
	}
	return false
}
