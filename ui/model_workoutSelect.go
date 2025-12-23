package ui

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	wodb "github.com/zmnpl/clift/db"
	coms "github.com/zmnpl/clift/ui/common"
)

type workoutSelect struct {
	workoutList list.Model
	workoutName textinput.Model
	datum       time.Time
}

func NewWorkoutSelectModel(datum time.Time) workoutSelect {
	workoutName := textinput.New()
	workoutName.Placeholder = "a nice name"
	workoutName.Width = 100

	return workoutSelect{
		workoutList: list.New(make([]list.Item, 0), coms.ListItemStyle(), 0, 0),
		workoutName: workoutName,
		datum:       datum,
	}
}

func (m workoutSelect) Init() tea.Cmd {
	return tea.Batch(coms.ReloadWorkouts, textinput.Blink)
}

func (m workoutSelect) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case coms.MsgDate:
		m.datum = time.Time(msg)
		return m, cmd

	case coms.MsgWorkoutsReload:
		if msg.Err != nil {
			//m.errMsg = msgErr(fmt.Sprintf("Error loading workouts: %v", msg.err))
		}
		m.refreshWorkoutList(msg.Workouts)
		return m, tea.Batch(cmd, tea.WindowSize())

	case tea.WindowSizeMsg:
		m.workoutList.SetHeight(coms.GetContentHeight(msg.Height) - 1)

	case coms.MsgWorkoutAddEdit:
		return m, coms.ReloadWorkouts

	case tea.KeyMsg:
		if m.workoutName.Focused() {
			switch msg.String() {
			case "enter":
				m.workoutName.Blur()
				cmd = coms.NewWorkout(m.workoutName.Value())
				m.workoutName.SetValue("")
				return m, cmd

			case "esc":
				m.workoutName.SetValue("")
				m.workoutName.Blur()
				return m, cmd

			default:
				m.workoutName, cmd = m.workoutName.Update(msg)
				return m, cmd
			}
		}

		if m.workoutList.FilterState() == list.Filtering {
			break
		}

		switch msg.String() {
		case "enter":
			return m, coms.GoTo(NewWorkoutModel(m.workoutList.SelectedItem().(coms.WorkoutItem).ID, m.datum))

		case "f2":
			return m, coms.GoTo(NewWorkoutModelEDIT(m.workoutList.SelectedItem().(coms.WorkoutItem).ID, m.datum))

		case "esc":
			return m, coms.Back

		case "n":
			return m, tea.Batch(m.workoutName.Focus(), textinput.Blink)

		case "delete":
			item := m.workoutList.SelectedItem()
			wi, ok := item.(coms.WorkoutItem)
			if !ok {
				return m, cmd
			}
			return m, coms.RemoveWorkout(wi.ID)

		default:

		}
	}

	if m.workoutName.Focused() {
		m.workoutName, cmd = m.workoutName.Update(msg)
		return m, cmd
	}

	m.workoutList, cmd = m.workoutList.Update(msg)

	return m, cmd
}

func (m workoutSelect) View() string {
	sb := &strings.Builder{}

	if m.workoutName.Focused() {
		sb.WriteString(m.workoutName.View() + "\n")
		//m.workoutList.SetHeight(m.workoutList.Height() - 1)
		sb.WriteString(m.workoutList.View())
		return sb.String()
	}

	sb.WriteString(coms.FocusedStyle.Render("Date: ") + m.datum.Format("2006-01-02") + "\n")
	sb.WriteString(m.workoutList.View())

	return sb.String()
}

func (m *workoutSelect) refreshWorkoutList(workouts []wodb.Workout) {
	items := make([]list.Item, len(workouts))
	for i := range workouts {
		items[i] = coms.WorkoutItem{&workouts[i]}
	}
	workoutsList := list.New(items, coms.ListItemStyle(), 0, 0)
	workoutsList.Title = "Select a Workout"
	workoutsList.SetSize(ListWidth, coms.WINDOW_HEIGHT-coms.HEADER_FOOTER_HEIGHT)
	workoutsList.SetShowTitle(false)
	//workoutsList.SetShowStatusBar(false)
	workoutsList.FilterInput.Cursor.Style = coms.FilterCursorStyle
	workoutsList.FilterInput.PromptStyle = coms.FilterPromptStyle

	workoutsList.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			workoutSelectKeys.enter,
			workoutSelectKeys.addWorkout,
			workoutSelectKeys.editWorkout,
			workoutSelectKeys.selectDate,
			workoutSelectKeys.back,
		}
	}

	m.workoutList = workoutsList
}

func (m workoutSelect) BreadCrumb() string {
	return "workouts"
}

// ---------------------------------------------------------------
type workoutSelectKeymap struct {
	enter       key.Binding
	back        key.Binding
	addWorkout  key.Binding
	editWorkout key.Binding
	selectDate  key.Binding
}

var workoutSelectKeys = workoutSelectKeymap{
	enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "do"),
	),
	back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	addWorkout: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "new"),
	),
	editWorkout: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit"),
	),
	selectDate: key.NewBinding(
		key.WithKeys("f5"),
		key.WithHelp("f5", "change date"),
	),
}
