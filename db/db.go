package db

import (
	"context"
	"time"

	"github.com/tidwall/gjson"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// func (m *Workout) AfterDelete(tx *gorm.DB) (err error) {
// 	tx.Clauses(clause.Returning{}).Where("model_a_id = ?", m.ID).Delete(&ModelB{})
// 	return
// }

type Workout struct {
	ID               uint              `gorm:"primaryKey;not null"`
	Name             string            `gorm:"not null"`
	WorkoutExercises []WorkoutExercise `gorm:"foreignKey:WorkoutID;constraint:OnDelete:CASCADE"`
	PerformedSets    []PerformedSet    `gorm:"foreignKey:WorkoutID"`
	Deleted          gorm.DeletedAt
}

type Exercise struct {
	ID               string `gorm:"primaryKey;not null"`
	Data             string
	WorkoutExercises []WorkoutExercise `gorm:"foreignKey:ExerciseID"`
	PerformedSets    []PerformedSet    `gorm:"foreignKey:ExerciseID"`
}

type WorkoutExercise struct {
	ID         uint     `gorm:"primaryKey;not null"`
	WorkoutID  uint     `gorm:"not null"`
	Workout    Workout  `gorm:"foreignKey:WorkoutID;references:ID"`
	ExerciseID string   `gorm:"not null"`
	Exercise   Exercise `gorm:"foreignKey:ExerciseID;references:ID"`
	//Exercise   Exercise `gorm:"foreignKey:ExerciseID;references:ID;constraint:OnDelete:CASCADE"`
	Sets    []Set `gorm:"constraint:OnDelete:CASCADE"`
	Note    string
	Deleted gorm.DeletedAt
}

type Set struct {
	ID                uint    `gorm:"primaryKey;not null"`
	WorkoutExerciseID uint    `gorm:"not null"`
	Reps              int     `gorm:"not null"`
	Weight            float64 `gorm:"not null"`
}

type PerformedSet struct {
	ID            uint `gorm:"primaryKey;not null"`
	WorkoutID     uint
	ExerciseID    string
	PerformedDate time.Time
	SetNo         int
	Reps          int
	Weight        float64
}

// Exercise

func (e Exercise) GetName() string {
	return gjson.Get(e.Data, "name").String()
}

func (e Exercise) GetMusclesString() []gjson.Result {
	// TODO: Maybe merge secondary muscles
	primary := gjson.Get(e.Data, "primaryMuscles").Array()
	secondary := gjson.Get(e.Data, "secondaryMuscles").Array()
	return append(primary, secondary...)
}

func (e Exercise) GetPrimaryMuscles() []string {
	j := gjson.Get(e.Data, "primaryMuscles").Array()
	result := make([]string, len(j))
	for i := range j {
		result[i] = j[i].String()
	}
	return result
}

func (e Exercise) GetSecondaryMuscles() []string {
	j := gjson.Get(e.Data, "secondaryMuscles").Array()
	result := make([]string, len(j))
	for i := range j {
		result[i] = j[i].String()
	}
	return result
}

// --- DB manager functions ---

func (t *TrainingDB) CreateWorkout(name string) (*Workout, error) {
	w := &Workout{Name: name}
	result := t.db.Create(w)
	return w, result.Error
}

func (t *TrainingDB) RemoveWorkout(id uint) error {
	var workout Workout
	if err := t.db.First(&workout, id).Error; err != nil {
		return err
	}

	// // 💡 Step 1: Check PRAGMA status on the current connection
	// var fkStatus int
	// // Use db.Raw() to run a direct query
	// if err := t.db.Raw("PRAGMA foreign_keys").Scan(&fkStatus).Error; err != nil {
	// 	log.Printf("Error checking PRAGMA status: %v", err)
	// } else {
	// 	log.Printf("PRAGMA foreign_keys status before delete: %d (should be 1)", fkStatus)
	// }

	// // 💡 Step 2: Ensure PRAGMA is ON one last time (redundant but crucial test)
	// //t.db.Exec("PRAGMA foreign_keys = ON")

	// Step 3: Execute the delete
	return t.db.Select(clause.Associations).Delete(&workout).Error
	//return t.db.Delete(&workout).Error
}

func (t *TrainingDB) GetAllWorkouts() ([]Workout, error) {
	var ws []Workout
	err := t.db.Preload("WorkoutExercises").
		Preload("WorkoutExercises.Exercise").
		Preload("WorkoutExercises.Sets").
		Find(&ws).Error
	return ws, err
}

func (t *TrainingDB) GetWorkoutWithExercises(workoutID uint) (Workout, error) {
	var workout Workout
	err := t.db.
		Preload("WorkoutExercises").
		Preload("WorkoutExercises.Exercise").
		Preload("WorkoutExercises.Sets").
		First(&workout, workoutID).Error

	if err != nil {
		return workout, err
	}

	return workout, nil
}

func (t *TrainingDB) GetAllExercises() ([]Exercise, error) {
	var es []Exercise
	err := t.db.Find(&es).Error
	return es, err
}

func (t *TrainingDB) AddExerciseToWorkout(workoutID uint, exerciseID string, note string) (*WorkoutExercise, error) {
	we := &WorkoutExercise{
		WorkoutID:  workoutID,
		ExerciseID: exerciseID,
		Note:       note,
	}
	if err := t.db.Create(we).Error; err != nil {
		return nil, err
	}
	return we, nil
}

// TODO - revise AI code
func (t *TrainingDB) UpdateWorkoutExerciseSets(weID uint, newSets []Set) error {
	return t.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Unscoped().
			Where("workout_exercise_id = ?", weID).
			Delete(&Set{})

		if result.Error != nil {
			return result.Error // Return error to rollback the transaction
		}

		// // 2. Prepare the new sets for insertion.
		// // Before bulk creation, we must ensure each new Set has the correct foreign key.
		// setsToCreate := make([]Set, len(newSets))
		// for i, set := range newSets {
		// 	// Copy the set data and explicitly set the foreign key.
		// 	setsToCreate[i] = Set{
		// 		WorkoutExerciseID: weID,
		// 		Reps:              set.Reps,
		// 		Weight:            set.Weight,
		// 	}
		// }

		// 3. Bulk create the new sets.
		if len(newSets) > 0 {
			result = tx.CreateInBatches(&newSets, len(newSets))
			if result.Error != nil {
				return result.Error // Return error to rollback the transaction
			}
		}

		return nil
	})
}

func (t *TrainingDB) RemoveExerciseFromWorkout(weID uint) error {
	return t.db.Delete(&WorkoutExercise{}, weID).Error
}

func (t *TrainingDB) AddSetTemplate(weID uint, reps int, weight float64) (*Set, error) {
	s := &Set{
		WorkoutExerciseID: weID,
		Reps:              reps,
		Weight:            weight,
	}
	if err := t.db.Create(s).Error; err != nil {
		return nil, err
	}
	return s, nil
}

func (t *TrainingDB) UpdateSetTemplate(setID uint, reps int, weight float64) error {
	return t.db.Model(&Set{}).Where("id = ?", setID).
		Updates(Set{Reps: reps, Weight: weight}).Error
}

func (t *TrainingDB) GetSetsForWorkoutExercise(weID uint) ([]Set, error) {
	var ss []Set
	err := t.db.Where("workout_exercise_id = ?", weID).Find(&ss).Error
	return ss, err
}

func (t *TrainingDB) GetAllPerformedSets() ([]PerformedSet, error) {
	var s []PerformedSet
	err := t.db.Order("performed_date desc, exercise_id asc, set_no asc").Find(&s).Error
	return s, err
}

func (t *TrainingDB) LogPerformedSet(workoutID uint, exerciseID string, setNo, reps int, weight float64, performedDate time.Time) error {
	if reps <= 0 {
		return nil
	}

	p := &PerformedSet{
		WorkoutID:     workoutID,
		ExerciseID:    exerciseID,
		SetNo:         setNo,
		Reps:          reps,
		Weight:        weight,
		PerformedDate: performedDate,
	}
	return t.db.Create(p).Error
}

func (t *TrainingDB) LogSet(set PerformedSet) error {
	if set.Reps <= 0 {
		return nil
	}

	return t.db.Create(&set).Error
}

func (t *TrainingDB) LogSetsTransaction(sets []PerformedSet) error {
	return t.db.Transaction(func(tx *gorm.DB) error {
		for _, set := range sets {
			// zero reps are not logged
			if set.Reps <= 0 {
				continue
			}

			err := tx.Create(&set).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// raw stuff

func (t *TrainingDB) MyRaw() [][]string {
	rows, err := gorm.G[any](t.db).Raw("select exercise_id, set_count, calendar_week from vw_weekly_volume where exercise_id = ?", "Pullups").Rows(context.Background())
	if err != nil {

	}
	defer rows.Close()

	foo := make([][]string, 0, 100)
	for rows.Next() {
		var exercise_id, set_count, calendar_week string
		r := make([]string, 3)

		rows.Scan(&exercise_id, &set_count, &calendar_week)
		r[0] = exercise_id
		r[1] = set_count
		r[2] = calendar_week

		foo = append(foo, r)
		// do something
	}
	return foo
}
