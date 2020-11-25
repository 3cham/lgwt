package maps

var (
	ErrKeyNotFound          = DictionaryErr("could not find the key you provide")
	ErrKeyAlreadyExists     = DictionaryErr("key already exists")
	ErrUpdateNonExistingKey = DictionaryErr("update non existing key")
	ErrDeleteNonExistingKey = DictionaryErr("delete non existing key")
)

type DictionaryErr string

func (e DictionaryErr) Error() string {
	return string(e)
}

type Dictionary map[string]string

func (d Dictionary) Search(key string) (string, error) {
	value, ok := d[key]
	if !ok {
		return "", ErrKeyNotFound
	}
	return value, nil
}

func (d Dictionary) Add(key, value string) error {
	_, err := d.Search(key)
	if err == nil {
		return ErrKeyAlreadyExists
	}

	d[key] = value
	return nil
}

func (d Dictionary) Update(key, value string) error {
	_, err := d.Search(key)
	if err != nil {
		return ErrUpdateNonExistingKey
	}
	d[key] = value
	return nil
}

func (d Dictionary) Delete(key string) error {
	_, err := d.Search(key)
	if err != nil {
		return ErrDeleteNonExistingKey
	}

	delete(d, key)
	return nil
}
