package documents

import (
	"b2/components/managed/docexmappings"
	"b2/manager"
	"errors"
	"sync"
)

type Document struct {
	ID       uint64 `json:"id"`
	Filename string `json:"filename"`
	Deleted  bool   `json:"deleted"`
	Date     string `json:"date"`
	Text     string `json:"text"`
	Filesize uint64 `json:"filesize"`
	Starred  bool   `json:"starred"`
	Archived bool   `json:"archived"`
	deleted  bool   `json:-`
	sync.RWMutex
	Expenses []*docexmappings.Mapping `json:"expenses"`
}

func (doc *Document) Type() string {
	return "document"
}

func (doc *Document) GetID() uint64 {
	return doc.ID
}

func (doc *Document) Merge(newThing manager.Thing) error {
	return doc.Overwrite(newThing)
}

func (doc *Document) Overwrite(newThing manager.Thing) error {
	document, ok := newThing.(*Document)
	if !ok {
		return errors.New("Non document passed to overwrite function")
	}
	document.RLock()
	doc.Lock()
	doc.Filename = document.Filename
	doc.Deleted = document.Deleted
	doc.Date = document.Date
	doc.Text = document.Text
	doc.Starred = document.Starred
	document.RUnlock()
	doc.Unlock()
	return nil
}

func (doc *Document) Check() error {
	doc.RLock()
	defer doc.RUnlock()
	if doc.deleted {
		return errors.New("Document deleted")
	}
	return nil
}
