package main

import (
	"errors"
	"fmt"
	"net/http"

	"bookstore.MirasKabykenov/internal/data"
)

func (app *application) addMovieInCart(w http.ResponseWriter, r *http.Request) {

	var input struct {
		BookID int32 `json:"bookID"`
	}

	// if there is error with decoding, we are sending corresponding message
	err := app.readJSON(w, r, &input) //non-nil pointer as the target decode destination
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
	}

	cart := &data.Cart{
		BookID: int64(input.BookID),
	}

	err = app.models.Cart.Ins(cart)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/cart/%d", cart.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"cart": cart}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	// // Dump the contents of the input struct in a HTTP response.
	// fmt.Fprintf(w, "%+v\n", input) //+v here is adding the field name of a value // https://pkg.go.dev/fmt
}

func (app *application) addBookToCart(w http.ResponseWriter, r *http.Request) {
	// Declare an anonymous struct to hold the information that we expect to be in the HTTP request body.
	var input struct {
		BookID int `json:"bookID"`
	}

	// Decode the request body into the input struct.
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	// Create a new CartItem with the given book ID.
	cart := &data.Cart{
		BookID: int64(input.BookID),
	}

	// Add the cart item to the cart.
	err = app.models.Cart.Ins(cart)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Send a JSON response with the added cart item.
	err = app.writeJSON(w, http.StatusCreated, envelope{"cart": cart}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) removeFrom(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Cart.RemoveFromCart(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "product sucsecfully removed from your cart"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) orderMovieFromCart(w http.ResponseWriter, r *http.Request) {
	err := app.models.Cart.Order()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "your order sucsecfully done"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
