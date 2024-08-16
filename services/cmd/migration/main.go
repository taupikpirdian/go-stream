package main

import (
	"database/sql"
	"fmt"
	"os"

	config2 "cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-pkg/config"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {

	dbConf := config2.DatabasePGSQL()
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbConf.Host, dbConf.Port, dbConf.User, dbConf.Password, dbConf.Database)

	fmt.Println(connStr)

	pgDb, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer pgDb.Close()

	err = pgDb.Ping()
	if err != nil {
		panic(err)
	}

	_, err = pgDb.Exec(fmt.Sprintf(`set search_path to "%s"`, dbConf.Schema))
	if err != nil {
		panic(err)
	}

	driver, err := postgres.WithInstance(pgDb, &postgres.Config{
		SchemaName: dbConf.Schema,
	})
	if err != nil {
		panic(err)
	}
	defer driver.Close()

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+os.Getenv("MIGRATION_DIR"),
		"postgres",
		driver,
	)
	if err != nil {
		panic(err)
	}

	err = m.Up()
	if err != nil {
		panic(err)
	}

	fmt.Println("success migrate")
	os.Exit(0)
}
