package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const (
	// URL for accessing RabbitMQ
	AMQP_URL = "amqp://sift:sift@rabbitmq:5672/sift"
	// URL for accessing Redis
	REDIS_URL = "redis:6379"
	// Max file size to store in memory. 100MB
	MAX_FILE_SIZE = 6 << 24
)

var db_host = flag.String("dbhost", "127.0.0.1", "The address at which the db listens")

// Configures the databse with user, password, host, name, and SSL encryption
// type
type DBConfig struct {
	DBUser     string
	DBPassword string
	DBHost     string
	DBName     string
	DBSSLType  string
}

func (cfg DBConfig) createDBQueryString() string {
	return fmt.Sprintf("user=%s password=%s host=%s dbname=%s sslmode=%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBName, cfg.DBSSLType)
}

// `main` function is the entry point, just like in C
func main() {
	flag.Parse()
	// Database configuration
	cfg := DBConfig{
		DBUser:     "test",
		DBPassword: "testpw",
		DBHost:     *db_host,
		DBName:     "sift_user_data",
		DBSSLType:  "disable", // switch to 'require' in production
	}
	// Lazily open a connection to the database. The database will only
	// be opened when the first query/exec statement is made against it
	// while service client requests
	db, err := gorm.Open("postgres", cfg.createDBQueryString())
	// Create a DataManager with our lazy connection (see datalayer.go)
	if err != nil {
		log.Fatal("gorm.Open: ", err)
	}
	// Close the connection on main() exit
	defer db.Close()
	dm := NewDataManager(db)
	// Migration of native types, which can be added as arguments as needed
	dm.AutoMigrate(&Profile{})
	dm.AutoMigrate(&Session{})
	// Create a new router, routers handle sets of logically related routes
	router := mux.NewRouter()
	// Handler for the feedback upload route
	router.HandleFunc("/feedback", FeedbackFormHandler).Methods("POST")
	// Handlers for profile CRUD operations
	router.HandleFunc("/profile", dm.IndexNewProfile).Methods("POST")
	router.HandleFunc("/profile/{company_name}/{user_name}", dm.GetExistingProfile).Methods("GET")
	router.HandleFunc("/profile/{company_name}/{user_name}", dm.UpdateExistingProfile).Methods("PUT")
	router.HandleFunc("/profile/{company_name}/{user_name}", dm.DeleteExistingProfile).Methods("DELETE")
	// Handlers for logins
	router.HandleFunc("/login", dm.Login).Methods("POST")
	router.HandleFunc("/logout", dm.Logout).Methods("POST")
	http.Handle("/", router)
	// Create an http server on port 9090 and start serving using our router.
	fmt.Println("Sift API running on port 9090...")
	if err := http.ListenAndServe(":9090", router); err != nil {
		log.Fatal("http.ListenAndServe: ", err)
	}

}
