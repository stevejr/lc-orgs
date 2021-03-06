// package function

// import (
// 	"fmt"
// 	"io/ioutil"
// 	"net/http"
// )

// func Handle(w http.ResponseWriter, r *http.Request) {
// 	var input []byte

// 	if r.Body != nil {
// 		defer r.Body.Close()

// 		body, _ := ioutil.ReadAll(r.Body)

// 		input = body
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte(fmt.Sprintf("Hello world, input was: %s", string(input))))
// }

package function

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"

	"github.com/openfaas/openfaas-cloud/sdk"
)

var db *sql.DB
var cors string

// init establishes a persistent connection to the remote database
// the function will panic if it cannot establish a link and the
// container will restart / go into a crash/back-off loop
func init() {
	password, _ := sdk.ReadSecret("mysql-password")
	user, _ := sdk.ReadSecret("mysql-username")
	host, _ := sdk.ReadSecret("mysql-host")

	dbName := os.Getenv("mysql_db")
	port := os.Getenv("mysql_port")

	connStr := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbName

	var err error
	db, err = sql.Open("mysql", connStr)

	if err != nil {
		panic(err.Error())
	}

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	if val, ok := os.LookupEnv("allow_cors"); ok && len(val) > 0 {
		cors = val
	}
}

// Handle a HTTP request as a middleware processor.
func Handle(w http.ResponseWriter, r *http.Request) {

	log.Printf("GET params were: %s", r.URL.Query())

	u := r.URL
	log.Printf("Path: %s", u.Path)
	log.Printf("RawPath: %s", u.RawPath)
	log.Printf("EscapedPath: %s", u.EscapedPath())

	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)

	log.Printf("Request body: %s", body)

	rows, getErr := db.Query(`select * from orgs;`)

	if getErr != nil {
		log.Printf("get error: %s", getErr.Error())
		http.Error(w, errors.Wrap(getErr, "unable to get from orgs").Error(),
			http.StatusInternalServerError)
		return
	}

	results := []Orgs{}
	defer rows.Close()
	for rows.Next() {
		result := Orgs{}
		scanErr := rows.Scan(&result.ID, &result.UniqueID, &result.ShortName, &result.LongName, &result.OwnerEmail)
		if scanErr != nil {
			log.Println("scan err:", scanErr)
		}
		results = append(results, result)
	}

	if len(cors) > 0 {
		w.Header().Set("Access-Control-Allow-Origin", cors)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	res, _ := json.Marshal(results)
	w.Write(res)
}

// Orgs is the struct for the Orgs table
type Orgs struct {
	ID         int
	UniqueID   string
	ShortName  string
	LongName   string
	OwnerEmail string
}
