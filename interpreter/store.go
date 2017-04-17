package interpreter

type Store interface {
	Save(*Scope) error
}
