package account

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
	"test/config"
)

const (
	INSERT         string = "INSERT INTO `accounts`(name, password) VALUES (?, ?)"
	COUNT_BY_NAME  string = "SELECT COUNT(*) FROM `accounts` WHERE `name` = ?"
	SELECT_ALL     string = "SELECT * FROM `accounts`"
	DELETE_BY_NAME string = "DELETE FROM `accounts` WHERE `name` = ?"
)

type Account struct {
	Name     string `json:"name"`
	Password string `json:"password" bson:"password"`
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func GetAllAccounts(config config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var response Response
		var listOfAccounts []Account
		dataSourceName := config.DBUsername + ":" + config.DBPassword + "@tcp(" + config.Source + ":" + config.DBPort + ")/" + config.DBName
		db, err := sql.Open(config.DBDriver, dataSourceName)
		if err != nil {
			response.Code = http.StatusInternalServerError
			response.Message = err.Error()
			response.Data = nil
			c.JSON(response.Code, response)
			return
		}
		defer func(db *sql.DB) {
			err := db.Close()
			if err != nil {
				response.Code = http.StatusInternalServerError
				response.Message = err.Error()
				response.Data = nil
				c.JSON(response.Code, response)
				return
			}
		}(db)
		res, err := db.Query(SELECT_ALL)
		if err != nil {
			response.Code = http.StatusInternalServerError
			response.Message = err.Error()
			response.Data = nil
			c.JSON(response.Code, response)
			return
		}
		for res.Next() {
			var id int
			var name string
			var password string
			err = res.Scan(&id, &name, &password)
			if err != nil {
				response.Code = http.StatusInternalServerError
				response.Message = err.Error()
				response.Data = nil
				c.JSON(response.Code, response)
				return
			}
			listOfAccounts = append(listOfAccounts, Account{
				Name:     name,
				Password: password,
			})
		}
		response.Code = http.StatusOK
		response.Message = "successfully"
		response.Data = listOfAccounts
		c.JSON(response.Code, response)
	}
}

func CreateAccount(config config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var response Response
		var account Account
		if err := c.ShouldBindBodyWith(&account, binding.JSON); err != nil {
			response.Code = http.StatusInternalServerError
			response.Message = err.Error()
			response.Data = nil
			c.JSON(response.Code, response)
			return
		}
		dataSourceName := config.DBUsername + ":" + config.DBPassword + "@tcp(" + config.Source + ":" + config.DBPort + ")/" + config.DBName
		db, err := sql.Open(config.DBDriver, dataSourceName)
		if err != nil {
			response.Code = http.StatusInternalServerError
			response.Message = err.Error()
			response.Data = nil
			c.JSON(response.Code, response)
			return
		}
		defer func(db *sql.DB) {
			err := db.Close()
			if err != nil {
				response.Code = http.StatusInternalServerError
				response.Message = err.Error()
				response.Data = nil
				c.JSON(response.Code, response)
				return
			}
		}(db)
		isNameExisted, err := db.Query(COUNT_BY_NAME, account.Name)
		if err != nil {
			response.Code = http.StatusInternalServerError
			response.Message = err.Error()
			response.Data = nil
			c.JSON(response.Code, response)
			return
		}
		isNameExisted.Next()
		var numOfNames int
		if err = isNameExisted.Scan(&numOfNames); err != nil {
			response.Code = http.StatusInternalServerError
			response.Message = err.Error()
			response.Data = nil
			c.JSON(response.Code, response)
			return
		}
		if numOfNames != 0 {
			response.Code = http.StatusBadRequest
			response.Message = "name already exists"
			response.Data = nil
			c.JSON(response.Code, response)
			return
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
		if err != nil {
			response.Code = http.StatusInternalServerError
			response.Message = err.Error()
			response.Data = nil
			c.JSON(response.Code, response)
			return
		}
		res, err := db.Exec(INSERT, account.Name, string(hashedPassword))
		if err != nil {
			response.Code = http.StatusInternalServerError
			response.Message = err.Error()
			response.Data = nil
			c.JSON(response.Code, response)
			return
		}
		lastInsertedID, _ := res.LastInsertId()
		rowsAffected, _ := res.RowsAffected()
		if rowsAffected == 0 {
			response.Code = http.StatusInternalServerError
			response.Message = "failed to insert new account"
			response.Data = nil
			c.JSON(response.Code, response)
			return
		}
		response.Code = http.StatusOK
		response.Message = "successfully added. id = " + strconv.Itoa(int(lastInsertedID))
		response.Data = account
		c.JSON(response.Code, response)
	}
}

func DeleteAccount(config config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var response Response
		name := c.Param("name")
		if name == "" {
			response.Code = http.StatusBadRequest
			response.Message = "missing name"
			response.Data = nil
			c.JSON(response.Code, response)
		}
		dataSourceName := config.DBUsername + ":" + config.DBPassword + "@tcp(" + config.Source + ":" + config.DBPort + ")/" + config.DBName
		db, err := sql.Open(config.DBDriver, dataSourceName)
		if err != nil {
			response.Code = http.StatusInternalServerError
			response.Message = err.Error()
			response.Data = nil
			c.JSON(response.Code, response)
		}
		defer func(db *sql.DB) {
			err := db.Close()
			if err != nil {
				response.Code = http.StatusInternalServerError
				response.Message = err.Error()
				response.Data = nil
				c.JSON(response.Code, response)
				return
			}
		}(db)
		res, err := db.Exec(DELETE_BY_NAME, name)
		if err != nil {
			response.Code = http.StatusInternalServerError
			response.Message = err.Error()
			response.Data = nil
			c.JSON(response.Code, response)
			return
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			response.Code = http.StatusInternalServerError
			response.Message = err.Error()
			response.Data = nil
			c.JSON(response.Code, response)
			return
		}
		if rowsAffected == 0 {
			response.Code = http.StatusInternalServerError
			response.Message = "no account such that " + name
			response.Data = nil
			c.JSON(response.Code, response)
			return
		}
		response.Code = http.StatusOK
		response.Message = "successfully removed"
		response.Data = struct {
			Name string
		}{Name: name}
		c.JSON(response.Code, response)
	}
}
