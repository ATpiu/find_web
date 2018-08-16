package main
import (
	"fmt"
	"sync"
	"io/ioutil"
	"strings"
	"syscall"
	"github.com/gocolly/colly"
)
var return_string [] string
func read_file(file_name string,wg *sync.WaitGroup) []string{
	b, err := ioutil.ReadFile(file_name)
	if err!=nil{
		fmt.Println(err)
		syscall.Exit(1)
	}
	str := string(b)
	url_list := strings.Split(str, "\r\n")
	wg.Add(len(url_list))
	return url_list
}
func main() {
	var wg sync.WaitGroup
	c := colly.NewCollector() 
	file_name:=""
	url_list:=read_file(file_name,&wg)
	var urlmap map[string]string
	urlmap = make(map[string]string)
	find_web(c,&wg,urlmap )
	for _,url := range url_list {
		go c.Visit(url)
	}
	wg.Wait()
	stringByte := strings.Join(return_string, "\x0D\x0A")
	ioutil.WriteFile("2.txt",[]byte(stringByte),0644)
	fmt.Println("done")

}

func find_web(c *colly.Collector,wg *sync.WaitGroup,urlmap map[string]string) {
	c.OnHTML("title", func(e *colly.HTMLElement) {
		//fmt.Println(e.Response.StatusCode)
		fmt.Println(e.Text,e.Request.URL)
		if e.Text !=""{
			_, ok := urlmap [ e.Text ]
			flag:=web_invalid(e)
			if ok==false && flag==false {
				urlmap [ e.Text ]=e.Request.URL.String()
				fmt.Println(e.Text,e.Request.URL.String())
				return_string=append(return_string,e.Request.URL.String())
			}
		}

	})
	c.OnError(func(response *colly.Response, e error) {
		defer wg.Done()
	})
	c.OnScraped(func(response *colly.Response) {
		defer wg.Done()
	})

}

func web_invalid(e *colly.HTMLElement) bool{
	if e.Text==""{
		return true
	}
	return false

}



