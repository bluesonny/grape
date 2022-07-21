package handler

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"grap-data/config"
	"grap-data/models"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
)

var (
	wg sync.WaitGroup
)

func Deal() {
	urls := config.ViperConfig.Urls.List
	//chs := make([]chan int, len(urls))
	for _, url := range urls {
		// 开启一个goroutine
		wg.Add(1)
		//chs[i] = make(chan int)
		go fetch(url)
	}
	wg.Wait()
	/*
		for _, ch := range chs {
			v, ok := <-ch
			log.Printf("通道数据：--------v:%v,ok:%v", v, ok)
		}
	*/
	log.Printf("--------全部解析完毕------")
}

// 根据URL获取资源内容
func fetch(url string) {
	defer wg.Done()
	//defer close(ch)
	res, err := http.Get(url)
	if err != nil {
		log.Printf("请求出错:---%v", err)
		return
	}
	dom, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Printf("dom获取出错:---%v", err)
		return
	}
	site := getFileName(url)
	log.Printf("---------开始解析和入库--------：%s", site)
	var list []*models.Grape
	var gr models.Grape
	switch site {

	case "www.autosport.com", "www.motorsport.com", "it.motorsport.com":
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

			fmt.Println("ID:" + id)
			fmt.Println("标题:" + title)
			fmt.Println("摘要:" + desc)
			fmt.Println("时间:" + time)
			fmt.Println("链接:" + link)

			g.ArticleId, _ = strconv.Atoi(id)
			g.ArticleTime = time
			g.Link = link
			g.Title = title
			g.Abstract = desc
			where := fmt.Sprintf("article_id=%s and site='%s'", id, site)
			ret := gr.Get(where)
			if ret {
				return
			}
			list = append(list, &g)

		})
		gr.Insert(list)

	case "www.skysports.com":
		dom.Find("div.news-list__item").Each(func(i int, selection *goquery.Selection) {
			var g models.Grape
			g.Site = site
			time := selection.Find("span.label__timestamp").Text()

			id, _ := selection.Attr("data-id")
			a := selection.Find("a.news-list__headline-link")
			link, _ := a.Attr("href")
			title := a.Text()
			d := selection.Find("p.news-list__snippet")
			desc := d.Text()

			fmt.Println("ID:" + id)
			fmt.Println("标题:" + strings.TrimSpace(title))
			fmt.Println("摘要:" + desc)
			fmt.Println("时间:" + time)
			fmt.Println("连接:" + cutStr(link, site) + "\n")

			g.ArticleId, _ = strconv.Atoi(id)
			g.ArticleTime = time
			g.Link = cutStr(link, site)
			g.Title = strings.TrimSpace(title)
			g.Abstract = desc
			where := fmt.Sprintf("article_id=%s and site='%s'", id, site)
			ret := gr.Get(where)
			if ret {
				return
			}
			list = append(list, &g)
		})
		gr.Insert(list)

	case "www.racefans.net":
		dom.Find("article.post").Each(func(i int, selection *goquery.Selection) {
			var g models.Grape
			g.Site = site
			id, _ := selection.Attr("id")
			id = cutStr(id, "post-")
			_, err := strconv.Atoi(id)
			if err != nil {
				return
			}

			time := "" // selection.Find("span.label__timestamp").Text()

			a := selection.Find("[rel='bookmark']")
			link, _ := a.Attr("href")
			link = cutStr(link, site)
			time = link[1:11]
			title := a.Text()

			d := selection.Find("header>p")
			desc := d.Text()
			fmt.Println("ID:" + id)
			fmt.Println("标题:" + strings.TrimSpace(title))
			fmt.Println("摘要:" + desc)
			fmt.Println("时间:" + time)
			fmt.Println("连接:" + link + "\n")

			g.ArticleId, _ = strconv.Atoi(id)
			g.ArticleTime = time
			g.Link = link
			g.Title = strings.TrimSpace(title)
			g.Abstract = desc
			where := fmt.Sprintf("article_id=%s and site='%s'", id, site)
			ret := gr.Get(where)
			if ret {
				return
			}
			list = append(list, &g)
		})
		gr.Insert(list)
	case "racingnews365.com":
		dom.Find("a.card--default").Each(func(i int, selection *goquery.Selection) {
			var g models.Grape
			g.Site = site
			time, _ := selection.Find("time.postdate").Attr("datetime")
			if time == "" {
				return
			}
			id, _ := selection.Attr("data-id")
			link, _ := selection.Attr("href")
			link = cutStr(link, site)
			title := selection.Find("span.card__title").Find("span").Text()

			desc := ""

			fmt.Println("ID:" + id)
			fmt.Println("标题:" + title)
			fmt.Println("摘要:" + desc)
			fmt.Println("时间:" + time)
			fmt.Println("连接:" + link)

			g.ArticleId, _ = strconv.Atoi(id)
			g.ArticleTime = time
			g.Link = link
			g.Title = title
			g.Abstract = desc
			where := fmt.Sprintf("article_id=%s and site='%s'", id, site)
			ret := gr.Get(where)
			if ret {
				return
			}
			list = append(list, &g)
		})
		gr.Insert(list)
	case "www.crash.net":
		dom.Find("div.views-row").Each(func(i int, selection *goquery.Selection) {
			var g models.Grape
			g.Site = site
			t := selection.Find("div.views-field-created")
			time := t.Find("span.field-content").Text()
			if time == "" {
				return
			}
			id := "0"
			a := selection.Find("div.views-field-title")
			l := a.Find("span>a")
			link, _ := l.Attr("href")
			title := l.Text()
			d := selection.Find("div.views-field-body")
			desc := d.Find("div>p").Text()

			fmt.Println("ID:" + id)
			fmt.Println("标题:" + title)
			fmt.Println("摘要:" + desc)
			fmt.Println("时间:" + time)
			fmt.Println("连接:" + link)

			g.ArticleId, _ = strconv.Atoi(id)
			g.ArticleTime = time
			g.Link = link
			g.Title = title
			g.Abstract = desc
			where := fmt.Sprintf("%s'%s'", "link=", link)
			ret := gr.Get(where)
			if ret {
				return
			}
			list = append(list, &g)
		})
		gr.Insert(list)

	case "www.planetf1.com":
		ar := dom.Find("header.articleList__header").Next()
		ar.Find("li.articleList__item").Each(func(i int, selection *goquery.Selection) {
			var g models.Grape
			g.Site = site
			time, _ := selection.Find("time").Attr("datetime")
			if time == "" {
				return
			}
			id := "0"
			a := selection.Find("a")
			link, _ := a.Attr("href")
			link = cutStr(link, site)
			title := selection.Find("h3").Text()
			desc := selection.Find("p").Text()

			fmt.Printf("--------------第%d条--------------\n:", i+1)
			fmt.Println("ID:" + id)
			fmt.Println("标题:" + title)
			fmt.Println("摘要:" + desc)
			fmt.Println("时间:" + time)
			fmt.Println("连接:" + link)

			g.ArticleId, _ = strconv.Atoi(id)
			g.ArticleTime = time
			g.Link = link
			g.Title = title
			g.Abstract = desc
			where := fmt.Sprintf("%s'%s'", "link=", link)
			ret := gr.Get(where)
			if ret {
				return
			}
			list = append(list, &g)
		})
		gr.Insert(list)

	case "www.gpfans.com":
		dom.Find("div.nieuwereeks").Each(func(i int, selection *goquery.Selection) {
			var g models.Grape
			g.Site = site
			time, _ := selection.Attr("data-datum")
			t := selection.Find("li.headlinelabel").Next()
			time = time + " " + t.Text()
			id := "0"
			a := selection.Find("a")
			link, _ := a.Attr("href")
			link = cutStr(link, site)
			title := selection.Find("h3").Text()
			d := selection.Find("li.headlinelabel")
			desc := d.Text()
			if link == "" || title == "" {
				return
			}

			fmt.Println("ID:" + id)
			fmt.Println("标题:" + title)
			fmt.Println("摘要:" + desc)
			fmt.Println("时间:" + time)
			fmt.Println("连接:" + link)

			g.ArticleId, _ = strconv.Atoi(id)
			g.ArticleTime = time
			g.Link = link
			g.Title = title
			g.Abstract = desc
			where := fmt.Sprintf("%s'%s'", "link=", link)
			ret := gr.Get(where)
			if ret {
				return
			}
			list = append(list, &g)
		})
		gr.Insert(list)
	case "www.grandprix247.com":
		dom.Find("article.cb-blog-style-c").Each(func(i int, selection *goquery.Selection) {
			var g models.Grape
			g.Site = site

			id, _ := selection.Attr("id")
			id = cutStr(id, "post-")

			a := selection.Find("h2>a")
			link, _ := a.Attr("href")
			link = cutStr(link, site)
			time := link[1:11]
			title := a.Text()
			title = strings.TrimSpace(title)
			d := selection.Find("div.cb-excerpt")
			desc := d.Text()
			desc = strings.TrimSpace(desc)

			fmt.Printf("-------------------第 %d 条-------------------\n", i+1)
			fmt.Println("ID:" + id)
			fmt.Println("标题:" + title)
			fmt.Println("摘要:" + desc)
			fmt.Println("时间:" + time)
			fmt.Println("连接:" + link)

			g.ArticleId, _ = strconv.Atoi(id)
			g.ArticleTime = time
			g.Link = link
			g.Title = title
			g.Abstract = desc
			where := fmt.Sprintf("article_id=%s and site='%s'", id, site)
			ret := gr.Get(where)
			if ret {
				return
			}
			list = append(list, &g)
		})
		gr.Insert(list)

	case "the-race.com":

		//抓取头条
		var gr models.Grape
		h := dom.Find("div.post-category").Next()
		a := h.Find("h1>a")
		link, _ := a.Attr("href")
		title := a.Text()
		d := h.Find("p")
		desc := d.Text()
		desc = strings.TrimSpace(desc)
		t := dom.Find("span.time")
		time := strings.TrimSpace(t.Text())
		id := "0"
		fmt.Printf("--------------头条--------------\n")
		fmt.Println("ID:" + id)
		fmt.Println("标题:" + title)
		fmt.Println("摘要:" + desc)
		fmt.Println("时间:" + time)
		fmt.Println("连接:" + link)

		gr.Site = site
		gr.ArticleId, _ = strconv.Atoi(id)
		gr.ArticleTime = time
		gr.Link = link
		gr.Title = title
		gr.Abstract = desc
		where := fmt.Sprintf("%s'%s'", "link=", link)
		ret := gr.Get(where)
		if !ret {
			list = append(list, &gr)
		}

		dom.Find("div.related_group").Each(func(i int, selection *goquery.Selection) {
			var g models.Grape
			g.Site = site

			t := selection.Find("div>span")
			time := t.Text()
			id := "0"
			a := selection.Find("h3>a")
			link, _ := a.Attr("href")
			title := a.Text()
			title = strings.TrimSpace(title)
			d := selection.Find("p")
			desc := d.Text()
			desc = strings.TrimSpace(desc)
			if title == "" || link == "" {
				return
			}

			fmt.Printf("-----------------第%d条-------------------\n", i+1)
			fmt.Println("ID:" + id)
			fmt.Println("标题:" + title)
			fmt.Println("摘要:" + desc)
			fmt.Println("时间:" + time)
			fmt.Println("连接:" + link)

			g.ArticleId, _ = strconv.Atoi(id)
			g.ArticleTime = time
			g.Link = link
			g.Title = title
			g.Abstract = desc
			where := fmt.Sprintf("%s'%s'", "link=", link)
			ret := gr.Get(where)
			if ret {
				return
			}
			list = append(list, &g)
		})
		gr.Insert(list)

	case "www.formula1.com":

		dom.Find("a.f1-cc").Each(func(i int, selection *goquery.Selection) {
			var g models.Grape
			g.Site = site
			time := ""
			id := "0"
			link, _ := selection.Attr("href")
			ti := selection.Find("p.no-margin")
			title := ti.Text()
			title = strings.TrimSpace(title)
			desc := ""
			if title == "" || link == "" {
				return
			}

			//fmt.Printf("序号:%d\n", i)
			fmt.Println("ID:" + id)
			fmt.Println("标题:" + title)
			fmt.Println("摘要:" + desc)
			fmt.Println("时间:" + time)
			fmt.Println("连接:" + link + "\n")

			g.ArticleId, _ = strconv.Atoi(id)
			g.ArticleTime = time
			g.Link = link
			g.Title = title
			g.Abstract = desc
			where := fmt.Sprintf("%s'%s'", "link=", link)
			ret := gr.Get(where)
			if ret {
				return
			}
			list = append(list, &g)
		})
		gr.Insert(list)

	}
	//ch <- 1
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
func cutStr(ori string, sub string) string {
	//ur := "https://www.skysports.com/f1/news/12433/12623916/monaco-gp-sergio-perez-and-max-verstappen-avoid-post-race-penalties-after-ferrari-protest"
	//o := "www.skysports.com"
	if ori == "" {
		return ""
	}
	pos := strings.Index(ori, sub)
	return ori[pos+len(sub) : len(ori)]
}
