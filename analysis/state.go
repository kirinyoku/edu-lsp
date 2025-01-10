package analysis

type State struct {
	Documents map[string]string
}

func NewState() State {
	return State{
		Documents: make(map[string]string),
	}
}

func (s *State) OpenDocument(document, text string) {
	s.Documents[document] = text
}

func (s *State) UpdateDocument(document, text string) {
	s.Documents[document] = text
}
