package uuid

import "github.com/google/uuid"

type Generator struct{}

func NewGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) NewString() string {
	return uuid.NewString()
}
