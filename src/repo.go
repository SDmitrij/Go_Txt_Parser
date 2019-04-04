package main

import (
	"database/sql"
)

type filesRepo struct {
	dbTblParams map[string]string
	dbConnection *sql.DB
}

/**
Init files repository
 */
func (fr *filesRepo) initFilesRepo() {
	fr.createEntryDatabase()
	fr.createFilesInfoTable()
}

/**
Create an entry database
 */
func (fr *filesRepo) createEntryDatabase() {
	_, err := fr.dbConnection.Exec("CREATE DATABASE IF NOT EXISTS " + fr.dbTblParams["db_name"] + " CHARSET=utf8")
	if err != nil {
		panic(err)
	}
}

/**
Create main files info table
 */
func (fr *filesRepo) createFilesInfoTable() {
	_, err := fr.dbConnection.Exec("CREATE TABLE IF NOT EXISTS " + fr.dbTblParams["db_name"] + "." +
		fr.dbTblParams["tbl_idx"] +
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
func (fr *filesRepo) createTableStrings(fileUniqueKey string, tblPref string) {
	_, err := fr.dbConnection.Exec("CREATE TABLE IF NOT EXISTS " + fr.dbTblParams["db_name"] + "." +
		fr.dbTblParams[tblPref] + fileUniqueKey +
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
func (fr *filesRepo) createTableWords(fileUniqueKey string, tblPref string) {
	_, err := fr.dbConnection.Exec("CREATE TABLE IF NOT EXISTS " + fr.dbTblParams["db_name"] + "." +
		fr.dbTblParams[tblPref] + fileUniqueKey +
		"(id INT(10) UNSIGNED AUTO_INCREMENT PRIMARY KEY," +
		"word_of_file VARCHAR(50) NOT NULL," +
		"INDEX wrd_idx (word_of_file))")
	if err != nil {
		panic(err)
	}
}

func (fr *filesRepo) deleteFileInfo(fileUniqueKey string) {

	sqlDelete := "DELETE FROM " + fr.dbTblParams["db_name"] + "." + fr.dbTblParams["tbl_idx"] +
		" WHERE file_unique_key = ?;"
	sqlDelete += "DROP TABLE IF EXISTS " + fr.dbTblParams["db_name"] + "." + fr.dbTblParams["tbl_str_pref"] +
		fileUniqueKey + ";"
	sqlDelete += "DROP TABLE IF EXISTS " + fr.dbTblParams["db_name"] + "." + fr.dbTblParams["tbl_wrd_pref"] +
		fileUniqueKey + ";"

	_, err := fr.dbConnection.Exec(sqlDelete, fileUniqueKey)

	if err != nil {
		panic(err)
	}
}

/**
Insert file's data into main info table
 */
func (fr *filesRepo) insIntoMainInfoFileTable(file File) {
	_, err := fr.dbConnection.Exec("INSERT INTO " + fr.dbTblParams["db_name"] + "." + fr.dbTblParams["tbl_idx"] +
		"(file_path, file_hash, file_unique_key, file_size, is_index) VALUES (?, ?, ?, ?, ?)",
		file.filePath, file.fileHash, file.fileUniqueKey, file.fileSize, 0)
	if err != nil {
		panic(err)
	}
}

/**
Insert data into table with strings of current file
 */
func (fr *filesRepo) insIntoTableStrings(stringAndKey map[string]string, tblPref string) {
	_, err := fr.dbConnection.Exec("INSERT INTO " + fr.dbTblParams["db_name"] + "." + fr.dbTblParams[tblPref] +
		stringAndKey["file_key"] + "(string_of_file) VALUES (?)",
		stringAndKey["str_of_file"])
	if err != nil {
		panic(err)
	}
}

/**
Insert data into table with words of current file
 */
func (fr *filesRepo) insIntoTableWords(wordAndKey map[string]string, tblPref string) {
	_, err := fr.dbConnection.Exec("INSERT INTO " + fr.dbTblParams["db_name"] + "." + fr.dbTblParams[tblPref] +
		wordAndKey["file_key"] + "(word_of_file) VALUES (?)",
		wordAndKey["wrd_of_file"])
	if err != nil {
		panic(err)
	}
}

/**
Get current file info as object
 */
func (fr *filesRepo) getFileInfoAsObj(fileUniqueKey string) File {
	file := File{}
	// Getting current file
	err := fr.dbConnection.QueryRow("SELECT file_path, file_unique_key, file_hash, file_size FROM " +
		fr.dbTblParams["db_name"] + "." + fr.dbTblParams["tbl_idx"] +
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
func (fr *filesRepo) getAllTermsOfFile(fileUniqueKey string, tblPref string) *[]string {

	var terms []string

	rawRes := make([][]byte, 1)
	temp := make([]interface{}, 1)
	temp[0] = & rawRes[0]

	result, err := fr.dbConnection.Query("SELECT word_of_file FROM "+ fr.dbTblParams["db_name"] + "." +
		fr.dbTblParams[tblPref] + fileUniqueKey)

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