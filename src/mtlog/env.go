package mtlog

type Env int

const (
	DEVELOP Env = iota
	ABTEST
	PRODUCT
)

var envNames = []string{
	"develop",
	"abtest",
	"product",
}

func (o Env) String() string {
	return envNames[o]
}
