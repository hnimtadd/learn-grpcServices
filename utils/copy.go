package utils

import "github.com/jinzhu/copier"

func DeepCopy[T any](i *T) (*T, error) {
	other := new(T)
	err := copier.Copy(other, i)
	if err != nil {
		return nil, err
	}
	return other, nil
}
