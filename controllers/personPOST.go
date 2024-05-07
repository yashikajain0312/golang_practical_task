package controllers

import (
	"log"
	"net/http"
	"golang_practical_task/db"

	"github.com/gin-gonic/gin"
)

type PersonCreateRequest struct {
	Name        string `json:"name" binding:"required"`
	Age         int    `json:"age" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	City        string `json:"city" binding:"required"`
	State       string `json:"state" binding:"required"`
	Street1     string `json:"street1" binding:"required"`
	Street2     string `json:"street2"`
	ZipCode     string `json:"zip_code" binding:"required"`
}

func PersonPOSTHandler(c *gin.Context) {
	var request PersonCreateRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dbConn, err := db.InitDB()
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}

	defer dbConn.Close()

	// Inserting data into tables.
	personRow := dbConn.QueryRow("INSERT INTO person (name, age) VALUES ($1, $2) RETURNING id", request.Name, request.Age)
	var personID int
	err = personRow.Scan(&personID)
	if err != nil {
		log.Fatal("Error inserting data into person table: ", err)
	}

	_, err = dbConn.Exec("INSERT INTO phone (person_id, number) VALUES ($1, $2)", personID, request.PhoneNumber)
	if err != nil {
		log.Fatal("Error inserting data into phone table: ", err)
	}

	_, err = dbConn.Exec("INSERT INTO address (city, state, street1, street2, zip_code) VALUES ($1, $2, $3, $4, $5)", request.City, request.State, request.Street1, request.Street2, request.ZipCode)
	if err != nil {
		log.Fatal("Error inserting data into address table: ", err)
	}

	// Retrieving the ID of the inserted address record.
	var addressID int
	
	row := dbConn.QueryRow("SELECT id FROM address WHERE city = $1 AND state = $2 AND street1 = $3 AND street2 = $4 AND zip_code = $5", request.City, request.State, request.Street1, request.Street2, request.ZipCode)
	err = row.Scan(&addressID)
	if err != nil {
		log.Fatal("Error retrieving address ID: ", err)
	}

	_, err = dbConn.Exec("INSERT INTO address_join (person_id, address_id) VALUES ($1, $2)", personID, addressID)
	if err != nil {
		log.Fatal("Error inserting data into address_join table: ", err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Person created successfully!"})
}
