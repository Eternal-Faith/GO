package main
import (
	"net/http"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	// "github.com/axgle/mahonia"
	)

func main (){
requestUrl := "http://work.muxi-tech.xyz/teamMember"   
// 发送Get请求
rsp, err := http.Get(requestUrl)    
if err != nil {
    log.Println(err.Error())
    return
}

body, err := ioutil.ReadAll(rsp.Body)
if err != nil {
    log.Println(err.Error())
    return
}
content := string(body)
defer rsp.Body.Close()

buf := content
// fmt.Println(buf)

// 解释正则表达式
reg := regexp.MustCompile(`<title>(?s:(.*?))</title>`)
if reg == nil{
	fmt.Println("MustCompile err")
	return
}
// 提取关键信息
result := reg.FindAllStringSubmatch(buf, -1)

// 过滤<> </>
for i, text := range result {
	fmt.Println("text",i+1, "=", text[i])
}
}
