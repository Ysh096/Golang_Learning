package mydict

import "errors"

// Dictionary type
type Dictionary map[string]string // map[keytype]valuetype

var (
	errNotFound   = errors.New("Not Found")
	errCantUpdate = errors.New("Cant update non-existing word")
	errWordExists = errors.New("That word already exists")
	errCantDelete = errors.New("Cant delete non-existing word")
)

// Search for a word
func (d Dictionary) Search(word string) (string, error) {
	value, exists := d[word] // value는 찾은 string값, exists는 존재 여부(boolean)
	if exists {
		return value, nil // zero value error nil
	}
	return "", errNotFound
}

// func (d Dictionary) Add(word, def string) error {
// 	// 아직 word가 dictionary에 없으면 추가 가능
// 	_, err := d.Search(word) // word가 이미 d에 있는지 파악
// 	if err == errNotFound {  // NotFound error는 word가 d에 없다는 뜻
// 		d[word] = def // 추가
// 	} else if err == nil { // word가 이미 d에 있으면
// 		return errWordExists // error를 return
// 	}
// 	return nil
// }

// Add a word to the dictionary
func (d Dictionary) Add(word, def string) error {
	_, err := d.Search(word)
	switch err {
	case errNotFound:
		d[word] = def
	case nil:
		return errWordExists
	}
	return nil
}

// Update a word
func (d Dictionary) Update(word, definition string) error {
	_, err := d.Search(word)
	switch err {
	case nil:
		d[word] = definition
	case errNotFound:
		return errCantUpdate
	}
	return nil
}

// Delete a word
func (d Dictionary) Delete(word string) error {
	// map document 참고
	_, err := d.Search(word)
	switch err {
	case nil:
		delete(d, word)
	case errNotFound:
		return errCantDelete
	}
	return nil
}
