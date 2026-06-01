package migration

import (
	"log"

	"github.com/farzadamr/go-backend-api/internal/domain"
	"github.com/farzadamr/go-backend-api/internal/infra/database"
	"gorm.io/gorm"
)

func Up_1() {
	database := database.GetDb()

	createTables(database)
	//createDefaultInformation(database)
}

func Down_1() {

}

func createTables(database *gorm.DB) {
	tables := []interface{}{}
	//User
	tables = addNewTable(database, domain.Item{}, tables)

	err := database.Migrator().CreateTable(tables...)
	if err != nil {
		log.Printf("Error creating tables: %v", err)
	}
	log.Println("Tables created successfully")
}

func addNewTable(database *gorm.DB, model interface{}, tables []interface{}) []interface{} {
	if !database.Migrator().HasTable(model) {
		tables = append(tables, model)
	}
	return tables
}
