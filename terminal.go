package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"strconv"
	//"github.com/jinzhu/gorm"
	//_ "github.com/jinzhu/gorm/dialects/mysql"
	//"reflect"

	_ "github.com/go-sql-driver/mysql"
)

//http://localhost:9090/?/url_long=111&url_long=222								//从？开始识别内容
/*func sayhelloName(w http.ResponseWriter, r *http.Request) { //http://localhost:9090/?url_long=111&url_long=222
	r.ParseForm()                       // 解析参数，默认是不会解析的
	fmt.Println(r.Form)                 // 这些信息是输出到服务器端的打印信息 //map[url_long:[111 222]]
	fmt.Println("path", r.URL.Path)     //path /
	fmt.Println("scheme", r.URL.Scheme) //scheme
	fmt.Println(r.Form["url_long"])     //[111 222]
	for k, v := range r.Form {
		fmt.Println("key:", k)                   //key: url_long
		fmt.Println("val:", strings.Join(v, "")) //val: 111222
	}
	fmt.Fprintf(w, "Hello astaxie!") // 这个写入到 w 的是输出到客户端的]
	fmt.Println("")
}
*/
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
	TYPE        int
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func polling_post(w http.ResponseWriter, r *http.Request) { //终端模拟设备发送设备号，温度湿度浓度，时间
	//http://127.0.0.1:9090/terminal/polling/post?id=111&temperature=12&humidty=30&density=13&time=202301121230
	r.ParseForm() // 解析参数，默认是不会解析的
/*
	fmt.Println(r.Form["id"])
	fmt.Println(r.Form["temperature"])
	fmt.Println(r.Form["humidty"])
	fmt.Println(r.Form["density"])
	fmt.Println(r.Form["time"]) //time

	polling_id :=r.Form["id"]
	polling_temperture := r.Form["temperature"]
	polling_humidty := r.Form["humidty"]
	polling_density := r.Form["density"]
	polling_time := r.Form["time"]
	

	fmt.Println(reflect.TypeOf(polling_id),polling_id)
	fmt.Println(reflect.TypeOf(polling_temperture),polling_temperture)
	fmt.Println(reflect.TypeOf(polling_humidty),polling_humidty)
	fmt.Println(reflect.TypeOf(polling_density),polling_density)
	fmt.Println(reflect.TypeOf(polling_time),polling_time)
	*/
	ing_id := strings.Join(r.Form["id"],"")
	ing_temperture := strings.Join(r.Form["temperature"],"")
	ing_humidty := strings.Join(r.Form["humidty"],"")
	ing_density := strings.Join(r.Form["density"],"")
	ing_time := strings.Join(r.Form["time"],"")
/*
	fmt.Println(reflect.TypeOf(ing_id),ing_id)
	fmt.Println(reflect.TypeOf(ing_temperture),ing_temperture)
	fmt.Println(reflect.TypeOf(ing_humidty),ing_humidty)
	fmt.Println(reflect.TypeOf(ing_density),ing_density)
	fmt.Println(reflect.TypeOf(ing_time),ing_time)
	*/
//	apple:="123"
//	fmt.Println(reflect.TypeOf(apple),apple)

	id,_ :=strconv.Atoi(ing_id)
	temperture,_ :=strconv.Atoi(ing_temperture)
	humidty,_ := strconv.Atoi(ing_humidty)
	density,_ :=strconv.Atoi(ing_density)
	time,_ :=strconv.Atoi(ing_time)
/*
	fmt.Println(reflect.TypeOf(id),id)
	fmt.Println(reflect.TypeOf(temperture),temperture)
	fmt.Println(reflect.TypeOf(humidty),humidty)
	fmt.Println(reflect.TypeOf(density),density)
	fmt.Println(reflect.TypeOf(time),time)
	*/
	polling:=terminal{
		ID: id,
		TEMPERATURE: temperture,
		HUMIDTY: humidty,
		DENSITY: density,
		TIME: time,
	}

	db, err := sql.Open("mysql", "root:m@/terminal?charset=utf8")
	checkErr(err)

	// 插入数据
	stmt, err := db.Prepare("INSERT terminal SET id=?,temperature=?,humidty=?,density=?,time=?")
	checkErr(err)

	res, err := stmt.Exec(polling.ID,polling.TEMPERATURE,polling.HUMIDTY,polling.DENSITY,polling.TIME)
	checkErr(err)

	iid, err := res.LastInsertId()
	checkErr(err)

	fmt.Println(iid)
	db.Close()


}

func state_post(w http.ResponseWriter, r *http.Request) { //wehcat设置设备号对应状态
	//http://127.0.0.1:9090/terminal/state/post?id=111&state=1
	r.ParseForm() // 解析参数，默认是不会解析的
	fmt.Println(r.Form["id"])
	fmt.Println(r.Form["state"])
}

func threshold_post(w http.ResponseWriter, r *http.Request) { //wechat设置设备号对应阈值
	//http://127.0.0.1:9090/terminal/threshold/post?id=111&temperature=12&humidty=30&density=13
	r.ParseForm() // 解析参数，默认是不会解析的
	fmt.Println(r.Form["id"])
	fmt.Println(r.Form["temperature"])
	fmt.Println(r.Form["humidty"])
	fmt.Println(r.Form["density"])
}

func polling_get(w http.ResponseWriter, r *http.Request) {
	//http://127.0.0.1:9090/terminal/polling/get?id=11
}

func state_get(w http.ResponseWriter, r *http.Request) {
	//http:127.0.0.1:9090/terminal/state/get?id=111
}

func main() {
	http.HandleFunc("/terminal/polling/post", polling_post)     // 模拟终端上传设备
	http.HandleFunc("/terminal/polling/get", polling_get)       //wechat警报轮询
	http.HandleFunc("/terminal/state/get", state_get)           //wechat单一设备历史记录
	http.HandleFunc("/terminal/threshold/post", threshold_post) //wechat设置门限
	http.HandleFunc("/terminal/state/post", state_post)         //wechat设置设备状态
	err := http.ListenAndServe(":9090", nil)                    // 设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
