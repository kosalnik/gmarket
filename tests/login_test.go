package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
)

func (s *FunctionalSuite) TestLogin() {
	s.LoadFixtures(DBDataset{
		{
			table: "user",
			row: DBDatasetRow{
				"id":       "11ef1e0b-5bc5-6257-bfaf-74563c32efde",
				"login":    "test",
				"password": "202cb962ac59075b964b07152d234b70",
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
			    "Login": "test",
			    "Password": "123"
		    }
	    `)
		resp, err := srv.Client().Post(srv.URL+`/api/user/login`, "application/json", req)
		s.Require().NoError(err)
		defer func() {
			s.Assert().NoError(resp.Body.Close())
		}()
		authHeader := resp.Header.Get("Authorization")
		s.Require().NotEmpty(authHeader)
	})
	s.Run("Register: Wrong password", func() {
		req := strings.NewReader(`
		    {
			    "Login": "test",
			    "Password": "asd"
		    }
	    `)
		resp, err := srv.Client().Post(srv.URL+`/api/user/login`, "application/json", req)
		s.Require().NoError(err)
		defer func() {
			s.Assert().NoError(resp.Body.Close())
		}()
		s.Require().Equal(http.StatusUnauthorized, resp.StatusCode)
		s.Require().Empty(resp.Header.Get("Authorization"))
	})
	s.Run("Register: Wrong login", func() {
		req := strings.NewReader(`
		    {
			    "Login": "test_wrong",
			    "Password": "123"
		    }
	    `)
		resp, err := srv.Client().Post(srv.URL+`/api/user/login`, "application/json", req)
		s.Require().NoError(err)
		defer func() {
			s.Assert().NoError(resp.Body.Close())
		}()
		s.Require().Equal(http.StatusUnauthorized, resp.StatusCode)
		s.Require().Empty(resp.Header.Get("Authorization"))
	})
}
