package populator

import (
	"fmt"
	"git.garena.com/frieda.hasanah/user_service/internal/model"
	"git.garena.com/frieda.hasanah/user_service/utils/hash"
	"git.garena.com/frieda.hasanah/user_service/utils/log"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"runtime"
	"sync"
	"time"
)

const DefaultPassword = "password123"

// PopulateUsers generates and inserts the specified number of users into the database
func PopulateUsers(db *sqlx.DB, numUsers int, batchSize int) {
	start := time.Now()
	numWorkers := runtime.NumCPU() * 2
	userBatches := make(chan []model.User, numWorkers*2)
	var wg sync.WaitGroup
	hashedPassword, _ := hash.EncryptPassword(DefaultPassword)

	startWorkers(db, userBatches, numWorkers, &wg)
	processBatches(numUsers, batchSize, userBatches, hashedPassword)

	wg.Wait()
	fmt.Printf("User population completed in %v\n", time.Since(start))
}

func startWorkers(db *sqlx.DB, userBatches chan []model.User, numWorkers int, wg *sync.WaitGroup) {
	fmt.Println("Starting workers")
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			fmt.Println("Starting worker", w)
			defer wg.Done()
			for users := range userBatches {
				if err := insertUsers(db, users); err != nil {
					log.Println("Failed to insert users:", err)
				} else {
					fmt.Printf("Successfully inserted %d users\n", len(users))
				}
			}
		}()
	}
}

func processBatches(numUsers, batchSize int, userBatches chan []model.User, defaultHashedPassword string) {
	var batchWg sync.WaitGroup
	for i := 0; i < numUsers; i += batchSize {
		end := i + batchSize
		if end > numUsers {
			end = numUsers
		}
		batchWg.Add(1)
		go func(start, end int) {
			//fmt.Printf("Generating users %d-%d\n", start, end)
			defer batchWg.Done()
			batch := generateUsers(start, end, defaultHashedPassword)
			userBatches <- batch
		}(i, end)
	}
	batchWg.Wait()
	close(userBatches)
}

func generateUsers(start, end int, defaultHashedPassword string) []model.User {
	users := make([]model.User, end-start)
	for i := range users {
		users[i] = generateUser(defaultHashedPassword)
	}
	return users
}

func generateUser(defaultHashedPassword string) model.User {
	id := uuid.New()
	identifier := id.String()
	fullName := fmt.Sprintf("fullname_%s", identifier)
	email := fmt.Sprintf("email_%s@garena.com", identifier)
	username := fmt.Sprintf("u_%s", identifier)
	//password, _ := hash.EncryptPassword(DefaultPassword)
	password := defaultHashedPassword

	return model.User{
		ID:       id,
		FullName: fullName,
		Email:    email,
		Username: username,
		Password: password,
	}
}

// Insert users into the database within a transaction
func insertUsers(db *sqlx.DB, users []model.User) error {
	tx, err := db.Begin() // Start transaction
	if err != nil {
		return err
	}

	query, values := prepareInsertQuery(users)

	stmt, err := tx.Prepare(query)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = stmt.Exec(values...)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	err = tx.Commit()
	return err
}

func prepareInsertQuery(users []model.User) (string, []interface{}) {
	query := "INSERT INTO users (id, full_name, email, username, password) VALUES "
	values := []interface{}{}
	for _, user := range users {
		query += "(?, ?, ?, ?, ?),"
		values = append(values, user.ID, user.FullName, user.Email, user.Username, user.Password)
	}
	query = query[:len(query)-1] // Trim the last comma
	return query, values
}
