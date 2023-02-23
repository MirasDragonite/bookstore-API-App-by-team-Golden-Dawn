package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/books", app.requireForAdmin(app.addBookInDataBase))         // добавить книгу в
	router.HandlerFunc(http.MethodGet, "/v1/books/:id", app.requireActivatedUser(app.showInfoAboutBook)) ///Посмотреть определенную книгу
	router.HandlerFunc(http.MethodPatch, "/v1/books/:id", app.requireForAdmin(app.updateBookInfo))       // Обновить данные в определенном кинге
	router.HandlerFunc(http.MethodDelete, "/v1/books/:id", app.requireForAdmin(app.deleteBookFromDB))    // удалить книгу из базы данных
	router.HandlerFunc(http.MethodPost, "/v1/cart", app.addBookInCart)
	router.HandlerFunc(http.MethodDelete, "/v1/cart/:id", app.RemoveFromCart)
	router.HandlerFunc(http.MethodPost, "/v1/upload", app.uploadFile)
	router.HandlerFunc(http.MethodPost, "/v1/users", app.userRegistration)
	router.HandlerFunc(http.MethodDelete, "/v1/cart", app.orderBookFromCart)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUser)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.authentication)
	return app.recoverPanic(app.rateLimit(app.authenticate(router)))
}
