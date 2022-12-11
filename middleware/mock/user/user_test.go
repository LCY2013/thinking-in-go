package user

import (
	"github.com/LCY2013/thinking-in-go/middleware/mock/mock"
	"github.com/LCY2013/thinking-in-go/middleware/mock/person"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUser_GetUserInfo(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	var id int64 = 1
	maleMock := mock.NewMockMale(ctl)
	gomock.InOrder(
		maleMock.EXPECT().Get(id).Return(nil),
	)

	tests := []struct {
		name string

		Person person.Male
		id     int64

		wantErr error
	}{
		{
			name:   "normal case",
			Person: NewUser(maleMock),
			id:     1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.Person.Get(tt.id); err != nil {
				assert.Equal(t, tt.wantErr, err)
			}
		})
	}
}
