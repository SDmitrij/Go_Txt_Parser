package main

import (
	"./lsa"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

// TODO refactor to SQLite
func main() {
	paths := lsa.GetMainFilesInfo("/Go_parser_core/texts/")
	files := lsa.InitFileObjects(paths)

	db, err := sql.Open("mysql", "root@tcp(127.0.0.1)/")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	dbParams := map[string]string {
		"db_name": "go_latent_semantic_analysis",
		"tbl_idx": "already_indexed_files",
		"tbl_str_pref": "strings_of__",
		"tbl_term_pref": "terms_of__"}

	fr := lsa.FilesRepo{Params: dbParams, Connection: db}
	idx := lsa.Indexing{Files: files, Repo: fr}
	idx.InvokeIndexing()
	analysis := lsa.LatentSemanticAnalysis{ Indexer: idx, Fm: &lsa.FrequencyMatrix{}}
	analysis.InvokeLsa()
}