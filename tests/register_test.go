package tests

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"
)

func (s *FunctionalSuite) TestRegister() {
	s.LoadFixtures(DBDataset{
		{
			table: "user",
			row: DBDatasetRow{
				"id":       "11ef1e0b-5bc5-6257-bfaf-74563c32efde",
				"login":    "test",
				"password": "8f4268c427658cffb225bdb3b0549a36",
			},
		},
		{
			table: "account",
			row: DBDatasetRow{
				"user_id": "11ef1e0b-5bc5-6257-bfaf-74563c32efde",
			},
		},
	})
	defer func() {
		s.Require().NoError(s.CleanupFixtures(DBDataset{
			{
				table: "account",
				row:   DBDatasetRow{"user_id": "11ef1e0b-5bc5-6257-bfaf-74563c32efde"},
			},
			{
				table: "user",
				row:   DBDatasetRow{"id": "11ef1e0b-5bc5-6257-bfaf-74563c32efde"},
			},
		}))
	}()

	routes := s.app.GetRoutes(context.Background())
	srv := httptest.NewServer(routes)
	s.Run("Register: Success", func() {
		req := strings.NewReader(`
		    {
			    "Login": "mytestuser` + fmt.Sprintf("%d", time.Now().UnixMilli()) + `",
			    "Password": "123"
		    }
	    `)
		resp, err := srv.Client().Post(srv.URL+`/api/user/register`, "application/json", req)
		s.Require().NoError(err)
		defer func() {
			s.Assert().NoError(resp.Body.Close())
		}()
		authHeader := resp.Header.Get("Authorization")
		s.Require().NotEmpty(authHeader)
	})
	s.Run("Register: Already exists", func() {
		req := strings.NewReader(`
		    {
			    "Login": "test",
			    "Password": "123"
		    }
	    `)
		resp, err := srv.Client().Post(srv.URL+`/api/user/register`, "application/json", req)
		s.Require().NoError(err)
		defer func() {
			s.Assert().NoError(resp.Body.Close())
		}()
		s.Require().Equal(http.StatusBadRequest, resp.StatusCode)
		s.Require().Empty(resp.Header.Get("Authorization"))
	})
}
