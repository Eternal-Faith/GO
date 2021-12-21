package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

var db, err = sql.Open("mysql", "root:036068@(127.0.0.1)/gms") //注意：1.数据库的用户名一般都是root 2.gms：Game management system

func main() {
	if err != nil {
		fmt.Println("连接失败", err)
	} else {
		fmt.Println("连接成功")
	}
	defer db.Close()
	http.HandleFunc("/register", Register)
	http.HandleFunc("/login", Login)
	http.HandleFunc("/login/insert", Insertgame)
	http.HandleFunc("/login/viewall", ViewAll)
	http.HandleFunc("/login/appoint", Appoint)
	http.HandleFunc("/login/viewown", ViewOwn)
	http.HandleFunc("/login/registered-player", RegisteredPlayer)
	http.HandleFunc("/login/registered-team", RegisteredTeam)
	http.HandleFunc("/login/modification", Modification)
	http.HandleFunc("/login/queryplayer", QueryPlayer)
	http.HandleFunc("/login/queryteam", QueryTeam)
	http.HandleFunc("/login/authorization", Authorization)
	err = http.ListenAndServe("localhost:8918", nil) //端口值尽量大些
	if err != nil {
		fmt.Println(err)
	}
}

//不能直接在if语句中初始化
func Insert(name string, password string, sex string, phone string, avatar string, role string) string {
	fmt.Println("开始插入数据")
	sql := "insert into user(name,password,sex,phone,avatar,role) values(?,?,?,?,?,?)"
	_, err := db.Exec(sql, name, password, sex, phone, avatar, role)
	var result string
	if err != nil {
		fmt.Println(err)
		result = "该手机号码已经被绑定"
	} else {
		result = "注册成功"
	}
	return result
}

func Register(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	password := r.FormValue("password")
	sex := r.FormValue("sex")
	phone := r.FormValue("phone")
	avatar := r.FormValue("avatar")
	role := r.FormValue("role")

	result := Insert(name, password, sex, phone, avatar, role)
	fmt.Fprintln(w, result)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var role string
	phone := r.FormValue("phone")
	password := r.FormValue("password")
	//用x-www-form-urlencoded才能正常返回
	sql1 := "select role from user where phone = ? and password = ?"
	err := db.QueryRow(sql1, phone, password).Scan(&role)
	fmt.Println(err)
	if err == sql.ErrNoRows {
		fmt.Fprintln(w, "登陆失败") //失败时err有信息
	} else {
		sql2 := "select role from user where phone = ?"
		row2 := db.QueryRow(sql2, phone) //直接打印row2是一大窜字符串，但可以拿来和真实值一样用
		err = row2.Scan(&role)
		if err != nil {
			fmt.Println(err)
		}

		cookie1 := &http.Cookie{
			Name:  "phone",
			Value: phone,
			// Path:    "/",
			// Expires: time.Time{},
			// MaxAge:  600,
		}
		cookie2 := &http.Cookie{
			Name:  "role",
			Value: role,
		}
		w.Header().Add("Set-Cookie", cookie1.String()) //设置两个cookie用Add函数取代Set
		w.Header().Add("Set-Cookie", cookie2.String()) //cookie不能设置为中文文字
		fmt.Fprintln(w, "登陆成功")                        //必须先对header操作，如果对body操作以后，默认header已经设置完成！！
	}
}

func Authorization(w http.ResponseWriter, r *http.Request) {
	role, _ := r.Cookie("role")
	role_tmp := role.Value
	if role_tmp != "3" {
		fmt.Fprintln(w, "对不起，您没有权限为其它用户授权")
		return
	}
	phone := r.FormValue("phone")
	NewRole := r.FormValue("newrole")
	sql := "update user set role = ? where phone = ?"
	_, err := db.Exec(sql, NewRole, phone)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Fprintln(w, "更改权限成功")
	}
}

func Insertgame(w http.ResponseWriter, r *http.Request) {
	role, _ := r.Cookie("role")
	role_tmp := role.Value
	if role_tmp == "1" {
		fmt.Fprintln(w, "对不起，您没有权限修改球赛数据")
		return
	}
	name := r.FormValue("name")
	data := r.FormValue("data")
	place := r.FormValue("place")
	info := r.FormValue("info")
	appointnum := r.FormValue("appointnum")
	teamA := r.FormValue("teamA")
	teamB := r.FormValue("teamB")
	fmt.Println("开始插入球赛信息")
	sql := "insert into list (name,data,place,info,appointnum,teamA,teamB) values (?,?,?,?,?,?,?)"
	_, err := db.Exec(sql, name, data, place, info, appointnum, teamA, teamB)
	fmt.Println(err)
	if err == nil {
		fmt.Fprintln(w, "插入成功")
	} else {
		fmt.Fprintln(w, "插入失败")
	}
}

//预约比赛/一次send只能预约一场比赛,不过在该路由内还可以进行第二次send
//一定要输入页数，不然会报错
func ViewAll(w http.ResponseWriter, r *http.Request) {
	//先查看：
	var (
		id         string
		name       string
		data       string
		place      string
		info       string
		appointnum string
		teamA      string
		teamB      string
	)
	//分页查询
	var pageno_tmp int
	pageno := r.FormValue("pageno")
	pageno_tmp, _ = strconv.Atoi(pageno)
	pageSize := 10
	startIndex := (pageno_tmp - 1) * pageSize

	//根据需求查看热度排序及筛选
	var rows *sql.Rows
	var sqlx string
	var team_name string
	fliter := r.FormValue("fliter")
	fliter_name := r.FormValue("fliter_name")
	if fliter == "heat" {
		sqlx = "select * from list order by appointnum desc limit ?,? "
		rows, err = db.Query(sqlx, startIndex, pageSize)
		Error(err)
	} else if fliter_name != "" {
		sqlx = "select team from player where name = ? "
		err := db.QueryRow(sqlx, fliter_name).Scan(&team_name)
		Error(err)
		sqlx = "select * from list where teamA=? or teamB = ? limit ?,?"
		rows, err = db.Query(sqlx, team_name, team_name, startIndex, pageSize)
		Error(err)
	} else if fliter == "" {
		sqlx = "select * from list limit ? ,? "
		rows, err = db.Query(sqlx, startIndex, pageSize)
		Error(err)
	}

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &name, &data, &place, &info, &appointnum, &teamA, &teamB)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Fprintln(w, id, name, data, place, info, appointnum, teamA, teamB)
	}
}

func Appoint(w http.ResponseWriter, r *http.Request) {
	//选择比赛并计入用户球赛表
	var tmp string
	Serialnum := r.FormValue("serialnum")

	//计入用户球赛表
	var phonetmp string
	phone, _ := r.Cookie("phone")
	phonetmp = phone.Value
	sql3 := "insert into user_gamelist(user_phone, list_id) values(?,?) "
	_, err4 := db.Exec(sql3, phonetmp, Serialnum)
	if err4 != nil {
		fmt.Fprintln(w, "您已预约过该球赛")
		return
	}

	sql1 := "select appointnum from list where id = ?"
	rows, err1 := db.Query(sql1, Serialnum)
	if err1 != nil {
		fmt.Println(err1)
	}
	defer rows.Close()
	// for rows.Next() { //查到的数据只有一条不需要遍历
	err2 := rows.Scan(&tmp)
	if err2 != nil {
		fmt.Println(err)
	}
	num, _ := strconv.Atoi(tmp)
	num++
	sql2 := "update list set appointnum = ? where id = ? "
	_, err3 := db.Exec(sql2, num, Serialnum)
	if err3 != nil {
		fmt.Println(err3)
	} else {
		fmt.Fprintln(w, "预约成功")
	}

}

func ViewOwn(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "您已预约的球赛如下")
	var list_id string
	var list_idtmp int
	phone, _ := r.Cookie("phone")
	phonenum := phone.Value
	sql := "select list_id from user_gamelist where user_phone = ?"
	rows, _ := db.Query(sql, phonenum)
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&list_id)
		if err != nil {
			fmt.Println(err)
		}
		list_idtmp, _ = strconv.Atoi(list_id)
		var (
			id         string
			name       string
			data       string
			place      string
			info       string
			appointnum string
			teamA      string
			teamB      string
		)
		sql := "select * from list where id = ?"
		rows, err := db.Query(sql, list_idtmp)
		if err != nil {
			fmt.Println(err)
		}
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&id, &name, &data, &place, &info, &appointnum, &teamA, &teamB)
			if err != nil {
				fmt.Println(err)
			}

			fmt.Fprintln(w, id, name, data, place, info, appointnum, teamA, teamB)
		}
	}
}

func RegisteredTeam(w http.ResponseWriter, r *http.Request) {
	role, _ := r.Cookie("role")

	role_tmp := role.Value
	if role_tmp == "1" {
		fmt.Fprintln(w, "对不起，您没有权限注册球队")
		return
	}
	name := r.FormValue("name")
	logo := r.FormValue("logo")
	info := r.FormValue("info")
	sql := "insert into team(name,logo,info) values (?,?,?)"
	_, err := db.Exec(sql, name, logo, info)
	if err != nil {
		fmt.Fprintln(w, "该队名已被注册")
	} else {
		fmt.Fprintln(w, "注册成功")
	}
}

func RegisteredPlayer(w http.ResponseWriter, r *http.Request) {
	role, _ := r.Cookie("role")
	role_tmp := role.Value
	if role_tmp == "1" {
		fmt.Fprintln(w, "对不起，您没有权限登记球员信息")
		return
	}
	name := r.FormValue("name")
	avatar := r.FormValue("avatar")
	team := r.FormValue("team")
	num := r.FormValue("num")
	position := r.FormValue("position")
	age := r.FormValue("age")
	sql := "insert into player(name,avatar,team,num,position,age) values (?,?,?,?,?,?)"
	_, err := db.Exec(sql, name, avatar, team, num, position, age)
	if err != nil {
		fmt.Println(err)
		fmt.Fprintln(w, "注册失败") //
	} else {
		fmt.Fprintln(w, "注册成功")
	}
	//如果团队不为空，上传到成员-团队信息表
	if team != "" {
		var player_id string
		sql1 := "select id from player where name = ?"
		err = db.QueryRow(sql1, name).Scan(&player_id)
		Error(err)
		var team_id string
		sql2 := "select id from team where name = ?"
		err = db.QueryRow(sql2, team).Scan(&team_id)
		Error(err)
		sql3 := "insert into player_team(player_id,team_id) values(?,?)"
		db.Exec(sql3, player_id, team_id)
	}
}

func Error(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func Modification(w http.ResponseWriter, r *http.Request) {
	player_name := r.FormValue("player_name")
	new_team := r.FormValue("new_team")
	var player_id string
	sql := "select id from player where name = ? "
	err := db.QueryRow(sql, player_name).Scan(&player_id)
	Error(err)
	var team_id string
	sql = "select id from team where name = ? "
	err = db.QueryRow(sql, new_team).Scan(&team_id)
	Error(err)
	//先改父表再改子表
	sql = "update player set team = ? where id = ?"
	_, err = db.Exec(sql, new_team, player_id)
	if err != nil {
		fmt.Fprintln(w, "修改失败")
		return
	}

	sql = "update player_team set team_id = ? where player_id = ?"
	_, err = db.Exec(sql, team_id, player_id)
	if err != nil {
		fmt.Fprintln(w, "修改失败")
	} else {
		fmt.Fprintln(w, "修改成功")
	}
}

func QueryPlayer(w http.ResponseWriter, r *http.Request) {
	var (
		id       string
		avatar   string
		team     string
		num      string
		position string
		age      string
	)
	name := r.FormValue("name")
	sql := "select * from player where name = ?"
	err := db.QueryRow(sql, name).Scan(&id, &name, &avatar, &team, &num, &position, &age)
	Error(err)
	fmt.Fprintln(w, id, name, avatar, team, num, position, age)
}

func QueryTeam(w http.ResponseWriter, r *http.Request) {
	var (
		id   string
		logo string
		info string
	)
	name := r.FormValue("name")
	sql := "select * from team where name = ?"
	err := db.QueryRow(sql, name).Scan(&id, &name, &logo, &info)
	Error(err)
	fmt.Fprintln(w, id, name, logo, info)
}
