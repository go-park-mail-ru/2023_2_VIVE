package main

import (
	"HnH/internal/repository"
	"database/sql"

	// "encoding/json"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/stdlib"
)

func main() {
	urlExample := "postgres://vive_admin:123@localhost:8054/hnh"
	// conn, err := pgx.Connect(context.Background(), urlExample)
	conn, err := sql.Open("pgx", urlExample)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	vacModel := repository.NewPsqlVacancyRepository(conn)

	vacancies, err := vacModel.GetAllVacancies()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(vacancies)




}
