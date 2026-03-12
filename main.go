package main

import (
	// "encoding/json"
	// "charm.land/bubbles/v2/list"
	// "charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-rod/rod"
	"log"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

const listHeight = 10

type queryModel struct {
	Id        []int
	title     []string
	currency  []string
	price     []string
	link      []string
	sortPrice bool
	converted bool
}

type queryResultMsg queryModel

type canvasT struct {
	width int
	flip  bool
}
type model struct {
	spinner    spinner.Model
	canvas     canvasT
	QModel     queryModel
	QSorted    queryModel
	QModelCopy queryModel
	sort       bool
	cursor     int // which to-do list item our cursor is pointing at
	startIndex int
	height     int
	selected   map[int]struct{} // which to-do items are selected
	viewport   viewport.Model
}

func sortPrice(qm *queryModel) queryModel {

	// bm := qm
	// newQm := make([]string, len(qm))
	temp := []float64{}
	// sorted := make([]string, len(qm.price))
	if qm.converted != true {
		for index, value := range qm.price {
			cleanValue := strings.ReplaceAll(value, ".", "")
			val, _ := strconv.ParseFloat(cleanValue, 64)

			var finalEuro float64
			if qm.currency[index] == "din" {
				finalEuro = val / 117.0
				qm.currency[index] = "€" // Update label
			} else {
				finalEuro = val // Already in Euro
			}
			temp = append(temp, finalEuro)
		}
	}

	for n := 0; n < len(temp); n++ {
		for i := 0; i < len(temp)-1; i++ {
			j := i + 1
			if temp[i] > temp[j] {
				temp[i], temp[j] = temp[j], temp[i]
				qm.Id[i], qm.Id[j] = qm.Id[j], qm.Id[i]
				qm.title[i], qm.title[j] = qm.title[j], qm.title[i]
				qm.currency[i], qm.currency[j] = qm.currency[j], qm.currency[i]
				qm.link[i], qm.link[j] = qm.link[j], qm.link[i]

			}
		}

	}
	for k := 0; k < len(qm.price); k++ {
		qm.price[k] = strconv.FormatFloat(temp[k], 'f', 0, 64)
	}

	qm.converted = true
	return *qm
}

func convertCurrencyToEUR(price float64) float64 {
	results := price / 117.0
	return results
}

// func convertCurrencyToRSD()
func initModel() model {
	// code providing spinner (loading screen tui)
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("144"))

	// code providing listings and viewport:

	// l := list.New()
	return model{

		QModel:   queryModel{sortPrice: false, converted: false},
		QSorted:  queryModel{sortPrice: true, converted: false},
		spinner:  sp,
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
	query := runQuery(keyword)

	return tea.Batch(
		m.spinner.Tick,
		query,
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.canvas.width = msg.Width
		m.height = msg.Height - 6
		return m, nil
	case queryResultMsg:
		m.QModel = queryModel(msg)
		// m.QSorted = queryModel(msg)
		// m.QSorted = m.QModel
		m.viewport.Update(msg)
		// m.QSorted.Id = append([]int{}, m.QModel.Id...)
		// m.QSorted.price = append([]string{}, m.QModel.price...)
		// m.QSorted.link = append([]string{}, m.QModel.link...)
		// m.QSorted.title = append([]string{}, m.QModel.title...)
		// m.QSorted.currency = append([]string{}, m.QModel.currency...)

		m.QSorted.Id = append([]int{}, m.QModel.Id...)
		m.QSorted.price = append([]string{}, m.QModel.price...)
		m.QSorted.link = append([]string{}, m.QModel.link...)
		m.QSorted.title = append([]string{}, m.QModel.title...)
		m.QSorted.currency = append([]string{}, m.QModel.currency...)
		m.QSorted = sortPrice(&m.QSorted)
		m.QModelCopy.Id = append([]int{}, m.QModel.Id...)
		m.QModelCopy.price = append([]string{}, m.QModel.price...)
		m.QModelCopy.link = append([]string{}, m.QModel.link...)
		m.QModelCopy.title = append([]string{}, m.QModel.title...)
		m.QModelCopy.currency = append([]string{}, m.QModel.currency...)
		return m, m.spinner.Tick

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--

				if m.cursor < m.startIndex {
					m.startIndex = m.cursor
				}
			}
		case "down", "j":
			if m.cursor < len(m.QModel.title)-1 {
				m.cursor++

				// m.cursor = 0
				// m.startIndex = 0
			}
			if m.cursor >= m.startIndex+m.height {
				m.startIndex = m.cursor - m.height + 1
			}
		case "s":
			m.sort = !m.sort
		case "l":
			m.QModel.sortPrice = !m.QModel.sortPrice
			if m.QModel.sortPrice == true {
				m.QModel = m.QSorted

			} else {
				m.QModel = m.QModelCopy
			}

		case "enter", " ":
			if err := openURL(m.QModel.link[m.cursor]); err != nil {
				fmt.Printf("error openining url %v\n", err)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		default:
			m.canvas.flip = !m.canvas.flip

		}
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m model) View() tea.View {
	strAnim := fmt.Sprintf("\n %s  Searching KupujemProdajem... Please wait.\n", m.spinner.View())
	if len(m.QModel.Id) == 0 {
		return tea.NewView(strAnim)
	}

	var b strings.Builder
	b.WriteString("What are you buying:\n\n")

	viewHeight := m.height
	if viewHeight <= 0 {
		viewHeight = 10
	}
	m.viewport.View()
	maxVisible := m.startIndex + m.height
	if maxVisible > len(m.QModel.title) || maxVisible < len(m.QModel.title) {
		maxVisible = len(m.QModel.title)
	}
	end := min(m.startIndex+maxVisible+10, len(m.QModel.title))
	for i := m.startIndex; i < end; i++ {
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
		gray := lipgloss.NewStyle().Foreground(lipgloss.Color("144"))
		// Build the line piece by piece

		if !m.sort {
			// 1. Format the string
			str1 := fmt.Sprintf("%s [%s] %s - %s%s", cursor, pink.Render(checked), gray.Render(m.QModel.title[i]), gray.Render(m.QModel.currency[i]), m.QModel.price[i])
			// 2. Render it and write it to the builder
			b.WriteString(str1 + "\n")
		} else {
			str2 := fmt.Sprintf("%s [%s] %s%s - %s", cursor, pink.Render(checked), gray.Render(m.QModel.currency[i]), m.QModel.price[i], gray.Render(m.QModel.title[i]))
			b.WriteString(str2 + "\n")
		}

	}

	b.WriteString("\n PRESS: [(q) quit] [(s) sort-order] [(l) sort-price] \n")
	return tea.NewView(b.String())
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
		if title != "" && price != "Kontakt" {
			Qbase.Id = append(Qbase.Id, i)
			Qbase.title = append(Qbase.title, title)
			parts := strings.Split(strings.TrimSpace(price), " ")
			if len(parts) == 2 {
				Qbase.price = append(Qbase.price, parts[0])
				Qbase.currency = append(Qbase.currency, parts[1])

			}
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

	// results := sendQuery(keyword)
	//
	// results.sortPrice = true
	// printable := sortPrice(&results)
	// for _, value := range printable.price {
	// 	fmt.Printf("\nprice: %s\n", value)
	//
	// }

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
