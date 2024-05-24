package v1

import (
	"EffectiveMobile/pkg/postgres"
	"EffectiveMobile/pkg/postgres/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
)

type Router struct {
	log *zap.Logger
	db  postgres.Service
}

func newAutoRoutes(handler *gin.Engine, log *zap.Logger, db postgres.Service) {
	r := &Router{
		log: log,
		db:  db,
	}

	handler.GET("/info", r.Info)

	handler.GET("/cars", r.GetCars)
	handler.POST("/cars", r.AddCars)
	handler.PATCH("/cars/:id", r.UpdateCar)
	handler.DELETE("/cars/:id", r.DeleteCarByRegNum)

	handler.POST("/people", r.CreatePeople)
}

// Info ShowAccount
// @Summary      Get car info
// @Description  Get cars info by regNum
// @Tags         cars
// @Accept       json
// @Produce      json
// @Param regNum query string true "Registration number" example(X123XX150)
// @Success      200 {object} models.Car
// @Failure      400 {object} string "regNum is required"
// @Failure      404 {object} string "car not found"
// @Router       /info [get]
func (r *Router) Info(c *gin.Context) {
	reqNum := c.Query("regNum")
	if reqNum == "" {
		c.String(400, "regNum is required")
		return
	}
	car, err := r.db.GetCar(reqNum)
	if err != nil {
		c.String(404, "car not found")
		return
	}

	c.JSON(200, &car)
}

// GetCars
// @Summary Get cars with filter and pagination
// @Description get cars by filter
// @Tags cars
// @Accept  json
// @Produce  json
// @Param page query int false "page number"
// @Param pageSize query int false "items per page"
// @Param regNum query string false "registration number"
// @Param mark query string false "mark"
// @Param model query string false "model"
// @Param year query int false "year"
// @Param createdAt query string false "created at"
// @Param updatedAt query string false "updated at"
// @Param deletedAt query string false "deleted at"
// @Success 200 {array} models.Car
// @Router /cars [get]
func (r *Router) GetCars(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	regNum := c.Query("regNum")
	mark := c.Query("mark")
	model := c.Query("model")
	year, _ := strconv.Atoi(c.Query("year"))
	createdAt, _ := time.Parse(time.RFC3339, c.Query("createdAt"))
	updatedAt, _ := time.Parse(time.RFC3339, c.Query("updatedAt"))
	deletedAt, _ := time.Parse(time.RFC3339, c.Query("deletedAt"))

	filter := models.Car{
		RegNum:    regNum,
		Mark:      mark,
		Model:     model,
		Year:      year,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		DeletedAt: &deletedAt,
	}

	cars, err := r.db.GetCars(filter, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cars)
}

// DeleteCarByRegNum
// @Summary Delete a car
// @Description delete car by registration number
// @Tags cars
// @Accept  json
// @Produce  json
// @Param regNum path string true "Registration number"
// @Success 204 {object} string "Successfully deleted"
// @Router /cars/{regNum} [delete]
func (r *Router) DeleteCarByRegNum(c *gin.Context) {
	regNum := c.Param("regNum")

	err := r.db.DeleteCarByRegNum(regNum)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "Successfully deleted"})
}

// UpdateCar
// @Summary Update a car
// @Description update car by ID
// @Tags cars
// @Accept  json
// @Produce  json
// @Param id path int true "Car ID"
// @Param car body models.Car true "Car object"
// @Success 200 {object} models.Car
// @Router /cars/{id} [patch]
func (r *Router) UpdateCar(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var updatedCar models.Car
	if err := c.ShouldBindJSON(&updatedCar); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := r.db.UpdateCar(uint(id), updatedCar)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedCar)
}

type PeopleRequest struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
}

// CreatePeople
// @Summary Create a new people
// @Description create new people
// @Tags people
// @Accept  json
// @Produce  json
// @Param people body PeopleRequest true "People object"
// @Success 201 {object} models.People
// @Router /people [post]
func (r *Router) CreatePeople(c *gin.Context) {
	var request PeopleRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	people, err := r.db.CreatePeople(request.Name, request.Surname)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, people)
}

type CarRequest struct {
	RegNums  []string `json:"regNums"`
	PeopleID uint     `json:"peopleId"`
}

// AddCars
// @Summary Add new cars
// @Description add new cars by registration numbers
// @Tags cars
// @Accept  json
// @Produce  json
// @Param cars body CarRequest true "Register numbers of cars to add"
// @Success 201 {object} string "Successfully added"
// @Router /cars [post]
func (r *Router) AddCars(c *gin.Context) {
	var request CarRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := r.db.AddCars(request.RegNums, request.PeopleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Successfully added"})
}
