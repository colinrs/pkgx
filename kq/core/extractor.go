package core

type Extractor interface {
	Unmarshal(*InputMessage) (Message, error)
}
