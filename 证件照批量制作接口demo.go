package main
//购买地址 https://market.aliyun.com/products/57124001/cmapi030523.html?spm=5176.2020520132.101.3.7a347218HunZ0T#sku=yuncode2452300001
import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
)
/*
返回数据
{
    "type": "jpg", # 图片类型，目前支持"jpg"和"png"两种类型
    "photo": '图片数据BASE64编码',
    "spec": "证件照规格ID",  // 见证件照规格列表
    "bk": "背景颜色" // 见证件照颜色列表
}
参考文档常用数据
1 一寸
2 小一寸
3 大一寸
4 二寸
5 小二寸
6 大二寸
19 三寸
20 四寸
21 五寸

blue 蓝底
red 红底
white 白底
blue_gradual 蓝色渐变
red_gradual 红色渐变
gray_gradual 灰色渐变
*/


func main()  {
	fmt.Println("输入照片尺寸代码\n"+"1 一寸 \n2 小一寸 \n3 大一寸  \n4 二寸 \n5 小二寸 \n6 大二寸 \n19 三寸 \n20 四寸 \n21 五寸")
	var spec,bkId,bk string
	fmt.Scan(&spec)
	fmt.Println("背景颜色代码\n 1  蓝底 \n 2  红底 \n 3  白底 \n 4  蓝色渐变 \n 5  红色渐变 \n 6  灰色渐变")
	fmt.Scan(&bkId)

	switch bkId {
	case "1":
		bk="blue"
	case "2":
		bk="red"
	case "3":
		bk="white"
	case "4":
		bk="blue_gradual"
	case "5":
		bk="red_gradual"
	case "6":
		bk="gray_gradual"
	}
	fmt.Println(bk,spec)
	finfo,err:=ioutil.ReadDir("./")
	if err!=nil{
		fmt.Println("读取文件错误",err)
	}
	bufs :=make(map[string][]byte)
	//获取文件下所有jpg和png图片
	for _,f:=range finfo{
		if f.Size()>8*1024*1024*4{
			fmt.Println(f.Name(),"文件过大")
		}else{
			ext:=path.Ext(f.Name())
			if ext==".jpeg"||ext==".jpg"||ext==".png"{
				file,err:=os.Open(f.Name())
				defer file.Close()
				if err!=nil{
					fmt.Println(err)
					return
				}
				buff,_:=ioutil.ReadAll(file)
				bufs[f.Name()]=buff
			}

		}
	}
	//对数据进行编码
	nameAndBase64:=make(map[string]string)
	for k,v:=range bufs{
		encodeString := base64.StdEncoding.EncodeToString(v)
		nameAndBase64[k]=encodeString
	}
	//调用接口
	for k,v:= range nameAndBase64{
		ext:=path.Ext(k)
		if ext==".jpeg"{
			ext=".jpg"
		}
		var JS = map[string]string{}
		JS["type"]=ext[1:]
		JS["photo"]=v
		JS["spec"]=spec
		JS["bk"]=bk
		url:="http://alidphoto.aisegment.com/idphoto/make"
		js,_:=json.Marshal(JS)
		body:=httpDo(url,&js)
		data:=make(map[string]interface{})
		json.Unmarshal(body,&data)
		fmt.Println(data)
		//处理数据
		if data["errmsg"].(string)!="SUCCESSFULLY"{
			fmt.Println("出现错误",err)
			continue
		}
		//获取照片地址
		tmp:=data["data"].(map[string]interface{})
		Url:=tmp["result"].(string)
		//写入数据到文件
		resp,err:=http.Get(Url)
		if err!=nil{
			return
		}
		defer resp.Body.Close()
		Body,_ := ioutil.ReadAll(resp.Body)
		fw,_:=os.Create("./ok"+k)
		defer fw.Close()
		fw.WriteString(string(Body))

	}

}



func httpDo(url string ,js *[]byte) []byte {
	client := &http.Client{}
	req, err := http.NewRequest("POST", url,bytes.NewReader(*js))
	if err != nil {
		// handle error
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Authorization", "APPCODE 你自己的AppCode ")
	resp, err := client.Do(req)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}
	return body
}
