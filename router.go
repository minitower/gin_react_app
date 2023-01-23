package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"net/http"
)

func PGSQLConnection() (*sqlx.DB, error) {
	conn, err := sqlx.Connect("postgres",
		"user=minitower dbname=vera_db sslmode=disable")
	return conn, err
}

func ReturnMain(c *gin.Context) {
	db, err := PGSQLConnection()
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": "delete",
		})
	}
	greating, _ := SelectLastGreatingVal(db)
	c.SetCookie("greating", greating, 1000, "/", "localhost",
		false, false)
	c.HTML(http.StatusOK, "index.html", gin.H{
		"err": "",
	})
}

func SaveMain(c *gin.Context) {
	conn, err := PGSQLConnection()
	fmt.Println(c)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "index.html",
			gin.H{
				"page": "err_db",
				"err":  "can't connect",
			})
	} else {
		q1, _ := c.Cookie("q1")
		q2, _ := c.Cookie("q2")
		q3, _ := c.Cookie("q3")
		q4, _ := c.Cookie("q4")
		q5, _ := c.Cookie("q5")
		q6, _ := c.Cookie("q6")
		q7, _ := c.Cookie("q7")
		q8, _ := c.Cookie("q8")

		n := Note{
			q1: q1,
			q2: q2,
			q3: q3,
			q4: q4,
			q5: q5,
			q6: q6,
			q7: q7,
			q8: q8,
		}
		fmt.Println(n)
		status := InsertNewNote(conn, &n)
		if status != "OK" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"err": "insert",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"err": "",
			})
		}
	}
}

func ShowNote(c *gin.Context) {
	conn, err := PGSQLConnection()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "index.html",
			gin.H{
				"page": "err_db",
				"err":  "can't connect",
			})
	}
	year, _ := c.GetPostForm("year")
	month, _ := c.GetPostForm("month")
	day, _ := c.GetPostForm("day")

	fmt.Println(year)
	fmt.Println(month)
	fmt.Println(day)
	n := Note{
		year:  year,
		month: month,
		day:   day,
	}
	fmt.Println(n)
	if err != nil {
		fmt.Println("Error parse date!")
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "date",
		})
	}

	ns, err := SelectNotes(conn, &n)
	if err != nil {
		fmt.Println("Error to select notes!")
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": "select",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"result": ns,
		})
	}
}

func UpdateNote(c *gin.Context) {
	method := c.Query("method")
	db, err := PGSQLConnection()
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": "delete",
		})
	}

	n := Note{
		ID: c.Query("id"),
		q1: c.Query("q1"),
		q2: c.Query("q2"),
		q3: c.Query("q3"),
		q4: c.Query("q4"),
		q5: c.Query("q5"),
		q6: c.Query("q6"),
		q7: c.Query("q7"),
		q8: c.Query("q8"),
	}

	if method == "delete" {
		err := DeleteNote(db, &n)
		if err != nil {
			fmt.Println(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"err": "delete",
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"err": "",
		})
	} else if method == "update" {

		err := UpdateNotes(db, &n)
		if err != nil {
			fmt.Println(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"err": "update",
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"err": "",
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "no-method",
		})
	}
}

func QueryTest(c *gin.Context) {
	db, err := PGSQLConnection()
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": "delete",
		})
	}
	res, err := SelectQuestions(db)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": "select",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"questions": res,
	})
}

func DeleteGreating(c *gin.Context) {
	db, err := PGSQLConnection()
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": "delete",
		})
	}
	err2 := UpdateGreating(db, c.Query("greating"))
	if err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err2.Error(),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"err": "",
	})
}

func GetGreating(c *gin.Context) {
	db, err := PGSQLConnection()
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": "delete",
		})
	}

	greating, err := SelectLastGreatingVal(db)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err.Error(),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"err":      "",
		"greating": greating,
	})
}

func BuildRouter(r *gin.Engine) {
	r.GET("/", ReturnMain)
	r.GET("vera/notes", ShowNote)
	r.POST("vera/notes", UpdateNote)
	r.POST("vera/g", DeleteGreating)
	r.GET("vera/g", GetGreating)
	r.POST("vera/diary_helper", SaveMain)
	r.GET("vera/test", QueryTest)
}
