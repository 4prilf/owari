package main

//hello
import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	//"github.com/jinzhu/gorm"
	//_ "github.com/jinzhu/gorm/dialects/mysql"
	//"reflect"

	_ "github.com/go-sql-driver/mysql"
)

type terminal struct {
	ID          int
	TEMPERATURE int
	HUMIDTY     int
	DENSITY     int
	TIME        int
}
type thresh struct {
	ID          int
	TEMPERATURE int
	HUMIDTY     int
	DENSITY     int
	STATE       int
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func polling_post(w http.ResponseWriter, r *http.Request) { //终端模拟设备发送设备号，温度湿度浓度，时间
	//http://127.0.0.1:9090/terminal/polling/post?id=1&temperature=12&humidty=30&density=13&time=202301121230
	r.ParseForm() // 解析参数，默认是不会解析的
	polling_id := strings.Join(r.Form["id"], "")
	polling_temperture := strings.Join(r.Form["temperature"], "")
	polling_humidty := strings.Join(r.Form["humidty"], "")
	polling_density := strings.Join(r.Form["density"], "")
	polling_time := strings.Join(r.Form["time"], "")

	id, _ := strconv.Atoi(polling_id)
	temperture, _ := strconv.Atoi(polling_temperture)
	humidty, _ := strconv.Atoi(polling_humidty)
	density, _ := strconv.Atoi(polling_density)
	time, _ := strconv.Atoi(polling_time)

	polling := terminal{
		ID:          id,
		TEMPERATURE: temperture,
		HUMIDTY:     humidty,
		DENSITY:     density,
		TIME:        time,
	}

	db, err := sql.Open("mysql", "root:m@/terminal?charset=utf8")
	checkErr(err)

	// 插入数据
	stmt, err := db.Prepare("INSERT terminal SET id=?,temperature=?,humidty=?,density=?,time=?")
	checkErr(err)

	res, err := stmt.Exec(polling.ID, polling.TEMPERATURE, polling.HUMIDTY, polling.DENSITY, polling.TIME)
	checkErr(err)

	iid, err := res.LastInsertId()
	checkErr(err)

	fmt.Println("polling_post:iid", iid)
	//新表项 当前state=1时进行比对，state=2说明有问题，state=3时用户已经设置了
	row := db.QueryRow("SELECT temperature,humidty,density,state FROM state where id =?", polling.ID)
	checkErr(err)
	var temperature int
	var humidaty int
	var densiaty int
	var state int
	err = row.Scan(&temperature, &humidaty, &densiaty, &state)
	checkErr(err)
	if state == 1 {
		if (polling.TEMPERATURE > temperature) || (polling.HUMIDTY > humidaty) || (polling.DENSITY > densiaty) {
			db, err := sql.Open("mysql", "root:m@/terminal?charset=utf8")
			checkErr(err)
			astmt, err := db.Prepare("update state set state=? where id=?")
			checkErr(err)
			ares, err := astmt.Exec(2, polling.ID)
			checkErr(err)
			aaffect, err := ares.RowsAffected()
			checkErr(err)
			fmt.Println("polling_post:aaffect", aaffect)
			db.Close()
		}
	}

}

func state_post(w http.ResponseWriter, r *http.Request) { //wehcat设置设备号对应状态
	//http://127.0.0.1:9090/terminal/state/post?id=1&state=3
	r.ParseForm() // 解析参数，默认是不会解析的
	alter_id := strings.Join(r.Form["id"], "")
	alter_state := strings.Join(r.Form["state"], "")
	id, _ := strconv.Atoi(alter_id)
	state, _ := strconv.Atoi(alter_state)
	polling := thresh{
		ID:    id,
		STATE: state,
	}
	db, err := sql.Open("mysql", "root:m@/terminal?charset=utf8")
	checkErr(err)
	stmt, err := db.Prepare("update state set state=? where id=?")
	checkErr(err)

	res, err := stmt.Exec(polling.STATE, polling.ID)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Println("state_post:affect", affect)
	db.Close()
}

func threshold_post(w http.ResponseWriter, r *http.Request) { //wechat设置设备号对应阈值
	//http://127.0.0.1:9090/terminal/threshold/post?id=1&temperature=100&humidty=100&density=100
	r.ParseForm() // 解析参数，默认是不会解析的
	alter_id := strings.Join(r.Form["id"], "")
	alter_temperature := strings.Join(r.Form["temperature"], "")
	alter_humidty := strings.Join(r.Form["humidty"], "")
	alter_density := strings.Join(r.Form["density"], "")
	id, _ := strconv.Atoi(alter_id)
	temperature, _ := strconv.Atoi(alter_temperature)
	humidty, _ := strconv.Atoi(alter_humidty)
	density, _ := strconv.Atoi(alter_density)
	polling := thresh{
		ID:          id,
		TEMPERATURE: temperature,
		HUMIDTY:     humidty,
		DENSITY:     density,
	}
	db, err := sql.Open("mysql", "root:m@/terminal?charset=utf8")
	checkErr(err)
	stmt, err := db.Prepare("update state set temperature=?,humidty=?,density=? where id=?")
	checkErr(err)

	res, err := stmt.Exec(polling.TEMPERATURE, polling.HUMIDTY, polling.DENSITY, polling.ID)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Println("threshold_post:affect", affect)
	db.Close()
}

func polling_get(w http.ResponseWriter, r *http.Request) { //警报
	//http://127.0.0.1:9090/terminal/polling/get?id=1
	r.ParseForm() // 解析参数，默认是不会解析的
	polling_id := strings.Join(r.Form["id"], "")
	id, _ := strconv.Atoi(polling_id)
	polling := terminal{
		ID: id,
	}
	db, err := sql.Open("mysql", "root:m@/terminal?charset=utf8")
	row := db.QueryRow("SELECT state FROM state where id =?", polling.ID)
	checkErr(err)
	var state int
	_ = row.Scan(&state)
	println("polling_get:state", state)
	db.Close()
}

func state_get(w http.ResponseWriter, r *http.Request) { //设备历史状态
	//http://127.0.0.1:9090/terminal/state/get?id=1
	r.ParseForm() // 解析参数，默认是不会解析的
	polling_id := strings.Join(r.Form["id"], "")
	id, _ := strconv.Atoi(polling_id)
	polling := terminal{
		ID: id,
	}
	db, err := sql.Open("mysql", "root:m@/terminal?charset=utf8")
	checkErr(err)
	rows, err := db.Query("SELECT * FROM terminal where id = ? ORDER BY time desc LIMIT 10 ", polling.ID)
	checkErr(err)
	for rows.Next() {
		var id int
		var temperature int
		var humidty int
		var density int
		var time int
		err = rows.Scan(&id, &temperature, &humidty, &density, &time)
		fmt.Println("state_get:",id, temperature, humidty, density, time)
		checkErr(err)
	}
	db.Close()

}

func main() {
	http.HandleFunc("/terminal/polling/post", polling_post)     // 模拟终端上传设备
	http.HandleFunc("/terminal/polling/get", polling_get)       //wechat警报轮询
	http.HandleFunc("/terminal/state/get", state_get)           //wechat单一设备历史记录
	http.HandleFunc("/terminal/state/post", state_post)         //wechat设置设备状态
	http.HandleFunc("/terminal/threshold/post", threshold_post) //wechat设置门限
	err := http.ListenAndServe(":9090", nil)                    // 设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
