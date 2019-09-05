package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"

	_ "github.com/lib/pq"
	fdk "github.com/fnproject/fdk-go"
)
type Message struct {
	Action Action `json:"action"`
	Data json.RawMessage `json:"data"`
}

type Action string
type UserDb map[uint64]*UserData

func (db UserDb) Create(user *UserData) {
	var id uint64
	err := userDb.QueryRow("INSERT INTO users (name, lastname) VALUES($1, $2) RETURNING id",
			user.Name, user.LastName).Scan(&id)

	if err != nil {
		log.Fatal("Could not insert into DB ", err)
	}
	user.Id = id
	db[user.Id] = user
}

func (db UserDb) Get(uid uint64) (*UserData, error) {
	user := UserData{}
	row := userDb.QueryRow("SELECT id, name, lastname FROM users WHERE id = $1", uid)

	err := row.Scan(&user.Id, &user.Name, &user.LastName)
	return &user, err
}

func (db UserDb) Delete(uid uint64) {
	_, err := userDb.Exec("DELETE FROM users WHERE id = $1", uid)

	if err != nil {
		log.Fatal("Could not delete user ", err)
	}
}

func (db UserDb) Update(user *UserData) {
	_, err := userDb.Exec("UPDATE users SET name=$1, lastname=$2 WHERE id=$3", user.Name, user.LastName, user.Id)

	if err != nil {
		log.Fatal("Could not update user ", err)
	}
}

type UserData struct {
	Id uint64 `json:"id"`
	Name string `json:"name"`
	LastName string `json:"last_name"`
}

const (
	CREATE Action = "create"
	DELETE = "delete"
	UPDATE = "update"
	RETRIEVE = "retrieve"
	LIST = "list"
)

var memDb = UserDb{}
var userDb *sql.DB

var actionHandlers = map[Action]func(*UserData) *UserData {
	CREATE: createUser,
	DELETE: deleteUser,
	UPDATE: updateUser,
	RETRIEVE: retrieveUser,
	LIST: listUser,
}

func main() {
	fdk.Handle(fdk.HandlerFunc(myHandler))
}

func createUser(user *UserData) *UserData {
	memDb.Create(user)
	return memDb[user.Id]
}

func deleteUser(user *UserData) *UserData {
	memDb.Delete(user.Id)
	return &UserData{Id: user.Id}
}

func retrieveUser(user *UserData) *UserData {
	user, err := memDb.Get(user.Id)

	if err == sql.ErrNoRows {
		return &UserData{Name: "Error, User not found"}
	} else if err != nil {
		log.Fatal("Could not retrieve user", err)
	}
	return user
}

func updateUser(user *UserData) *UserData {
	memDb.Update(user)
	return &UserData{Name: "OK"}
}

func initDb() {
	connStr := "user=postgres dbname=test host=postgresdb password=12345 port=5432 sslmode=disable"
	db_, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	userDb = db_
}

func myHandler(ctx context.Context, in io.Reader, out io.Writer) {
	initDb()
	defer userDb.Close()

	msg := Message{}

	err := json.NewDecoder(in).Decode(&msg)
	if err != nil {
		log.Fatal("The was an error while parsing the request ", err)
	} else {
		// Get the function to handle the action
		handler := actionHandlers[msg.Action]
		userData := UserData{}
		err = json.Unmarshal(msg.Data, &userData)
		if err != nil {
			log.Fatal("The was error while parsing the input data ", err)
		} else {
			response := handler(&userData)
			err = json.NewEncoder(out).Encode(response)
			if err != nil {
				log.Fatal("The was an error encoding the response ", err)
			}
		}
	}
}
