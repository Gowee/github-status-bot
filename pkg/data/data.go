package data

type Database struct {
	storage Storage
}

func NewDBFromFilePath(filePath string) Database {
	return Database{Storage{FilePath: filePath}}
}

// var storage = Storage{
// 	FilePath: "./data.json", // TODO: specifying via options
// }

func (db *Database) EnsureInitialized() error {
	return db.storage.EnsureInitialized(EmtpyRootDocument)
}

func (db *Database) Load() (*RootDocument, error) {
	document := RootDocument{}
	err := db.storage.Load(&document)
	if err != nil {
		return nil, err
	}
	return &document, nil
}

func (db *Database) Store(document *RootDocument) error {
	err := db.storage.Store(document)
	if err != nil {
		return err
	}
	return nil
}

// WTF: why manually use reference type here and there?
