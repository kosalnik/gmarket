package tests

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/kosalnik/gmarket/internal/application"
	"github.com/stretchr/testify/suite"

	"github.com/kosalnik/gmarket/internal/config"
	"github.com/kosalnik/gmarket/internal/infra/logger"
	"github.com/kosalnik/gmarket/internal/infra/postgres"
)

type FunctionalSuite struct {
	suite.Suite

	db  *sql.DB
	app *application.Application
}

func (s *FunctionalSuite) SetupSuite() {
	cfg := config.NewConfig()
	var err error
	s.db, err = postgres.NewDB(context.Background(), cfg.Database)
	s.Require().NoError(err)
	s.app = application.New(cfg)
	s.app.InitServices(context.Background())
}

func TestFunctionalSuite(t *testing.T) {
	suite.Run(t, new(FunctionalSuite))
}

type DBDataset []struct {
	table string
	row   DBDatasetRow
}
type DBDatasetRow map[string]any

func (s *FunctionalSuite) LoadFixtures(fixtures DBDataset) (err error) {
	tx, err := s.db.Begin()
	s.Require().NoError(err)
	defer func() {
		if err != nil {
			if er := tx.Rollback(); er != nil {
				err = fmt.Errorf("failed to rollback transaction %w: %w", er, err)
			}
			s.Fail(err.Error())
			return
		}
		if er := recover(); er != nil {
			err = tx.Rollback()
			s.Fail("fixture panic")
			return
		}
	}()
	for _, d := range fixtures {
		tbl := d.table
		row := d.row
		i := 1
		keys := []string{}
		placeholders := []string{}
		values := []any{}
		for k, v := range row {
			placeholders = append(placeholders, fmt.Sprintf("$%d", i))
			keys = append(keys, fmt.Sprintf(`%q`, k))
			values = append(values, v)
			i++
		}
		q := fmt.Sprintf(
			`INSERT INTO %q (%s) VALUES (%s)`,
			tbl, strings.Join(keys, ","), strings.Join(placeholders, ","),
		)
		fmt.Println(q)
		fmt.Println(values)
		if _, err := tx.Exec(q, values...); err != nil {
			logger.Error("failed load fixtures", "err", err)
			return err
		}
	}
	return tx.Commit()
}

func (s *FunctionalSuite) CleanupFixtures(fixtures DBDataset) error {
	tx := s.db
	for _, d := range fixtures {
		tbl := d.table
		row := d.row
		i := 1
		where := []string{}
		values := []any{}
		for k, v := range row {
			where = append(where, fmt.Sprintf(`%q = $%d`, k, i))
			values = append(values, v)
			i++
		}
		q := fmt.Sprintf(
			`DELETE FROM %q WHERE %s`,
			tbl, strings.Join(where, " AND "),
		)
		if _, err := tx.Exec(q, values...); err != nil {
			logger.Error("failed cleanup fixtures", "err", err, "q", q, "v", values)
			continue
		}
	}
	return nil
}

func (s *FunctionalSuite) RequireDBContains(fixtures DBDataset) {
	tx := s.db
	for _, d := range fixtures {
		tbl := d.table
		row := d.row
		i := 1
		where := []string{}
		values := []any{}
		for k, v := range row {
			where = append(where, fmt.Sprintf(`%q = $%d`, k, i))
			values = append(values, v)
			i++
		}
		q := fmt.Sprintf(
			`SELECT true AS ok FROM %q WHERE %s LIMIT 1`,
			tbl, strings.Join(where, " AND "),
		)
		var ok bool
		res := tx.QueryRow(q, values...)
		s.Require().NoError(res.Err())
		s.Require().NoError(res.Scan(&ok))
		s.Require().True(ok)
	}
}
