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
	MODE_DO = iota
	MODE_EDIT
)

type workout struct {
	workoutID uint
	workout   wodb.Workout

	datum        time.Time
	exerciseList list.Model

	sessionSets map[uint][]coms.SetInput

	mode           int
	escapeUnlocked bool
	deleteUnlocked bool

	status string
}

func NewWorkoutModel(workoutID uint, datum time.Time) workout {
	if datum.IsZero() {
		datum = time.Now()
	}

	items := make([]list.Item, 0)
	l := list.New(items, list.NewDefaultDelegate(), 0, 0)

	return workout{
		datum:        datum,
		workoutID:    workoutID,
		exerciseList: l,
		sessionSets:  make(map[uint][]coms.SetInput),
		mode:         MODE_DO,
	}
}

func NewWorkoutModelEDIT(workoutID uint, datum time.Time) workout {
	if datum.IsZero() {
		datum = time.Now()
	}

	items := make([]list.Item, 0)
	l := list.New(items, list.NewDefaultDelegate(), 0, 0)

	return workout{
		datum:        datum,
		workoutID:    workoutID,
		exerciseList: l,
		sessionSets:  make(map[uint][]coms.SetInput),
		mode:         MODE_EDIT,
	}
}

func (m workout) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, coms.ReloadWorkoutSingle(m.workoutID))
}

func (m workout) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.exerciseList.SetHeight(coms.GetContentHeight(msg.Height) - 3)

	case coms.LockCriticalKey:
		m.deleteUnlocked = false
		m.escapeUnlocked = false
		return m, coms.SendStatus("", nil)

	case coms.MsgExerciseID:
		return m, coms.AddWorkoutExercise(m.workoutID, string(msg))

	case coms.MsgDate:
		m.datum = time.Time(msg)
		return m, cmd

	case coms.MsgPerformedSets:
		if m.mode == MODE_DO {
			m.updateSetsForWE(msg.Sets, msg.Weid)
			m.refreshWEList()
			return m, tea.Batch(cmd, tea.WindowSize())
		}
		if m.mode == MODE_EDIT {
			return m, tea.Batch(coms.UpdateWorkoutExerciseSets(msg.Weid, msg.Sets), tea.WindowSize())
		}

	case coms.MsgUpdatedWorkoutExercise:
		if msg.Err != nil {
			m.status = msg.Err.Error()
		}

		return m, coms.ReloadWorkoutSingle(m.workout.ID)

	case coms.MsgExerciseAddedToWorkout:
		return m, coms.ReloadWorkoutSingle(m.workout.ID)

	case coms.MsgWorkoutSingleReload:
		if msg.Err != nil {
			m.status = msg.Err.Error()
		} else {
			m.workout = msg.Workout
			m.refreshWEList()
			return m, tea.Batch(cmd, tea.WindowSize())
		}

	case tea.KeyMsg:
		if m.exerciseList.FilterState() == list.Filtering {
			break
		}

		switch msg.String() {
		case "enter":
			if m.mode == MODE_DO && len(m.exerciseList.SelectedItem().(coms.WeItem).SetInputs) < 1 {
				m.status = "No sets defined"
				return m, cmd
			}

			weitem := m.exerciseList.SelectedItem().(coms.WeItem)
			return m, coms.GoTo(NewWorkoutExerciseEntryModel(m.datum, &m.workout, weitem.WorkoutExercise, &weitem.WorkoutExercise.Exercise, weitem.SetInputs, MODE_RETURN_SETS))

		case "+":
			return m, coms.GoTo(NewSelectExercise(m.workout.ID, m.datum))

		case "f1":
			if m.mode == MODE_DO {
				weItems := make([]coms.WeItem, len(m.exerciseList.Items()))
				for i, v := range m.exerciseList.Items() {
					weItems[i] = v.(coms.WeItem)
				}
				return m, coms.LogWorkout(weItems, m.datum)
			}
			if m.mode == MODE_EDIT {

			}

		case "delete":
			if m.deleteUnlocked {
				return m, coms.RemoveWorkoutExercise(m.exerciseList.SelectedItem().(coms.WeItem).WorkoutExercise.ID)
			}
			m.deleteUnlocked = true
			return m, tea.Batch(coms.SendStatus("Yo! Press delete one more time and that workout is gone.", nil), coms.SleepToLockKey(2000*time.Millisecond))

		case "esc":
			if m.escapeUnlocked {
				return m, tea.Batch(coms.Back, coms.SendStatus("", nil))
			}
			m.escapeUnlocked = true
			return m, tea.Batch(coms.SendStatus("Next time you press that, we go back without saving...", nil), coms.SleepToLockKey(2000*time.Millisecond))
		}
	}

	m.exerciseList, cmd = m.exerciseList.Update(msg)
	return m, cmd
}

func (m workout) View() string {
	sb := &strings.Builder{}
	sb.WriteString(coms.FocusedStyle.Render("Date: ") + m.datum.Format("2006-01-02") + "\n")
	sb.WriteString(m.exerciseList.View() + "\n\n")
	sb.WriteString(m.exerciseList.Help.View(m.exerciseList))
	return sb.String()
}

func (m *workout) updateSetsForWE(sets []coms.SetInput, weid uint) {
	m.sessionSets[weid] = sets
}

func (m *workout) refreshWEList() {
	wes := m.workout.WorkoutExercises

	// workout list
	maxSets := 1
	items := make([]list.Item, len(wes))
	for i := range wes {
		templates := coms.CreateSetTemplatesForWE(wes[i])

		// overwrite with user entered sessoin sets
		sessionTemplates, ok := m.sessionSets[wes[i].ID]
		if ok {
			templates = sessionTemplates
		}

		items[i] = coms.WeItem{
			&wes[i],
			templates,
		}
		if len(templates) > maxSets {
			maxSets = len(templates)
		}
	}
	maxSets = maxSets + 1

	d := coms.ListItemStyle()
	d.SetHeight(maxSets)

	l := list.New(items, d, 0, 0)
	l.Title = "Train hard"
	l.SetSize(ListWidth, 10)
	l.SetShowTitle(false)
	l.SetShowHelp(false)
	//l.SetShowStatusBar(false)

	l.FilterInput.Cursor.Style = coms.FilterCursorStyle
	l.FilterInput.PromptStyle = coms.FilterPromptStyle

	// TODO keys
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			workoutKeys.submit,
			workoutKeys.enter,
			workoutKeys.addExercise,
			workoutKeys.changedate,
			workoutKeys.back,
		}
	}

	m.exerciseList = l
}

func (m workout) BreadCrumb() string {
	edit := ""
	if m.mode == MODE_EDIT {
		edit = " (edit)"
	}
	return m.workout.Name + edit
}

func (m workout) Help() string {
	return ""
}

func (m workout) summary() string {
	sb := &strings.Builder{}
	sb.WriteString(m.workout.Name + "\n")
	// for _, item := range m.exerciseList.Items() {

	// }
	return sb.String()
}

// ------------------------------------------------------------------------------
type workoutKeymap struct {
	enter       key.Binding
	addExercise key.Binding
	submit      key.Binding
	changedate  key.Binding
	back        key.Binding
}

var workoutKeys = workoutKeymap{
	enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "do"),
	),
	addExercise: key.NewBinding(
		key.WithKeys("+"),
		key.WithHelp("+", "add exercise"),
	),
	submit: key.NewBinding(
		key.WithKeys("f1"),
		key.WithHelp("f1", "submit"),
	),
	changedate: key.NewBinding(
		key.WithKeys("f5"),
		key.WithHelp("f5", "change date"),
	),
	back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
}
