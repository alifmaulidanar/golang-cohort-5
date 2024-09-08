package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Book struct {
	BookID string `json:"book_id"`
	Title  string `json:"title"`
	Stock  int    `json:"stock"`
	Author string `json:"author"`
}

var BookDatas = []Book{}

func CreateBook(ctx *gin.Context) {
	var newBook Book

	if err := ctx.ShouldBindJSON(&newBook); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	newBook.BookID = fmt.Sprintf("%d", len(BookDatas)+1)
	BookDatas = append(BookDatas, newBook)

	ctx.JSON(http.StatusOK, gin.H{
		"data":    newBook,
		"message": "succeed create book",
	})
}

func GetBook(ctx *gin.Context) {
	bookID := ctx.Param("bookID")
	condition := false
	var bookData Book

	for i, book := range BookDatas {
		if bookID == book.BookID {
			condition = true
			bookData = BookDatas[i]
			break
		}
	}

	if !condition {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"data":    nil,
			"message": fmt.Sprintf("book with id %v not found", bookID),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":    bookData,
		"message": "succeed get book",
	})

}

func GetAllBooks(ctx *gin.Context) {
	if len(BookDatas) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"data":    nil,
			"message": "no books found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":    BookDatas,
		"message": "succeed get all books",
	})
}

func UpdateBook(ctx *gin.Context) {
	bookID := ctx.Param("bookID")
	condition := false
	var updatedBook Book

	if err := ctx.ShouldBindJSON(&updatedBook); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	for i, book := range BookDatas {
		if bookID == book.BookID {
			condition = true
			BookDatas[i] = updatedBook
			BookDatas[i].BookID = bookID
			break
		}
	}

	if !condition {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"data":    nil,
			"message": fmt.Sprintf("book with id %v not found", bookID),
		})
		return
	}

	convertBookId, _ := strconv.Atoi(bookID)
	ctx.JSON(http.StatusOK, gin.H{
		"data":    BookDatas[convertBookId-1],
		"message": fmt.Sprintf("Book with id %v has been successfully updated", bookID),
	})
}

func DeleteBook(ctx *gin.Context) {
	bookID := ctx.Param("bookID")
	condition := false
	var bookIndex int

	for i, book := range BookDatas {
		if bookID == book.BookID {
			condition = true
			bookIndex = i
			break
		}
	}

	if !condition {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"data":    nil,
			"message": fmt.Sprintf("book with id %v not found", bookID),
		})
		return
	}

	copy(BookDatas[bookIndex:], BookDatas[bookIndex+1:])
	BookDatas[len(BookDatas)-1] = Book{}
	BookDatas = BookDatas[:len(BookDatas)-1]

	ctx.JSON(http.StatusOK, gin.H{
		"data":    nil,
		"message": fmt.Sprintf("Book with id %v has been successfully deleted", bookID),
	})
}
