package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	paths := getMainFilesInfo("/Go_parser_core/texts/")
	files := initFileObjects(paths)

	db, err := sql.Open("mysql", "root@tcp(127.0.0.1)/")
	defer db.Close()
	if err != nil {
		panic(err)
	}

	dbTableParams := map[string]string {
		"db_name": "go_parser_core",
		"tbl_idx": "already_indexed_files",
		"tbl_str_pref": "strings_of__",
		"tbl_wrd_pref": "words_of__"}

	fr := filesRepo{dbTableParams, db}
	idx := indexing{&files, &fr}
	idx.filesRepo.initFilesRepo()
	idx.initFilesInfo()
	lsa := latentSemanticAnalysis{&files, &fr, &idx}
	lsa.setFrequencyMatrix(true)
}