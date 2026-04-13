package keygenerator

type (
	KeyGenerator interface {
		Generate(identifier, key string) (string, error)
		Parse(message, key string) (string, error)
		GenerateBase64(message string) (string, error)
		DecryptBase64(base64 string) (string, error)
	}
)
