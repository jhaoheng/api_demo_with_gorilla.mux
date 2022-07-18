package modules

type EmailTemplate string

type Email struct {
	Title    string
	Body     string
	To       []string
	From     string
	Template EmailTemplate
}
