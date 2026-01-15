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

const (
	MODE_EXERCISE_DO = iota
	MODE_EXERCISE_RETURNID
)

type exerciseSelect struct {
	exerciseList list.Model
	mode         int

	workoutID uint
	datum     time.Time

	ws tea.WindowSizeMsg

	status string
}

func NewDoExercise(workoutID uint, datum time.Time) exerciseSelect {

	if datum.IsZero() {
		datum = time.Now()
	}

	m := exerciseSelect{
		mode:         MODE_EXERCISE_DO,
		exerciseList: list.New(make([]list.Item, 0), coms.ListItemStyle(), 0, 0),
		workoutID:    workoutID,
		//ws:           ws,
		datum: datum,
	}
	return m
}

func NewSelectExercise(workoutID uint, datum time.Time) exerciseSelect {

	if datum.IsZero() {
		datum = time.Now()
	}

	m := exerciseSelect{
		mode:         MODE_EXERCISE_RETURNID,
		exerciseList: list.New(make([]list.Item, 0), coms.ListItemStyle(), 0, 0),
		workoutID:    workoutID,
		//ws:           ws,
		datum: datum,
	}

	return m
}

func (m *exerciseSelect) refreshExerciseList(exercises []wodb.Exercise) {
	// exercise list
	items := make([]list.Item, len(exercises))
	for i := range exercises {
		items[i] = coms.ExerciseItem{&exercises[i]}
	}
	exerciseList := list.New(items, coms.ListItemStyle(), 0, 0)
	exerciseList.Title = "Select Exercise"
	exerciseList.SetSize(ListWidth, 10) // -2: height of date
	exerciseList.SetShowTitle(false)
	exerciseList.FilterInput.Cursor.Style = coms.FilterCursorStyle
	exerciseList.FilterInput.PromptStyle = coms.FilterPromptStyle
	exerciseList.SetShowHelp(false)
	//exerciseList.Help.Styles.ShortDesc = exerciseList.Help.Styles.ShortDesc.Padding(0)

	exerciseList.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			exerciseSelectKeys.logExercise,
			exerciseSelectKeys.selectDate,
			exerciseSelectKeys.back,
		}
	}

	m.exerciseList = exerciseList
}

func (m exerciseSelect) Init() tea.Cmd {
	return tea.Batch(coms.ReloadExercises, textinput.Blink)
}

func (m exerciseSelect) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.exerciseList.SetHeight(coms.GetContentHeight(msg.Height) - 1 - 3) // - 1 for date

	case coms.MsgDate:
		m.datum = time.Time(msg)
		return m, cmd

	case coms.MsgExercisesReload:
		if msg.Err != nil {
			m.status = msg.Err.Error()
			return m, cmd
		}
		m.refreshExerciseList(msg.Exercises)
		return m, tea.Batch(cmd, tea.WindowSize())

	case coms.MsgExerciseAddedToWorkout:
		if msg.Err != nil {
			m.status = msg.Err.Error()
			return m, cmd
		}
		return m, coms.GoTo(NewWorkoutModel(m.workoutID, m.datum))

	case tea.KeyMsg:
		if m.exerciseList.FilterState() == list.Filtering {
			break
		}
		switch msg.String() {
		case "enter":
			switch m.mode {
			case MODE_EXERCISE_DO:
				return m, coms.GoTo(NewExerciseEntry(m.datum, m.exerciseList.SelectedItem().(coms.ExerciseItem).Exercise))
			case MODE_EXERCISE_RETURNID:
				return m, coms.Ret(coms.SendExerciseID(m.exerciseList.SelectedItem().(coms.ExerciseItem).ID))
			}

		case "esc":
			return m, coms.Back
		}
	}
	m.exerciseList, cmd = m.exerciseList.Update(msg)
	return m, cmd
}

func (m exerciseSelect) View() string {
	sb := &strings.Builder{}
	sb.WriteString(coms.FocusedStyle.Render("Date: ") + m.datum.Format("2006-01-02") + "\n")
	sb.WriteString(m.exerciseList.View() + "\n\n\n")
	sb.WriteString(m.exerciseList.Help.View(m.exerciseList))

	return sb.String()
}

func (m exerciseSelect) BreadCrumb() string {
	return "exercises"
}

func (m exerciseSelect) Help() string {
	return ""
}

// --------------------------------------------------------------------------------------

type exerciseSelectKeymap struct {
	logExercise key.Binding
	back        key.Binding
	selectDate  key.Binding
}

var exerciseSelectKeys = exerciseSelectKeymap{
	logExercise: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "log exercise"),
	),
	selectDate: key.NewBinding(
		key.WithKeys("f5"),
		key.WithHelp("f5", "change date"),
	),
	back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
}
