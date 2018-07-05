package app

import (
	"fmt"
	"encoding/json"
	"html/template"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"

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

	// templateに埋める内容をrequestのFormValueから用意する。
	content := Page{
		A: result,
	}

	// patcartaxi.htmlというtemplateをcontentの内容を使って、{{.A}}などのとこ
	// ろを実行して、内容を埋めて、wに書き込む。
	tmpl.ExecuteTemplate(w, "patcartaxi.html", content)
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
}



// LineはJSONに入ってくる線路の情報をtypeとして定義している。このJSON
// にこの名前にこういうtypeのデータが入ってくるということを表している。
type Line struct {
	Name     string
	Stations []string
}

// Pokemon
// TransitNetworkは http://pokemon.fantasy-transit.appspot.com/net?format=json
// の一番外側のリストのことを表しています。
type TransitNetwork []Line

func handleNorikae(w http.ResponseWriter, r *http.Request) {
	// Appengineの「Context」を通してAppengineのAPIを利用する。
	ctx := appengine.NewContext(r)

	// clientはAppengine用のHTTPクライエントで、他のウェブページを読み込
	// むことができる。
	client := urlfetch.Client(ctx)

	// JSONとしての路線グラフ内容を読み込む
	resp, err := client.Get("http://pokemon.fantasy-transit.appspot.com/net?format=json")
	if err != nil {
		panic(err)
	}

	// 読み込んだJSONをパースするJSONのDecoderを作る。
	decoder := json.NewDecoder(resp.Body)

	// JSONをパースして、「network」に保存する。
	var network TransitNetwork
	if err := decoder.Decode(&network); err != nil {
		panic(err)
	}

	//city1 := r.FormValue("c")
	//city2 := r.FormValue("d")

	// templateに埋める内容をrequestのFormValueから用意する。
	//content2 := Page{
	//	A: city1,
	//	B: city2,
	// }

	fmt.FPrint(w, network[0])

	// handleExampleと同じようにtemplateにテンプレートを埋めて、出力する。
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// tmpl.ExecuteTemplate(w, "norikae.html", network)

}
