package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "swaggo-item-api/docs"

	"github.com/gorilla/mux"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	httpSwagger "github.com/swaggo/http-swagger"
)

//Order Struct
type Order struct {
	OrderID      uint      `json:"orderId" gorm:"primary_key"`
	CustomerName string    `json:"customerName"`
	OrderAt      time.Time `json:"orderedAt"`
	Items        []Item    `json:"items" gorm:"foreginkey:OrderID"`
}

//Item Struct
type Item struct {
	//gorm Model
	LineItemID  uint   `json:"lineItemId" gorm:"primary_key"`
	ItemCode    string `json:"itemCode"`
	Description string `json:"description"`
	Quantity    uint   `json:"quantity"`
	OrderID     uint   `json:"-"`
}

var db *gorm.DB

func initDB() {
	var err error
	dataSourceName := "root:root@tcp(localhost:3306)/?parseTime=True"
	db, err = gorm.Open("mysql", dataSourceName)

	if err != nil {
		fmt.Println(err)
		panic("failed to connect to the database:try Again...")
	}

	//create the database
	db.Exec("CREATE DATABASE orders_db")
	db.Exec("USE orders_db")
	//Migartion to create the tables for Order and Item schema
	db.AutoMigrate(&Order{}, &Item{})
}

// @title items API
// @version 1.0
// @description This is a sample serice for managing orders
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email soberkoder@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8082
// @BasePath /

func main() {
	router := mux.NewRouter()
	//create
	router.HandleFunc("/orders", createOrder).Methods("POST")
	//Read
	router.HandleFunc("/orders/{orderId}", getOrder).Methods("GET")
	//Read-All
	router.HandleFunc("/orders", getOrders).Methods("GET")
	//Update
	router.HandleFunc("/orders/{orderId}", updateOrder).Methods("PUT")
	//delete
	router.HandleFunc("/orders/{orderId}", deleteOrder).Methods("DELETE")

	//initialze DB connection
	initDB()
	// Swagger
	router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	log.Fatal(http.ListenAndServe(":8082", router))

}

// CreateOrder godoc
// @Summary Create a new order
// @Description Create a new order with the input paylod
// @Tags orders
// @Accept  json
// @Produce  json
// @Param order body Order true "Create order"
// @Success 200 {object} Order
// @Router /orders [post]

func createOrder(w http.ResponseWriter, r *http.Request) {
	var order Order
	json.NewDecoder(r.Body).Decode(&order)
	//create new orders by inserting the records int orders and item table
	db.Create(&order)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)

}

// GetOrders godoc
// @Summary Get details of all orders
// @Description Get details of all orders
// @Tags orders
// @Accept  json
// @Produce  json
// @Success 200 {array} Order
// @Router /orders [get]
func getOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var orders []Order
	db.Preload("Items").Find(&orders)
	json.NewEncoder(w).Encode(orders)
}

func getOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	inputOrderID := params["orderId"]

	var order Order
	db.Preload("Items").First(&order, inputOrderID)
	json.NewEncoder(w).Encode(order)
}

func updateOrder(w http.ResponseWriter, r *http.Request) {
	var updatedOrder Order
	json.NewDecoder(r.Body).Decode(&updatedOrder)
	db.Save(&updatedOrder)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedOrder)
}

func deleteOrder(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	inputOrderID := params["orderId"]

	//converting 'orderId' String param to unit64
	id64, _ := strconv.ParseUint(inputOrderID, 10, 64)
	//convert uint64  to int
	idToDelete := uint(id64)

	db.Where("order_id=?", idToDelete).Delete(&Item{})
	db.Where("order_id=?", idToDelete).Delete(&Order{})
	w.WriteHeader(http.StatusNoContent)
}
