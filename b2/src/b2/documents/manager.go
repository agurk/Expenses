package documents

import (
    "database/sql"
    "sync"
    "errors"
)

type DocManager struct {
    db *sql.DB
    documents docMap
}

type docMap struct {
    sync.RWMutex
    m map[uint64]*Document
}

func (manager *DocManager) Initalize (db *sql.DB) error {
    manager.db = db
    manager.documents.m = make(map[uint64]*Document)
    return nil
}

func (manager *DocManager) GetDocument(did uint64) (*Document, error) {
    manager.documents.RLock()
    if document, ok := manager.documents.m[did]; ok {
        manager.documents.RUnlock()
        return document, nil
    }
    manager.documents.RUnlock()
    document, err := loadDocument(did, manager.db)
    if (err != nil ) {
        return nil, err
    }
    err = loadExpenses(document, manager.db)
    if err == nil && document != nil {
        manager.documents.Lock()
        defer manager.documents.Unlock()
        // check someone hasn't already inserted it while we were creating it
        if  newDoc, ok := manager.documents.m[did]; ok {
            return newDoc, nil
        }
        manager.documents.m[did] = document
    }
    return document, err
}

func (manager *DocManager) GetDocuments(from, to string) ([]*Document, error) {
    // create empty array so we return [] not null
    documents := []*Document{}
    dids, err := findDocuments(from, to, manager.db)
    for _, did := range dids {
        document, err := manager.GetDocument(did)
        if (err == nil ) {
            documents = append (documents, document)
        }
    }

    return documents, err
}

func (manager *DocManager) SaveDocument(document *Document) error {
    oldDoc, err := manager.GetDocument(document.ID)
    if err != nil {
        if err.Error() == "404" {
            err := createDocument(document, manager.db)
            if err != nil && document.ID > 0 {
                manager.documents.Lock();
                defer manager.documents.Unlock()
                manager.documents.m[document.ID] = document
            }
            return err
        }
        return errors.New("Error loading existing document")
    } else if document == oldDoc {
        return updateExpenes(document, manager.db)
    } else {
        return errors.New("Document pointer different to one in manager")
    }
}

func (manager *DocManager) OverwriteDocument(document *Document) (*Document, error) {
    oldDoc, err := manager.GetDocument(document.ID)
    if err != nil {
        return nil, errors.New("Error loading existing document")
    }
    document.RLock()
    oldDoc.Lock()
    oldDoc.Filename = document.Filename
    oldDoc.Deleted = document.Deleted
    oldDoc.Date = document.Date
    oldDoc.Text = document.Text
    document.RUnlock()
    oldDoc.Unlock()
    return oldDoc, manager.SaveDocument(oldDoc)
}

