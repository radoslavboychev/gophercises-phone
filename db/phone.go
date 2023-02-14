package db

import "database/sql"

/////////////////////////
///// EXPORTED TYPES ///
////////////////////////

// Phone represents a phone number object
type Phone struct {
	ID     int
	Number string
}

// DB represents a database object
type DB struct {
	db *sql.DB
}

// Close shuts down a database connection
func (db *DB) Close() error {
	return db.db.Close()
}

/////////////////////////
/// SET UP FUNCTIONS ///
////////////////////////

// Open
func Open(driverName, dataSource string) (*DB, error) {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

// Migrate
func Migrate(driverName, datasource string) error {
	db, err := sql.Open(driverName, datasource)
	if err != nil {
		return err
	}

	err = createPhoneNumbersTable(db)
	if err != nil {
		return err
	}
	return db.Close()
}

// Reset
func Reset(driverName, datasource, dbName string) error {
	db, err := sql.Open(driverName, datasource)
	if err != nil {
		return err
	}

	err = resetDB(db, dbName)
	if err != nil {
		return err
	}

	return db.Close()

}

// Seed inserts random values into the numbers table
func (db *DB) Seed() error {
	data := []string{
		"1234567890",
		"123 456 7891",
		"(123) 456 7892",
		"(123) 456-7893",
		"123-456-7894",
		"123-456-7890",
		"1234567892",
		"(123)456-7892",
	}
	for _, number := range data {
		if _, err := insertPhone(db.db, number); err != nil {
			return err
		}
	}
	return nil
}

/////////////////////////
//// CRUD FUNCTIONS ////
////////////////////////

// AllPhones returns all records for phones from the database
func (db *DB) AllPhones() ([]Phone, error) {
	rows, err := db.db.Query("SELECT id, value FROM phone_numbers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ret []Phone
	for rows.Next() {
		var p Phone
		if err := rows.Scan(&p.ID, &p.Number); err != nil {
			return nil, err
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return ret, nil
}

// FindPhones returns a single phone record based on the number searched for
func (db *DB) FindPhone(number string) (*Phone, error) {
	var p Phone
	row := db.db.QueryRow("SELECT * FROM phone_numbers WHERE value=$1", number)
	err := row.Scan(&p.ID, &p.Number)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &p, nil
}

// UpdatePhone
func (db *DB) UpdatePhone(p *Phone) error {
	statement := `UPDATE phone_numbers SET value=$2 WHERE id=$1`
	_, err := db.db.Exec(statement, p.ID, p.Number)
	if err != nil {
		return err
	}
	return nil
}

// DeletePhone
func (db *DB) DeletePhone(id int) error {
	statement := `DELETE FROM phone_numbers WHERE id=$1`
	_, err := db.db.Exec(statement, id)
	if err != nil {
		return err
	}
	return nil
}

/////////////////////////
// IMPORTED FUNCTIONS //
////////////////////////

// return a fan from ID
func getPhone(db *sql.DB, id int) (string, error) {
	var number string
	row := db.QueryRow("SELECT * FROM phone_numbers WHERE id =$1", id)
	err := row.Scan(&id, &number)
	if err != nil {
		return "", err
	}
	return number, nil
}

// creates a new SQL DB with the provided name
func createDB(db *sql.DB, name string) error {
	_, err := db.Exec("CREATE DATABASE " + name)
	if err != nil {
		return err
	}
	return nil
}

// deletes a DB with the provided name if it exists
func resetDB(db *sql.DB, name string) error {
	_, err := db.Exec("DROP DATABASE IF EXISTS " + name)
	if err != nil {
		return err
	}
	return nil
}

// creates the phone numbers table in the provided DB
func createPhoneNumbersTable(db *sql.DB) error {
	statement := `CREATE TABLE IF NOT EXISTS phone_numbers (
		id SERIAL,
		value VARCHAR(255)
	)`

	_, err := db.Exec(statement)
	if err != nil {
		return err
	}

	return nil
}

// inserts a phone number record into the DB
func insertPhone(db *sql.DB, phone string) (int, error) {
	statement := `INSERT INTO phone_numbers(value) VALUES($1) RETURNING id`
	var id int
	err := db.QueryRow(statement, phone).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}
