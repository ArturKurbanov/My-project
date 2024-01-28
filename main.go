package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

// ПОЛУЧАЕМ ДАННЫЕ ИЗ ТАБЛИЦЫ БД И ВЫВОДИМ НА ГЛАВНОЙ СТРАНИЦЕ
type Article struct {
	Id                     uint16
	Title, Anons, FullText string
}

var posts = []Article{} //СОЗДАЕМ СПИСОК КОТОРЫЙ БУДЕМ НАПОЛНЯТЬ
var showPost = Article{}

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")

	if err != nil {
		fmt.Fprint(w, err.Error())
	}

	//---------------------------------------------------
	// ПОДКЛЮЧЕНИЕ К БАЗЕ ДАННЫХ
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/golang")
	if err != nil {
		panic(err)
	}

	defer db.Close()
	//--------------------------------------------------
	//Выборка данных

	res, err := db.Query("SELECT * FROM `articles`")
	if err != nil {
		panic(err)
	}
	posts = []Article{}
	for res.Next() {
		var post Article
		err = res.Scan(&post.Id, &post.Title, &post.Anons, &post.FullText)
		if err != nil {
			panic(err)
		}

		posts = append(posts, post) // append позволяет внутрь какого либо списка добавить новые элементы

		//Каждый пост помещаем в некий список, а дальше этот список будет передавать в сам шаблон
	}

	t.ExecuteTemplate(w, "index", posts) // posts передали в сам шаблон
}

func create(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/create.html", "templates/header.html", "templates/footer.html")

	if err != nil {
		fmt.Fprint(w, err.Error())
	}

	t.ExecuteTemplate(w, "create", nil)
}

func save_article(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")         //Получаем данные из формочки
	anons := r.FormValue("anons")         //Получаем данные из формочки
	full_text := r.FormValue("full_text") //Получаем данные из формочки

	if title == "" || anons == "" || full_text == "" { // ПРОВЕРКИ ПРИ ДОБАВЛЕНИИ ДАННЫХ
		fmt.Fprintf(w, "Не все данные заполнены")
	} else {
		//---------------------------------------------------
		// ПОДКЛЮЧЕНИЕ К БАЗЕ ДАННЫХ
		db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/golang")
		if err != nil {
			panic(err)
		}

		defer db.Close()
		//--------------------------------------------------

		//------------------------------------------------------
		// ДОБАВЛЕНИЕ ЗАПИСЕЙ ВНУТРЬ САМОЙ ТАБЛИЦЫ В БД ЧЕРЕЗ СТРАНИЦУ
		// Установка данных

		insert, err := db.Query(fmt.Sprintf("INSERT INTO `articles`(`title`, `anons`, `full_text`) VALUES('%s', '%s', '%s')", title, anons, full_text))
		if err != nil {
			panic(err)
		}
		defer insert.Close()

		http.Redirect(w, r, "/", http.StatusSeeOther) //ПЕРЕАДРЕСАЦИЯ НА ДРУГУЮ СТРАНИЧКУ ПОСЛЕ ДОБАВЛЕНИЯ ДАННЫХ В БД

		//---------------------------------------------------
	}
}

// Получаем параметр из URL АДРЕСА
func show_post(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Создаем объект vars и передаем параметр "r"
	//---------------------------------------------------
	//Подключаем определенные шаблоны
	t, err := template.ParseFiles("templates/show.html", "templates/header.html", "templates/footer.html")

	if err != nil {
		fmt.Fprint(w, err.Error())
	}
	//---------------------------------------------------
	//---------------------------------------------------
	// ПОДКЛЮЧЕНИЕ К БАЗЕ ДАННЫХ
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/golang")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	//--------------------------------------------------
	//Выборка данных
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		panic(err)
	}

	res, err := db.Query("SELECT * FROM `articles` WHERE `id` = ?", id)
	if err != nil {
		panic(err)
	}

	showPost = Article{}
	for res.Next() {
		var post Article
		err = res.Scan(&post.Id, &post.Title, &post.Anons, &post.FullText)
		if err != nil {
			panic(err)
		}

		showPost = post
	}

	t.ExecuteTemplate(w, "show", showPost) // posts передали в сам шаблон
}

func handleFunc() {
	rtr := mux.NewRouter()                                        //Отслеживание URL адресов
	rtr.HandleFunc("/", index).Methods("GET")                     // ".Methods("GET")"-Указываем метод передачи данных за счет которого полюьзователь сможет подключаться по этому URL адресу
	rtr.HandleFunc("/create", create).Methods("GET")              // "GET"- ЭТО ТАКОЙ МЕТОД ПЕРЕДАЧИ ДАННЫХ КОГДА МЫ НАПРЯМУЮ ЗАХОДИМ НА САМУ СТРАНИЦУ. Например по ссылке или URL АДРЕСУ
	rtr.HandleFunc("/save_article", save_article).Methods("POST") //"POST"- Это такое метод передачи данных который срабатывает при отправки данных из какой либо формы
	rtr.HandleFunc("/post/{id:[0-9]+}", show_post).Methods("GET") // Отслеживаем URL адрес "post" c уникальным идентификатором "{id:[0-9]+}".Указываем какой метод будет все это обрабатывать "show_post"

	http.Handle("/", rtr) // Указываем  что обработка всех URL адресов "/", будет происходить через такой обЪект "rtr"
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.ListenAndServe(":8080", nil)
}

func main() {
	handleFunc()
}
