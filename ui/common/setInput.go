package common

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	wodb "github.com/zmnpl/clift/db"
)

type SetInput struct {
	SetNo      int
	Reps       textinput.Model
	Weight     textinput.Model
	WorkoutId  uint
	ExerciseId string
	Datum      time.Time
}

func (i *SetInput) FocusReps() tea.Cmd {
	cmd := i.Reps.Focus()
	i.Reps.PromptStyle = FocusedStyle
	i.Reps.TextStyle = FocusedStyle

	i.Weight.Blur()
	i.Weight.PromptStyle = NoStyle
	i.Weight.TextStyle = NoStyle

	return cmd
}

func (i *SetInput) FocusWeight() tea.Cmd {
	cmd := i.Weight.Focus()
	i.Weight.PromptStyle = FocusedStyle
	i.Weight.TextStyle = FocusedStyle

	i.Reps.Blur()
	i.Reps.PromptStyle = NoStyle
	i.Reps.TextStyle = NoStyle

	return cmd
}

func (i *SetInput) Unfocus() {
	i.Reps.Blur()
	i.Reps.PromptStyle = NoStyle
	i.Reps.TextStyle = NoStyle

	i.Weight.Blur()
	i.Weight.PromptStyle = NoStyle
	i.Weight.TextStyle = NoStyle
}

func CreateSetTemplatesForWE(we wodb.WorkoutExercise) []SetInput {
	inputs := make([]SetInput, 0, 999)
	// sets of workout exercise
	for i, s := range we.Sets {
		inputs = append(inputs, CreateSetTemplate(i+1, s.Reps, int(s.Weight), we.WorkoutID, we.ExerciseID))
	}

	return inputs
}

func CreateEmptySetTemplate(eid string, set_cnt int) []SetInput {
	inputs := make([]SetInput, 0, set_cnt)
	// sets of workout exercise
	for i := 0; i < set_cnt; i++ {
		inputs = append(inputs, CreateSetTemplate(i+1, 10, 0, 0, eid))
	}

	return inputs
}

func CreateSetTemplate(setno int, reps int, weight int, wrokoutId uint, exerciseId string) SetInput {
	repTextIn := textinput.New()
	repTextIn.Placeholder = fmt.Sprintf("%v", reps)
	repTextIn.CharLimit = 4
	repTextIn.Width = 20
	//repTextIn.Cursor.SetMode(cursor.CursorBlink)

	weightTextIn := textinput.New()
	weightTextIn.Placeholder = fmt.Sprintf("%v", weight)
	weightTextIn.CharLimit = 50
	weightTextIn.Width = 100

	template := SetInput{
		SetNo:      setno,
		Reps:       repTextIn,
		Weight:     weightTextIn,
		WorkoutId:  wrokoutId,
		ExerciseId: exerciseId,
		Datum:      time.Now(),
	}

	return template
}
