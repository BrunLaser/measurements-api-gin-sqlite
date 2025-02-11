package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// define Measurement point here bc we dont have that many models
type Measurement struct {
	ID int `json:"id"`
	//Quantity  string  `json:"quantity"`
	//Timestamp string  `json:"timestamp"`
	Value float64 `json:"value"`
}

type Database struct {
	db *sql.DB
}

// Open the data base
func (d *Database) Open() error {
	var err error
	d.db, err = sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Println("Error opening the database: ", err)
	}
	return err
}

func (d *Database) CreateBasicTable() error {
	createTableSQL := `CREATE TABLE IF NOT EXISTS messpunkte (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        wert REAL
    );`

	_, err := d.db.Exec(createTableSQL)
	if err != nil {
		log.Println("Error creating table: ", err)
	}
	return err
}

func (d *Database) InsertPoint(m Measurement) error {
	insertSQL := `INSERT INTO messpunkte (wert) VALUES (?);`
	_, err := d.db.Exec(insertSQL, m.Value)
	if err != nil {
		log.Println("Error inserting point: ", err)
	}
	return err
}

func (d *Database) GetAllPoints() ([]Measurement, error) {
	rows, err := d.db.Query(`SELECT id, wert FROM messpunkte;`)
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
	if d.db != nil {
		return d.db.Close()
	}
	return nil
}
