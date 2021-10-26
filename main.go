package main

import (
	"context"
	"errors"
	"myapp/pkg/db"

	"github.com/kataras/iris/v12"
	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	db.ConnectMongoDB()
	app := iris.New()

	booksAPI := app.Party("/books")
	{
		booksAPI.Use(iris.Compression)

		// GET: http://localhost:8080/books
		booksAPI.Get("/", list)
		// POST: http://localhost:8080/books
		booksAPI.Post("/", create)
	}

	app.Listen(":8080")
}

// Book example.
type Book struct {
	Title string `json:"title"`
}

func reportError(ctx iris.Context, e error, text string) {
	ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
		Title(text).DetailErr(e))
}

func list(ctx iris.Context) {
	c_books := db.GetMongoDBClient().Database("library").Collection("books")

	filter := bson.D{}
	cursor, e := c_books.Find(context.TODO(), filter)
	if e != nil {
		reportError(ctx, e, "Find failure")
		return
	}
	var books []Book
	for cursor.Next(context.TODO()) {
		var res Book
		cursor.Decode(&res)
		books = append(books, res)
	}

	ctx.JSON(books)
	// TIP: negotiate the response between server's prioritizes
	// and client's requirements, instead of ctx.JSON:
	// ctx.Negotiation().JSON().MsgPack().Protobuf()
	// ctx.Negotiate(books)
}

func create(ctx iris.Context) {
	var b Book
	e := ctx.ReadJSON(&b, iris.JSONReader{true, false, false})

	// TIP: use ctx.ReadBody(&b) to bind
	// any type of incoming data instead.
	if e != nil {
		reportError(ctx, e, "Book parsing failure")
		return
	}

	if b.Title == "" {
		reportError(ctx, errors.New(""), "No title")
		return
	}
	c_books := db.GetMongoDBClient().Database("library").Collection("books")
	_, e = c_books.InsertOne(context.TODO(), b)
	if e != nil {
		reportError(ctx, e, "Book insert failure")
		return
	}
	println("Received Book: " + b.Title)

	ctx.StatusCode(iris.StatusCreated)
}
