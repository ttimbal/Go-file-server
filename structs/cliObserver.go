package structs

type CliObserver interface {
	OnEntry(options []string)
	Identifier() string
}
