package documents

import (
	"b2/components/managed/docexmappings"
	"b2/errors"
	"b2/manager"
	"sync"
)

// Document represents a document in this system, typically a receipt for an expense
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

// Type returns a string description of a document type
func (doc *Document) Type() string {
	return "document"
}

// GetID returns the id for a document
func (doc *Document) GetID() uint64 {
	return doc.ID
}

// Merge is a synonym for Overwrite
func (doc *Document) Merge(newThing manager.Thing) error {
	return doc.Overwrite(newThing)
}

// Overwrite replaces key fields (not id) in the existing document
// with values from that passed it
func (doc *Document) Overwrite(newThing manager.Thing) error {
	document, ok := newThing.(*Document)
	if !ok {
		panic("Non document passed to overwrite function")
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

// Check returns an error if the document has been deleted i.e. you have
// a stale reference
func (doc *Document) Check() error {
	doc.RLock()
	defer doc.RUnlock()
	if doc.deleted {
		return errors.New("Document deleted", nil, "documents.Check", true)
	}
	return nil
}
