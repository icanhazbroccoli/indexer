package document

import (
	"io"
)

type Document struct {
	name   string
	reader io.Reader
}

func New(name string, reader io.Reader) *Document {
	return &Document{
		name:   name,
		reader: reader,
	}
}

func (d *Document) Name() string {
	return d.name
}

func (d *Document) Reader() io.Reader {
	return d.reader
}
