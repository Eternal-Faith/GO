 package main
import (
	"fmt"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"strings"
)

type a struct {
	Email    string `json:"username"`
	Password string `json:"password"`
}
func main () {
// 请求url
requestUrl := "http://pass.muxi-tech.xyz/auth/api/signin"

// 加入表单数据
var user a
// data.Set("username", "2295616516@qq.com")
// data.Set("password", "aGpqMDkxOCsrKw==")
user = a{"2295616516@qq.com", "aGpqMDkxOCsrKw=="}
buf, err := json.MarshalIndent(user, "", " ")
	if err != nil {
		panic(err)
	}
payload := strings.NewReader(string(buf))

req, err := http.NewRequest("POST", requestUrl, payload)
if err != nil {
	panic(err)
	return
}

req.Header.Add("Accept", "*/*")
req.Header.Add("Accept-Encoding", "*gzip, deflate*")
req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
req.Header.Add("Connection", "keep-alive")
req.Header.Add("Content-Length", "62")
req.Header.Add("Content-Type", "text/plain;charset=UTF-8")
req.Header.Add("Host", "pass.muxi-tech.xyz")
req.Header.Add("Origin", "http://pass.muxi-tech.xyz")
req.Header.Add("Referer", "http://pass.muxi-tech.xyz/intro")
req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.87 Safari/537.36 SE 2.X MetaSr 1.0")

res, err := http.DefaultClient.Do(req)
if err != nil {
	panic(err)
	return
}
defer res.Body.Close()
body, err := ioutil.ReadAll(res.Body)
if err != nil {
	panic(err)
	return
}

fmt.Println(string(body))
}