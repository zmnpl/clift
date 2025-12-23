package common

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	wodb "github.com/zmnpl/clift/db"
)

type StatusMsg struct {
	Status string
	Err    error
}

type MsgWorkoutAddEdit string

type MsgExercisesReload struct {
	Exercises []wodb.Exercise
	Err       error
}

type MsgWorkoutsReload struct {
	Workouts []wodb.Workout
	Err      error
}

type MsgWorkoutSingleReload struct {
	Workout wodb.Workout
	Err     error
}

type MsgPerformedSetsLoaded struct {
	PerformedSets []wodb.PerformedSet
	Err           error
}

type MsgExerciseAddedToWorkout struct {
	Exercise wodb.WorkoutExercise
	Err      error
}

type MsgUpdatedWorkoutExercise struct {
	Err error
}

type MsgExerciseLogged struct {
	Err error
}

type MsgExerciseID string

type MsgDate time.Time

type MsgPerformedSets struct {
	Weid uint
	Sets []SetInput
}

type MsgGoTo tea.Model

type MsgBack bool

type LockCriticalKey bool

type MsgBackCmd struct {
	Cmd tea.Cmd
}

func GoTo(model tea.Model) func() tea.Msg {
	return func() tea.Msg {
		return MsgGoTo(model)
	}
}

func Back() tea.Msg {
	return MsgBack(true)
}

func Ret(cmd tea.Cmd) func() tea.Msg {
	return func() tea.Msg {
		return MsgBackCmd{
			Cmd: cmd,
		}
	}
}

func SleepToLockKey(duration time.Duration) func() tea.Msg {
	return func() tea.Msg {
		time.Sleep(duration)
		return LockCriticalKey(true)
	}
}

func SendStatus(status string, err error) func() tea.Msg {
	return func() tea.Msg {
		return StatusMsg{
			Status: status,
			Err:    err,
		}
	}
}

func SendDate(d time.Time) func() tea.Msg {
	return func() tea.Msg {
		return MsgDate(d)
	}
}

func SendExerciseID(id string) func() tea.Msg {
	return func() tea.Msg {
		return MsgExerciseID(id)
	}
}

func SendPerformedSets(sets []SetInput, weid uint) func() tea.Msg {
	return func() tea.Msg {
		return MsgPerformedSets{
			Weid: weid,
			Sets: sets,
		}
	}
}

func ReloadExercises() tea.Msg {
	exercises, err := wodb.Instance().GetAllExercises()
	if err != nil {
		log.Fatalf("Failed to get all exercises: %v", err) // TODO: maybe pass error to ui
	}

	return MsgExercisesReload{
		Exercises: exercises,
		Err:       err,
	}
}

func ReloadWorkouts() tea.Msg {
	workouts, err := wodb.Instance().GetAllWorkouts()
	if err != nil {
		log.Fatalf("Failed to get all workout: %v", err)
	}

	return MsgWorkoutsReload{
		Workouts: workouts,
		Err:      err,
	}
}

func ReloadWorkoutSingle(wid uint) func() tea.Msg {
	return func() tea.Msg {
		workout, err := wodb.Instance().GetWorkoutWithExercises(wid)
		if err != nil {
			log.Fatalf("Could not reload workout ID: %v (%v)", wid, err)
		}

		return MsgWorkoutSingleReload{
			Workout: workout,
			Err:     err,
		}
	}
}

func AddWorkoutExercise(workoutID uint, eid string) func() tea.Msg {
	return func() tea.Msg {
		ex, err := wodb.Instance().AddExerciseToWorkout(workoutID, eid, "")
		if err != nil {
			log.Fatalf("Could not add exercise (%v) to workout (%v): %v", eid, workoutID, err)
		}
		return MsgExerciseAddedToWorkout{
			Exercise: *ex,
			Err:      err,
		}
	}
}

func RemoveWorkoutExercise(weID uint) func() tea.Msg {
	return func() tea.Msg {
		return MsgUpdatedWorkoutExercise{
			Err: wodb.Instance().RemoveExerciseFromWorkout(weID),
		}
	}
}

func NewWorkout(name string) func() tea.Msg {
	return func() tea.Msg {
		_, err := wodb.Instance().CreateWorkout(name)
		if err != nil {
			return StatusMsg{Status: "", Err: fmt.Errorf("Error creating workout: %v", err)}
		}

		return MsgWorkoutAddEdit("Created " + name)
	}
}

func RemoveWorkout(id uint) func() tea.Msg {
	return func() tea.Msg {
		err := wodb.Instance().RemoveWorkout(id)
		if err != nil {
			return StatusMsg{Status: "", Err: fmt.Errorf("Error removing workout: %v", err)}
		}

		return MsgWorkoutAddEdit("Removed workout")
	}
}

func MakeDBPerformedSet(set SetInput, setno int, datum time.Time) wodb.PerformedSet {
	reps, err := strconv.Atoi(set.Reps.Value())
	if err != nil {
		reps = 0
		// TODO - pass error to user maybe
	}
	weight, err := strconv.ParseFloat(set.Weight.Value(), 64)
	if err != nil {
		weight = 0
	}

	foo := wodb.PerformedSet{
		WorkoutID:     set.WorkoutId,
		ExerciseID:    set.ExerciseId,
		SetNo:         setno,
		Reps:          reps,
		Weight:        weight,
		PerformedDate: datum,
	}

	return foo
}

func LoadPerformedSets() tea.Msg {
	performedSets, err := wodb.Instance().GetAllPerformedSets()
	if err != nil {
		log.Printf("Could not get performed sets: %v", err)
	}

	return MsgPerformedSetsLoaded{
		PerformedSets: performedSets,
		Err:           err,
	}
}

func MakeJournal(performedSets []wodb.PerformedSet) table.Model {

	columns := []table.Column{
		{Title: "Exercise", Width: 30},
		{Title: "Date", Width: 12},
		{Title: "Set", Width: 5},
		{Title: "Reps", Width: 8},
		{Title: "Weight", Width: 8},
	}

	rows := make([]table.Row, 0, len(performedSets))
	for _, s := range performedSets {
		r := table.Row{
			s.ExerciseID,
			fmt.Sprintf("%v", s.PerformedDate),
			fmt.Sprintf("%v", s.SetNo),
			fmt.Sprintf("%v", s.Reps),
			fmt.Sprintf("%v", s.Weight),
		}
		rows = append(rows, r)
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(20),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(Theme.Foreground).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color(Theme.Red)).
		//Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return t
}

// model commands

func LogSingleExercise(datum time.Time, setInputs []SetInput) func() tea.Msg {
	return func() tea.Msg {
		sets := make([]wodb.PerformedSet, 0, 30)
		for i, s := range setInputs {
			foo := MakeDBPerformedSet(s, i, datum)
			sets = append(sets, foo)
		}

		err := wodb.Instance().LogSetsTransaction(sets)
		if err != nil {
			return MsgExerciseLogged{Err: fmt.Errorf("error logging your sets: %v", err.Error())}
		}
		return MsgExerciseLogged{}
	}
}

func UpdateWorkoutExerciseSets(weid uint, setInputs []SetInput) func() tea.Msg {
	return func() tea.Msg {
		sets := make([]wodb.Set, 0, len(setInputs))
		for _, v := range setInputs {

			// parse reps and weight
			// try value first, then placeholder
			// if 0 reps -> skip
			// 0 weight is ok
			reps, err := strconv.Atoi(v.Reps.Value())
			if err != nil {
				reps, err = strconv.Atoi(v.Reps.Placeholder)
				if err != nil || reps == 0 {
					continue
				}
			}

			weight, err := strconv.ParseFloat(v.Weight.Value(), 64)
			if err != nil {
				weight, err = strconv.ParseFloat(v.Weight.Placeholder, 64)
				if err != nil {
					weight = 0
				}
			}

			sets = append(sets, wodb.Set{
				WorkoutExerciseID: weid,
				Reps:              reps,
				Weight:            weight,
			})
		}

		err := wodb.Instance().UpdateWorkoutExerciseSets(weid, sets)
		if err != nil {
			log.Fatalf("Could not update sets for  weid: %v: %v", weid, err)
		}

		return MsgUpdatedWorkoutExercise{
			Err: err,
		}
	}
}

func LogWorkout(weItems []WeItem, datum time.Time) func() tea.Msg {
	return func() tea.Msg {
		sets := make([]wodb.PerformedSet, 0, 30)
		for _, we := range weItems {
			for i, s := range we.SetInputs {
				foo := MakeDBPerformedSet(s, i, datum)
				sets = append(sets, foo)
			}
		}

		err := wodb.Instance().LogSetsTransaction(sets)
		if err != nil {
			return StatusMsg{Status: "", Err: fmt.Errorf("Error logging your sets: %v", err.Error())}
		}
		return StatusMsg{Status: "Good job, logged workout 💪"}
	}
}
