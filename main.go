package main

import (
	"fmt"
	"log"
	"golang_practical_task/db"
	"golang_practical_task/controllers"
	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
)

func main() {
	// Initializing database connection.
	dbConn, err := db.InitDB()
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}

	defer dbConn.Close()

	// Creating required tables.
	_, err = dbConn.Exec(`CREATE TABLE IF NOT EXISTS person (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255),
		age INT
	)`)
	if err != nil {
		log.Fatal("Error creating person table: ", err)
	}

	_, err = dbConn.Exec(`CREATE TABLE IF NOT EXISTS phone (
		id SERIAL PRIMARY KEY,
		number VARCHAR(255),
		person_id INT
	)`)
	if err != nil {
		log.Fatal("Error creating phone table: ", err)
	}

	_, err = dbConn.Exec(`CREATE TABLE IF NOT EXISTS address (
		id SERIAL PRIMARY KEY,
		city VARCHAR(255),
		state VARCHAR(255),
		street1 VARCHAR(255),
		street2 VARCHAR(255),
		zip_code VARCHAR(255)
	)`)
	if err != nil {
		log.Fatal("Error creating address table: ", err)
	}

	_, err = dbConn.Exec(`CREATE TABLE IF NOT EXISTS address_join (
		id SERIAL PRIMARY KEY,
		person_id INT,
		address_id INT
	)`)
	if err != nil {
		log.Fatal("Error creating address_join table: ", err)
	}

	// Checking whether table is empty or not.
	var count int
	err = dbConn.QueryRow("SELECT COUNT(*) FROM person").Scan(&count)
	if err != nil {
		log.Fatal("Error checking person table: ", err)
	}
	if count == 0 {
		// Since the table is empty, inserting sample data into tables.
		_, err = dbConn.Exec(`INSERT INTO person (name, age) VALUES
			('mike', 31),
			('John', 20),
			('Joseph', 20)
		`)
		if err != nil {
			log.Fatal("Error inserting sample data into person table: ", err)
		}

		_, err = dbConn.Exec(`INSERT INTO phone (person_id, number) VALUES
			(1, '444-444-4444'),
			(2, '123-444-7777'),
			(3, '445-222-1234')
		`)
		if err != nil {
			log.Fatal("Error inserting sample data into phone table: ", err)
		}

		_, err = dbConn.Exec(`INSERT INTO address (city, state, street1, street2, zip_code) VALUES
			('Eugene', 'OR', '111 Main St', '', '98765'),
			('Sacramento', 'CA', '432 First St', 'Apt 1', '22221'),
			('Austin', 'TX', '213 South 1st St', '', '78704')
		`)
		if err != nil {
			log.Fatal("Error inserting sample data into address table: ", err)
		}

		_, err = dbConn.Exec(`INSERT INTO address_join (person_id, address_id) VALUES
			(1, 3),
			(2, 1),
			(3, 2)
		`)
		if err != nil {
			log.Fatal("Error inserting sample data into address_join table: ", err)
		}

		fmt.Println("Tables created and sample data inserted successfully!")
	}

    r := gin.Default()

    r.GET("/person/:person_id/info", controllers.PersonGETHandler)
	r.POST("/person/create", controllers.PersonPOSTHandler)

    // Running server on 9000 port.
    r.Run(":9000")
}
