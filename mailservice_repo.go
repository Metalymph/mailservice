package main

import (
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

// MailItem represents an email data structure
type MailItem struct {
	Name    string `json:"name"`
	Mail    string `json:"mail"`
	Message string `json:"message"`
}

// saveMail saves a complete mail details into db
func saveMail(a *App, mail *MailItem) error {
	switch a.dbtype {
	case Sqlite:
		if _, err := a.db.Exec("INSERT INTO mails(name, address, message) VALUES (?, ?, ?)", mail.Name, mail.Mail, mail.Message); err != nil {
			return err
		}
	case Postgresql:
		if _, err := a.db.Exec(`insert into "Mails"("Name", "Address", "Message") values($1, $2, $3)`, mail.Name, mail.Mail, mail.Message); err != nil {
			return err
		}
	}
	return nil
}

// getMails returns all MailItem saved into db
func getMails(a *App) ([]MailItem, error) {
	rows, err := a.db.Query("SELECT * FROM mails")
	if err != nil {
		return nil, err
	}

	var mail MailItem
	var mails []MailItem
	for rows.Next() {
		err = rows.Scan(&mail.Name, &mail.Mail, &mail.Message)
		if err != nil {
			return nil, err
		}
		mails = append(mails, mail)
	}
	return mails, nil
}
