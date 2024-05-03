package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/andreibuiciuc/qr"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type EncodingPayload struct {
	Input           string `json:"input"`
	CorrectionLevel string `json:"correction_level"`
}

type EncodedResponse struct {
	EncodedMatrix []string `json:"encoded_matrix"`
}

func newResponse() *EncodedResponse {
	return &EncodedResponse{
		EncodedMatrix: nil,
	}
}

func encode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var encodingPayload EncodingPayload

	json.NewDecoder(r.Body).Decode(&encodingPayload)
	encoded, _ := qr.New().Encode(encodingPayload.Input, rune(encodingPayload.CorrectionLevel[0]), "test.png")
	encodedResponse := newResponse()
	encodedMatrix := encoded.GetMatrix()

	for row := range encodedMatrix {
		rowValues := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(encodedMatrix[row])), " "), "[]")
		encodedResponse.EncodedMatrix = append(encodedResponse.EncodedMatrix, rowValues)
	}

	json.NewEncoder(w).Encode(encodedResponse)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/encode", encode).Methods(http.MethodPost)

	corsOpts := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "Access-Control-Allow-Origin"},
		AllowedMethods:   []string{"POST"},
		AllowCredentials: true,
		Debug:            true,
	})

	handler := corsOpts.Handler(r)
	http.ListenAndServe(":8000", handler)
}
