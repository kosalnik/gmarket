package tests

import (
	"context"
	"database/sql"
	"testing"

	"github.com/kosalnik/gmarket/internal/config"
	"github.com/kosalnik/gmarket/internal/infra/postgres"
	"github.com/stretchr/testify/suite"
)

type FunctionalSuite struct {
	suite.Suite

	db *sql.DB
}

func (s *FunctionalSuite) SetupSuite() {
	cfg := config.NewConfig()
	var err error
	s.db, err = postgres.NewDB(context.Background(), cfg.Database)
	s.Require().NoError(err)
}

func TestFunctionalSuite(t *testing.T) {
	suite.Run(t, new(FunctionalSuite))
}
