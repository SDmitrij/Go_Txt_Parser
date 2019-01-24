package main

import (
	"database/sql"
	"fmt"
)

/**
This struct describes indexing process
 */
type indexing struct {
	filesToIndex *[]File
	filesRepo *sql.DB
}

/**
Create an entry database
 */
func (idx *indexing) createEntryDatabase() {
	_, err := idx.filesRepo.Exec("CREATE DATABASE IF NOT EXISTS go_parser_core")
	if err != nil {
		panic(err)
	}
}

func (idx *indexing) createFilesInfoTable() {

}

/**
Indexing current directory files
 */
func (idx *indexing) indexing() {
	for _, files := range *idx.filesToIndex {
		fmt.Println(files)
	}
}

func (idx indexing) searching() {

}

func (idx indexing) excludeOrIncludeFiles() {

}