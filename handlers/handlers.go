package handlers

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"

	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Values struct {
	ID        uint
	FirstName string
	LastName  string
	City      string
	Email     string
}
type Name struct {
	ID           uint
	Name         string
	WithEmail    int
	WithoutEmail int
}

func Home(c *gin.Context) {

	name := c.PostForm("name")
	f, err := c.FormFile("file")

	if err != nil {
		c.AbortWithStatusJSON(404, gin.H{"err": err.Error()})
		return
	}
	extension := filepath.Ext(f.Filename)
	fmt.Print(extension)
	if extension != ".csv" {
		c.JSON(http.StatusBadRequest, gin.H{"err": "unsupported file format"})
		return
	}
	cs := uuid.New().String() + extension
	c.SaveUploadedFile(f, "./"+cs)
	file, err := os.Open(cs)
	if err != nil {
		c.AbortWithStatusJSON(404, gin.H{"err": err.Error()})
		return
	}
	defer file.Close()
	read := csv.NewReader(file)
	red, err := read.ReadAll()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	fmt.Print(name)
	withemail := 0
	withoutEmail := 0
	var s []Values
	for i, v := range red {
		if i > 0 {
			var rec Values
			for j, field := range v {
				if j == 1 {
					rec.FirstName = field

				} else if j == 2 {
					rec.LastName = field
				} else if j == 3 {
					if field != "" {
						withemail++
						rec.Email = field
					} else {
						withoutEmail++
						break
					}
				} else if j == 4 {
					rec.City = field
				} else {
					print(field)
				}

			}
			s = append(s, rec)

			result := Db.Create(&s)
			if result.Error != nil {
				c.AbortWithStatusJSON(400, gin.H{
					"error": result.Error.Error(),
				})
				return
			}
		}

	}
	var nme Name
	nme.Name = name
	nme.WithEmail = withemail
	nme.WithoutEmail = withoutEmail
	r := Db.Create(&nme)
	if r.Error != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"error": r.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"res": "successfully added",
	})

}

func LoadEnvVariables() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error while loading env file")
	}
}

var Db *gorm.DB

func Connect() {
	var err error
	dbhost := os.Getenv("dbHost")
	dbuser := os.Getenv("dbUser")
	dbpassword := os.Getenv("dbPassword")
	dbname := os.Getenv("dbName")
	dbport := os.Getenv("dbPort")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable port=%s", dbhost, dbuser, dbpassword, dbname, dbport)
	Db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Db not connected")
	}
	Db.AutoMigrate(
		&Name{},
		&Values{},
	)

}
