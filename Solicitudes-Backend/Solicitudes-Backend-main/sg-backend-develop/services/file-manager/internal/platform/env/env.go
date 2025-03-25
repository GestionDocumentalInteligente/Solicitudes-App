package env

import (
	"os"
)

var credentials struct {
	email    string
	password string
}

var db struct {
	user     string
	password string
	name     string
	host     string
	sslmode  string
}

var recordData struct {
	treatmentCode string
	reason        string
}

func LoadConfigs() {
	credentials.email = os.Getenv("EMAIL")
	credentials.password = os.Getenv("PASS")

	db.user = os.Getenv("DB_USER")
	db.password = os.Getenv("DB_PASS")
	db.name = os.Getenv("DB_NAME")
	db.host = os.Getenv("DB_HOST")
	db.sslmode = os.Getenv("SSL_MODE")

	recordData.treatmentCode = os.Getenv("RECORD_CODE")
	recordData.reason = os.Getenv("RECORD_REASON")
}

func GetBaseURLGDE() string {
	return os.Getenv("BASE_URL_GDE")
}

func GetEmail() string {
	return credentials.email
}

func GetPassword() string {
	return credentials.password
}

func GetDBUser() string {
	return db.user
}

func GetDBPass() string {
	return db.password
}

func GetDBHost() string {
	return db.host
}

func GetDBName() string {
	return db.name
}

func GetDBSSLMode() string {
	return db.sslmode
}

func GetTreatmentCode() string {
	return recordData.treatmentCode
}

func GetReason() string {
	return recordData.reason
}
