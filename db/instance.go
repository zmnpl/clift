package db

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"sync"

	//"gorm.io/driver/sqlite"
	"github.com/glebarez/sqlite" // use this for pure go (no c)
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var lock = &sync.Mutex{}

type TrainingDB struct {
	db *gorm.DB

	// err = db.SetupJoinTable(&Workout{}, "WorkoutExercises", &WorkoutExercise{})
	// if err != nil {
	// 	log.Fatalf("failed to migrate: %v", err)
	// }
}

var instance *TrainingDB

func Instance() *TrainingDB {
	return instance
}

func Init(path string) {
	dir := filepath.Dir(path)

	// if db not existing
	firstTime := false
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		// create directory tree
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			log.Fatalf("Could not create folder structure: %v", err)
		}
		firstTime = true
	}

	if _, err := os.Open(dir); os.IsNotExist(err) {
		log.Fatalf("Cannot access db path (%v): %v", path, err)
	}

	if instance == nil {
		lock.Lock()
		defer lock.Unlock()
		if instance == nil {
			dsn := path + "?_pragma=foreign_keys(1)"

			db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
			if err != nil {
				log.Fatalf("failed to open DB: %v", err)
			}

			err = db.AutoMigrate(&Workout{}, &Exercise{}, &PerformedSet{}, &WorkoutExercise{}, &Set{})
			if err != nil {
				log.Fatalf("failed to migrate: %v", err)
			}

			// if first time, fill up with basic data
			if firstTime {
				// use session here to create silent logger; somehow prints all statements to stdout if using blunt exec
				err := db.Session(&gorm.Session{Logger: db.Logger.LogMode(logger.Silent)}).Exec(INSERT_DATA).Error

				if err != nil {
					log.Fatalf("failed to insert initial data: %v", err)
				}
				//db.Exec(INSERT_DATA)
			}

			instance = &TrainingDB{db: db}
		}
	}

}

func NewTrainingDB(path string) *TrainingDB {
	//dsn := path + "?_pragma=foreign_keys(1)"
	//dsn := path + "?_foreign_keys=on"
	dsn := path + "?_pragma=foreign_keys(1)"

	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to open DB: %v", err)
	}

	// // 💡 Step 1: Get the underlying database connection pool
	// sqlDB, err := db.DB()
	// if err != nil {
	// 	log.Fatalf("failed to get sql.DB: %v", err)
	// }

	// // SQLite usually runs fine with a single connection (default).
	// // If you explicitly set MaxOpenConns to 1, you guarantee
	// // that the subsequent Exec will target the only active connection.
	// //sqlDB.SetMaxOpenConns(1)

	// // 💡 Step 2: Manually enforce the PRAGMA on the primary connection.
	// // This is run once during initialization.
	// if _, err := sqlDB.Exec("PRAGMA foreign_keys = ON"); err != nil {
	// 	log.Fatalf("failed to enforce PRAGMA foreign_keys = ON: %v", err)
	// }

	// Step 3: Run migrations (ensure foreign keys are created correctly)
	err = db.AutoMigrate(&Workout{}, &Exercise{}, &PerformedSet{}, &WorkoutExercise{}, &Set{})
	if err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	return &TrainingDB{db: db}
}
