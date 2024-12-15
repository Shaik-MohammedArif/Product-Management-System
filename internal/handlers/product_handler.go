package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
)

func GetAllProductsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse query parameters
		userID := r.URL.Query().Get("user_id")
		minPrice := r.URL.Query().Get("min_price")
		maxPrice := r.URL.Query().Get("max_price")
		name := r.URL.Query().Get("product_name")

		// Validate required parameters
		if userID == "" {
			http.Error(w, "`user_id` is required", http.StatusBadRequest)
			return
		}

		// Build SQL query with optional filters
		query := `SELECT id, user_id, product_name, product_description, product_images, product_price
				  FROM products WHERE user_id = $1`
		args := []interface{}{userID}

		// Add filters for price range
		if minPrice != "" && maxPrice != "" {
			query += " AND product_price BETWEEN $2 AND $3"
			min, _ := strconv.ParseFloat(minPrice, 64)
			max, _ := strconv.ParseFloat(maxPrice, 64)
			args = append(args, min, max)
		}

		// Add filter for product name
		if name != "" {
			query += " AND product_name ILIKE $4"
			args = append(args, "%"+name+"%")
		}

		// Execute query
		rows, err := db.Query(query, args...)
		if err != nil {
			http.Error(w, "Failed to fetch products: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Parse results into a slice of products
		type Product struct {
			ID          int      `json:"id"`
			UserID      int      `json:"user_id"`
			ProductName string   `json:"product_name"`
			Description string   `json:"product_description"`
			Images      []string `json:"product_images"`
			Price       float64  `json:"product_price"`
		}

		var products []Product
		for rows.Next() {
			var product Product
			var images pq.StringArray

			err := rows.Scan(&product.ID, &product.UserID, &product.ProductName, &product.Description, &images, &product.Price)
			if err != nil {
				http.Error(w, "Error reading product data: "+err.Error(), http.StatusInternalServerError)
				return
			}
			product.Images = images
			products = append(products, product)
		}

		// Return the products as JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(products)
	}
}

// GetProductByIDHandler - Get a product by its ID
// GetProductByIDHandler - Get a product by its ID
func GetProductByIDHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract product ID from URL parameters (not query parameters)
		vars := mux.Vars(r)
		productID := vars["id"]
		if productID == "" {
			http.Error(w, "`id` is required", http.StatusBadRequest)
			return
		}

		// Prepare the query
		query := `SELECT id, user_id, product_name, product_description, product_images, product_price
				  FROM products WHERE id = $1`
		var product struct {
			ID          int      `json:"id"`
			UserID      int      `json:"user_id"`
			ProductName string   `json:"product_name"`
			Description string   `json:"product_description"`
			Images      []string `json:"product_images"`
			Price       float64  `json:"product_price"`
		}

		err := db.QueryRow(query, productID).Scan(&product.ID, &product.UserID, &product.ProductName, &product.Description, pq.Array(&product.Images), &product.Price)
		if err != nil {
			http.Error(w, "Failed to fetch product: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Return the product as JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(product)
	}
}



// CreateProductHandler - Create a new product
func CreateProductHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the request body into a product struct
		var product struct {
			ID          int      `json:"id"` // Add ID field here
			UserID      int      `json:"user_id"`
			ProductName string   `json:"product_name"`
			Description string   `json:"product_description"`
			Images      []string `json:"product_images"`
			Price       float64  `json:"product_price"`
		}

		// Decode the JSON body
		err := json.NewDecoder(r.Body).Decode(&product)
		if err != nil {
			http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Prepare the SQL query to insert the new product
		query := `INSERT INTO products (user_id, product_name, product_description, product_images, product_price)
				  VALUES ($1, $2, $3, $4, $5) RETURNING id`
		var newProductID int
		err = db.QueryRow(query, product.UserID, product.ProductName, product.Description, pq.Array(product.Images), product.Price).Scan(&newProductID)
		if err != nil {
			http.Error(w, "Failed to create product: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Set the ID of the created product
		product.ID = newProductID

		// Return the newly created product as JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(product)
	}
}



// UpdateProductHandler - Update a product by its ID
func UpdateProductHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract product ID from URL parameters
		vars := mux.Vars(r)
		productID := vars["id"]
		if productID == "" {
			http.Error(w, "`id` is required", http.StatusBadRequest)
			return
		}

		// Parse the request body into a product struct
		var product struct {
			UserID      int      `json:"user_id"`
			ProductName string   `json:"product_name"`
			Description string   `json:"product_description"`
			Images      []string `json:"product_images"`
			Price       float64  `json:"product_price"`
		}

		// Decode the JSON body
		err := json.NewDecoder(r.Body).Decode(&product)
		if err != nil {
			http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Prepare the SQL query to update the product
		query := `UPDATE products
				  SET user_id = $1, product_name = $2, product_description = $3, product_images = $4, product_price = $5
				  WHERE id = $6`
		_, err = db.Exec(query, product.UserID, product.ProductName, product.Description, pq.Array(product.Images), product.Price, productID)
		if err != nil {
			http.Error(w, "Failed to update product: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Return a success message
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Product updated successfully"))
	}
}
// DeleteProductHandler - Delete a product by its ID
func DeleteProductHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract product ID from URL parameters
		vars := mux.Vars(r)
		productID := vars["id"]
		if productID == "" {
			http.Error(w, "`id` is required", http.StatusBadRequest)
			return
		}

		// Prepare the SQL query to delete the product
		query := `DELETE FROM products WHERE id = $1`
		result, err := db.Exec(query, productID)
		if err != nil {
			http.Error(w, "Failed to delete product: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Check if any rows were affected (product exists)
		rowsAffected, err := result.RowsAffected()
		if err != nil || rowsAffected == 0 {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}

		// Return a success message
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Product deleted successfully"))
	}

	
}
