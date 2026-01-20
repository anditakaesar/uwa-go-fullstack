package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/anditakaesar/uwa-go-fullstack/internal/domain"
	"github.com/anditakaesar/uwa-go-fullstack/internal/env"
	"github.com/anditakaesar/uwa-go-fullstack/internal/infra"
	"github.com/anditakaesar/uwa-go-fullstack/internal/xlog"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()

	pool, err := pgxpool.New(context.Background(), env.Values.DBUrl)
	if err != nil {
		xlog.Logger.Error(fmt.Sprintf("Unable to connect database: %v", err))
		os.Exit(1)
	}
	defer pool.Close()

	err = pool.Ping(context.Background())
	if err != nil {
		xlog.Logger.Error(fmt.Sprintf("Unable to connect database: %v", err))
		os.Exit(1)
	}

	services := infra.NewInfra(pool)
	seedUsers(ctx, services)
}

type UserSeed struct {
	Username string
	Password string
}

func readUserFromCSV(path string) ([]UserSeed, error) {
	userSeedFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer userSeedFile.Close()

	r := csv.NewReader(userSeedFile)
	r.TrimLeadingSpace = true

	_, err = r.Read() // skip header
	if err != nil {
		return nil, err
	}

	var users []UserSeed

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		users = append(users, UserSeed{
			Username: record[0],
			Password: record[1],
		})
	}

	return users, nil
}

func seedUsers(ctx context.Context, services *infra.Services) {
	users, err := readUserFromCSV("./cmd/seed/users.csv")
	if err != nil {
		xlog.Logger.Error(fmt.Sprintf("error reading user seed file: %v", err))
		return
	}

	var seedErrs []string
	for _, u := range users {
		_, seedErr := services.UserService.CreateUser(ctx, domain.User{
			Username: u.Username,
			Password: u.Password,
		})
		if seedErr != nil {
			seedErrs = append(seedErrs, fmt.Sprintf("error processing %s: %v", u.Username, seedErr))
		}
	}

	if len(seedErrs) > 0 {
		xlog.Logger.Error(strings.Join(seedErrs, "\n"))
	}
}
