package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	coms "github.com/zmnpl/clift/ui/common"
)

type journal struct {
	journal table.Model
}

func NewReportModel() journal {
	return journal{}
}

func (m journal) Init() tea.Cmd {
	return tea.Batch(coms.LoadPerformedSets)
}

func (m journal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.journal.SetHeight(msg.Height - coms.HEADER_FOOTER_HEIGHT)

	case coms.MsgPerformedSetsLoaded:
		m.journal = coms.MakeJournal(msg.PerformedSets)
		return m, tea.Batch(cmd, tea.WindowSize())

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, coms.Back
		default:

		}
	}

	m.journal, cmd = m.journal.Update(msg)
	return m, cmd
}

func (m journal) View() string {
	sb := &strings.Builder{}
	sb.WriteString(m.journal.View() + "\n")
	return sb.String()
}

func (m journal) BreadCrumb() string {
	return "journal"
}

func (m journal) Help() string {
	return ""
}
