package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddr string
	store      Storage
}

func NewAPIserver(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (a *APIServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/account", makeHTTPHandleFunc(a.handleAccount))
	router.HandleFunc("/account/{id}", withJWTAuth(makeHTTPHandleFunc(a.handleAccountById), a.store))
	router.HandleFunc("/transfer", makeHTTPHandleFunc(a.handleTransfer))

	log.Println("Starting server on", a.listenAddr)
	http.ListenAndServe(a.listenAddr, router)
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodGet {
		return s.handleGetAccount(w, r)
	}
	if r.Method == http.MethodPost {
		return s.handleCreateAccount(w, r)
	}
	return fmt.Errorf("unsupported method %s", r.Method)
}

func (s *APIServer) handleAccountById(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodGet {
		return s.handleGetAccountById(w, r)
	}
	if r.Method == http.MethodDelete {
		return s.handleDeleteAccount(w, r)
	}
	if r.Method == http.MethodPatch || r.Method == http.MethodPut {
		return s.handleUpdateAccount(w, r)
	}
	return fmt.Errorf("unsupported method %s", r.Method)
}
func (s *APIServer) handleGetAccount(w http.ResponseWriter, _ *http.Request) error {
	accounts, err := s.store.GetAccounts()
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, &accounts)
}

func (s *APIServer) handleGetAccountById(w http.ResponseWriter, r *http.Request) error {
	id, err := GetId(r)
	if err != nil {
		return err
	}
	account, err := s.store.GetAccountById(id)
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleUpdateAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := GetId(r)
	if err != nil {
		return err
	}
	account, err := s.store.GetAccountById(id)
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	createAccountReq := new(createAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(createAccountReq); err != nil {
		return err
	}
	defer r.Body.Close()
	account := NewAccount(createAccountReq.FirstName, createAccountReq.LastName)
	accountRecord, err := s.store.CreateAccount(account)
	if err != nil {
		return err
	}

	tokenString, err := createJWT(account)
	if err != nil {
		return err
	}
	fmt.Printf("JWT string token: %s\n", tokenString)
	return WriteJSON(w, http.StatusOK, accountRecord)
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := GetId(r)
	if err != nil {
		return err
	}
	if err = s.store.DeleteAccount(id); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, map[string]string{"message": fmt.Sprintf("the Account with id: %d deleted", id)})
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	transFerReq := new(TransFerRequest)
	if err := json.NewDecoder(r.Body).Decode(transFerReq); err != nil {
		return err
	}
	defer r.Body.Close()
	return WriteJSON(w, http.StatusOK, transFerReq)
}

func WriteJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type APIError struct {
	Error string `json:"error"`
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusInternalServerError, APIError{Error: err.Error()})
			return
		}
	}
}

func GetId(r *http.Request) (int, error) {
	vars := mux.Vars(r)
	id := vars["id"]
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return idInt, fmt.Errorf("invalid id given %s", id)
	}
	return idInt, nil
}

func permissionDenied(w http.ResponseWriter) {
	WriteJSON(w, http.StatusForbidden, APIError{Error: "permission denied"})
}

func withJWTAuth(handlerFunc http.HandlerFunc, s Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("calling JWT auth middleware")
		tokenString := r.Header.Get("x-jwt-token")
		token, err := validateJWT(tokenString)
		if err != nil {
			permissionDenied(w)
			return
		}
		if !token.Valid {
			permissionDenied(w)
			return
		}

		userId, err := GetId(r)
		if err != nil {
			permissionDenied(w)
			return
		}

		account, err := s.GetAccountById(userId)
		if err != nil {
			permissionDenied(w)
			return
		}
		claims := token.Claims.(jwt.MapClaims)

		if account.Number != int64(claims["accountNumber"].(float64)) {
			permissionDenied(w)
			return
		}
		handlerFunc(w, r)
	}
}

// const JWTSecret = "secret"
func validateJWT(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWTSecret")
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secret), nil
	})
}

func createJWT(account *Account) (string, error) {
	secret := os.Getenv("JWTSecret")
	claims := &jwt.MapClaims{
		"ExpiresAt":     jwt.NewNumericDate(time.Unix(1516239022, 0)),
		"accountNumber": account.Number,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

//eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFeHBpcmVzQXQiOjE1MTYyMzkwMjIsImFjY291bnROdW1iZXIiOjY4NDIzOH0.99gpit7uS1UcO8TWBPiQQ3hBoP3uHgsF8PEJE8JIPWY
