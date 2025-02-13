package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// define Measurement point here bc we dont have that many models
type Measurement struct {
	ID int64 `json:"id"`
	//Quantity  string  `json:"quantity"`
	//Timestamp string  `json:"timestamp"`
	Value float64 `json:"value"`
}

type Database struct {
	dbConn *sql.DB
}

func InitDB() (*Database, error) {
	//Open Connection and Create BasicTable
	connection, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Println("Error opening the database: ", err)
		return nil, err
	}
	db := &Database{dbConn: connection}
	err = createBasicTable(db)
	if err != nil {
		log.Println("Error creating the basic tables: ", err)
		return nil, err
	}
	return db, nil
}

func createBasicTable(db *Database) error {
	createTableSQL := `CREATE TABLE IF NOT EXISTS messpunkte (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        wert REAL
    );`

	_, err := db.dbConn.Exec(createTableSQL)
	if err != nil {
		log.Println("Error creating table: ", err)
		return err
	}
	return err
}

func (d *Database) InsertPoint(m *Measurement) error {
	insertSQL := `INSERT INTO messpunkte (wert) VALUES (?);`
	result, err := d.dbConn.Exec(insertSQL, m.Value)
	if err != nil {
		log.Println("Error inserting point: ", err)
		return err
	}
	// Retrieve the last inserted ID
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		log.Println("Error retrieving last insert ID: ", err)
		return err
	}
	m.ID = lastInsertId //change the ID to the asserted one
	return err
}

func (d *Database) GetAllPoints() ([]Measurement, error) {
	rows, err := d.dbConn.Query(`SELECT id, wert FROM messpunkte;`)
	if err != nil {
		log.Println("Error getting all points: ", err)
		return nil, err
	}
	defer rows.Close()

	var points []Measurement
	for rows.Next() {
		var point Measurement
		if err := rows.Scan(&point.ID, &point.Value); err != nil {
			log.Println("Error scanning row: ", err)
			return nil, err
		}
		points = append(points, point) //Append points to the Measurement slice
	}

	// Check for errors from the row iteration
	if err := rows.Err(); err != nil {
		log.Println("Error during row iteration: ", err)
		return nil, err
	}

	return points, nil
}

func (d *Database) Close() error {
	if d.dbConn != nil {
		return d.dbConn.Close()
	}
	return nil
}
