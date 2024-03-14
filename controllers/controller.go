package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	// "strconv"

	m "Week3/models"
	// "github.com/gorilla/mux"
)

func Login(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		sendErrorResponse(w, "Failed Parse")
		return
	}

	platform := r.Header.Get("platform")
	if platform == "" {
		sendErrorResponse(w, "Failed Platform")
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	var user m.User

	errQuery := db.QueryRow("SELECT * FROM users WHERE email=? AND password=?",
		email,
		password,
	).Scan(&user.ID, &user.Name, &user.Age, &user.Address)

	var response m.UserResponse
	if errQuery == nil {
		response.Status = 200
		response.Message = "Success login from " + platform
		response.Data = user
	} else {
		response.Status = 400
		response.Message = "Login failed"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetAllUsers...
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	query := "SELECT * FROM users"

	// Read from Header
	// name := r.Header.Get("Name")
	// if name != "" {
	// 	query := " WHERE name='" + name + "'"
	// }

	// Read from Query Param
	name := r.URL.Query()["name"]
	age := r.URL.Query()["age"]
	if name != nil {
		fmt.Println(name[0])
		// fmt.Println(name[1])
		query += " WHERE name='" + name[0] + "'"
	}

	if age != nil {
		if name[0] != "" {
			query += " AND"
		} else {
			query += " WHERE"
		}
		query += " age='" + age[0] + "'"
	}

	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		// send error response
		return
	}
	var user m.User
	var users []m.User
	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Name, &user.Age, &user.Address); err != nil {
			log.Println(err)
			// send error response
			return
		} else {
			users = append(users, user)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	// var response UsersResponse
	// if len(users) < 5 {
	var response m.UsersResponse
	response.Status = 200
	response.Message = "Success"
	response.Data = users
	json.NewEncoder(w).Encode(response)
	// } else {
	// 	var response ErrorResponse
	// 	response.Status = 400
	// 	response.Message = "Error Array Size Not Correct"
	// 	json.NewEncoder(w).Encode(response)
	// }
}

func GetAllUsersGorm(w http.ResponseWriter, r *http.Request) {
	db := connectGorm()

	var users []m.User
	result := db.Find(&users)
	if result.Error != nil {
		sendErrorResponse(w, "Find Error")
		return
	}

	var response m.UsersResponse
	response.Status = 200
	response.Message = "Success"
	response.Data = users
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetAllProduct
func GetAllProduct(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	query := "SELECT * FROM products"

	// Read from Query Param
	name := r.URL.Query()["name"]
	price := r.URL.Query()["price"]
	if name != nil {
		fmt.Println(name[0])
		// fmt.Println(name[1])
		query += " WHERE name='" + name[0] + "'"
	}

	if price != nil {
		if name[0] != "" {
			query += " AND"
		} else {
			query += " WHERE"
		}
		query += " price='" + price[0] + "'"
	}

	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		// send error response
		var response m.ErrorResponse
		response.Status = 400
		response.Message = "Error Get Products"
		json.NewEncoder(w).Encode(response)
		return
	}
	var product m.Products
	var products []m.Products
	for rows.Next() {
		if err := rows.Scan(&product.ID, &product.Name, &product.Price); err != nil {
			log.Println(err)

			var response m.ErrorResponse
			response.Status = 400
			response.Message = "Error Get Products"
			json.NewEncoder(w).Encode(response)
			return
		} else {
			products = append(products, product)
		}
	}
	w.Header().Set("Content-Type", "application/json")

	var response m.ProductsResponse
	response.Status = 200
	response.Message = "Success"
	response.Data = products
	json.NewEncoder(w).Encode(response)
}

// GetAllTransactions
func GetAllTransactions(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	query := "SELECT * FROM transactions"

	// Read from Query Param
	userID := r.URL.Query()["userID"]
	productID := r.URL.Query()["productID"]
	if userID != nil {
		fmt.Println(userID[0])
		// fmt.Println(name[1])
		query += " WHERE userID='" + userID[0] + "'"
	}

	if productID != nil {
		if userID[0] != "" {
			query += " AND"
		} else {
			query += " WHERE"
		}
		query += " productID='" + productID[0] + "'"
	}

	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)

		var response m.ErrorResponse
		response.Status = 400
		response.Message = "Error Get Transactions"
		json.NewEncoder(w).Encode(response)
		return
	}
	var transaction m.Transactions
	var transactions []m.Transactions
	for rows.Next() {
		if err := rows.Scan(&transaction.ID, &transaction.UserID, &transaction.ProductID, &transaction.Quantity); err != nil {
			log.Println(err)

			var response m.ErrorResponse
			response.Status = 400
			response.Message = "Error Get Transactions"
			json.NewEncoder(w).Encode(response)
			return
		} else {
			transactions = append(transactions, transaction)
		}
	}
	w.Header().Set("Content-Type", "application/json")

	var response m.TransactionsResponse
	response.Status = 200
	response.Message = "Success"
	response.Data = transactions
	json.NewEncoder(w).Encode(response)
}

func GetDetailUserTransaction(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	params := r.URL.Query()

	userIDStr := params.Get("userID")

	var query string
	var args []interface{}

	if userIDStr != "" {
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			log.Println(err)

			var response m.ErrorResponse
			response.Status = 400
			response.Message = "Error User ID"
			json.NewEncoder(w).Encode(response)
			return
		}
		query = `SELECT t.id, t.quantity, u.id, u.name, u.age, u.address, p.id, p.name, p.price
            FROM transactions t 
            JOIN users u ON t.userid = u.id
            JOIN products p ON t.productid = p.id
            WHERE u.id = ?`
		args = append(args, userID)

	} else {
		query = `SELECT t.id, t.quantity, u.id, u.name, u.age, u.address, p.id, p.name, p.price
            FROM transactions t 
            JOIN users u ON t.userid = u.id
            JOIN products p ON t.productid = p.id`
	}

	detailTransactionRow, err := db.Query(query, args...)
	if err != nil {
		log.Println(err)

		var response m.ErrorResponse
		response.Status = 400
		response.Message = "Invalid query"
		json.NewEncoder(w).Encode(response)
		return
	}

	var detailTransaction m.DetailTransaction
	var detailTransactions []m.DetailTransaction
	for detailTransactionRow.Next() {
		if err := detailTransactionRow.Scan(&detailTransaction.ID, &detailTransaction.Quantity, &detailTransaction.User.ID,
			&detailTransaction.User.Name, &detailTransaction.User.Age, &detailTransaction.User.Address,
			&detailTransaction.Product.ID, &detailTransaction.Product.Name, &detailTransaction.Product.Price); err != nil {

			log.Println(err)
			var response m.ErrorResponse
			response.Status = 400
			response.Message = "Error Get Detail Transaction"
			json.NewEncoder(w).Encode(response)
			return
		} else {
			detailTransactions = append(detailTransactions, detailTransaction)
		}
	}

	var response m.DetailTransactionsResponse
	response.Status = 200
	response.Message = "Success!"
	response.Data = detailTransactions
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetDetailUserTransactionGorm(w http.ResponseWriter, r *http.Request) {
	db := connectGorm()

	params := r.URL.Query()

	userIDStr := params.Get("userID")

	if userIDStr == "" {
		sendErrorResponse(w, "ID kosong")
		return
	}

	var queryBuilder = db.Table("transactions").
		Select("transactions.id, transactions.quantity, users.id as user_id, users.name, users.age, users.address, products.id as product_id, products.name as product_name, products.price").
		Joins("JOIN users ON transactions.userid = users.id").
		Joins("JOIN products ON transactions.productid = products.id")

	if userIDStr != "" {
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			log.Println(err)

			sendErrorResponse(w, "ID harus angka")
			return
		}
		queryBuilder = queryBuilder.Where("users.id = ?", userID)
	}

	var detailTransactions []m.DetailTransaction
	if err := queryBuilder.Preload("User").Preload("Product").Scan(&detailTransactions).Error; err != nil {
		log.Println(err)

		sendErrorResponse(w, "Query Error")
		return
	}

	var response m.DetailTransactionsResponse
	response.Status = 200
	response.Message = "Success!"
	response.Data = detailTransactions
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func InsertUser(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	// Read from request body
	err := r.ParseForm()
	if err != nil {
		return
	}

	name := r.Form.Get("name")
	age, _ := strconv.Atoi(r.Form.Get("age"))
	address := r.Form.Get("address")

	if name == "" || address == "" || age == 0 {
		sendErrorResponse(w, "Missing input data")
		return
	}

	_, errQuery := db.Exec("INSERT INTO users(name, age, address) values (?,?,?)",
		name,
		age,
		address,
	)

	var response m.UserResponse
	if errQuery == nil {
		response.Status = 200
		response.Message = "Success"
	} else {
		response.Status = 400
		response.Message = "Insert Failed !"
	}
	var user m.User
	user.Name = name
	user.Age = age
	user.Address = address
	response.Data = user

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func InsertUserGorm(w http.ResponseWriter, r *http.Request) {
	db := connectGorm()

	err := r.ParseForm()
	if err != nil {
		sendErrorResponse(w, "Failed Parse")
		return
	}

	name := r.Form.Get("name")
	age, _ := strconv.Atoi(r.Form.Get("age"))
	address := r.Form.Get("address")
	email := r.Form.Get("email")
	passwords := r.Form.Get("password")

	if name == "" || address == "" || age <= 0 || passwords == "" || email == "" {
		sendErrorResponse(w, "Missing input data")
		return
	}

	user := m.User{Name: name, Address: address, Age: age, Password: passwords, Email: email}
	result := db.Create(&user)

	if result.Error != nil {
		sendErrorResponse(w, "Error")
	} else {
		sendSuccessResponse(w, "Success")
	}
}

func InsertProduct(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	// Read from request body
	err := r.ParseForm()
	if err != nil {
		return
	}

	name := r.Form.Get("name")
	price, _ := strconv.Atoi(r.Form.Get("price"))

	if name == "" || price == 0 {
		sendErrorResponse(w, "Missing input data")
		return
	}

	_, errQuery := db.Exec("INSERT INTO products(name, price) values (?,?)",
		name,
		price,
	)

	var response m.ProductResponse
	if errQuery == nil {
		response.Status = 200
		response.Message = "Success"
	} else {
		response.Status = 400
		response.Message = "Insert Failed !"
	}
	var product m.Products
	product.Name = name
	product.Price = price
	response.Data = product

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func InsertTransaction(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	// Read from request body
	err := r.ParseForm()
	if err != nil {
		return
	}

	userId, _ := strconv.Atoi(r.Form.Get("userId"))
	productId, _ := strconv.Atoi(r.Form.Get("productId"))
	quantity, _ := strconv.Atoi(r.Form.Get("quantity"))

	if userId == 0 || productId == 0 || quantity == 0 {
		sendErrorResponse(w, "Missing input data")
		return
	}

	if !userExists(userId) {
		sendErrorResponse(w, "user not exists")
		return
	}
	if !productExists(productId) {
		_, err := db.Exec("INSERT INTO products(id, name) VALUES (?,?)",
			productId,
			"")

		if err != nil {
			sendErrorResponse(w, "product not exists")
			return
		}
	}

	result, errQuery := db.Exec("INSERT INTO transactions(userid, productid, quantity) VALUES (?,?,?)",
		userId,
		productId,
		quantity,
	)

	var response m.TransactionResponse
	if errQuery == nil {
		transactionId, _ := result.LastInsertId()
		transaction := m.Transactions{
			ID:        int(transactionId),
			UserID:    userId,
			ProductID: productId,
			Quantity:  quantity,
		}
		response.Status = 200
		response.Message = "Success!"
		response.Data = transaction
	} else {
		response.Status = 400
		response.Message = "Insert Failed!"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		sendErrorResponse(w, "Failed Parse")
		return
	}

	userID := r.Form.Get("id")
	if userID == "" {
		sendErrorResponse(w, "user id kosong")
		return
	}

	_, err = strconv.Atoi(userID)
	if err != nil {
		sendErrorResponse(w, "Failed to convert string")
		return
	}

	name := r.Form.Get("name")
	age, _ := strconv.Atoi(r.Form.Get("age"))
	address := r.Form.Get("address")

	// Buat Query sesuai dengan value yang diberikan
	query := "UPDATE users SET"
	params := []interface{}{}

	if name != "" {
		query += " name=?,"
		params = append(params, name)
	}

	if age != 0 {
		query += " age=?,"
		params = append(params, age)
	}

	if address != "" {
		query += " address=?,"
		params = append(params, address)
	}

	query = strings.TrimSuffix(query, ",")
	query += " WHERE id=?"
	params = append(params, userID)
	_, errQuery := db.Exec(query, params...)

	if errQuery == nil {
		var user m.User
		user.Name = name
		user.Age = age
		user.Address = address

		var response m.UserResponse
		response.Status = 200
		response.Message = "Success"
		response.Data = user
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		sendErrorResponse(w, "Query Error")
	}
}

func UpdateUserGorm(w http.ResponseWriter, r *http.Request) {
	db := connectGorm()

	err := r.ParseForm()
	if err != nil {
		sendErrorResponse(w, "Failed Parse")
		return
	}

	userID, _ := strconv.Atoi(r.Form.Get("id"))

	if userID == 0 {
		sendErrorResponse(w, "ID kosong")
		return
	}

	name := r.Form.Get("name")
	age, _ := strconv.Atoi(r.Form.Get("age"))
	address := r.Form.Get("address")
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	if name == "" || age <= 0 || address == "" || password == "" || email == "" {
		sendErrorResponse(w, "Missing Data")
		return
	}

	user := m.User{ID: userID, Name: name, Age: age, Address: address, Password: password, Email: email}
	result := db.Save(&user)

	if result.Error == nil {
		var user m.User
		user.ID = userID
		user.Name = name
		user.Age = age
		user.Address = address
		user.Email = email
		user.Password = password

		var response m.UserResponse
		response.Status = 200
		response.Message = "Success"
		response.Data = user
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		sendErrorResponse(w, "Query Error")
	}
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		sendErrorResponse(w, "Failed Parse")
		return
	}

	productId := r.Form.Get("id")
	if productId == "" {
		sendErrorResponse(w, "Product Id not defined")
		return
	}

	_, err = strconv.Atoi(productId)
	if err != nil {
		sendErrorResponse(w, "Failed to convert string")
		return
	}

	name := r.Form.Get("name")
	price, _ := strconv.Atoi(r.Form.Get("price"))

	// Buat Query sesuai dengan value yang diberikan
	query := "UPDATE products SET"
	params := []interface{}{}

	if name != "" {
		query += " name=?,"
		params = append(params, name)
	}

	if price != 0 {
		query += " price=?,"
		params = append(params, price)
	}

	query = strings.TrimSuffix(query, ",")
	query += " WHERE id=?"
	params = append(params, productId)
	_, errQuery := db.Exec(query, params...)

	if errQuery == nil {
		var product m.Products
		product.Name = name
		product.Price = price

		var response m.ProductResponse
		response.Status = 200
		response.Message = "Success"
		response.Data = product
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		sendErrorResponse(w, "Query Error")
	}
}

func UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		sendErrorResponse(w, "Failed Parse")
		return
	}

	transactionId := r.Form.Get("id")
	if transactionId == "" {
		sendErrorResponse(w, "id is not defined")
		return
	}

	_, err = strconv.Atoi(transactionId)
	if err != nil {
		sendErrorResponse(w, "failed to convert string")
		return
	}

	userId, _ := strconv.Atoi(r.Form.Get("userId"))
	productId, _ := strconv.Atoi(r.Form.Get("productId"))
	quantity, _ := strconv.Atoi(r.Form.Get("quantity"))

	// Buat Query sesuai dengan value yang diberikan
	query := "UPDATE transactions SET"
	params := []interface{}{}

	if userId != 0 {
		query += " userId=?,"
		params = append(params, userId)
	}

	if productId != 0 {
		query += " productId=?,"
		params = append(params, productId)
	}

	if quantity != 0 {
		query += " quantity=?,"
		params = append(params, quantity)
	}

	query = strings.TrimSuffix(query, ",")
	query += " WHERE id=?"
	params = append(params, transactionId)
	_, errQuery := db.Exec(query, params...)

	if errQuery == nil {
		var transaction m.Transactions
		transaction.UserID = userId
		transaction.ProductID = productId
		transaction.Quantity = quantity

		var response m.TransactionResponse
		response.Status = 200
		response.Message = "Success"
		response.Data = transaction
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		sendErrorResponse(w, "query error")
	}
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		sendErrorResponse(w, "failed parse")
		return
	}

	userID, _ := strconv.Atoi(r.Form.Get("id"))

	_, errQuery := db.Exec("DELETE FROM users WHERE id=?",
		userID,
	)

	if errQuery == nil {
		sendSuccessResponse(w, "success")
	} else {
		sendErrorResponse(w, "query error")
	}
}

func DeleteUserGorm(w http.ResponseWriter, r *http.Request) {
	db := connectGorm()

	err := r.ParseForm()
	if err != nil {
		sendErrorResponse(w, "failed parse")
		return
	}

	userID, _ := strconv.Atoi(r.Form.Get("id"))

	result := db.Delete(&m.User{}, userID)

	if result.Error == nil {
		sendSuccessResponse(w, "success")
	} else {
		sendErrorResponse(w, "query error")
	}
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		sendErrorResponse(w, "failed parse")
		return
	}

	productId, _ := strconv.Atoi(r.Form.Get("id"))

	if !productExists(productId) {
		sendErrorResponse(w, "product not exists")
		return
	}

	_, errQuery := db.Exec("DELETE FROM transactions WHERE productid=?",
		productId,
	)
	if errQuery != nil {
		sendErrorResponse(w, "error query")
	}

	_, errQuery = db.Exec("DELETE FROM products WHERE id=?",
		productId,
	)

	if errQuery == nil {
		sendSuccessResponse(w, "success")
	} else {
		sendErrorResponse(w, "query error")
	}
}

func DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		sendErrorResponse(w, "failed parse")
		return
	}

	transactionId, _ := strconv.Atoi(r.Form.Get("id"))

	_, errQuery := db.Exec("DELETE FROM transactions WHERE id=?",
		transactionId,
	)

	if errQuery == nil {
		sendSuccessResponse(w, "success")
	} else {
		sendErrorResponse(w, "error query")
	}
}

func sendSuccessResponse(w http.ResponseWriter, pesan string) {
	var response m.SuccessResponse
	response.Status = 200
	response.Message = pesan
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func sendErrorResponse(w http.ResponseWriter, pesan string) {
	var response m.ErrorResponse
	response.Status = 400
	response.Message = pesan
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func userExists(userId int) bool {
	db := connect()
	defer db.Close()

	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE id=?",
		userId,
	).Scan(&count)

	if err != nil {
		log.Println("Error checking UserID:", err)
		return false
	}

	return count > 0
}

func productExists(productId int) bool {
	db := connect()
	defer db.Close()

	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM products WHERE id=?", productId).Scan(&count)

	if err != nil {
		log.Println("Error checking ProductID:", err)
		return false
	}

	return count > 0
}
