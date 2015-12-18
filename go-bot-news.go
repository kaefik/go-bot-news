// go-bot-news
package main

import (
	"fmt"
	"go-bot-news/pkg"
	"go-bot-news/pkg/html"
	"golang.org/x/net/html/charset"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type News struct {
	url     string //урл новости
	title   string // заголовок новости
	content string // содержимое новости
}

//инициализация лог файла
func InitLogFile(namef string) *log.Logger {
	file, err := os.OpenFile(namef, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file", os.Stderr, ":", err)
	}
	multi := io.MultiWriter(file, os.Stdout)
	LFile := log.New(multi, "Info: ", log.Ldate|log.Ltime|log.Lshortfile)
	return LFile
}

//получение страницы из урла url
func gethtmlpage(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("HTTP error:", err)
		panic("HTTP error")
	}
	defer resp.Body.Close()
	// вот здесь и начинается самое интересное
	utf8, err := charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
	if err != nil {
		fmt.Println("Encoding error:", err)
		panic("Encoding error")
	}
	body, err := ioutil.ReadAll(utf8)
	if err != nil {
		fmt.Println("IO error:", err)
		panic("IO error")
	}
	return body
}

//удаление повторных элементов в массиве
func delpovtor(s []string) []string {
	if len(s)==0 {
		return make([]string,0)
	}
	fl := false
	st := make([]string, 0)
	st = append(st, s[0])
	for i := 0; i < len(s); i++ {
		fl = true
		for j := 0; j < len(st); j++ {
			if s[i] == st[j] {
				fl = false
			}
		}
		if fl {
			st = append(st, s[i])
		}
	}
	return st
}

// вывод на пкчать массива строк
func printarray(s []string) {
	for i := 0; i < len(s); i++ {
		fmt.Println(s[i])
	}
	return
}

//--------------- парсинг Эха Москвы
func GetNewsEchoMsk() []News {
	url := "http://echo.msk.ru/"	
	n := make([]News, 0)
	ss := GetNewsUrlEchoMsk(url)
	for i := 0; i < len(ss); i++ {
		n = append(n, News{url: ss[i]})
	}
	for i := 0; i < len(n); i++ {
		n[i].ParserNewsEchoMsk()
	}
	return n
}

//получение урлы новостей с главной страницы
func GetNewsUrlEchoMsk(url string) []string {
	//	var ss []string
	if url == "" {
		return make([]string, 0)
	}
	body := gethtmlpage(url)
	shtml := string(body)

	// <a rel="nofollow" href="/likes/e1678230/" class="share" data-url="http://echo.msk.ru/news/1678230-echo.html" data-title="Новое уголовное дело о ремонте кораблей Северного флота поступило в суд">
	snewsmusor, _ := pick.PickAttr(&pick.Option{&shtml, "a", nil}, "data-url")
	snews := make([]string, 0)
	for i := 0; i < len(snewsmusor); i++ {
		if strings.Contains(snewsmusor[i], "-echo.htm") && (strings.Contains(snewsmusor[i], "/news/")) {
			snews = append(snews, snewsmusor[i])
		}
	}

	//	printarray(delpovtor(snews))

	return delpovtor(snews)
}

//парсер новостей с сайта Эха Москвы
func (this *News) ParserNewsEchoMsk() {

	if this.url == "" {
		return
	}
	body := gethtmlpage(this.url)
	shtml := string(body)

	//	<meta property="og:title" content="Новости / 17 декабря, 16:31 | Путин утверждает, что  никогда  не  обсуждал  с  региональными  лидерами расследование конкретных  уголовных  дел" />

	stitle, _ := pick.PickAttr(&pick.Option{&shtml, "meta", &pick.Attr{"property", "og:title"}}, "content")
	if len(stitle) > 0 {
		this.title = stitle[0]
	}

	//	<meta property="og:description" content="
	//В   том числе дела об убийстве    Бориса  Немцова. «Следствие должно установить, как бы долго оно ни продолжалось. Это преступление должно быть расследовано и участники должны быть наказаны, кто бы это ни был, — сказал глава государства." />
	scont, _ := pick.PickAttr(&pick.Option{&shtml, "meta", &pick.Attr{"property", "og:description"}}, "content")
	this.content = scont[0]

	return
}

//--------------- END парсинг Эха Москвы

//--------------- парсинг РБК

func GetNewsRbc() []News {
	url := "http://rt.rbc.ru/"	
	n := make([]News, 0)
	ss := 	GetNewsUrlRbc(url)
	
	for i := 0; i < len(ss); i++ {
		n = append(n, News{url: ss[i]})
	}
	for i := 0; i < len(n); i++ {
		n[i].ParserNewsRbc()
	}
	return n
}

//получение урлы новостей с главной страницы
func GetNewsUrlRbc(url string) []string {
	//	var ss []string
	if url == "" {
		return make([]string, 0)
	}
	body := gethtmlpage(url)
	shtml := string(body)

	// <a href="http://www.rbc.ru/politics/18/12/2015/5673fcd39a794764ce0cd14e" class="news-main-feed__item__link chrome" data-ati-item="item_1" data-ati-title="%D0%95%D0%B2%D1%80%D0%BE%D0%BA%D0%BE%D0%BC%D0%B8%D1%81%D1%81%D0%B8%D1%8F+%D1%80%D0%B5%D0%BA%D0%BE%D0%BC%D0%B5%D0%BD%D0%B4%D0%BE%D0%B2%D0%B0%D0%BB%D0%B0+%D0%BE%D1%82%D0%BC%D0%B5%D0%BD%D0%B8%D1%82%D1%8C+%D0%B2%D0%B8%D0%B7%D0%BE%D0%B2%D1%8B%D0%B9+%D1%80%D0%B5%D0%B6%D0%B8%D0%BC+%D0%B4%D0%BB%D1%8F%D0%A3%D0%BA%D1%80%D0%B0%D0%B8%D0%BD%D1%8B" data-ati-id="5673fcd39a794764ce0cd14e" data-ati-url="http://www.rbc.ru/politics/18/12/2015/5673fcd39a794764ce0cd14e">
	snewsmusor, _ := pick.PickAttr(&pick.Option{&shtml, "a", &pick.Attr{"class", "news-main-feed__item__link chrome"}}, "href")
	snews := snewsmusor
//	make([]string, 0)
	
//	fmt.Println(snewsmusor)
	
//	for i := 0; i < len(snewsmusor); i++ {
//		if strings.Contains(snewsmusor[i], "-echo.htm") && (strings.Contains(snewsmusor[i], "/news/")) {
//			snews = append(snews, snewsmusor[i])
//		}
//	}

	//	printarray(delpovtor(snews))

	return delpovtor(snews)
}

//парсер новостей с сайта РБК
func (this *News) ParserNewsRbc() {

	if this.url == "" {
		return
	}
	body := gethtmlpage(this.url)
	shtml := string(body)

	//    <div class="article__overview__text">Еврокомиссия констатировала выполнение Украиной всех требований плана действий визовой либерализации. Еврочиновники в своем новом отчете рекомендуют Евросовету и Европарламенту начать процесс отмены виз для украинцев</div>

	stitle, _ := pick.PickText(&pick.Option{ 
		&shtml,
		"div",
		&pick.Attr{
			"class",
			"article__overview__text",
		},
	})
	
//	fmt.Println(stitle)
	
	if len(stitle) > 0 {
		this.title = stitle[0]
	}

	//	<meta property="og:description" content="
	//В   том числе дела об убийстве    Бориса  Немцова. «Следствие должно установить, как бы долго оно ни продолжалось. Это преступление должно быть расследовано и участники должны быть наказаны, кто бы это ни был, — сказал глава государства." />
	scont, _ := pick.PickAttr(&pick.Option{&shtml, "meta", &pick.Attr{"property", "og:description"}}, "content")
	this.content = scont[0]

	return
}

//--------------- END парсинг РБК



// генерация html главной страницы
func Htmlpage(sn []News) string {
	zagol := "ГРАББЕР НОВОСТЕЙ"
	begstr := "<html>\n <head>\n <meta charset='utf-8'>\n <title>" + zagol + "</title>\n </head>\n <body>\n"
	//	<h3 id=”Razdel2”> Раздел2 </h3>
	bodystr := "<h1 align=\"center\"><a name=\"MainPage\"> ГРАББЕР НОВОСТЕЙ </a></h1><br>"
	bodystr += HtmlNews(sn,"EchoMSK")
	endstr := "</body>\n" + "</html>"
	return begstr + bodystr + endstr
}

// шаблон оформления новости из одного ресурса
func HtmlNews(sn []News,titlenews string) string{
	bodystr := "<h3 align=\"center\"><a name=\""+titlenews+"\"> "+titlenews+" </a></h3><br>" + "<TABLE align=\"center\" border=\"1\">"
	for i := 0; i < len(sn); i++ {
		bodystr += "<TR> <TD width=\"350\"> <b>" + genhtml.Link(sn[i].title, sn[i].url) + "</b></TD>" + "<TD width=\"550\"><br>" + sn[i].content + "" + "<br> <a href=\"#MainPage\"> В начало </a>" + " <a href=\"#"+titlenews+"\"> К "+titlenews+" </a> "+ "</TD> </TR>"
	}
	bodystr += "</TABLE>"	
	return bodystr
}


func main() {
//	fmt.Println("Starting программы")
	
    n:=GetNewsEchoMsk()	
	str := Htmlpage(n)	
	genhtml.Savestrtofile("news.html", str)
	
	rbc:=GetNewsRbc()
	s:=Htmlpage(rbc)
	genhtml.Savestrtofile("rbc.html", s)

//	fmt.Println("Ending программы")
}
