package log

type Option struct {
	DirPath        string
	MaxFileSize    string
	RotateDuration string
	Level          string
	BackCount      uint32
	BackTime       string
}
