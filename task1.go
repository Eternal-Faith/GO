package main
import (
	"fmt"
	"net/http"
)

type data struct {
	Username string
	Nickname string
	Password string
	Phone string
}
//n是data类型的切片，为这个类型分配了0个元素，预分配1000个元素
var n = make([]data, 0, 1000) //用来记录已注册用户数据
var details data
func register(w http.ResponseWriter,r *http.Request){
	//进行注册并检查是否用户名已经被用过
								//客户端输入栏名称
	details.Username = r.FormValue("username")
	details.Nickname = r.FormValue("nickname")
	details.Password = r.FormValue("password")
	details.Phone = r.FormValue("phone")
	var pass = true 	//pass是bool类型					
	for _, values := range  n{					//查重
		if values.Username == details.Username{
			pass = false
		}
	}
	if pass {			//pass为true时执行 
		n = append(n,details)      //进行记录，写入切片
		cookies := http.Cookie {   //发送 Cookie 到客户端
			Name: details.Username,
			Value: details.Password,
			HttpOnly: true,
		}
			w.Header().Set("Set-cookie", cookies.String())
			w.Write([]byte("您已注册成功"))
	}
	if !pass {			//(!true的值为 false),pass为false时执行？							
		w.Write([]byte("该用户名已存在"))
	}
}
func show(p *data, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%+v", *p)			// %+v先输出字段名字，再输出该字段的值
}
func edit(w http.ResponseWriter, r *http.Request) {
	var name = r.FormValue("username")
	var value = r.FormValue("password")
	cookie, err := r.Cookie(name)
	if cookie != nil {       //若cookie不是0值
		for i := 0;i < len(n); i++ {
			if name == n[i].Username && n[i].Password == value {  //验证身份
				w.Write([]byte("基本信息："))
				show(&n[i], w, r)
				n[i].Password =r.FormValue("password")
				n[i].Nickname =r.FormValue("nickname")
				n[i].Phone = r.FormValue("phone")
				w.Write([]byte("请查看您的修改："))
				show(&n[i],w,r)
			}
		}
	} else {
		fmt.Fprintln(w, err)
	}
}

	func login(w http.ResponseWriter, r *http.Request){
		for _, value := range n {
			if value.Username == r.FormValue("username") && value.Password == r.FormValue("password"){
				cookies := http.Cookie{
					Name:   details.Username,
					Value:  details.Password,
					HttpOnly: true,
				}
				w.Header().Set("Set-Cookie",cookies.String())
				w.Write([]byte("登陆成功"))
		}
		if value.Username == r.FormValue("username") && value.Password != r.FormValue("password"){
				w.Write([]byte("密码错误"))
		}
	}
}
    func All(w http.ResponseWriter, r *http.Request) {
		var t = r.FormValue("username")
		cookie, err := r.Cookie(t)
		if cookie != nil{
			if cookie.Name == t && cookie.Value == r.FormValue("password"){
				for ix, value := range n{
					fmt.Fprintf(w,"用户%d:%+v\n",ix+1,value)
				} 
			}
		} else {
			fmt.Fprintln(w,err)
		}
	}
func main() {
	fmt.Println("running")
	http.HandleFunc("/register",register) //注册
	http.HandleFunc("/login",login)       //登录
	http.HandleFunc("/login/edit",edit)   //登录后可修改用户信息
	http.HandleFunc("/register/All",All)  //查看所有用户信息
	http.ListenAndServe(":80",nil)	
}