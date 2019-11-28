package documents

import (
    "database/sql"
    "sync"
    "errors"
    "b2/manager"
)

type DocManager struct {
    db *sql.DB
    documents docMap
}

type docMap struct {
    sync.RWMutex
    m map[uint64]*Document
}

func (dm *DocManager) Initalize (db *sql.DB) error {
    dm.db = db
    dm.documents.m = make(map[uint64]*Document)
    return nil
}

func (dm *DocManager) Get(did uint64) (manager.Thing, error) {
    dm.documents.RLock()
    if document, ok := dm.documents.m[did]; ok {
        dm.documents.RUnlock()
        return document, nil
    }
    dm.documents.RUnlock()
    document, err := loadDocument(did, dm.db)
    if (err != nil ) {
        return nil, err
    }
    err = loadExpenses(document, dm.db)
    if err == nil && document != nil {
        dm.documents.Lock()
        defer dm.documents.Unlock()
        // check someone hasn't already inserted it while we were creating it
        if  newDoc, ok := dm.documents.m[did]; ok {
            return newDoc, nil
        }
        dm.documents.m[did] = document
    }
    return document, err
}

func (dm *DocManager) GetMultiple(from, to string) ([]manager.Thing, error) {
    // create empty array so we return [] not null
    documents := []manager.Thing{}
    dids, err := findDocuments(from, to, dm.db)
    for _, did := range dids {
        document, err := dm.Get(did)
        if (err == nil ) {
            documents = append (documents, document)
        }
    }

    return documents, err
}

func (dm *DocManager) Save(doc manager.Thing) error {
    document, ok := doc.(*Document)
    if !ok {
        return errors.New("Non document passed to function")
    }
    oldDoc, err := dm.Get(document.ID)
    if err != nil {
        if err.Error() == "404" {
            err := createDocument(document, dm.db)
            if err != nil && document.ID > 0 {
                dm.documents.Lock();
                defer dm.documents.Unlock()
                dm.documents.m[document.ID] = document
            }
            return err
        }
        return errors.New("Error loading existing document")
    } else if document == oldDoc {
        return updateExpenes(document, dm.db)
    } else {
        return errors.New("Document pointer different to one in manager")
    }
}

func (dm *DocManager) Overwrite(doc manager.Thing) (manager.Thing, error) {
    document, ok := doc.(*Document)
    if !ok {
        return nil, errors.New("Non document passed to function")
    }
    oldDocument, err := dm.Get(document.ID)
    if err != nil {
        return nil, errors.New("Error loading existing document")
    }
    oldDoc, ok := oldDocument.(*Document)
    if !ok {
        return nil, errors.New("Non document passed to function")
    }
    document.RLock()
    oldDoc.Lock()
    oldDoc.Filename = document.Filename
    oldDoc.Deleted = document.Deleted
    oldDoc.Date = document.Date
    oldDoc.Text = document.Text
    document.RUnlock()
    oldDoc.Unlock()
    return oldDoc, dm.Save(oldDoc)
}

