package person

//go:generate mockgen -destination=../mock/male_mock.go -package=mock github.com/LCY2013/thinking-in-go/middleware/mock/person Male

type Male interface {
	Get(id int64) error
}
