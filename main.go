package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

const (
	username = "root"
	password = ""
	portnum  = 3306
	dbname   = "todolist"
)

func main() {
	fmt.Println("Hello world")

	r := mux.NewRouter()

	type Todo struct {
		Id    int
		Title string
		Done  bool
	}

	type TodoPageData struct {
		PageTitle string
		Todos     []Todo
	}

	conf := username + ":" + password + "@(localhost:" + strconv.Itoa(portnum) + ")/" + dbname + "?parseTime=true"

	fmt.Println(conf)

	db, err := sql.Open("mysql", conf)

	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println("no error")
	}

	homeTempl := template.Must(template.ParseFiles("./templates/index.html"))

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		todos, err := db.Query("select * from todos")

		defer todos.Close()

		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("no error selecting todos")
		}

		var todo []Todo

		for todos.Next() {
			var t Todo
			err := todos.Scan(&t.Id, &t.Title, &t.Done)
			// fmt.Println(err)
			if err == nil {
				if t.Done != true {
					fmt.Println(t)
					todo = append(todo, t)
					// fmt.Println(todo)
				}
			} else {
				fmt.Println(err)
				return
			}
		}

		fmt.Println(todo)

		data := TodoPageData{
			PageTitle: "To do list",
			Todos:     todo,
		}

		fmt.Println(data)
		er := homeTempl.Execute(w, data)

		if er != nil {
			fmt.Println(er)
		}
	})

	r.HandleFunc("/remove-todo/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		//_, err := db.Exec(`Delete from todos where id = ?`, id)
		_, err := db.Exec(`UPDATE todos SET done = ? where id = ?`, true, id)
		if err != nil {
			fmt.Println(err)
		} else {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}

	})

	r.HandleFunc("/add-todo", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			title := r.FormValue("todotitle")
			_, err = db.Exec("insert into todos(title) values(?)", title)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("no errors")
				http.Redirect(w, r, "/", http.StatusSeeOther)
			}
		}
	})

	r.
		PathPrefix("/assets/").
		Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("."+"/assets/"))))

	http.ListenAndServe(":8090", r)
}
