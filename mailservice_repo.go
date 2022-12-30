package main

import "database/sql"

//MailItem represents an email data structure
type MailItem struct {
	Name string `json:"name"`
	Mail string `json:"mail"`
	Message string `json:"message"`
}

//saveMail saves a complete mail details into db
func saveMail(db *sql.DB, mail *MailItem) error {
	if _, err := db.Exec("INSERT INTO mails(name, address, message) VALUES (?, ?, ?)", mail.Name, mail.Mail, mail.Message); err != nil {
		return err
	}
	return nil
}

//getMails returns all MailItem saved into db
func getMails(db *sql.DB) ([]MailItem, error) {
	rows, err := db.Query("SELECT * FROM mails")
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