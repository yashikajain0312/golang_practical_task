package controllers

import (
	"database/sql"
	"log"
	"net/http"
	"golang_practical_task/db"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type PersonRequest struct {
	PersonID int `uri:"person_id" binding:"required"`
}

type PersonInfoResponse struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	City        string `json:"city"`
	State       string `json:"state"`
	Street1     string `json:"street1"`
	Street2     string `json:"street2"`
	ZipCode     string `json:"zip_code"`
}

func PersonGETHandler(c *gin.Context) {
	var request PersonRequest

	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request!"})
		return
	}

	// Initialize database connection
	dbConn, err := db.InitDB()
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}

	defer dbConn.Close()

	row := dbConn.QueryRow("SELECT p.name, ph.number, a.city, a.state, a.street1, a.street2, a.zip_code FROM person p JOIN phone ph ON p.id = ph.person_id JOIN address_join aj ON p.id = aj.person_id JOIN address a ON aj.address_id = a.id WHERE p.id = $1", request.PersonID)

	var resp PersonInfoResponse

	err = row.Scan(&resp.Name, &resp.PhoneNumber, &resp.City, &resp.State, &resp.Street1, &resp.Street2, &resp.ZipCode)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Person not found"})
			return
		}
		log.Fatal("Error scanning row: ", err)
	}

	c.JSON(http.StatusOK, resp)
}
