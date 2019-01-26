package main

import "database/sql"

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
	_, err := fr.dbConnection.Exec("CREATE DATABASE IF NOT EXISTS " + fr.dbTblParams["db_name"])
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
func (fr *filesRepo) createTableStrings(fileUniqueKey string) {
	_, err := fr.dbConnection.Exec("CREATE TABLE IF NOT EXISTS " + fr.dbTblParams["db_name"] + "." +
		fr.dbTblParams["tbl_str_pref"] + fileUniqueKey +
		"(id INT(10) UNSIGNED AUTO_INCREMENT PRIMARY KEY," +
		"string_of_file VARCHAR(200) NOT NULL," +
		"num_of_line INT(10) NOT NULL," +
		"INDEX str_idx (string_of_file))")
	if err != nil {
		panic(err)
	}
}

/**
Create table that keeps words of file
 */
func (fr *filesRepo) createTableWords(fileUniqueKey string) {
	_, err := fr.dbConnection.Exec("CREATE TABLE IF NOT EXISTS " + fr.dbTblParams["db_name"] + "." +
		fr.dbTblParams["tbl_wrd_pref"] + fileUniqueKey +
		"(id INT(10) UNSIGNED AUTO_INCREMENT PRIMARY KEY," +
		"word_of_file VARCHAR(50) NOT NULL," +
		"num_of_line INT(10) NOT NULL," +
		"INDEX wrd_idx (word_of_file))")
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
func (fr *filesRepo) insIntoTableStrings(stringAndKey map[string]string,  lineCounter int) {
	_, err := fr.dbConnection.Exec("INSERT INTO "+fr.dbTblParams["db_name"] + "." + fr.dbTblParams["tbl_str_pref"]+
		stringAndKey["file_key"] + "(string_of_file, num_of_line) VALUES (?, ?)",
		stringAndKey["str_of_file"], lineCounter)
	if err != nil {
		panic(err)
	}
}

/**
Insert data into table with words of current file
 */
func (fr *filesRepo) insIntoTableWords(wordAndKey map[string]string, lineCounter int) {
	_, err := fr.dbConnection.Exec("INSERT INTO " + fr.dbTblParams["db_name"] + "." + fr.dbTblParams["tbl_wrd_pref"] +
		wordAndKey["file_key"] + "(word_of_file, num_of_line) VALUES (?, ?)",
		wordAndKey["wrd_of_file"], lineCounter)
	if err != nil {
		panic(err)
	}
}

func (fr *filesRepo) getFilesInfoAsObj(fileUniqueKey string) File {
	_, err
}
