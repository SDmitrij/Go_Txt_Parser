package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	paths := getMainFilesInfo("/Go_parser_core/texts/")
	files := initFileObjects(paths)

	for _, file := range files {
		fmt.Printf("File key: %s, file hash: %s, file path: %s, file size: %d\n",
			file.fileUniqueKey, file.fileHash, file.filePath, file.fileSize)
	}

	db, err := sql.Open("mysql", "root@tcp(127.0.0.1)/")
	if err != nil {
		panic(err)
	}
	idx := indexing{&files, db}
	idx.indexing()
	idx.createEntryDatabase()

}