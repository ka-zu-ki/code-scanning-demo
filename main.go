package main

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"regexp"
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

func main() {
	http.HandleFunc("/vulnerable", vulnerableHandler)
	http.ListenAndServe(":8080", nil)

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

	http.ListenAndServe(":8080", nil) // ここではHTTPSを使用していない

	zero := new(big.Int)
	value := new(big.Int).Div(big.NewInt(10), zero)
	fmt.Println(value)

	http.HandleFunc("/", handler)

	fmt.Println(generateRandom())
	mul([]int{1, 2, 3, 4, 5})
	broken([]byte("forbidden.host.org"))

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

func mul(xs []int) int {
	res := 1
	for i := 0; i < len(xs); i++ {
		x := xs[i]
		res *= x
		if res == 0 {
		}
		return 0
	}
	return res
}

func broken(hostNames []byte) string {
	var hostRe = regexp.MustCompile("\bforbidden.host.org")
	if hostRe.Match(hostNames) {
		return "Must not target forbidden.host.org"
	} else {
		// This will be reached even if hostNames is exactly "forbidden.host.org",
		// because the literal backspace is not matched
		return ""
	}
}
