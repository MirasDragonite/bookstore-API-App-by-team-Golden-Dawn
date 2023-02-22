package main

import (
	"fmt"
	"net/http"

	"bookstore.MirasKabykenov/internal/data"
)

// return a JSON response.
func (app *application) addBookInDataBase(w http.ResponseWriter, r *http.Request) {
	//Declare an anonymous struct to hold the information that we expect to be in the
	// HTTP request body (note that the field names and types in the struct are a subset
	// of the Movie struct that we created earlier). This struct will be our *target
	// decode destination*.
	var input struct {
		Title    string   `json:"title"`
		Author   string   `json:"author"`
		Year     int32    `json:"year"`
		Language string   `json:"language"`
		Price    int32    `json:"price"`
		Quantity int32    `json:"quantity"`
		Genres   []string `json:"genres"`
	}

	// if there is error with decoding, we are sending corresponding message
	err := app.readJSON(w, r, &input) //non-nil pointer as the target decode destination
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
	}

	book := &data.Book{
		Title:    input.Title,
		Author:   input.Author,
		Year:     input.Year,
		Language: input.Language,
		Price:    input.Price,
		Quantity: input.Quantity,
		Genres:   input.Genres,
	}

	err = app.models.Books.AddMovieInDBB(book)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/book/%d", book.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"book": book}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	// // Dump the contents of the input struct in a HTTP response.
	// fmt.Fprintf(w, "%+v\n", input) //+v here is adding the field name of a value // https://pkg.go.dev/fmt
}
