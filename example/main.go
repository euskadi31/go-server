// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"net/http"

	"github.com/euskadi31/go-server"
	"github.com/euskadi31/go-server/response"
)

var userSchema = `{
	"properties": {
		"id": {
			"type": "number"
		},
		"name": {
			"type": "string",
			"pattern": "^[A-Za-z\\-]+$",
			"minLength": 2
		},
		"email": {
			"type": "string",
			"pattern": "^[a-zA-Z0-9_.+\\-]+@[a-zA-Z0-9\\-]+\\.[a-zA-Z0-9\\-.]+$"
		}
	},
	"required": [
		"name",
		"email"
	],
	"additionalProperties": false
}`

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	validator := server.NewValidator()

	if err := validator.AddSchemaFromJSON("user", []byte(userSchema)); err != nil {
		panic(err)
	}

	router := server.NewRouter()
	router.AddRouteFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		user := &User{}

		if err := json.NewDecoder(r.Body).Decode(user); err != nil {
			response.FailureFromError(w, http.StatusBadRequest, err)

			return
		}

		// Validate User struct by user schema
		if result := validator.Validate("user", user); !result.IsValid() {
			response.FailureFromValidator(w, result)

			return
		}

		response.JSON(w, http.StatusCreated, user)
	}).Methods(http.MethodPost)

	if err := http.ListenAndServe(":1337", router); err != nil {
		panic(err)
	}
}
