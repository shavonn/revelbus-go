package forms

type FAQForm struct {
	ID       string
	Question string
	Answer   string
	Category string
	Active   bool
	Order    string
	Errors   map[string]string
}

func (f *FAQForm) Valid() bool {
	v := newValidator()

	v.Required("Question", f.Question)
	v.Required("Answer", f.Answer)

	f.Errors = v.Errors
	return len(f.Errors) == 0
}
