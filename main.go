package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"grap-data/config"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

type Data struct {
	number string
	title  string
	url    string
	desc   string
	time   string
}

func main() {
	log.Printf("%v", config.C.Urls.List)
	return
	urls := []string{
		//"https://www.formula1.com/en/latest/all.html",
		"https://www.autosport.com/all/news/",
		//"https://www.motorsport.com/all/news/",
	}

	// 创建channel
	//ch := make(chan string)

	// 开始时间
	start := time.Now()

	for _, url := range urls {
		// 开启一个goroutine
		//go fetch(url, ch)
		fetch(url)
	}
	// 总消耗的时间
	elapsed := time.Since(start).Seconds()

	fmt.Printf("%.2fs elapsed\n", elapsed)
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

	// 读取资源数据
	//body, err := ioutil.ReadAll(res.Body)

	// 关闭资源
	//res.Body.Close()

	if err != nil {
		// 输出异常信息
		//ch <- fmt.Sprintf("while reading %s: %v", url, err)
		os.Exit(1)
	}

	dom, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	//var data []Data
	// 筛选class为td-01的元素
	dom.Find("[data-entity-type='article']").Each(func(i int, selection *goquery.Selection) {

		time, exists := selection.Find("time.ms-item_date-value").Attr("datetime")
		if !exists {
			return
		}
		id, _ := selection.Attr("data-entity-id")
		fmt.Println("ID:" + id)
		fmt.Println("时间:" + time)
		a := selection.Find("a.ms-item_link")
		link, _ := a.Attr("href")
		fmt.Println("连接:" + link)
		title, _ := a.Attr("title")
		fmt.Println("标题:" + title)
		d := selection.Find("p.ms-item_subheader")
		desc := d.Text()
		fmt.Println("摘要:" + desc)
	})

	//分析文本ms-item_subheader
	/*
		content := string(body)
		doc := soup.HTMLParse(content)
		subDocs := doc.FindAll("div", "data-entity-type", "article")
		for _, subDoc := range subDocs {
			//link := subDoc.Find("time", "class", "ms-item_date-value")

			//fmt.Println(link.Attrs()["datetime"])

			a := subDoc.Find("a", "class", "ms-item_link")
			fmt.Println(a.Attrs()["href"])
			fmt.Println(a.Attrs()["title"])
		}

	*/
	//log.Printf("匹配的内容:%v", str)
	// 写入文件
	//log.Printf("名字:%v", getFileName(url))
	//ioutil.WriteFile(getFileName(url), , 0644)

	// 消耗的时间
	//elapsed := time.Since(start).Seconds()

	// 输出单个URL消耗的时间
	//ch <- fmt.Sprintf("%.2fs %s", elapsed, url)
}

// 获取文件名
func getFileName(urls string) string {
	// 从URL中匹配域名部分
	u, err := url.Parse(urls)
	if err != nil {
		log.Fatal(err)
	}
	return u.Hostname() + ".html"
	//	return RE.FindString(url) + ".txt"
}
