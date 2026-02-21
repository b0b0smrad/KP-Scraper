package main

import (
	// "encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/PuerkitoBio/goquery"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-rod/rod"

	// "io"
	// "log"
	// "net/http"
	"net/url"
	"strings"
)

type queryModel struct {
	Id    []int
	title []string
	price []string
	link  []string
}

type model struct {
	QModel   queryModel
	cursor   int              // which to-do list item our cursor is pointing at
	selected map[int]struct{} // which to-do items are selected
}

func initModel() model {
	return model{

		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	// for _, i := range q.Id {
	// 	m.IdArray[i] = q.Id[i]
	// }
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.QModel.Id)-1 {
				m.cursor++
			}
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}

		}

	}
	return m, nil
}

func (m model) View() string {

	// note the _l means it's a locallized post-fix

	string_l := "what are you buying:\n\n"

	tea.ClearScreen()
	for i, choice := range m.QModel.Id {
		cursor_l := " "
		if m.cursor == i {
			cursor_l = ">"
		}
		checked_l := " " // not selected
		if _, ok := m.selected[i]; ok {
			checked_l = "o" // linked or selected
		}
		string_l += fmt.Sprintf("%s [%s] %s\n", cursor_l, checked_l, choice)
	}

	string_l += "\nPress q or C^c to quit\n"

	return string_l

}
func sendQuery(keyword string) queryModel {
	searchURL := "https://www.kupujemprodajem.com/pretraga?keywords=" + url.QueryEscape(keyword)

	var Qbase queryModel
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
			Qbase.Id = append(Qbase.Id, i)
			Qbase.title = append(Qbase.title, title)
			Qbase.price = append(Qbase.price, price)
			Qbase.link = append(Qbase.link, link)
			fmt.Printf("Title: %s\n", title) // Fixed: added title argument
			fmt.Printf("Price: %s\n", price)
			fmt.Printf("Link: https://www.kupujemprodajem.com%s\n", link)
		}
	})

	return Qbase
}
func main() {

	keyword := "TCL 4k"

	data := queryModel{}
	data = sendQuery(keyword)
	for _, i := range data.Id {
		fmt.Printf("title: %s; price: %s\n\n", data.title[i], data.price[i])
	}
	// Use the REGULAR search page, not API
	p := tea.NewProgram(initModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Errrrorrrr: %v", err)
		os.Exit(1)
	}
	fmt.Printf(" ")
}
func checkIFEmpty(value string) bool {

	if value != "" {
		log.Fatal("string is empty")
		return true
	}
	return false

}
