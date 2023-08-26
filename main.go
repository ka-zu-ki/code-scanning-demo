package main

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
)

const apiKey = "SECRET_API_KEY_12345"

var db *sql.DB

func vulnerableHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")

	// 脆弱なSQLクエリ
	query := fmt.Sprintf("SELECT * FROM users WHERE id = %s", userID)
	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<script>alert('XSS!');</script>"))
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	http.Redirect(w, r, url, http.StatusFound)
}

func main() {
	http.HandleFunc("/vulnerable", vulnerableHandler)
 
	http.HandleFunc("/header-injection", func(w http.ResponseWriter, r *http.Request) {
		value := r.URL.Query().Get("header_value")
		w.Header().Set("X-Custom-Header", value)
		w.Write([]byte("Header set!"))
	})

	http.HandleFunc("/path-traversal", func(w http.ResponseWriter, r *http.Request) {
		filename := r.URL.Query().Get("filename")
		data, err := ioutil.ReadFile("/safe/path/" + filename)
		fmt.Sprintf(string(data))
		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
	})

	http.HandleFunc("/change-password", func(w http.ResponseWriter, r *http.Request) {
		newPassword := r.URL.Query().Get("new_password")
		fmt.Sprintf(newPassword)
	})

	http.HandleFunc("/cors", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write([]byte("CORS headers set!"))
	})

	var counter int
	http.HandleFunc("/race-condition", func(w http.ResponseWriter, r *http.Request) {
		counter++
		w.Write([]byte(fmt.Sprintf("Counter: %d", counter)))
	})

	http.HandleFunc("/bad-error-handling", func(w http.ResponseWriter, r *http.Request) {
		_, err := someCriticalFunction()
		if err != nil {
			w.Write([]byte("Oops, something went wrong!"))
		}
	})

	zero := new(big.Int)
	value := new(big.Int).Div(big.NewInt(10), zero)
	fmt.Println(value)

	http.HandleFunc("/", handler)

	fmt.Println(generateRandom())
	http.HandleFunc("/redirect", redirectHandler)

	http.ListenAndServe(":8080", nil)
}

func someCriticalFunction() (string, error) {
	return "", fmt.Errorf("Critical error!")
}

func generateRandom() []byte {
	buffer := make([]byte, 10)
	rand.Read(buffer)
	return buffer
}
