package photo

type (
	Photo struct {
		FileType    string `json:"-"`
		ContentType string `json:"-"`
		Byte        []byte `json:"-"`
		Size        int64  `json:"-"`
		PhotoUrl    string `json:"-"`
		Filename    string `json:"-"`
	}
	Photos []Photo
)
