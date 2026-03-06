package main

import (
	// "encoding/json"
	"charm.land/lipgloss/v2"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-rod/rod"
	"log"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type queryModel struct {
	Id    []int
	title []string
	price []string
	link  []string
}
type queryResultMsg queryModel

type model struct {
	QModel   queryModel
	sort     bool
	cursor   int              // which to-do list item our cursor is pointing at
	selected map[int]struct{} // which to-do items are selected
}

func initModel() model {
	return model{

		selected: make(map[int]struct{}),
		sort:     false,
	}
}
func openURL(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start", url}
	case "darwin": // macOS
		cmd = "open"
		args = []string{url}
	default: // "linux", "freebsd", "openbsd", "netbsbsd"
		cmd = "xdg-open"
		args = []string{"https://kupujemprodajem.com/" + url}
	}
	return exec.Command(cmd, args...).Start()
}
func runQuery(keyword string) tea.Cmd {
	return func() tea.Msg {
		results := sendQuery(keyword)
		return queryResultMsg(results) // Send this back to Update
	}
}
func (m model) Init() tea.Cmd {
	keyword := "TCL 4k"
	return runQuery(keyword)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	// for _, i := range q.Id {
	// 	m.IdArray[i] = q.Id[i]
	// }
	switch msg := msg.(type) {
	case queryResultMsg:
		m.QModel = queryModel(msg)
		return m, nil
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
		case "s":
			m.sort = !m.sort
		case "enter", " ":
			if err := openURL(m.QModel.link[m.cursor]); err != nil {
				fmt.Printf("error openining url %v\n", err)
			} else {
				m.selected[m.cursor] = struct{}{}
			}

		}

	}
	return m, nil
}

func (m model) View() string {
	if len(m.QModel.Id) == 0 {
		return "\n  Searching KupujemProdajem... Please wait.\n"
	}

	var b strings.Builder
	b.WriteString("What are you buying:\n\n")

	for i, title := range m.QModel.title {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}
		// Using title instead of ID for better visibility
		pink := lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
		gray := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Background(lipgloss.BrightGreen)

		// Build the line piece by piece
		if !m.sort {
			// 1. Format the string
			str := fmt.Sprintf("%s [%s] %s - %s", cursor, pink.Render(checked), gray.Render(title), m.QModel.price[i])
			// 2. Render it and write it to the builder
			b.WriteString(str + "\n")
		} else {
			str := fmt.Sprintf("%s [%s] %s - %s", cursor, pink.Render(checked), m.QModel.price[i], gray.Render(title))
			b.WriteString(str + "\n")
		}
	}

	b.WriteString("\nPress q to quit\n")
	return b.String()
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
		// fmt.Printf("\n=== Item %d ===\n", i+1)
		if title != "" {
			Qbase.Id = append(Qbase.Id, i)
			Qbase.title = append(Qbase.title, title)
			Qbase.price = append(Qbase.price, price)
			Qbase.link = append(Qbase.link, link)
			// fmt.Printf("Title: %s\n", title) // Fixed: added title argument
			// fmt.Printf("Price: %s\n", price)
			// fmt.Printf("Link: https://www.kupujemprodajem.com%s\n", link)
		}
	})

	return Qbase
}
func main() {

	// keyword := "TCL 4k"
	//
	// data := queryModel{}
	// data = sendQuery(keyword)

	// fmt.Printf("title: %s; price: %s\n\n", data.title[1], data.price[1])
	// for _, i := range data.Id {
	// 	fmt.Printf("title: %s; price: %s\n\n", data.title[i], data.price[i])
	// }
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
