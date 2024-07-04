package steam_helper

import (
	"context"
	"encoding/json"
)

func MapToStruct(m map[string]interface{}, s interface{})error{
	data, err := json.Marshal(m)
	if err != nil{
		return err
	}

	return json.Unmarshal(data, &s)
}

func StructToMap(s interface{})(map[string]interface{}, error){
	data, err := json.Marshal(s)
	if err != nil{
		return nil, err
	}

	var m map[string]interface{}

	if err := json.Unmarshal(data, &m); err != nil{
		return nil, err
	}

	return m, nil
}

type (
	Cursor[T any] struct {
		Model T
		Error error
	}
	CursorCh[T any] chan Cursor[T]
)

// CursorModel - create Cursor[T] model item
func CursorModel[T any](model T) Cursor[T] {
	return Cursor[T]{
		Model: model,
	}
}

// CursorError - create Cursor[T] error item
func CursorError[T any](err error) Cursor[T] {
	return Cursor[T]{
		Error: err,
	}
}

// Write - attempt to write data to chan
func (a CursorCh[T]) Write(c context.Context, item Cursor[T]) error {
	select {
	case <-c.Done():
		return c.Err()
	case a <- item:
		return nil
	}
}

// WriteModel - attempt to write model to chan
func (a CursorCh[T]) WriteModel(c context.Context, model T) error {
	return a.Write(c, CursorModel(model))
}

// WriteError - attempt to write error to chan
func (a CursorCh[T]) WriteError(c context.Context, err error) error {
	return a.Write(c, CursorError[T](err))
}
