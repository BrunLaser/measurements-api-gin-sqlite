// General databse actions like connection opening and initialising
package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	dbConn *sql.DB
}

// Measurement point here bc we dont have that many models

type Measurement struct {
	ID        int64   `json:"id"`
	SensorsId int64   `json:"sensor_id"`
	Value     float64 `json:"value"`
	Unit      string  `json:"unit"`
	Timestamp string  `json:"timestamp"`
}

func (db *Database) withTransaction(fn func(transaction *sql.Tx) error) error {
	//Initialise connection
	tx, err := db.dbConn.Begin()
	if err != nil {
		log.Println("Failed to start transaction:", err)
		return err
	}
	defer tx.Rollback() //Rollback if we dont reach commit

	if err := fn(tx); err != nil {
		log.Println("Transaction failed", err)
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Println("Transaction commit failed", err)
		return err
	}

	return nil
}

func InitDB() (*Database, error) {
	//Open Connection and Create BasicTable
	connection, err := sql.Open("sqlite3", "./experiments.db")
	if err != nil {
		log.Println("Error opening the database: ", err)
		return nil, err
	}

	db := &Database{dbConn: connection}

	//Create tables with transaction
	err = db.withTransaction(db.createTables)
	if err != nil {
		return nil, fmt.Errorf("error creating tables: %w", err)
	}

	//Initialise tables with transaction
	err = db.withTransaction(db.initTables)
	if err != nil {
		return nil, fmt.Errorf("error initialising tables: %w", err)
	}

	return db, nil
}

func (db *Database) createTables(tx *sql.Tx) error {
	//experiments table

	tableStatements := []string{
		`CREATE TABLE IF NOT EXISTS experiments (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT,
            description TEXT,
            date DATE DEFAULT CURRENT_DATE
        );`,
		`CREATE TABLE IF NOT EXISTS sensors (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            experiment_id INTEGER,
            sensor_type TEXT,
            FOREIGN KEY (experiment_id) REFERENCES experiments(id)
        );`,
		`CREATE TABLE IF NOT EXISTS measurements (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            sensors_id INTEGER,
			value REAL,
			unit TEXT,
            timestamp TEXT DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (sensors_id) REFERENCES sensors(id)
        );`,
		`PRAGMA foreign_keys = ON;`, // Ensure foreign keys are enforced
	}

	for _, stmt := range tableStatements {
		if _, err := tx.Exec(stmt); err != nil {
			return fmt.Errorf("error executing statement: %s, error: %w", stmt, err)
		}
	}

	return nil
}

// Inserts 2 experiments and 2 sensors each, no measurements inserted
func (db *Database) initTables(tx *sql.Tx) error {
	//Creates two experiments, one yesterday, one today.
	createExpSQL := `INSERT OR IGNORE INTO experiments (id, name, description, date)
				     VAlUES (1, ?, ?, ?);`
	yesterday := time.Now().AddDate(0, 0, -1)
	_, err := tx.Exec(createExpSQL, "Exp1", "the first experiment", yesterday.Format("2006-01-02"))
	if err != nil {
		log.Println("Couldnt create first Experiment: ", err)
		return err
	}

	createExpSQL = `INSERT OR IGNORE INTO experiments (id, name, description, date)
				     VAlUES (2, ?, ?, ?);`
	_, err = tx.Exec(createExpSQL, "Exp2", "the second experiment", time.Now().Format("2006-01-02"))
	if err != nil {
		log.Println("Couldnt create second Experiment: ", err)
		return err
	}

	//create two sensors for each experiment
	createSensorSQL := `INSERT OR IGNORE INTO sensors (id, experiment_id, sensor_type) 
						VALUES (?, ?, "elektrisch")`
	for i := 1; i <= 2; i++ {
		for j := 1; j <= 2; j++ {
			_, err := tx.Exec(createSensorSQL, j, i)
			if err != nil {
				log.Println("Couldnt create sensor: ", err)
				return err
			}
		}
	}
	return nil
}

func (d *Database) Close() error {
	if d.dbConn != nil {
		return d.dbConn.Close()
	}
	return nil
}
