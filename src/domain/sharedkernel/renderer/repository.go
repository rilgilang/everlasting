package renderer

type Template string

const (
	ResetPasswordEmailTemplate Template = "email/reset_password.html"
)

var Templates []Template = []Template{ResetPasswordEmailTemplate}

type HTMLRenderer interface {
	Render(template Template, values map[string]any) (result string, err error)
}
