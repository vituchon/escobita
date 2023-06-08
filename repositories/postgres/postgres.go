// Routines for connecting to a postgres dbms

package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Connection = *sql.DB

type Tx = *sql.Tx

// Common interface for sql.Tx and sql.Db regarding Exec, // see : https://github.com/golang/go/issues/14468
type DbTx interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Server struct {
	Host         string `json:"host"`
	Port         int    `json:"port"`
	DatabaseName string `json:"databaseName"`
}

type Config struct {
	Credentials Credentials `json:"credentials"`
	Server      Server      `json:"server"`
}

func OpenConnection(config Config) (Connection, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Server.Host, config.Server.Port, config.Credentials.Username, config.Credentials.Password, config.Server.DatabaseName)
	return sql.Open("postgres", psqlInfo)
}

var Conn Connection = nil

// el init se invoca antes del main asi que queda establecida la conexion... y si hay tests que usan este paquete tmb se correrá el init ANTES de los tests
func init() {
	// por ahora no se usa postgres
	/*config := loadConfig()
	var err error
	Conn, err = OpenConnection(config)
	if err != nil {
		panic(err)
	}
	// verificamos que hay conexión...
	err = Conn.Ping()
	if err != nil {
		panic(err)
	}
	_, err = Conn.Exec("SELECT 1")
	if err != nil {
		panic(err)
	}*/
}

// Load the config for open a connection to a given db
func loadConfig() Config {
	// Esta "hardcodeado" y NO suele ser así, solo para empezar...
	// Luego, el día de mañana tranquilamente acá se puede agregr código que levante un archivo o variables de entorno donde este la configuración!
	var config Config = Config{
		Credentials: Credentials{
			Username: "username",
			Password: "********",
		},
		Server: Server{
			Host:         "locahost",
			Port:         5432,
			DatabaseName: "dbname",
		},
	}
	return config
}
