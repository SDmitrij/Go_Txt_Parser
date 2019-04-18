package lsa

import (
	"database/sql"
)

type FilesRepo struct {
	Params     map[string]string
	Connection *sql.DB
}

/**
Init files repository
 */
func (fr *FilesRepo) initFilesRepo() {
	fr.createEntryDatabase()
	fr.createFilesInfoTable()
}

/**
Create an entry database
 */
func (fr *FilesRepo) createEntryDatabase() {
	_, err := fr.Connection.Exec("CREATE DATABASE IF NOT EXISTS " + fr.Params["db_name"] + " CHARSET=utf8")
	if err != nil {
		panic(err)
	}
}

/**
Create main files info table
 */
func (fr *FilesRepo) createFilesInfoTable() {
	_, err := fr.Connection.Exec("CREATE TABLE IF NOT EXISTS " + fr.Params["db_name"] + "." +
		fr.Params["tbl_idx"] +
		"(id INT(10) UNSIGNED AUTO_INCREMENT PRIMARY KEY," +
		"file_path VARCHAR(100) NOT NULL," +
		"file_hash VARCHAR(32) NOT NULL," +
		"file_unique_key VARCHAR(32) NOT NULL," +
		"file_size INT(10) NOT NULL," +
		"is_index TINYINT(1) NOT NULL)")
	if err != nil {
		panic(err)
	}
}

/**
Create table that keeps strings of file
 */
func (fr *FilesRepo) createTableStrings(fileUniqueKey string, tblPref string) {
	_, err := fr.Connection.Exec("CREATE TABLE IF NOT EXISTS " + fr.Params["db_name"] + "." +
		fr.Params[tblPref] + fileUniqueKey +
		"(id INT(10) UNSIGNED AUTO_INCREMENT PRIMARY KEY," +
		"string_of_file VARCHAR(200) NOT NULL," +
		"INDEX str_idx (string_of_file))")
	if err != nil {
		panic(err)
	}
}

/**
Create table that keeps words of file
 */
func (fr *FilesRepo) createTableTerms(fileUniqueKey string, tblPref string) {
	_, err := fr.Connection.Exec("CREATE TABLE IF NOT EXISTS " + fr.Params["db_name"] + "." +
		fr.Params[tblPref] + fileUniqueKey +
		"(id INT(10) UNSIGNED AUTO_INCREMENT PRIMARY KEY," +
		"term_of_file VARCHAR(50) NOT NULL," +
		"INDEX term_idx (term_of_file))")
	if err != nil {
		panic(err)
	}
}

func (fr *FilesRepo) deleteFileInfo(fileUniqueKey string) {

	sqlDelete := "DELETE FROM " + fr.Params["db_name"] + "." + fr.Params["tbl_idx"] +
		" WHERE file_unique_key = ?;"
	sqlDelete += "DROP TABLE IF EXISTS " + fr.Params["db_name"] + "." + fr.Params["tbl_str_pref"] +
		fileUniqueKey + ";"
	sqlDelete += "DROP TABLE IF EXISTS " + fr.Params["db_name"] + "." + fr.Params["tbl_term_pref"] +
		fileUniqueKey + ";"

	_, err := fr.Connection.Exec(sqlDelete, fileUniqueKey)

	if err != nil {
		panic(err)
	}
}

/**
Insert file's data into main info table
 */
func (fr *FilesRepo) insIntoMainInfoFileTable(file File) {
	_, err := fr.Connection.Exec("INSERT INTO " + fr.Params["db_name"] + "." + fr.Params["tbl_idx"] +
		"(file_path, file_hash, file_unique_key, file_size, is_index) VALUES (?, ?, ?, ?, ?)",
		file.filePath, file.fileHash, file.fileUniqueKey, file.fileSize, 0)
	if err != nil {
		panic(err)
	}
}

/**
Insert data into table with strings of current file
 */
func (fr *FilesRepo) insIntoTableStrings(stringAndKey map[string]string, tblPref string) {
	_, err := fr.Connection.Exec("INSERT INTO " + fr.Params["db_name"] + "." + fr.Params[tblPref] +
		stringAndKey["file_key"] + "(string_of_file) VALUES (?)",
		stringAndKey["str_of_file"])
	if err != nil {
		panic(err)
	}
}

/**
Insert data into table with words of current file
 */
func (fr *FilesRepo) insIntoTableTerms(termAndKey map[string]string, tblPref string) {
	_, err := fr.Connection.Exec("INSERT INTO " + fr.Params["db_name"] + "." + fr.Params[tblPref] +
		termAndKey["file_key"] + "(term_of_file) VALUES (?)",
		termAndKey["term_of_file"])
	if err != nil {
		panic(err)
	}
}

/**
Get current file info as object
 */
func (fr *FilesRepo) getFileInfoAsObj(fileUniqueKey string) File {
	file := File{}
	// Getting current file
	err := fr.Connection.QueryRow("SELECT file_path, file_unique_key, file_hash, file_size FROM " +
		fr.Params["db_name"] + "." + fr.Params["tbl_idx"] +
		" WHERE file_unique_key = ?", fileUniqueKey).Scan(&file.filePath, &file.fileUniqueKey, &file.fileHash,
			&file.fileSize)
	// Handling error
	if err != nil && err != sql.ErrNoRows {
		panic(err)
	}

	return file
}

/**
Get all terms of current file
 */
func (fr *FilesRepo) getAllTermsOfFile(fileUniqueKey string, tblPref string) *[]string {

	var terms []string

	rawRes := make([][]byte, 1)
	temp :=   make([]interface{}, 1)
	temp[0] = & rawRes[0]

	result, err := fr.Connection.Query("SELECT term_of_file FROM "+ fr.Params["db_name"] + "." +
		fr.Params[tblPref] + fileUniqueKey)

	if err != nil {
		panic(err)
	}

	for result.Next() {
		err := result.Scan(temp...)
		if err != nil {
			panic(err)
		}

		terms = append(terms, string(rawRes[0]))
	}

	return &terms
}