package model

import "io"

type Image struct {
	File        io.Reader
	Name        string
	Size        int64
	ContentType string
	BucketName  string
}
