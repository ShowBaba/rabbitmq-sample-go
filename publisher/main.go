package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/go-rabbitmq-sample/shared"
)

var (
	ctx         = context.Background()
	err         error
	messageChan *amqp.Channel
)

type App struct {
	Router *mux.Router
}

func (a *App) Initialize() {
	a.Router = mux.NewRouter()
	a.setRouters()
}

func (a *App) setRouters() {
	a.Post("/message", SendMessage)
}

func (a *App) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("Post")
}

type MessageIn struct {
	Message string `json:"message"`
}

func SendMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	var input MessageIn
	if body, err := ioutil.ReadAll(r.Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(shared.WriteError("Invalid body: %s", err))
		return
	} else if err := json.Unmarshal(body, &input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(shared.WriteError("Invalid body: %s", err))
		return
	}
	validate := validator.New()
	err := validate.Struct(input)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(shared.WriteError(validationErrors.Error()))
		return
	}
	buf := new(bytes.Buffer)
	b, err := json.Marshal(input)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(shared.WriteError("%s", err))
		return
	}
	if err = binary.Write(buf, binary.BigEndian, &b); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(shared.WriteError("%s", err))
		return
	}
	// create a message to publish
	message := amqp.Publishing{
		ContentType: "application/json",
		Body:        []byte(b),
	}
	// publish message to queue
	if err = messageChan.PublishWithContext(
		ctx,
		"",
		shared.SERVICE_ONE,
		false,
		false,
		message,
	); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(shared.WriteError("%s", err))
		return
	}
	w.Write(shared.WriteInfo(fmt.Sprintf("Message \"%s\" have been sent.", input.Message)))
}

func (a *App) Run(port string) {
	log.Fatal(
		http.ListenAndServe(
			port,
			handlers.CORS(
				handlers.AllowCredentials(),
				handlers.AllowedMethods([]string{"POST"}),
				handlers.AllowedHeaders([]string{"Authorization", "Content-Type"}),
				handlers.MaxAge(3600),
			)(a.Router),
		),
	)
}

func main() {
	// setup rabbitmq
	connection, err := amqp.Dial(shared.RABBITMQ_SERVER_URL)
	if err != nil {
		panic(err)
	}
	defer connection.Close()
	// open a channel to the instance over the created connection
	messageChan, err = connection.Channel()
	if err != nil {
		panic(err)
	}
	defer messageChan.Close()
	//declare queue(s)
	_, err = messageChan.QueueDeclare(
		shared.SERVICE_ONE,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}
	app := App{}
	port := ":3000"
	app.Initialize()
	log.Printf("starting server on port: %s", port)
	app.Run(port)
}
