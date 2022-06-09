package models

import (
	"log"
)

type Grape struct {
	ArticleId   int    `db:"article_id"`
	Link        string `db:"link"`
	Title       string `db:"title"`
	Desc        string `db:"desc"`
	ArticleTime string `db:"article_time"`
	Site        string `db:"site"`
}

func (gp *Grape) Insert(list []*Grape) (err error) {
	log.Println("写数据库....")
	defer DB.Close()
	statement := "insert into grape (article_id, link, title,desc,article_time,site) values (:article_id, :link, :title,:desc,:article_time,:site)"
	_, err = DB.NamedExec(statement, list)
	if err != nil {
		log.Printf("数据库失败1%v\n", err)
		return
	}
	return
}
