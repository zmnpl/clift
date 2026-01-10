package ui

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	datepicker "github.com/ethanefung/bubble-datepicker"
	com "github.com/zmnpl/clift/ui/common"
)

type selectDate struct {
	datum      time.Time
	datepicker datepicker.Model
}

func NewDateSelectModel(datum time.Time) selectDate {

	if datum.IsZero() {
		datum = time.Now()
	}

	datepicker := datepicker.New(datum)
	datepicker.Styles.FocusedText = com.FocusedStyle
	datepicker.SelectDate()

	return selectDate{
		datum:      datum,
		datepicker: datepicker,
	}
}

func (m selectDate) Init() tea.Cmd {
	return nil
}

func (m selectDate) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return m, com.Ret(com.SendDate(m.datepicker.Time))
		case "esc":
			return m, com.Back
		}
		m.datepicker, cmd = m.datepicker.Update(msg)
		return m, cmd
	}

	return m, cmd
}

func (m selectDate) View() string {
	sb := &strings.Builder{}
	sb.WriteString(m.datepicker.View() + "\n")
	return sb.String()
}

func (m selectDate) BreadCrumb() string {
	return "select date"
}

func (m selectDate) Help() string {
	return ""
}
