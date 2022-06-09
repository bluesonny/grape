package handler

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"grap-data/config"
	"grap-data/models"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

func Deal() {
	urls := config.ViperConfig.Urls.List
	for _, url := range urls {
		// 开启一个goroutine
		//go fetch(url, ch)
		fetch(url)
	}

}

// 根据URL获取资源内容
func fetch(url string) {
	//start := time.Now()
	// 发送网络请求
	res, err := http.Get(url)
	if err != nil {
		// 输出异常信息
		//ch <- fmt.Sprint(err)
		os.Exit(1)
	}
	// 关闭资源
	//res.Body.Close()
	dom, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	site := getFileName(url)
	var list []*models.Grape
	var gr models.Grape
	switch site {

	case "www.autosport.com", "www.motorsport.com":
		dom.Find("[data-entity-type='article']").Each(func(i int, selection *goquery.Selection) {
			var g models.Grape
			g.Site = site
			time, exists := selection.Find("time.ms-item_date-value").Attr("datetime")
			if !exists {
				return
			}
			id, _ := selection.Attr("data-entity-id")
			a := selection.Find("a.ms-item_link")
			link, _ := a.Attr("href")
			title, _ := a.Attr("title")
			d := selection.Find("p.ms-item_subheader")
			desc := d.Text()

			fmt.Println("标题:" + title)
			fmt.Println("连接:" + link)
			fmt.Println("ID:" + id)
			fmt.Println("时间:" + time)
			fmt.Println("摘要:" + desc)

			g.ArticleId, _ = strconv.Atoi(id)
			g.ArticleTime = time
			g.Link = link
			g.Title = title
			g.Abstract = desc
			list = append(list, &g)
		})
		gr.Insert(list)
	}

}

// 获取文件名
func getFileName(urls string) string {
	// 从URL中匹配域名部分
	u, err := url.Parse(urls)
	if err != nil {
		log.Fatal(err)
	}
	return u.Hostname()
}
