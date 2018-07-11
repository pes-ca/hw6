package app

import (
	//"fmt"
	"encoding/json"
	"html/template"
	"net/http"

	"appengine"
	//"google.golang.org/appengine"
	"appengine/urlfetch"
	//"google.golang.org/appengine/urlfetch"

	"unicode/utf8"
	"container/list"
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
	C string
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

func getSelected(r *http.Request)(if_valid bool, time_transits string) {
	slice:=[]string{"1","2"}
	for _, v := range slice {
	    if v == r.Form.Get("time_or_transits") {
	        return true, v
	    }
	}
	return false, "0"
}

func handleNorikae(w http.ResponseWriter, r *http.Request) {
	// handleExampleと同じようにtemplateにテンプレートを埋めて、出力する。
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

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
		// fmt.Fprint(w, network)
		//  -> [{Outer Loop [Pallet Town Viridian City Pewter City Cerulean City Lavender Town Fuschia City Celadon City Saffron City Lavender Town]} {Inner Loop [Saffron City Vermillion City Fuschia City Celadon City Saffron City]} {Victory Road [Viridian City Mt. Silver Indigo Plateau]} {Seafoam Island Ferry [Fuschia City Cinnibar island Pallet Town]} {Route 11 [Vermillion City Lavender Town]} {Diglett Network [Vermillion Cave Viridian Cave Rock Tunnel Cave]}]
		// fmt.Fprint(w, network[0])
		//  -> {Outer Loop [Pallet Town Viridian City Pewter City Cerulean City Lavender Town Fuschia City Celadon City Saffron City Lavender Town]}
		// fmt.Fprint(w, network[0].Stations)
		//  -> [Pallet Town Viridian City Pewter City Cerulean City Lavender Town Fuschia City Celadon City Saffron City Lavender Town]
		// fmt.Fprint(w, network[0].Stations[0])
		//  -> Pallet Town

	cityToNum, numToCity, numToLoop, cityToCities := networkInterpreter(w, network)

	city1 := r.FormValue("c")
	city2 := r.FormValue("d")

	routes_min_depth, routes_min_transits := BFS(w, cityToCities, cityToNum[city1], cityToNum[city2], numToCity)

	if_valid, time_transits := getSelected(r)
	result := ""
	if if_valid{
	  if time_transits == "1"{
			if len(routes_min_transits) == 0{
				result += "not found"
			}
			for num, route := range routes_min_transits {
				for i, transit := range route.Transits {
					result += numToLoop[route.Lines[i]]
					result += " --> Change at "
					result += numToCity[transit]
					result += " to --> "
				}
				result +=  numToLoop[route.Lines[len(route.Lines)-1]]
				if num != len(routes_min_transits) -1{
					result += "\nNEXT ROUTE\n"
				}
			}
			result += "\n"

	  } else if time_transits == "2"{
			if len(routes_min_depth) == 0{
				result += "not found"
			}
			for num, route := range routes_min_depth {
				for i, transit := range route.Transits {
					result += numToLoop[route.Lines[i]]
					result += " --> Change at "
					result += numToCity[transit]
					result += " to --> "
				}
				result +=  numToLoop[route.Lines[len(route.Lines)-1]]
				if num != len(routes_min_depth) -1{
					result += "\nNEXT ROUTE\n"
				}
			}
			result += "\n"
	  }
	}

	// templateに埋める内容をrequestのFormValueから用意する。
	content2 := Page{
		A: city1,
		B: city2,
		C: result,
	}

	tmpl.ExecuteTemplate(w, "norikae.html", content2)

}

// delete overlap in [][]int that needs to be treated like set
func deleteOverlap2(array [][]int) [][]int {
	uniqArray := [][]int{}
	for _, arr := range array {
		uniq := []int{}
		m := make(map[int]bool)
		for _, ele := range arr {
				if !m[ele] {
						m[ele] = true
						uniq = append(uniq, ele)
				}
		}
		uniqArray = append(uniqArray, uniq)
	}
	return uniqArray
}


type Info struct {
	Transits     []int
	Lines []int
	City int
	Last int
	Depth int
}

type Route struct {
	Transits     []int
	Lines []int
}


func BFS(w http.ResponseWriter, l [][][]int, start int, goal int, numToCity []string) ([]Route, []Route){
	q := list.New()

	var first Info
	first.City = start
	first.Depth = 0

	q.PushBack(first)

	var routes_min_depth []Route
	var	routes_min_transits []Route
	min_depth := 1000000
	min_transits := 1000000

	count:=0

	for true {

	  count = count + 1
		if count == 50{
			break
		}
	  if q.Len() == 0 {
	    break
	  }

	  node_info := q.Front().Value.(Info)

	  q.Remove(q.Front())

		var new_info Info

		for _, city_info := range l[node_info.City] {
		  city := city_info[0]
		  line := city_info[1]
		  if city == node_info.Last {
		    continue
		  } else{
				if len(node_info.Lines) == 0{
					new_info.Transits = node_info.Transits
					new_info.Lines = append(node_info.Lines, line)
				} else {
					if line != node_info.Lines[len(node_info.Lines) - 1]{
			      new_info.Transits = append(node_info.Transits, node_info.City)
			      new_info.Lines = append(node_info.Lines, line)
			    } else {
						new_info.Transits = node_info.Transits
						new_info.Lines = node_info.Lines
					}
				}

				new_info.Depth = node_info.Depth + 1
				new_info.Last = node_info.City
				new_info.City = city

				if city == goal{
					//fmt.Fprint(w, "city, goal:", city, goal, "\n")
					if new_info.Depth == min_depth{
						var new_route Route
						new_route.Transits = new_info.Transits
						new_route.Lines = new_info.Lines
						routes_min_depth = append(routes_min_depth, new_route)
						min_depth = new_info.Depth
					} else if len(new_info.Transits) < min_transits {
						var new_route Route
						new_route.Transits = new_info.Transits
						new_route.Lines = new_info.Lines
						routes_min_depth = make([]Route, 0)
						routes_min_depth = append(routes_min_depth, new_route)
						min_depth = new_info.Depth
					}
					if len(new_info.Transits) == min_transits{
						//fmt.Fprint(w, "new_info.Depth:", new_info.Depth, "\n")
						var new_route Route
						new_route.Transits = new_info.Transits
						new_route.Lines = new_info.Lines
						routes_min_transits = append(routes_min_transits, new_route)
						min_transits = len(new_info.Transits)

					//fmt.Fprint(w, "routes_min_depth:", routes_min_depth, "\n")
					//fmt.Fprint(w, "routes_min_transits:", routes_min_transits, "\n")
					} else if len(new_info.Transits) < min_transits {
						var new_route Route
						new_route.Transits = new_info.Transits
						new_route.Lines = new_info.Lines
						routes_min_transits = make([]Route, 0)
						routes_min_transits = append(routes_min_transits, new_route)
						min_transits = len(new_info.Transits)
					}
				} else {
					q.PushBack(new_info)
				}
		  }
		}
	}
	return routes_min_depth, routes_min_transits
}



func networkInterpreter(w http.ResponseWriter, network TransitNetwork) (map[string]int, []string, []string, [][][]int){
	// initialisation
	cityToNum := make(map[string]int)
	numToCity := make([]string, 0)
	numToLoop := make([]string, 0)

	//cityToLoops, loopToCities: made but not used outside this function at the moment
	cityToLoops := make([][]int, 15)
	loopToCities := make([][]int, 6)

	cityToCities := make([][][]int, 15)	//l := new(Dic)

	last_city := -1

	// i < len(network)
	for i := 0; i < len(network) ; i++ {
		numToLoop = append(numToLoop, network[i].Name)

		for j := 0; j < len(network[i].Stations) ; j++ {
			// if the city has been seen before
			if cityID, ok := cityToNum[network[i].Stations[j]]; ok {
				cityToLoops[cityID] = append(cityToLoops[cityID], len(numToLoop)-1)
				loopToCities[i] = append(loopToCities[i], cityID)

				if j > 0 {
					cityToCities[last_city] = append(cityToCities[last_city], []int{cityID, len(numToLoop)-1})
					cityToCities[cityID] = append(cityToCities[cityID], []int{last_city, len(numToLoop)-1})
					last_city = cityID
				} else {
					last_city = cityID
				}

    	} else {
				numToCity = append(numToCity, network[i].Stations[j])
				cityToNum[network[i].Stations[j]] = len(numToCity)-1

				if j > 0 {
					cityToCities[last_city] = append(cityToCities[last_city], []int{len(numToCity)-1, len(numToLoop)-1})
					cityToCities[len(numToCity)-1] = append(cityToCities[len(numToCity)-1], []int{last_city, len(numToLoop)-1})
					last_city = len(numToCity)-1
				} else {
					last_city = len(numToCity)-1
				}

				cityToLoops[len(numToCity)-1] = append(cityToLoops[len(numToCity)-1], len(numToLoop)-1)
				loopToCities[i] = append(loopToCities[i], len(numToCity)-1)
    	}
		}
	}
	//cityToLoops, loopToCities: made but not used outside this function at the moment
	cityToLoops = deleteOverlap2(cityToLoops)
	loopToCities = deleteOverlap2(loopToCities)
	return cityToNum, numToCity, numToLoop, cityToCities
}
