package service

//
//import (
//	"goTSV/internal/shema"
//	"testing"
//)
//
//func TestService_ParseFile(t *testing.T) {
//	tests := []struct {
//		name      string
//		args      string
//		wantTsv   []shema.Tsv
//		wantGuids []string
//		wantErr   error
//	}{
//		{
//			name:    "OK1",
//			args:    "test.tsv",
//			wantTsv: []shema.Tsv{},
//			wantErr: nil,
//		},
//		{
//			name: "BAD1",
//			args: model.User{
//				Username: "dima",
//				Password: "test1",
//				Session:  "ahsjufil12-fk",
//			},
//			sessionMock: func(c *mocks.SessionService, user model.User) {
//				c.Mock.On("Generate").Return(user.Session, nil).Times(1)
//			},
//			repositoryMock: func(c *mocks.IRepository, user model.User) {
//				c.Mock.On("CreateUser", user).Return(invalidErr).Times(1)
//			},
//			wantErr: invalidErr,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			storage := mocks.NewIRepository(t)
//			session := mocks.NewSessionService(t)
//			tt.repositoryMock(storage, tt.args)
//			tt.sessionMock(session, tt.args)
//			service := Service{
//				database:       storage,
//				sessionService: session,
//			}
//			cook, err := service.SignUp(tt.args)
//			if !errors.Is(err, tt.wantErr) {
//				t.Errorf("got %v, want %v", err, tt.wantErr)
//			}
//			if tt.wantErr == nil && cook != tt.args.Session {
//				t.Errorf("got %s, want %s", cook, tt.args.Session)
//			}
//		})
//	}
//}
