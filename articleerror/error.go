package articleerror

import (
	"errors"
	"fmt"
)

// ArticleError error structure for article
type ArticleError struct {
	Code int
	Err error
}

// NewArticleError creates new ArticleError
func NewArticleError(code int, message string) error {
	return &ArticleError{
		Code: code,
		Err: errors.New(message),
	}
}

//ArticleDoesNotExistError returns not found error
func ArticleDoesNotExistError(id int64) error {
	return &ArticleError{
		Code: 404,
		Err: fmt.Errorf("Article with ID %d does not exist", id),
	}
}

// UnknownError returns an error representation of unknown error
func UnknownError(e error) error {
	return &ArticleError{
		Code: 500,
		Err: e,
	}
}

// InputError returns an error when input are invalid
func InputError(message string) error {
	return &ArticleError{
		Code: 400,
		Err: errors.New(message),
	}
}

func (e *ArticleError) Error() string {
	return e.Err.Error()
}

// IsArticleEmpty tells whether the article is empty
func (e *ArticleError) IsArticleEmpty() bool {
	return e.Code == 404
}

// IsUnknownError tells whether error is unknown
func (e *ArticleError) IsUnknownError() bool {
	return e.Code == 500
}

// IsInputError tells whether the error is related to invalid input
func (e *ArticleError) IsInputError() bool {
	return e.Code == 400
}