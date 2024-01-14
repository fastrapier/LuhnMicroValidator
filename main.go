package main

import (
	"encoding/base64"
	"encoding/json"
	"github.com/EClaesson/go-luhn"
	"log"
	"net/http"
	"time"
)

type Message struct {
	Status  bool
	Message string
	Time    time.Time
}

type RequestMessage struct {
	CreditCard string
}

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	log.Println(w.Header())

	message, err := json.Marshal(Message{
		Status:  true,
		Message: "Hello from my App",
		Time:    time.Now(),
	})

	if err != nil {
		log.Fatal(err)
	}

	w.Write(message)
}

func validateCreditCard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)

		message := Message{
			Status:  false,
			Message: "Method not allow",
			Time:    time.Now(),
		}

		jsonResponse, err := json.Marshal(message)

		if err != nil {
			log.Fatal("Unable to encode JSON")
		}

		http.Error(w, string(jsonResponse), 405)

		return
	}

	if !r.URL.Query().Has("payload") {
		message := Message{
			Status:  false,
			Message: "Query param payload is not exists",
			Time:    time.Now(),
		}
		jsonResponse, err := json.Marshal(message)

		if err != nil {
			log.Fatal("Unable to encode JSON")
		}

		http.Error(w, string(jsonResponse), 405)
	}

	payload := r.URL.Query().Get("payload")

	requestBase64DecodedJson, decodeRequestBase64Error := base64.StdEncoding.DecodeString(payload)

	if decodeRequestBase64Error != nil {
		decodeRequestBase64Message, decodeRequestBase64MessageError := json.Marshal(Message{
			Status:  false,
			Message: decodeRequestBase64Error.Error(),
			Time:    time.Now(),
		})

		if decodeRequestBase64MessageError != nil {
			log.Fatal(decodeRequestBase64MessageError)
		}

		http.Error(w, string(decodeRequestBase64Message), 400)
	}

	var requestMessage RequestMessage

	err := json.Unmarshal(requestBase64DecodedJson, &requestMessage)

	if err != nil {

		unmarshalErrorMessage, unmarshalErrorMessageError := json.Marshal(Message{
			Status:  false,
			Message: err.Error(),
			Time:    time.Now(),
		})

		if unmarshalErrorMessageError != nil {
			log.Fatal(unmarshalErrorMessageError)
		}

		http.Error(w, string(unmarshalErrorMessage), 400)
		return
	}

	valid, err := luhn.IsValid(requestMessage.CreditCard)

	if err != nil {
		invalidCharsMessage, marshalErrorMessageError := json.Marshal(Message{
			Status:  false,
			Message: err.Error(),
			Time:    time.Now(),
		})

		if marshalErrorMessageError != nil {
			log.Fatal(marshalErrorMessageError)
		}

		http.Error(w, string(invalidCharsMessage), 400)
		return
	}

	responseMessage, responseMessageEncodeError := json.Marshal(Message{
		Status:  valid,
		Message: "Validate success!",
		Time:    time.Now(),
	})

	if responseMessageEncodeError != nil {
		log.Fatal(responseMessageEncodeError)
	}

	w.Write(responseMessage)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/validateCreditCard", validateCreditCard)

	log.Println("Запуск веб-сервера на http://127.0.0.1:4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
