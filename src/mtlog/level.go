package mtlog

type Level int

const (
	TRACE Level = iota
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
	REPORT
	OFF
)

var levelNames = []string{
	"trace",
	"debug",
	"info",
	"warn",
	"error",
	"fatal",
	"report",
}

func (o Level) String() string {
	return levelNames[o]
}
