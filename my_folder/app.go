package app

import (
	"encoding/json"
	"html/template"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

func init() {
	http.HandleFunc("/", handleExample)
}

// このディレクトリーに入っているすべての「.html」終わるファイルをtemplateとして読み込む。
var tmpl = template.Must(template.ParseGlob("*.html"))

// Templateに渡す内容を分かりやすくするためのtypeを定義しておきます。
// （「Page」という名前などは重要ではありません）。
type Page struct {
	A string
	B string
}

func handleExample(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// templateに埋める内容をrequestのFormValueから用意する。
	content := Page{
		A: r.FormValue("a"),
		B: r.FormValue("b"),
	}

	// example.htmlというtemplateをcontentの内容を使って、{{.A}}などのとこ
	// ろを実行して、内容を埋めて、wに書き込む。
	tmpl.ExecuteTemplate(w, "test.html", content)
}
