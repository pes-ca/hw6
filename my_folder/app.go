package app

import (
	// "fmt"
	"encoding/json"
	"html/template"
	"net/http"

	//"google.golang.org/appengine"
	//"google.golang.org/appengine/urlfetch"

	"unicode/utf8"
)

func init() {
	http.HandleFunc("/", handleExample)
	http.HandleFunc("/norikae", handleNorikae)
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

	A := r.FormValue("a")
	B := r.FormValue("b")
	result := joinWords(w, A, B)

	// fmt.Fprintf(w, result)

	// templateに埋める内容をrequestのFormValueから用意する。
	content := Page{
		A: result,
	}

	// example.htmlというtemplateをcontentの内容を使って、{{.A}}などのとこ
	// ろを実行して、内容を埋めて、wに書き込む。
	tmpl.ExecuteTemplate(w, "test.html", content)
}

func getRuneAt(s string, i int) rune {
    rs := []rune(s)
    return rs[i]
}

func joinWords(w http.ResponseWriter, word1 string, word2 string) string{
		combined_word := ""

		len_1 := utf8.RuneCountInString(word1)
		len_2 := utf8.RuneCountInString(word2)
		max_len := 0

		if len_1 >= len_2 {
			max_len = len_1 * 2
		} else {
			max_len = len_2 * 2
		}

    for i := 0; i < max_len ; i++ {
      if i % 2 == 0 {
				if i/2 < len_1 {
					combined_word += string(getRuneAt(word1, i/2))
				}
			} else {
				if (i-1)/2 < len_2 {
					combined_word += string(getRuneAt(word2, (i-1)/2))
				}
			}
    }
		return combined_word + "\n"
    //for pos, c := range word1 {
			//	combined_word += string([]rune{c})
        // println("位置:", pos, "文字:", string([]rune{c}))


}
