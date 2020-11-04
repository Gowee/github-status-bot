package data

var storage = Storage{
	FilePath: "./data.json", // TODO: specifying via options
}

func EnsureInitialized() error {
	return storage.EnsureInitialized(EMTPY_ROOT_DOCUMENT)
}

func Load() (*RootDocument, error) {
	document := RootDocument{}
	err := storage.Load(&document)

	if err != nil {
		return nil, err
	}
	return &document, nil
}

func Store(document *RootDocument) error {
	err := storage.Store(document)
	if err != nil {
		return err
	}
	return nil
}

// WTF: why manually use reference type here and there?
