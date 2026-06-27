package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

    "github.com/joho/godotenv"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)


type User struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type App struct {
    db *pgxpool.Pool // connection pool
}

// var users = make(map[string]string)

func main() {
    if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, relying on system environment variables")
	}

    dsn := os.Getenv("DATABASE_URL") //with export
	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("could not ping database: %v", err)
	}

	app := &App{db: pool}

    http.HandleFunc("/register", app.registerHandler)
    http.HandleFunc("/login", app.loginHandler)
    log.Println("Server started on port 8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func (a *App) registerHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var user User
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    _, err = a.db.Exec(r.Context(),
		`INSERT INTO users (username, password_hash, email )
         VALUES ($1, $2, $1)`,
		user.Username, string(hashedPassword),
	)
	if err != nil {
		http.Error(w, "could not create user: "+err.Error(), http.StatusConflict)
		return
	}
    w.WriteHeader(http.StatusCreated)
}

func (a *App) loginHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var user User
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    var hashedPassword string
    err = a.db.QueryRow(r.Context(),
    `SELECT password_hash FROM users WHERE username = $1`,
    user.Username).Scan(&hashedPassword)

    if errors.Is(err, pgx.ErrNoRows){
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }
    if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
    if bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password)) != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}
    w.WriteHeader(http.StatusOK)
}