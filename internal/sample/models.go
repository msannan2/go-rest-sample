package sample

import (
	"unsafe"
)

type Article struct {
	Id     int64  `json:"id,omitempty"`                         // Id of the article.
	Views  int64  `json:"views,omitempty" validate:"required"`  // Views on the article.
	Title  string `json:"title,omitempty" validate:"required"`  // Title of the article.
	Author string `json:"author,omitempty" validate:"required"` // Author of the article.
}

type Category struct {
	Id   int64     `json:"id,omitempty"`                        // Id of the category.
	Name [150]byte `json:"views,omitempty" validate:"required"` // Name of the category.
	Type [150]byte `json:"title,omitempty" validate:"required"` // Subtype of the category.
}

const StrutSizeArticles = 186

func (a *Article) Size() uint64 {
	size := *(*uint64)(unsafe.Pointer(unsafe.Sizeof(*a)))
	return size
}
