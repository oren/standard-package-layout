package myapp

type User struct {
	ID      int
	Name    string
	Address Address
}

type UserService interface {
	User(id int) (*User, error)
	Users() ([]*User, error)
	CreateUser(u *User) error
	DeleteUser(id int) error
}

////

package postgres

import (
	"database/sql"

	"github.com/benbjohnson/myapp"
	_ "github.com/lib/pq"
)

// UserService represents a PostgreSQL implementation of myapp.UserService.
type UserService struct {
	DB *sql.DB
}

// User returns a user for a given id.
func (s *UserService) User(id int) (*myapp.User, error) {
	var u myapp.User
	row := db.QueryRow(`SELECT id, name FROM users WHERE id = $1`, id)
	if row.Scan(&u.ID, &u.Name); err != nil {
		return nil, err
	}
	return &u, nil
}

// implement remaining myapp.UserService interface...


////

package myapp

// UserCache wraps a UserService to provide an in-memory cache.
type UserCache struct {
        cache   map[int]*User
        service UserService
}

// NewUserCache returns a new read-through cache for service.
func NewUserCache(service UserService) *UserCache {
        return &UserCache{
                cache: make(map[int]*User),
                service: service,
        }
}

// User returns a user for a given id.
// Returns the cached instance if available.
func (c *UserCache) User(id int) (*User, error) {
	// Check the local cache first.
        if u := c.cache[id]]; u != nil {
                return u, nil
        }

	// Otherwise fetch from the underlying service.
        u, err := c.service.User(id)
        if err != nil {
        	return nil, err
        } else if u != nil {
        	c.cache[id] = u
        }
        return u, err
}


////

package http

import (
        "net/http"

        "github.com/benbjohnson/myapp"
)

type Handler struct {
        UserService myapp.UserService
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
        // handle request
			}

////

package mock

import "github.com/benbjohnson/myapp"

// UserService represents a mock implementation of myapp.UserService.
type UserService struct {
        UserFn      func(id int) (*myapp.User, error)
        UserInvoked bool

        UsersFn     func() ([]*myapp.User, error)
        UsersInvoked bool

        // additional function implementations...
}

// User invokes the mock implementation and marks the function as invoked.
func (s *UserService) User(id int) (*myapp.User, error) {
        s.UserInvoked = true
        return s.UserFn(id)
}

// additional functions: Users(), CreateUser(), DeleteUser()


////

package http_test

import (
	"testing"
	"net/http"
	"net/http/httptest"

	"github.com/benbjohnson/myapp/mock"
)

func TestHandler(t *testing.T) {
	// Inject our mock into our handler.
	var us mock.UserService
	var h Handler
	h.UserService = &us

	// Mock our User() call.
	us.UserFn = func(id int) (*myapp.User, error) {
		if id != 100 {
			t.Fatalf("unexpected id: %d", id)
		}
		return &myapp.User{ID: 100, Name: "susy"}, nil
	}

	// Invoke the handler.
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/users/100", nil)
	h.ServeHTTP(w, r)

	// Validate mock.
	if !us.UserInvoked {
		t.Fatal("expected User() to be invoked")
	}
}

////


package main

import (
	"log"
	"os"

	"github.com/benbjohnson/myapp"
	"github.com/benbjohnson/myapp/postgres"
	"github.com/benbjohnson/myapp/http"
)

func main() {
	// Connect to database.
	db, err := postgres.Open(os.Getenv("DB"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create services.
	us := &postgres.UserService{DB: db}

	// Attach to HTTP handler.
	var h http.Handler
	h.UserService = us

	// start http server...
}




