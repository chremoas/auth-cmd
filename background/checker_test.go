package background

import (
	proto "github.com/abaeve/auth-srv/proto"
	gomock "github.com/golang/mock/gomock"
	client "github.com/micro/go-micro/client"
	context "golang.org/x/net/context"
	"github.com/bwmarrin/discordgo"
	"testing"
	"time"
)

//<editor-fold desc="Generated Mocks">
type MockUserAuthenticationClient struct {
	ctrl     *gomock.Controller
	recorder *_MockUserAuthenticationClientRecorder
}

type _MockUserAuthenticationClientRecorder struct {
	mock *MockUserAuthenticationClient
}

func NewMockUserAuthenticationClient(ctrl *gomock.Controller) *MockUserAuthenticationClient {
	mock := &MockUserAuthenticationClient{ctrl: ctrl}
	mock.recorder = &_MockUserAuthenticationClientRecorder{mock}
	return mock
}

func (_m *MockUserAuthenticationClient) EXPECT() *_MockUserAuthenticationClientRecorder {
	return _m.recorder
}

func (_m *MockUserAuthenticationClient) Confirm(_param0 context.Context, _param1 *proto.AuthConfirmRequest, _param2 ...client.CallOption) (*proto.AuthConfirmResponse, error) {
	_s := []interface{}{_param0, _param1}
	for _, _x := range _param2 {
		_s = append(_s, _x)
	}
	ret := _m.ctrl.Call(_m, "Confirm", _s...)
	ret0, _ := ret[0].(*proto.AuthConfirmResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockUserAuthenticationClientRecorder) Confirm(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	_s := append([]interface{}{arg0, arg1}, arg2...)
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Confirm", _s...)
}

func (_m *MockUserAuthenticationClient) Create(_param0 context.Context, _param1 *proto.AuthCreateRequest, _param2 ...client.CallOption) (*proto.AuthCreateResponse, error) {
	_s := []interface{}{_param0, _param1}
	for _, _x := range _param2 {
		_s = append(_s, _x)
	}
	ret := _m.ctrl.Call(_m, "Create", _s...)
	ret0, _ := ret[0].(*proto.AuthCreateResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockUserAuthenticationClientRecorder) Create(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	_s := append([]interface{}{arg0, arg1}, arg2...)
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Create", _s...)
}

func (_m *MockUserAuthenticationClient) GetRoles(_param0 context.Context, _param1 *proto.GetRolesRequest, _param2 ...client.CallOption) (*proto.AuthConfirmResponse, error) {
	_s := []interface{}{_param0, _param1}
	for _, _x := range _param2 {
		_s = append(_s, _x)
	}
	ret := _m.ctrl.Call(_m, "GetRoles", _s...)
	ret0, _ := ret[0].(*proto.AuthConfirmResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockUserAuthenticationClientRecorder) GetRoles(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	_s := append([]interface{}{arg0, arg1}, arg2...)
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetRoles", _s...)
}

type MockClient struct {
	ctrl     *gomock.Controller
	recorder *_MockClientRecorder
}

type _MockClientRecorder struct {
	mock *MockClient
}

func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &_MockClientRecorder{mock}
	return mock
}

func (_m *MockClient) EXPECT() *_MockClientRecorder {
	return _m.recorder
}

func (_m *MockClient) GetAllMembers(_param0 string, _param1 string, _param2 int) ([]*discordgo.Member, error) {
	ret := _m.ctrl.Call(_m, "GetAllMembers", _param0, _param1, _param2)
	ret0, _ := ret[0].([]*discordgo.Member)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockClientRecorder) GetAllMembers(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetAllMembers", arg0, arg1, arg2)
}

func (_m *MockClient) GetAllRoles(_param0 string) ([]*discordgo.Role, error) {
	ret := _m.ctrl.Call(_m, "GetAllRoles", _param0)
	ret0, _ := ret[0].([]*discordgo.Role)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockClientRecorder) GetAllRoles(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetAllRoles", arg0)
}

func (_m *MockClient) UpdateMember(_param0 string, _param1 string, _param2 []string) error {
	ret := _m.ctrl.Call(_m, "UpdateMember", _param0, _param1, _param2)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockClientRecorder) UpdateMember(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "UpdateMember", arg0, arg1, arg2)
}

type MockClientFactory struct {
	ctrl     *gomock.Controller
	recorder *_MockClientFactoryRecorder
}
type _MockClientFactoryRecorder struct {
	mock *MockClientFactory
}

func NewMockClientFactory(ctrl *gomock.Controller) *MockClientFactory {
	mock := &MockClientFactory{ctrl: ctrl}
	mock.recorder = &_MockClientFactoryRecorder{mock}
	return mock
}

func (_m *MockClientFactory) EXPECT() *_MockClientFactoryRecorder {
	return _m.recorder
}

func (_m *MockClientFactory) NewClient() proto.UserAuthenticationClient {
	ret := _m.ctrl.Call(_m, "NewClient")
	ret0, _ := ret[0].(proto.UserAuthenticationClient)
	return ret0
}

func (_mr *_MockClientFactoryRecorder) NewClient() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "NewClient")
}

type MockRoleMap struct {
	ctrl     *gomock.Controller
	recorder *_MockRoleMapRecorder
}

type _MockRoleMapRecorder struct {
	mock *MockRoleMap
}

func NewMockRoleMap(ctrl *gomock.Controller) *MockRoleMap {
	mock := &MockRoleMap{ctrl: ctrl}
	mock.recorder = &_MockRoleMapRecorder{mock}
	return mock
}

func (_m *MockRoleMap) EXPECT() *_MockRoleMapRecorder {
	return _m.recorder
}

func (_m *MockRoleMap) GetRoleId(_param0 string) string {
	ret := _m.ctrl.Call(_m, "GetRoleId", _param0)
	ret0, _ := ret[0].(string)
	return ret0
}

func (_mr *_MockRoleMapRecorder) GetRoleId(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetRoleId", arg0)
}

func (_m *MockRoleMap) GetRoles() map[string]*discordgo.Role {
	ret := _m.ctrl.Call(_m, "GetRoles")
	ret0, _ := ret[0].(map[string]*discordgo.Role)
	return ret0
}

func (_mr *_MockRoleMapRecorder) GetRoles() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetRoles")
}

func (_m *MockRoleMap) UpdateRoles() error {
	ret := _m.ctrl.Call(_m, "UpdateRoles")
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockRoleMapRecorder) UpdateRoles() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "UpdateRoles")
}

//</editor-fold>

func TestPoll(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockClient := NewMockClient(mockCtrl)
	mockAuthSvc := NewMockUserAuthenticationClient(mockCtrl)
	mockFactory := NewMockClientFactory(mockCtrl)
	mockRoleMap := NewMockRoleMap(mockCtrl)
	defer mockCtrl.Finish()

	checker := NewChecker("g1234567890", mockClient, mockFactory, mockRoleMap, time.Millisecond*500)

	mockRoleMap.EXPECT().UpdateRoles().Return(nil).Times(10)
	mockClient.EXPECT().GetAllMembers("g1234567890", "", 1000).Return(
		[]*discordgo.Member{
			{
				User: &discordgo.User{
					ID:       "u1234567890",
					Username: "Test User 1",
				},
			},
		},
		nil,
	).Times(10)
	mockFactory.EXPECT().NewClient().Return(mockAuthSvc).Times(10)
	mockAuthSvc.EXPECT().GetRoles(
		context.Background(),
		&proto.GetRolesRequest{UserId: "u1234567890"},
	).Return(
		&proto.AuthConfirmResponse{
			Success:       true,
			CharacterName: "Test Character Name 1",
			Roles:         []string{"ROLE1", "ROLE2"},
		},
		nil,
	).Times(10)
	mockRoleMap.EXPECT().GetRoleId("ROLE1").Return("r1234567890").Times(10)
	mockRoleMap.EXPECT().GetRoleId("ROLE2").Return("r2234567890").Times(10)
	mockClient.EXPECT().UpdateMember("g1234567890", "u1234567890", []string{"r1234567890", "r2234567890"}).Times(10)

	checker.Start()

	//Sleep a little longer than 10 ticks so we get all the calls we want to happen
	time.Sleep(time.Millisecond * 501 * 10)

	checker.Stop()
}

func TestUpdate(t *testing.T) {

}
