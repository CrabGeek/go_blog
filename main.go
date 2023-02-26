package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func defaultHandlerFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Hello, 欢迎来到 goblog! </h1>")
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "此博客是用以记录编程笔记，如您有反馈或建议，请联系 "+
		"<a href=\"mailto:summer@example.com\">summer@example.com</a>")
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>请求页面未找到 :(</h1><p>如有疑惑，请联系我们。</p>")
}

func articleShowHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Fprint(w, "文章ID: "+id)
}

func articlesIndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "访问文章列表")
}

func articlesStoreHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "创建新的文章")
}

// 中间件
func forceHTMLMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		// 继续处理请求
		next.ServeHTTP(w, r)
	})
}

func main() {
	// route := http.NewServeMux()
	route := mux.NewRouter()

	route.HandleFunc("/", defaultHandlerFunc).Methods("GET").Name("home")
	route.HandleFunc("/about", aboutHandler).Methods("GET").Name("about")

	route.HandleFunc("/articles/{id:[0-9]+}", articleShowHandler).Methods("GET").Name("articles.show")
	route.HandleFunc("/articles", articlesIndexHandler).Methods("GET").Name("articles.index")
	route.HandleFunc("/articles", articlesStoreHandler).Methods("POST").Name("articles.store")

	route.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	// 使用中间件
	route.Use(forceHTMLMiddleware)

	homeURL, _ := route.Get("home").URL()
	fmt.Println("homeURL: ", homeURL)

	articleURL, _ := route.Get("articles.show").URL("id", "23")
	fmt.Println("articleURL: ", articleURL)

	http.ListenAndServe(":3000", route)
}
