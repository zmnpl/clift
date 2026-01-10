package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	wodb "github.com/zmnpl/clift/db"
	coms "github.com/zmnpl/clift/ui/common"
)

const (
	MODE_LOGSETS = iota
	MODE_EDITSETS
	MODE_RETURN_SETS
)

type exerciseEntry struct {
	focusIndex int
	setInputs  []coms.SetInput
	datum      time.Time

	workout         *wodb.Workout
	workoutExercise *wodb.WorkoutExercise
	exercise        *wodb.Exercise

	mode int

	help help.Model
}

func NewExerciseEntry(datum time.Time, exercise *wodb.Exercise) exerciseEntry {
	if datum.IsZero() {
		datum = time.Now()
	}

	m := exerciseEntry{
		datum:    datum,
		exercise: exercise,
		help:     help.New(),
	}

	if m.exercise != nil {
		m.setInputs = coms.CreateEmptySetTemplate(m.exercise.ID, 3)
	} else if m.workoutExercise != nil {
		m.setInputs = coms.CreateSetTemplatesForWE(*m.workoutExercise)
		m.exercise = &m.workoutExercise.Exercise
	}

	m.setInputs[0].FocusReps()

	return m
}

func NewWorkoutExerciseEntryModel(datum time.Time, workout *wodb.Workout, workoutExercise *wodb.WorkoutExercise, exercise *wodb.Exercise, setInputs []coms.SetInput, mode int) exerciseEntry {
	if datum.IsZero() {
		datum = time.Now()
	}

	m := exerciseEntry{
		datum:           datum,
		workout:         workout,
		workoutExercise: workoutExercise,
		exercise:        exercise,
		setInputs:       setInputs,
		mode:            mode,
		help:            help.New(),
	}

	if len(setInputs) > 0 {
		m.setInputs[0].FocusReps()
	}

	return m
}

func (m exerciseEntry) Init() tea.Cmd {
	return tea.Batch(textinput.Blink)
}

func (m exerciseEntry) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case coms.MsgExerciseLogged:
		if msg.Err == nil {
			return m, coms.Back
		} else {
			// TODO: bubble up error to user
		}

	case coms.MsgDate:
		m.datum = time.Time(msg)
		return m, cmd

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "q":
			// TODO - apply reps / weight from selected inputs plan to actual
			return m, coms.SendStatus("foo", nil)

		case "f9":
			if m.focusIndex < len(m.setInputs) {
				m.setInputs[m.focusIndex].PlaceholderToValue()
			}
			return m, cmd

		case "f10":
			for i := 0; i < len(m.setInputs); i++ {
				m.setInputs[i].PlaceholderToValue()
			}
			return m, cmd

		case "enter":
			if m.focusIndex == len(m.setInputs) {
				if m.mode == MODE_RETURN_SETS {
					return m, coms.Ret(coms.SendPerformedSets(m.setInputs, m.workoutExercise.ID))
				}
				return m, coms.LogSingleExercise(m.datum, m.setInputs)
			}

		case "+":
			if len(m.setInputs) == 0 {
				m.setInputs = make([]coms.SetInput, 0, 10)
			}

			// start with dummy wid
			var mywid uint
			if m.workout != nil {
				mywid = m.workout.ID
			}
			m.setInputs = append(m.setInputs, coms.CreateSetTemplate(
				len(m.setInputs)+1,
				10,
				0,
				mywid,
				m.exercise.ID))

			if len(m.setInputs) > 0 {
				m.setInputs[m.focusIndex].Unfocus()
				m.focusIndex = len(m.setInputs) - 1
			}

			return m, m.setInputs[m.focusIndex].FocusReps()

		case "-":
			if len(m.setInputs) > 0 {
				m.setInputs[m.focusIndex].Unfocus()
				m.setInputs = m.setInputs[:len(m.setInputs)-1]
			}
			if len(m.setInputs) > 0 {
				m.focusIndex = len(m.setInputs) - 1
				return m, m.setInputs[m.focusIndex].FocusReps()
			}
			return m, cmd

		case "tab", "shift+tab", "up", "down":
			s := msg.String()

			same, up, down := 0, 1, -1
			direction := same
			switch s {
			case "up", "shift+tab":
				if m.focusIndex == len(m.setInputs) || m.setInputs[m.focusIndex].Reps.Focused() {
					m.focusIndex--
					direction = up
				}
			case "down", "tab":
				if m.focusIndex == len(m.setInputs) || m.setInputs[m.focusIndex].Weight.Focused() {
					m.focusIndex++
					direction = down
				}
			}

			if m.focusIndex > len(m.setInputs) {
				m.focusIndex = len(m.setInputs)
			}
			if m.focusIndex < 0 {
				m.focusIndex = 0
			}

			// Set/unset focus
			cmds := make([]tea.Cmd, len(m.setInputs))
			for i := 0; i <= len(m.setInputs)-1; i++ {
				if i == m.focusIndex {
					// up top
					if direction == up && i == 0 && m.setInputs[i].Reps.Focused() {
						continue
					}
					// within line
					if m.setInputs[i].Reps.Focused() {
						cmds = append(cmds, m.setInputs[i].FocusWeight())
						continue
					}
					if m.setInputs[i].Weight.Focused() {
						cmds = append(cmds, m.setInputs[i].FocusReps())
						continue
					}
					//moved up or down
					if direction == up {
						cmds = append(cmds, m.setInputs[i].FocusWeight())
						continue
					}
					if direction == down {
						cmds = append(cmds, m.setInputs[i].FocusReps())
						continue
					}
				}
				// remove focus from all others
				m.setInputs[i].Unfocus()
			}

			return m, tea.Batch(cmds...)

		case "esc":
			return m, coms.Back

		}
	}
	return m, m.updateInputs(msg)
}

func (m exerciseEntry) View() string {
	sb := &strings.Builder{}

	sb.WriteString(coms.FocusedStyle.Render("Date: ") + m.datum.Format("2006-01-02") + "\n\n")
	for _, v := range m.setInputs {
		sb.WriteString(fmt.Sprintf("%v | Reps %s Weight %s\n", v.SetNo, v.Reps.View(), v.Weight.View()))
	}

	button := blurredButton
	if m.focusIndex == len(m.setInputs) {
		button = focusedButton
	}
	sb.WriteString(fmt.Sprintf("\n%v\n", button))

	return sb.String()
}

func (m exerciseEntry) Help() string {
	sb := &strings.Builder{}

	helpView := m.help.View(exerciseEntryKeys)
	sb.WriteString(helpView)

	return sb.String()
}

func (m *exerciseEntry) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.setInputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.setInputs {
		var cmd tea.Cmd

		m.setInputs[i].Reps, cmd = m.setInputs[i].Reps.Update(msg)
		cmds = append(cmds, cmd)

		m.setInputs[i].Weight, cmd = m.setInputs[i].Weight.Update(msg)
		cmds = append(cmds, cmd)
	}
	return tea.Batch(cmds...)
}

func (m exerciseEntry) BreadCrumb() string {
	return m.exercise.GetName()
}

//------------------------------------------------------

type exerciseEntryKeymap struct {
	nav                 key.Binding
	confirm             key.Binding
	back                key.Binding
	changedate          key.Binding
	applyPlaceholder    key.Binding
	applyPlaceholderAll key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k exerciseEntryKeymap) ShortHelp() []key.Binding {
	return []key.Binding{k.nav, k.changedate, k.applyPlaceholder, k.applyPlaceholderAll, k.confirm, k.back}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k exerciseEntryKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.nav, k.applyPlaceholder, k.changedate}, // first column
		{k.confirm, k.back},                       // second column
	}
}

var exerciseEntryKeys = exerciseEntryKeymap{
	nav: key.NewBinding(
		key.WithKeys("up", "k", "down", "j"),
		key.WithHelp("↑/k/↓/j", "navigate"),
	),
	changedate: key.NewBinding(
		key.WithKeys("f5"),
		key.WithHelp("f5", "change date"),
	),
	applyPlaceholder: key.NewBinding(
		key.WithKeys("f9"),
		key.WithHelp("f9", "apply placeholder"),
	),
	applyPlaceholderAll: key.NewBinding(
		key.WithKeys("f10"),
		key.WithHelp("f10", "apply placeholder all"),
	),
	confirm: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "confirm"),
	),
	back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
}
