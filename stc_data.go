package main

import (
	"fmt"
	"regexp"
	"log"
	"os"
	"strconv"
    "net/http"
    "net/url"
    "io/ioutil"
)

//指定代理ip
func getTransportFieldURL(proxy_addr *string) (transport *http.Transport) {
    url_i := url.URL{}
    url_proxy, _ := url_i.Parse(*proxy_addr)
    transport = &http.Transport{Proxy : http.ProxyURL(url_proxy)}
    return
}
//从环境变量$http_proxy或$HTTP_PROXY中获取HTTP代理地址
func getTransportFromEnvironment() (transport *http.Transport) {
    transport = &http.Transport{Proxy : http.ProxyFromEnvironment}
    return
}

func fetch(url , proxy_addr *string) (html string) {
    var client *http.Client
    if *proxy_addr == "" {
        client = &http.Client{}
    }else{
      transport := getTransportFieldURL(proxy_addr)
      client = &http.Client{Transport : transport}
    }
    req, err := http.NewRequest("GET", *url, nil)
    if err != nil {
        log.Fatal(err.Error())
    }
    resp, err := client.Do(req)
    if err != nil {
        log.Fatal(err.Error())
    } 
    if resp.StatusCode == 200 {
        robots, err := ioutil.ReadAll(resp.Body);
        resp.Body.Close()
        if err != nil {
            log.Fatal(err.Error())
        }
        html = string(robots);
    } else {
        html = ""
    }
    return
}

func gethis( stockid string ,y string,q string){
	proxy_addr := ""
    url := "http://money.finance.sina.com.cn/corp/go.php/vMS_MarketHistory/stockid/"+stockid+".phtml?year="+y+"&jidu="+q;
    html := fetch(&url, &proxy_addr)
/*
	text := "<td><div align=\"center\">17.980</div></td>'"
	text += "<td><div align=\"center\">18.980</div></td>'"
	text += "<td><div align=\"center\">18.980</div></td>'"
	text += "<td><div align=\"center\">18.980</div></td>'"
	text += "<td><div align=\"center\">18.980</div></td>'"
	text += "<td><div align=\"center\">18.980</div></td>'"
	text += "<td><div align=\"center\">18.980</div></td>'"
	text += "<td><div align=\"center\">18.980</div></td>'"
*/
	regtext := ""
	regtext += "<td><div align=\"center\">(.*?)(\\d+.\\d+)(.*?)</td>"
	regtext += "|<td class=\"tdr\"><div align=\"center\">(.*?)(\\d+.\\d+)(.*?)</td>"
	regtext += "|<td><div align=\"center\">(.*?)\\d{4}-\\d{2}-\\d{2}(.*?)</td>"
	regtext += "|(.*?)\t\\d{4}-\\d{2}-\\d{2}(.*?)</a>"
	regtext += "|(.*?)\t\\d{4}-\\d{2}-\\d{2}(.*?)\t\t\t</div>"

	regdigtext := ""
	regdigtext += "\\d{4}-\\d{2}-\\d{2}"
	regdigtext += "|\\d+.\\d+"

	reg := regexp.MustCompile(regtext)
	ret := reg.FindAllString(html, -1)

	// fmt.Printf("%q\n", ret)

	regdig := regexp.MustCompile(regdigtext)

	for ind,tt := range ret {
		if (ind % 7 == 0){			
			if (ind != 0){ 
				fmt.Println(" ");
			}
			fmt.Print(stockid+" ")
		}		

		rets := regdig.FindAllString(tt, -1)
		
		for _,tts := range rets {
			fmt.Printf("%s", tts)
		}
		
		fmt.Print(" ");	
	}
}

func main() {
	todo := len(os.Args)
	if (todo < 2 ){
		fmt.Print("need id ");	
		return 
	}
	stockid :=  os.Args[1]
	//gethis("600036","2016","2")
	
	for n := 2016; n>=2010 ; n-- {
		for m := 4; m>=1; m-- { 
			gethis(stockid,strconv.Itoa(n),strconv.Itoa(m))
			fmt.Print("\n");	
		}
	}
}