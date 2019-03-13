package database

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/hashhar/go-hammertime/pkg/models"

	"github.com/hashhar/go-hammertime/pkg/constants"
	// postgres sql driver
	_ "github.com/lib/pq"
)

var (
	db *sql.DB
)

// LoadManagementTable creates an in-memory cache from fclkp_managementlookup table
func LoadManagementTable() ([]int, map[int]string) {
	managementTable := make(map[int]string)
	managementKeys := make([]int, 0)
	rows, err := db.Query("SELECT id, fulfillment_center_sk FROM fclkp_managementlookup;")
	checkErr(err)
	for rows.Next() {
		var id int
		var fcSk string
		err = rows.Scan(&id, &fcSk)
		checkErr(err)
		managementTable[id] = fcSk
		managementKeys = append(managementKeys, id)
	}
	return managementKeys, managementTable
}

// InitDB connects to the database and returns the database object
func InitDB(dbHost string, dbPort int, dbName string) {
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		os.Getenv(constants.UserEnvVar),
		os.Getenv(constants.PassEnvVar),
		dbHost, dbPort, dbName,
	)
	var err error
	db, err = sql.Open("postgres", dbURL)
	checkErr(err)
}

// Close closes the database connection
func Close() {
	db.Close()
}

// Count returns ims_item row count
func Count() int64 {
	var count int64
	res := db.QueryRow("SELECT count(1) FROM ims_item")
	err := res.Scan(&count)
	checkErr(err)
	return count
}

// CreateAndUpdateStatus creates an ims_item entry and then updates it's status
func CreateAndUpdateStatus(row models.Work) {
	var err error
	// Create item
	_, err = db.Exec(constants.CreateQuery, row.Serial, row.FcSk, row.ManagementID)
	checkErr(err)

	// Update item
	_, err = db.Exec(fmt.Sprintf(constants.QueryMap[row.Status], row.Serial))
	checkErr(err)
}

// BatchCreateAndUpdateStatus creates multiple ims_items entries in a single
// query and then updates their status using one query for updating all items
// for a single status
func BatchCreateAndUpdateStatus(rows []models.Work) {
	// Create items
	numQueryParams := 3
	valueStrings := make([]string, 0, len(rows))
	valueArgs := make([]interface{}, 0, len(rows)*numQueryParams)
	statusRowsMap := make(map[string][]string)
	for i := range rows {
		statusRowsMap[rows[i].Status] = append(statusRowsMap[rows[i].Status], rows[i].Serial)
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, '{}', 'BS', false, false, NOW(), NOW())", i*numQueryParams+1, i*numQueryParams+2, i*numQueryParams+3))
		valueArgs = append(valueArgs, rows[i].Serial)
		valueArgs = append(valueArgs, rows[i].FcSk)
		valueArgs = append(valueArgs, rows[i].ManagementID)
	}
	query := fmt.Sprintf(constants.BatchCreateQuery, strings.Join(valueStrings, ","))
	_, err := db.Exec(query, valueArgs...)
	checkErr(err)

	// Update items
	for status, serials := range statusRowsMap {
		_, err := db.Exec(fmt.Sprintf(constants.QueryMap[status], strings.Join(serials, `','`)))
		checkErr(err)
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
