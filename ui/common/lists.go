package common

import (
	"fmt"
	"strconv"
	"strings"

	wodb "github.com/zmnpl/clift/db"
)

// ------------------------------------------

type WorkoutItem struct {
	*wodb.Workout
}

func (wi WorkoutItem) Title() string { return wi.Name }
func (wi WorkoutItem) Description() string {
	return fmt.Sprintf("%v Exercices", len(wi.WorkoutExercises))
}
func (wi WorkoutItem) FilterValue() string { return wi.Name }

// ------------------------------------------

type WeItem struct {
	*wodb.WorkoutExercise
	SetInputs []SetInput
}

func (we WeItem) Title() string { return we.Exercise.ID }
func (we WeItem) Description() string {
	sb := &strings.Builder{}

	for _, setInput := range we.SetInputs {
		reps := setInput.Reps.Placeholder
		weight := setInput.Weight.Placeholder

		doneReps, _ := strconv.Atoi(setInput.Reps.Value())
		doneWeight, _ := strconv.Atoi(setInput.Weight.Value())

		sb.WriteString(fmt.Sprintf("%v (%v) reps @ %v (%v) kg\n", doneReps, reps, doneWeight, weight))
	}

	return sb.String()
}
func (we WeItem) FilterValue() string { return we.Exercise.ID }

// ------------------------------------------

type ExerciseItem struct {
	*wodb.Exercise
}

func (ei ExerciseItem) Title() string { return ei.GetName() }
func (ei ExerciseItem) Description() string {
	primary := ei.GetPrimaryMuscles()
	secondary := ei.GetSecondaryMuscles()
	return fmt.Sprintf("[%v](%v)", strings.Join(primary, " "), strings.Join(secondary, " "))
}
func (ei ExerciseItem) FilterValue() string {
	return ei.GetName() + fmt.Sprintf(" | %v", ei.GetMusclesString())
}
