package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

var route = mux.NewRouter()

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
	err := r.ParseForm()
	if err != nil {
		// 解析错误，这里应该有错误处理
		fmt.Fprint(w, "请提供正确的数据！")
		return
	}

	fmt.Fprintf(w, "r.Form 中 title的值为: %v <br>", r.FormValue("title"))
	fmt.Fprintf(w, "r.PostForm 中 title的值为: %v <br>", r.PostFormValue("title"))
	fmt.Fprintf(w, "f.Form 中 test 的值为: %v <br>", r.FormValue("test"))
	fmt.Fprintf(w, "r.PostForm 中 test的值为: %v <br>", r.PostFormValue("test"))
}

// 中间件
func forceHTMLMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		// 继续处理请求
		next.ServeHTTP(w, r)
	})
}

func removeTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		}
		next.ServeHTTP(w, r)
	})
}

func articleCreateHandler(w http.ResponseWriter, r *http.Request) {
	html := `
  <!DOCTYPE html>
	<html lang="en">
	<head>
			<title>创建文章 —— 我的技术博客</title>
	</head>
	<body>
			<form action="%s?test=data" method="post">
					<p><input type="text" name="title"></p>
					<p><textarea name="body" cols="30" rows="10"></textarea></p>
					<p><button type="submit">提交</button></p>
			</form>
	</body>
	</html>
	`
	storeURL, _ := route.Get("articles.store").URL()
	fmt.Fprintf(w, html, storeURL)
}

func main() {
	// route := http.NewServeMux()

	route.HandleFunc("/", defaultHandlerFunc).Methods("GET").Name("home")
	route.HandleFunc("/about", aboutHandler).Methods("GET").Name("about")

	route.HandleFunc("/articles/{id:[0-9]+}", articleShowHandler).Methods("GET").Name("articles.show")
	route.HandleFunc("/articles", articlesIndexHandler).Methods("GET").Name("articles.index")
	route.HandleFunc("/articles", articlesStoreHandler).Methods("POST").Name("articles.store")
	route.HandleFunc("/articles/create", articleCreateHandler).Methods("GET").Name("articles.create")

	route.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	// 使用中间件
	route.Use(forceHTMLMiddleware)

	http.ListenAndServe(":3000", removeTrailingSlash(route))
}
