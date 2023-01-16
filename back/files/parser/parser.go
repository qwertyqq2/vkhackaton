package parser

type Parser interface {
	Body() string

	Head() string

	Add(comm string) string
}
