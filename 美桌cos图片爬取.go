package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

func HTTPCL(str string) (shuju []string) {
	ret:=regexp.MustCompile(`meinv([0-9]{6}).html`)// <a href="http://www.win4000.com/meinv175542.html" target="_blank">
	shuju1:=ret.FindAllStringSubmatch(str,-1)
	shuju=make([]string,0)
	for _,v:=range shuju1{
		shuju=append(shuju,v[1])
	}

	return
}
func ZWYCL(str string) (string) {
	ret:=regexp.MustCompile(`em>([0-9]*?)</em`)// <a href="http://www.win4000.com/meinv175542.html" target="_blank">
	shuju:=ret.FindAllStringSubmatch(str,1)

	return shuju[0][1]
}
func zymtp(str string) (page string) {
	ret:=regexp.MustCompile(`data-original="(.*?)"`)//data-original="http://pic1.win4000.com/pic/d/ea/61d85feb7f.jpg"
	shuju:=ret.FindAllStringSubmatch(str,1)
	page=shuju[0][1]
	return
}
func FANDWJM(str string) (page string) {
	ret:=regexp.MustCompile(`ptitle"><h1>(.*?)</h1>`)//data-original="http://pic1.win4000.com/pic/d/ea/61d85feb7f.jpg"
	shuju:=ret.FindAllStringSubmatch(str,1)
	page=shuju[0][1]
	return
}
func baocun(str string,pathname string){
	f,err:=os.Create(pathname)
	if err!=nil{
		fmt.Println(err)
		os.Exit(-1)
	}
	f.WriteString(str)
	f.Close()
}
func HTTPGET(url string)(back string,err error)  {

	//获取网页内容
	resp,err:=http.Get(url)
	if err!=nil{
		return
	}
defer resp.Body.Close()
	buf:=make([]byte,4096)
	//循环读取内容
	for  {
		n,err:=resp.Body.Read(buf)
		if n==0 {
			break
		}
		if err!=nil&&err!=io.EOF{
			break
		}
		back+=string(buf[:n])
	}
	return
}
func ERR(s string,err error)  {
	if err!=nil{
		fmt.Println(s,err)
	}
}
func main()  {
	for i:=1;i<=5 ;i++  {
		 PAGE(i)
	}

}
func PAGE(i int) {

	url:=`http://www.win4000.com/meinvtag26_`+strconv.Itoa(i)+`.html`
	STR,err:=HTTPGET(url)
	ERR("get",err)
	//fmt.Println(back)
	wzshuju:=HTTPCL(STR)
	fmt.Println(wzshuju)
//循环取出共有多少页面
for i:=0;i<len(wzshuju);i++{
	zymurl:="http://www.win4000.com/meinv"+wzshuju[i]+".html"
	PAGEwjm(zymurl,wzshuju,i)
	}
for i:=0;i<len(wzshuju);i++{


}

}
func PAGEwjm(zymurl string,wzshuju []string ,i int){

	zymsj,err:=HTTPGET(zymurl)
	ERR("get",err)
	//取出总页数和图片名称
	Page:=ZWYCL(zymsj)
	//fmt.Println(Page)
	WJM:=FANDWJM(zymsj)
	//fmt.Println(WJM)
	paend,err:=strconv.Atoi(Page)
	for j:=1;j<=paend;j++{
		URL:="http://www.win4000.com/meinv"+wzshuju[i]+"_"+strconv.Itoa(j)+".html"
		wjianxiazai(URL,WJM,j)

	}

}
func wjianxiazai(URL,WJM string ,j int){

	YMXX,_:=HTTPGET(URL)
	WJXX:=zymtp(YMXX)
	//	保存文件
	WJ,_:=HTTPGET(WJXX)
	pathname:="/Users/lee/Desktop/go/数据/"+WJM+strconv.Itoa(j)+".jpg"//WJM+"/"+strconv.Itoa(j)+".jpg"
	fmt.Println(pathname)
	baocun(WJ,pathname)

}