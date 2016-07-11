package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// テンプレートをコンパイルするのは初回のみで良い
	// sync.Onceを使うと複数のgoroutineから呼び出されても引数として渡した関数は一度しか実行されないことが保証される。
	// また、ServerHTTP内でコンパイルされるため一度も呼び出されなければコンパイルされることはない
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	if err := t.templ.Execute(w, nil); err != nil {
		log.Fatal("Template.Execute:", err)
	}
}

func main() {
	http.Handle("/", &templateHandler{filename: "chat.html"})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServer:", err)
	}
}
