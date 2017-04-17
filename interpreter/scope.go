package interpreter

func NewScope(name string, parent *Scope) *Scope {
	return &Scope{Name: name, Parent: parent, Variables: make(map[string]*Object)}
}

type Scope struct {
	Name      string
	Parent    *Scope
	Variables map[string]*Object
	Store     Store
}

func (s *Scope) Lookup(name string) *Object {
	if _, ok := s.Variables[name]; !ok {
		// Todo: Look up in parents, but create in self.
		s.Variables[name] = &Object{Name: name}
	}
	return s.Variables[name]
}

func (s *Scope) Update(name string, o *Object) (*Object, error) {
	obj := s.Lookup(name)
	obj.Type = o.Type
	obj.Value = o.Value
	return obj, nil
}
