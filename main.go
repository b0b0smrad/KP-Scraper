package main

import (
	// "encoding/json"
	"fmt"
	"log"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-rod/rod"

	// "io"
	// "log"
	// "net/http"
	"net/url"
	"strings"
)

func main() {
	keyword := "TCL 4k"

	// Use the REGULAR search page, not API
	searchURL := "https://www.kupujemprodajem.com/pretraga?keywords=" + url.QueryEscape(keyword)

	browser := rod.New().MustConnect()
	defer browser.MustClose()

	page := browser.MustPage(searchURL)
	page.MustWaitLoad()
	page.MustWaitIdle()

	html := page.MustHTML()
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))

	// Parse HTML with goquery
	doc.Find(".AdItem_adHolder__rKT82").Each(func(i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Find(".AdItem_name__iOZvA").Text())
		price := strings.TrimSpace(s.Find(".AdItem_price__VZ_at").Text())
		link, _ := s.Find("a.Link_link__cqSOS.Link_inherit__05Kzh").Attr("href")

		fmt.Printf("\n=== Item %d ===\n", i+1)
		if title != "" {
			fmt.Printf("Title: %s\n", title) // Fixed: added title argument
			fmt.Printf("Price: %s\n", price)
			fmt.Printf("Link: https://www.kupujemprodajem.com%s\n", link)
		}
		// if checkIFEmpty(price) {
		// }
		// if checkIFEmpty(link) {
		// }

	})
}
func checkIFEmpty(value string) bool {

	if value != "" {
		log.Fatal("string is empty")
		return true
	}
	return false

}

// Add the cookies from your browser
// In Firefox: F12 -> Storage -> Cookies -> copy the Cookie header value
// Replace this entire string with YOUR cookie value:
// {
// 	req.Header.Set("Cookie", `g_state={"i_l":1,"i_ll":1765497211069,"i_b":"d0STXBIlIwhf95/jWDNzj0r9e3wtEE7MXFHJjF1tnwU","i_e":{"enable_itp_optimization":0},"i_p":1765504453558}; machine_id=54fd5d21e97e2ad49e0cb0f629795381; cookie[emailSSL]=markoVIII1999%40gmail.com; cookie[password_hashSSL]=9db5013f23a6ef5ec36664b3ddc003d0; recentSearchFilterIds=[8262832945%2C8435707930%2C8493105390%2C8493105795%2C8493106119%2C8493107055%2C8493145581%2C8493146133%2C8493346749%2C8493387138%2C8492499840]; screenWidth=2288; cookie_consent_v2=1; KP-DEVICE-J…d7FDzx7sciddS4QGFpgwmaqk5aLrqtFcwz7G4BKD8; KP-TEMP-AUTH-JWT=eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjo3OTE4NTksImlzcyI6Imh0dHBzOi8vd3d3Lmt1cHVqZW1wcm9kYWplbS5jb20vIiwiaWF0IjoxNzcxMjU5OTU5LCJqdGkiOiI0OTZkNjEzMDJiYTk0YWZjOWQxMDVlZjI1YWE1NGUyZiJ9.1El07YECH5lx19Q75RB9bP3uO56JIJvP8DYvA0lTT68; KUPUJEMPRODAJEM=oggqsb12hs1hr2rsjddn7d4uvm; cookie[user_idSSL]=a60bd86e802e55a9c4e628c96e168ca4; imageExpandInfo={%22adId%22:%22114055238%22%2C%22count%22:3}; zoomInfo={%22adId%22:%22187150520%22%2C%22count%22:2}`)
// 	req.Header.Set("Authorization", "f416a7586a8ef31a2d6c63a94ae901202e6e0c84fa9b8175c3cd5d2904c94245")
// 	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:147.0) Gecko/20100101 Firefox/147.0")
// 	req.Header.Set("x-kp-channel", "desktop_react")
// 	req.Header.Set("x-kp-machine-id", "54fd5d21e97e2ad49e0cb0f629795381")
// 	req.Header.Set("x-kp-session", "oggqsb12hs1hr2rsjddn7d4uvm")
// 	req.Header.Set("x-kp-signature", "1baf8a098c90f29a3f9553428ef647ae2592474e")
// 	req.Header.Set("Accept", "application/json, text/plain, */*")
// 	req.Header.Set("Referer", "https://www.kupujemprodajem.com/pretraga?keywords=gtx%201060")
// 	req.Header.Set("User-Agent", "Mozilla/5.0")
// 	req.Header.Set("x-kp-channel", "desktop_react")
//
// }
//
// res, err := client.Do(req)
// if err != nil {
// 	log.Fatal(err)
// }
// defer res.Body.Close()
//
// body, _ := io.ReadAll(res.Body)
//
// if res.StatusCode != 200 {
// 	fmt.Println("Status:", res.Status)
// 	fmt.Println("Response:", string(body))
// 	return
// }
//
//
// fmt.Printf("Found %d results\n\n", searchRes.Results.Total)
//
// for i, ad := range searchRes.Results.Ads {
// 	fmt.Printf("=== %d ===\n", i+1)
// 	fmt.Printf("%s - %d %s\n", ad.Name, ad.Price, ad.Currency)
// 	fmt.Printf("https://www.kupujemprodajem.com%s\n\n", ad.AdURL)
// }
//

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
