package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"strconv"
	"strings"
	"time"
)

type Note struct {
	ID    string
	year  string
	month string
	day   string
	q1    string
	q2    string
	q3    string
	q4    string
	q5    string
	q6    string
	q7    string
	q8    string
}

type NoteSelector struct {
	ID string
	Q1 string
	Q2 string
	Q3 string
	Q4 string
	Q5 string
	Q6 string
	Q7 string
	Q8 string
}

type QuestionSelector struct {
	Question string
	Q_type   string
}

func InsertNewNote(db *sqlx.DB, n *Note) string {
	tx := db.MustBegin()
	date := time.Now()
	fmt.Println(date.String())

	date_insert := strings.Split(date.String(), ".")[0]
	mt_insert := time.Now()
	mt_update := time.Now()

	_, err := tx.Exec("INSERT INTO diary_helper.notes("+
		"q1, q2, q3, q4, q5, q6, q7, q8, mt_insert, mt_update, date) "+
		"VALUES (CAST($1 AS varchar), CAST($2 AS varchar), CAST($3 AS varchar), "+
		"CAST($4 AS varchar), CAST($5 AS varchar), CAST($6 AS varchar), CAST($7 AS varchar), "+
		"CAST($8 AS varchar), CAST($9 AS timestamp), CAST($10 AS timestamp), CAST($11 as date))",
		n.q1, n.q2, n.q3, n.q4, n.q5, n.q6, n.q7, n.q8, mt_insert, mt_update, date_insert)
	if err != nil {
		fmt.Println("err with insert operation!")
		fmt.Println(err.Error())
		return ""
	}
	err2 := tx.Commit()
	if err2 != nil {
		fmt.Println("err with commit operation!")
		fmt.Println(err2.Error())
		return ""
	}
	return "OK"
}

func DeleteNote(db *sqlx.DB, n *Note) error {
	tx := db.MustBegin()
	q := "delete from diary_helper.notes WHERE id=" + n.ID
	_, err := tx.Exec(q)
	if err != nil {
		return err
	}
	err1 := tx.Commit()
	if err1 != nil {
		return err1
	}
	return nil
}

func UpdateNotes(db *sqlx.DB, n *Note) error {
	tx := db.MustBegin()
	var update_arr []string
	update_arr = append(update_arr, n.q1)
	update_arr = append(update_arr, n.q2)
	update_arr = append(update_arr, n.q3)
	update_arr = append(update_arr, n.q4)
	update_arr = append(update_arr, n.q5)
	update_arr = append(update_arr, n.q6)
	update_arr = append(update_arr, n.q7)
	update_arr = append(update_arr, n.q8)
	fmt.Println(update_arr)
	for i, up := range update_arr {
		if up != "" {
			q := "update diary_helper.notes set q" + strconv.Itoa(i+1) +
				"=" + up + " where id=" + n.ID + ";"
			_, err := tx.Query(q)
			if err != nil {
				fmt.Println(err.Error())
				return err
			}
			err2 := tx.Commit()
			if err2 != nil {
				fmt.Println(err2.Error())
				return err2
			}
		}
	}
	return nil
}

func SelectNotes(db *sqlx.DB, n *Note) ([]NoteSelector, error) {
	answer := []NoteSelector{}
	date := n.year + "-" + n.month + "-" + n.day
	q := "SELECT id, q1, q2, q3, q4, q5, q6, q7, q8 " +
		"FROM diary_helper.notes " +
		"WHERE q1 is not NULL AND q2 is not null AND " +
		"q3 is not NULL and q4 is not NULL and q5 is not NULL " +
		"AND q6 is not NULL AND q7 is not NULL AND q8 is not null AND " +
		"date=CAST( '" + date + "' AS date)"
	row, err := db.Queryx(q)
	if err != nil {
		fmt.Println(err)
		return []NoteSelector{}, err
	}
	for row.Next() {
		ns := NoteSelector{}
		err2 := row.StructScan(&ns)
		if err2 != nil {
			fmt.Println("Scan error!")
			fmt.Println(err2)
			return []NoteSelector{}, err
		}
		answer = append(answer, ns)
	}
	return answer, nil
}

func UpdateGreating(db *sqlx.DB, greating string) error {
	tx := db.MustBegin()
	mt_insert := time.Now()
	mt_update := time.Now()

	_, err := tx.Exec("insert into diary_helper.service "+
		"values ($1, $2, $3)", greating, mt_insert, mt_update)
	if err != nil {
		return err
	}
	err2 := tx.Commit()
	return err2
}

func SelectLastGreatingVal(db *sqlx.DB) (string, error) {
	q := "SELECT greating " +
		"FROM diary_helper.service " +
		"group by greating, mt_insert " +
		"HAVING MAX(mt_insert)=mt_insert " +
		"LIMIT 1"
	var g string
	err := db.Get(&g, q)
	if err != nil {
		return "", err
	}
	return g, nil
}

func SelectQuestions(db *sqlx.DB) ([]QuestionSelector, error) {
	q := "SELECT question, q_type FROM diary_helper.questions"
	row, err := db.Queryx(q)
	res := []QuestionSelector{}
	if err != nil {
		return []QuestionSelector{}, err
	} else {
		for row.Next() {
			qs := QuestionSelector{}
			err2 := row.StructScan(&qs)
			if err2 != nil {
				fmt.Println("Scan error!")
				fmt.Println(err2)
				return []QuestionSelector{}, err
			}
			res = append(res, qs)
		}
	}
	return res, nil
}
