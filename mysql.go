package main

import (
	"fmt"
	"net/http"
	"database/sql"
	_"github.com/go-sql-driver/mysql"
	"log"
)


var(
	Db *sql.DB
	err error
)

func init() {
	Db,err = sql.Open("mysql","root:036068@/muxi") //数据库类型；数据库用户名:数据库密码@[tcp(localhost:3306)]/数据库名
	if err !=nil{									
		log.Println(err)
	}
	log.Println(Db)
}


type User struct{
	Uid int
	Uname string
	Upwd string
}

func CreatUser(name string,pwd string) error {
	sqlStr :="insert into sim(userName,pwd)values(?,?)"
	_,err :=Db.Exec(sqlStr,name,pwd)
	if err !=nil{
		fmt.Println("error",err)
		return err
	} else {
		return nil
	}
}


func IsTrue(name string,pwd string)(string,string,error){
	var username string
	var password string
	var err error
	sqlStr :="select userName,Pwd from sim where userName = ? and pwd = ?"
	row,_:=Db.Query(sqlStr,name,pwd)
	for row.Next(){
		err = row.Scan(&username,&password)
	}
	if err != nil{
		return "","",err
	}
	return username,password,nil
}

func Register(w http.ResponseWriter,r *http.Request){
	Uname :=r.FormValue("username")
	Upwd :=r.FormValue("password")
	err := CreatUser(Uname,Upwd)
	if err != nil{
		fmt.Println(err)
	}else{
		fmt.Fprintln(w,"success")
	}
}

func Login(w http.ResponseWriter,r *http.Request){
	Uname :=r.FormValue("username")
	Upwd :=r.FormValue("password")
	name,_,_ := IsTrue(Uname,Upwd)
	if name == ""{
		fmt.Fprintln(w,"fail")
	}else{
		fmt.Fprintln(w,"success")
		fmt.Println(name)
	}
}

func main(){
	http.HandleFunc("/register",Register)
	http.HandleFunc("/login",Login)
	http.ListenAndServe(":8080",nil)
}