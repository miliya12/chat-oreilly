package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/gorilla/pat"

	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/gplus"

	"github.com/miliya12/chat-oreilly/trace"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, r)
}

func main() {
	var addr = flag.String("addr", ":8080", "アプリケーションのアドレス")
	flag.Parse()

	goth.UseProviders(
		gplus.New("653742812102-6u91isjhgmria659j0309aj5m1v137c7.apps.googleusercontent.com", "HxFqoyT0YOIKSAnPOE1TsTXU", "http://localhost:8080/auth/callback/gplus"),
		github.New("690885fbbe3c254463e8", "907abad6ebc82d6e736cdae74e52ba3a90d7faa6", "http://localhost:8080/auth/callback/github"),
	)

	r := newRoom()
	r.tracer = trace.New(os.Stdout)

	p := pat.New()
	p.Get("/auth/callback/{provider}", callbackHandler)
	p.Get("/auth/{provider}", gothic.BeginAuthHandler)
	p.Add("GET", "/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	p.Add("GET", "/login", &templateHandler{filename: "login.html"})
	p.Add("GET", "/room", r)

	// start chatroom
	go r.run()
	// launch web server
	log.Println("Webサーバーを開始します。ポート: ", *addr)
	if err := http.ListenAndServe(*addr, p); err != nil {
		log.Fatal("Listen and Serve:", err)
	}
}
