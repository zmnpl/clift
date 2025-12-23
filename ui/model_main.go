package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	coms "github.com/zmnpl/clift/ui/common"
)

type Screen interface {
	BreadCrumb() string
}

const (
	ListWidth  = 85
	ListHeight = 30
)

type model struct {
	screenStack   []tea.Model
	currentScreen tea.Model

	datum time.Time

	// ui stuff
	help help.Model
	ws   tea.WindowSizeMsg

	// signals / error
	statusMsg coms.StatusMsg
}

func NewModel() model {
	screnStack := make([]tea.Model, 0, 100)

	return model{
		screenStack: screnStack,
		datum:       time.Now(),
		help:        help.New(),
		statusMsg:   coms.StatusMsg{Status: "Do you even lift, bro?"},
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c": // universal quit across screens
			return m, tea.Quit
		case "f5": // change date works from everywhere as well
			return m, coms.GoTo(NewDateSelectModel(m.datum))
		}

	case coms.MsgDate:
		m.datum = time.Time(msg)
		for i, _ := range m.screenStack {
			m.screenStack[i], _ = m.screenStack[i].Update(msg)
		}

	case coms.MsgGoTo:
		if m.currentScreen != nil {
			m.screenStack = append(m.screenStack, m.currentScreen)
		}
		m.currentScreen = msg
		m.currentScreen.Update(m.ws) // send window size first time
		return m, m.currentScreen.Init()

	case coms.MsgBack:
		m.currentScreen = m.popScreen()
		return m, cmd

	case coms.MsgBackCmd:
		m.currentScreen = m.popScreen()
		return m, msg.Cmd

	case tea.WindowSizeMsg:
		coms.WINDOW_HEIGHT = msg.Height
		coms.WINDOW_WIDTH = msg.Width
		m.ws = msg
		m.help.Width = msg.Width
		for i, _ := range m.screenStack {
			m.screenStack[i], _ = m.screenStack[i].Update(msg)
		}

	case coms.StatusMsg:
		m.statusMsg = msg
		//case coms.MsgErr, coms.MsgCool:
		//	m.statusMsg = msg // TODO: differentiate between good and bad
	}

	if m.currentScreen == nil {
		return m.upateScreenMain(msg)
	} else {
		m.currentScreen, cmd = m.currentScreen.Update(msg)
		return m, cmd
	}
}

func (m model) upateScreenMain(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "1":
			return m, coms.GoTo(NewWorkoutSelectModel(m.datum))

		case "2":
			return m, coms.GoTo(NewDoExercise(0, time.Now()))

		case "3":
			return m, coms.GoTo(NewReportModel())

		case "esc":
			m.statusMsg = coms.StatusMsg{}
		}
	}
	return m, cmd

}

var (
	focusedButton = coms.FocusedStyle.Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", coms.BlurredStyle.Render("Submit"))
)

func (m model) View() string {
	sb := &strings.Builder{}

	// content height
	ch := coms.GetContentHeight(m.ws.Height)

	content := m.viewScreenMain()
	if m.currentScreen != nil {
		content = m.currentScreen.View()
	}

	breadcrumb := m.makeBreadCrumb() + "\n"
	status := "\n" + m.makeStatus() //+ fmt.Sprintf(" - %v", ch)

	//status = status  + fmt.Sprintf("bc: %v - st: %v", lipgloss.Height(breadcrumb), lipgloss.Height(status))

	sb.WriteString(breadcrumb)

	// content = "\n"
	// for i := 1; i <= ch; i++ {
	// 	content = content + fmt.Sprintf("%v\n", i)
	// }

	sb.WriteString("\n" +
		lipgloss.PlaceVertical(ch,
			lipgloss.Top,
			content) + "\n")

	sb.WriteString(status)

	return coms.Margin.Render(sb.String())
}

func (m model) makeBreadCrumb() string {
	bc := &strings.Builder{}

	if len(m.screenStack) == 0 && m.currentScreen == nil {
		bc.WriteString("/")
	}

	// breadcrumb for screen stack
	for _, s := range m.screenStack {
		if si, ok := s.(Screen); ok {
			bc.WriteString("/" + si.BreadCrumb())
		} else {
			bc.WriteString("/" + "stack not screen")
		}
	}

	// add current screen
	if m.currentScreen != nil {
		// current screen breadcrumb
		if si, ok := m.currentScreen.(Screen); ok {
			bc.WriteString("/" + si.BreadCrumb())
		} else {
			bc.WriteString("/" + "current not screen")
		}
	}

	return coms.HeaderStyle.Render(bc.String())
}

func (m model) makeStatus() string {
	status := &strings.Builder{}

	statusRenderer := coms.StatusGood
	statusText := m.statusMsg.Status

	if m.statusMsg.Err != nil {
		statusRenderer = coms.StatusBad
		statusText = m.statusMsg.Err.Error()
	}

	status.WriteString(statusRenderer.Render("Coach") + coms.StatusCenter.Render(fmt.Sprintf("\"%v\"", statusText)))

	return status.String()
}

func (m model) viewScreenMain() string {
	sb := &strings.Builder{}
	sb.WriteString(coms.FocusedStyle.Render("1) ") + "workouts" + "\n")
	sb.WriteString(coms.FocusedStyle.Render("2) ") + "exercises" + "\n")
	sb.WriteString(coms.FocusedStyle.Render("3) ") + "journal" + "\n")
	return sb.String()
}

func (m *model) popScreen() tea.Model {
	var s tea.Model
	if len(m.screenStack) > 0 {
		s = m.screenStack[len(m.screenStack)-1]
		m.screenStack = m.screenStack[:len(m.screenStack)-1]
	}
	return s
}
