// methods for the measurement table
package database

import (
	"fmt"
	"log"
	"reflect"
)

func (d *Database) InsertMeasurement(m *Measurement) error {
	insertSQL := `INSERT INTO measurements (
		sensors_id,
		value,
		unit) VALUES (?, ?, ?);`
	result, err := d.dbConn.Exec(insertSQL, m.SensorsId, m.Value, m.Unit)
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
	return nil
}

func (d *Database) GetAllMeasurements() ([]Measurement, error) {
	rows, err := d.dbConn.Query(`SELECT * FROM measurements;`)
	if err != nil {
		log.Println("Error getting all points: ", err)
		return nil, err
	}
	defer rows.Close()

	var points []Measurement
	for rows.Next() {
		var p Measurement
		if err := rows.Scan(&p.ID, &p.SensorsId, &p.Value, &p.Unit, &p.Timestamp); err != nil {
			log.Println("Error scanning row: ", err)
			return nil, err
		}
		points = append(points, p) //Append points to the Measurement slice
	}

	// Check for errors from the row iteration
	if err := rows.Err(); err != nil {
		log.Println("Error during row iteration: ", err)
		return nil, err
	}

	return points, nil
}

func (d *Database) GetMeasurementById(queryId int) (*Measurement, error) {
	row := d.dbConn.QueryRow(`SELECT * FROM measurements WHERE id = ? LIMIT 1;`, queryId)
	p := &Measurement{}

	if err := row.Scan(&p.ID, &p.SensorsId, &p.Value, &p.Unit, &p.Timestamp); err != nil {
		log.Println("Error getting single row: ", err)
		return nil, err
	}
	return p, nil
}

func (d *Database) DeleteMeasurement(id int) error {
	_, err := d.dbConn.Exec(`DELETE FROM measurements WHERE id = ?;`, id)
	if err != nil {
		return fmt.Errorf("error deleting measurement(id=%v): %w", id, err)
	}
	return nil
}

// this should fix the problem with updating not supported types, but not finished
func (d *Database) measurementToMap(updateData Measurement) map[string]any {
	var updateDataMap map[string]any
	val := reflect.ValueOf(updateData)
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if !val.IsZero() {
			// Get the field name
			fieldName := val.Type().Field(i).Name
			// Add to the map
			updateDataMap[fieldName] = field.Interface()
		}
	}
	return updateDataMap
}

/*
curl -X PUT http://localhost:8080/measurement/1 \
     -H "Content-Type: application/json" \
     -d '{"unit": "volt"}'
*/

func (d *Database) UpdateMeasurement(id int, updateData map[string]any) error {
	//should check with types but not now

	// Build SQL query dynamically
	query := "UPDATE measurements SET "
	args := make([]any, 0, len(updateData)+1) //max cap is one more than updateData (+id)
	i := 0
	for key, value := range updateData {
		if i > 0 {
			//no comma in first iteration
			query += ", "
		}
		query += key + " = ?"
		args = append(args, value)
		i++
	}
	query += " WHERE id = ?"
	fmt.Printf("final update query: %s", query)

	args = append(args, id)
	// Execute the query
	res, err := d.dbConn.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update: %w", err)
	}

	// Check if the row exists
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("record not found")
	}

	return nil
}
