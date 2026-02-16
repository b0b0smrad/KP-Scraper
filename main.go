package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

func main() {

	keyword := "gtx%201060"

	baseURL := "https://www.kupujemprodajem.com/search.php"
	params := url.Values{}
	params.Add("action", "list")
	params.Add("data[search_type]", "2")
	params.Add("data[keywords]", keyword)
	params.Add("data[page]", "1")

	finalURL := baseURL + "?" + params.Encode()
	fmt.Println("URL:", finalURL)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", finalURL, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0")

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	if err := os.WriteFile("file.txt", []byte(body), 0666); err != nil {
		log.Fatal(err)
	}
	fmt.Println("successfuly dumped")

}

// package main
//
// import (
// 	"fmt"
// 	"github.com/PuerkitoBio/goquery"
// 	tea "github.com/charmbracelet/bubbletea"
// 	// "io"
// 	"log"
// 	"net/http"
// 	"net/url"
// )
//
// //	type model struct {
// //		choices  []string
// //		cursor   int
// //		selected map[int]struct{}
// //	}
// //
// //	func initialModel() model {
// //		return model{
// //			// Our to-do list is a grocery list
// //			choices: []string{"empty for now"},
// //
// //			// A map which indicates which choices are selected. We're using
// //			// the  map like a mathematical set. The keys refer to the indexes
// //			// of the `choices` slice, above.
// //			selected: make(map[int]struct{}),
// //		}
// //	}
// func main() {
//
// 	// var keyword string
// 	// fmt.Print("search name: ")
// 	// fmt.Scanf("%s", &keyword)
//
// 	fixed_keyword := "gtx%201060"
// 	url_base := "https://www.kupujemprodajem.com/pretraga"
// 	params := url.Values{}
// 	params.Add("ignoreUserId", "no")
// 	// params.Add("keywords", keyword)
//
// 	url_con := url_base + "?" + "keywords=" + fixed_keyword + "&ignoreUserId=no"
// 	// url_con := url_base + "?" + "keywords=" + keyword + "&ignoreUserId=no"
// 	println(url_con)
// 	res, err := http.Get(url_con)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer res.Body.Close()
// 	// body, _ := io.ReadAll(res.Body)
// 	// fmt.Println(string(body))
// 	if res.StatusCode != 200 {
// 		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
// 	}
// 	doc, err := goquery.NewDocumentFromReader(res.Body)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	container := doc.Find(".AdItem_adOuterHolder__hb5N_")
// 	fmt.Printf("Container found: %d\n", container.Length())
// 	query := "class=\"AdItem_descriptionHolder__La9qE\""
// 	items := doc.Find(query)
// 	fmt.Printf("Total items found: %d\n", items.Length())
// 	doc.Find(query).Each(func(i int, s *goquery.Selection) {
// 		title := s.Find("a").Text()
// 		fmt.Printf("Review %d: %s\n", i, title)
// 	})
//
// 	fmt.Println("hello world")
// 	tea.ClearScreen()
//
// }
