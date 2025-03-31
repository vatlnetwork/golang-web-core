package domain

import "github.com/google/uuid"

type TestObject struct {
	Id      string `json:"id"`
	Number  int    `json:"number"`
	Boolean bool   `json:"boolean"`
}

func NewTestObject(number int, boolean bool) TestObject {
	return TestObject{
		Id:      uuid.NewString(),
		Number:  number,
		Boolean: boolean,
	}
}
