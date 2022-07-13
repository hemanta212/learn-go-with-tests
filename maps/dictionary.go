package main

type Dictionary map[string]string

type DictionaryErr string

const (
	ErrorNotFound  = DictionaryErr("unknown word")
	ErrorKeyExists = DictionaryErr("Key already exists")
)

func (e DictionaryErr) Error() string {
	return string(e)
}

func (d Dictionary) Search(word string) (string, error) {
	value, ok := d[word]
	if !ok {
		return "", ErrorNotFound
	}
	return value, nil
}

func (d Dictionary) Add(key, value string) error {
	_, err := d.Search(key)
	switch err {
	case ErrorNotFound:
		d[key] = value
	case nil:
		return ErrorKeyExists
	default:
		return err
	}
	return nil
}

func (d Dictionary) Update(key, value string) error {
	_, err := d.Search(key)
	switch err {
	case ErrorNotFound:
		return err
	case nil:
		d[key] = value
	default:
		return err
	}
	return nil
}

func (d Dictionary) Delete(key string) {
	delete(d, key)
}
