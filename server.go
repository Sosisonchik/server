package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type Girl struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Imgurl  string `json:"imgurl"`
	Like    int    `json:"like"`
	Dislike int    `json:dislike`
}

type Girls struct {
	List []struct {
		Id      int    `json:"id"`
		Name    string `json:"name"`
		Imgurl  string `json:"imgurl"`
		Like    int    `json:"like"`
		Dislike int    `json:dislike`
	} `json:"list"`
}

type Adder struct {
	Name   string `json:"name"`
	Imgurl string `json:"imgurl"`
}

// Индентификатор для методов: like,dislike,unlike,undislike
type Ind struct {
	ID int `json:"id"`
}

func main() {
	http.HandleFunc("/getGirls", getGirls)
	http.HandleFunc("/add", add)
	http.HandleFunc("/like", like)
	http.HandleFunc("/dislike", dislike)
	http.HandleFunc("/unlike", unlike)
	http.HandleFunc("/undislike", undislike)

	http.ListenAndServe(":8080", nil)
}

func getGirls(w http.ResponseWriter, r *http.Request) {
	fmt.Println("New query (getGirls)")
	db, er := sql.Open("sqlite3", "girlsdb")
	if er != nil {
		fmt.Println(er.Error())
		return
	}

	rows, err := db.Query("SELECT * FROM Girl")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var girl Girl
	var girls Girls
	for rows.Next() {
		err := rows.Scan(&girl.Id, &girl.Name, &girl.Imgurl, &girl.Like, &girl.Dislike)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		girls.List = append(girls.List, girl)
	}
	rows.Close()

	data, e := json.Marshal(&girls)
	if e != nil {
		fmt.Println(e.Error())
	}
	w.Write(data)
}

func add(w http.ResponseWriter, r *http.Request) {
	fmt.Println("New query (add)")
	db, er := sql.Open("sqlite3", "girlsdb")
	if er != nil {
		fmt.Println(er.Error())
		return
	}

	defer r.Body.Close()
	var adder Adder
	data, e := ioutil.ReadAll(r.Body)
	if e != nil {
		fmt.Println(e.Error())
		return
	}

	err := json.Unmarshal(data, &adder)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	res, e := db.Exec("INSERT INTO Girl (name, imgurl, like_field, dislike_field) VALUES(?,?,0,0)", adder.Name, adder.Imgurl)
	if e != nil {
		fmt.Println(e.Error())
	}

	i, _ := res.RowsAffected()

	fmt.Println(i)

}

func like(w http.ResponseWriter, r *http.Request) {
	fmt.Println("New query (like)")
	defer r.Body.Close()

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var ind Ind
	er := json.Unmarshal(data, &ind)
	if er != nil {
		fmt.Println(er.Error())
		return
	}

	db, e := sql.Open("sqlite3", "girlsdb")
	if e != nil {
		fmt.Println(e.Error())
		return
	}

	db.Exec("UPDATE Girl SET like_field=like_field+1 WHERE id=?", ind.ID)
	fmt.Println(ind.ID, "was incremented")

}

func dislike(w http.ResponseWriter, r *http.Request) {
	fmt.Println("New query (dislike)")
	defer r.Body.Close()

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var ind Ind
	er := json.Unmarshal(data, &ind)
	if er != nil {
		fmt.Println(er.Error())
		return
	}

	db, e := sql.Open("sqlite3", "girlsdb")
	if e != nil {
		fmt.Println(e.Error())
		return
	}

	db.Exec("UPDATE Girl SET dislike_field=dislike_field+1 WHERE id=?", ind.ID)
	fmt.Println(ind.ID, "was incremented")
}

func unlike(w http.ResponseWriter, r *http.Request) {
	fmt.Println("New query (unlike)")
	defer r.Body.Close()

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var ind Ind
	er := json.Unmarshal(data, &ind)
	if er != nil {
		fmt.Println(er.Error())
		return
	}

	db, e := sql.Open("sqlite3", "girlsdb")
	if e != nil {
		fmt.Println(e.Error())
		return
	}

	var like int
	row := db.QueryRow("SELECT like_field FROM Girl WHERE id=?", ind.ID)
	row.Scan(&like)

	if like > 0 {
		db.Exec("UPDATE Girl SET like_field=like_field-1 WHERE id=?", ind.ID)
	}

	fmt.Println(ind.ID, "was unincremented")
}

func undislike(w http.ResponseWriter, r *http.Request) {
	fmt.Println("New query (undislike)")
	defer r.Body.Close()

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var ind Ind
	er := json.Unmarshal(data, &ind)
	if er != nil {
		fmt.Println(er.Error())
		return
	}

	db, e := sql.Open("sqlite3", "girlsdb")
	if e != nil {
		fmt.Println(e.Error())
		return
	}

	var dislike int
	row := db.QueryRow("SELECT dislike_field FROM Girl WHERE id=?", ind.ID)
	row.Scan(&dislike)

	if dislike > 0 {
		db.Exec("UPDATE Girl SET dislike_field=dislike_field-1 WHERE id=?", ind.ID)
	}

	fmt.Println(ind.ID, "was unincremented")
}
