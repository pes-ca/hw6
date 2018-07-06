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

	cityToNum, numToCity, numToLoop, cityToLoops, loopToCities, cityToCities := networkInterpreter(w, network)
	fmt.Fprint(w, "cityToNum:", cityToNum, "\n")
	fmt.Fprint(w, "numToCity:", numToCity, "\n")
	fmt.Fprint(w, "cityToCities:", cityToCities, "\n")
	fmt.Fprint(w, "numToLoop:", numToLoop, "\n")
	fmt.Fprint(w, "cityToLoops:", cityToLoops, "\n")
	fmt.Fprint(w, "loopToCities:", loopToCities, "\n")

	//city1 := r.FormValue("c")
	//city2 := r.FormValue("d")

	// templateに埋める内容をrequestのFormValueから用意する。
	//content2 := Page{
	//	A: city1,
	//	B: city2,
	// }

	// fmt.Fprint(w, network)
	// -> [{Outer Loop [Pallet Town Viridian City Pewter City Cerulean City Lavender Town Fuschia City Celadon City Saffron City Lavender Town]} {Inner Loop [Saffron City Vermillion City Fuschia City Celadon City Saffron City]} {Victory Road [Viridian City Mt. Silver Indigo Plateau]} {Seafoam Island Ferry [Fuschia City Cinnibar island Pallet Town]} {Route 11 [Vermillion City Lavender Town]} {Diglett Network [Vermillion Cave Viridian Cave Rock Tunnel Cave]}]

	// fmt.Fprint(w, network[0])
	// -> {Outer Loop [Pallet Town Viridian City Pewter City Cerulean City Lavender Town Fuschia City Celadon City Saffron City Lavender Town]}

	// fmt.Fprint(w, network[0].Stations)
	// [Pallet Town Viridian City Pewter City Cerulean City Lavender Town Fuschia City Celadon City Saffron City Lavender Town]

	// fmt.Fprint(w, network[0].Stations[0])
	// Pallet Town


	// handleExampleと同じようにtemplateにテンプレートを埋めて、出力する。
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// tmpl.ExecuteTemplate(w, "norikae.html", network)


}




//  []string, []string, []Dic, []Dic
func networkInterpreter(w http.ResponseWriter, network TransitNetwork) (map[string]int, []string, []string, [][]int, [][]int, [][]int){
	// make a list of
	// city1 [ list of loops that contain city1 ]
	// city2 [ list of loops that contain city2 ]
	// ...

	cityToNum := make(map[string]int)
	// map[]
	// fmt.Fprint(w, cityToNum)

	numToCity := make([]string, 0)
	numToLoop := make([]string, 0)

	// adjacency list
	cityToCities := make([][]int, 15)
	fmt.Fprint(w, "\ncityToCities:", cityToCities, "\n")

	last_city := -1
	flag := 0
	flag1, flag2 := 0, 0



	cityToLoops := make([][]int, 15)
	fmt.Fprint(w, "cityToLoops:", cityToLoops, "\n")
	// fmt.Fprint(w, "cityToLoops[0]:", cityToLoops[0], "\n") -> cityToLoops[0]:[]
	loopToCities := make([][]int, 6)


	//l := new(Dic)
	// &{0 []}
	//fmt.Fprint(w, l)


	// i < len(network)
	for i := 0; i < 2 ; i++ {
		numToLoop = append(numToLoop, network[i].Name)
		fmt.Fprint(w, "numToLoop:", numToLoop, "\n")

		// j < len(network[i].Stations
		for j := 0; j < len(network[i].Stations) ; j++ {
			// city: network[i].Stations[j]

			// if the city has been seen before
			if cityID, ok := cityToNum[network[i].Stations[j]]; ok {
				flag = 0
				for n := 0; n < len(cityToLoops[cityID]) ; n ++ {
					if cityToLoops[cityID][n] == len(numToLoop)-1 {
						flag = 1
					}
				}
				if flag == 0{
					cityToLoops[cityID] = append(cityToLoops[cityID], len(numToLoop)-1)
				}

				fmt.Fprint(w, "cityToLoops:", cityToLoops, "\n")
				loopToCities[i] = append(loopToCities[i], cityID)
				fmt.Fprint(w, "loopToCities:", loopToCities, "\n")

				if j > 0 {

					flag1, flag2 = 0, 0
					for n := 0; n < len(cityToCities[last_city]) ; n ++ {
						if cityToCities[last_city][n] == cityID {
							flag1 = 1
						}
					}
					for n := 0; n < len(cityToCities[cityID]) ; n ++ {
						if cityToCities[cityID][n] == last_city {
							flag2 = 1
						}
					}

					if flag1 == 0 {
						cityToCities[last_city] = append(cityToCities[last_city], cityID)
					}
					if flag2 == 0 {
						cityToCities[cityID] = append(cityToCities[cityID], last_city)
					}

					last_city = cityID
				} else {
					last_city = cityID
				}
				fmt.Fprint(w, "cityToCities:", cityToCities, "\n")

    	} else {
				numToCity = append(numToCity, network[i].Stations[j])
				fmt.Fprint(w, "numToCity:", numToCity, "\n")
				cityToNum[network[i].Stations[j]] = len(numToCity)-1
				fmt.Fprint(w, "cityToNum:", cityToNum, "\n")

				if j > 0 {

					flag1, flag2 = 0, 0
					for n := 0; n < len(cityToCities[last_city]) ; n ++ {
						if cityToCities[last_city][n] == len(numToCity)-1 {
							flag1 = 1
						}
					}
					for n := 0; n < len(cityToCities[len(numToCity)-1]) ; n ++ {
						if cityToCities[len(numToCity)-1][n] == last_city {
							flag2 = 1
						}
					}

					if flag1 == 0 {
						cityToCities[last_city] = append(cityToCities[last_city], len(numToCity)-1)
					}
					if flag2 == 0 {
						cityToCities[len(numToCity)-1] = append(cityToCities[len(numToCity)-1], last_city)
					}

					last_city = len(numToCity)-1
				} else {
					last_city = len(numToCity)-1
				}
				fmt.Fprint(w, "cityToCities:", cityToCities, "\n")

				cityToLoops[len(numToCity)-1] = append(cityToLoops[len(numToCity)-1], len(numToLoop)-1)
				fmt.Fprint(w, "cityToLoops:", cityToLoops, "\n")
				loopToCities[i] = append(loopToCities[i], len(numToCity)-1)
				fmt.Fprint(w, "loopToCities:", loopToCities, "\n")
    	}
		}
	}


	return cityToNum, numToCity, numToLoop, cityToLoops, loopToCities, cityToCities
}
