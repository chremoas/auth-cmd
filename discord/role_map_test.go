package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/golang/mock/gomock"
	"testing"
)

//<editor-fold desc="Generated Mocks">
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

func (_m *MockClient) GetUser(_param0 string) (*discordgo.User, error) {
	ret := _m.ctrl.Call(_m, "GetUser", _param0)
	ret0, _ := ret[0].(*discordgo.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockClientRecorder) GetUser(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetUser", arg0)
}

func (_m *MockClient) RemoveMemberRole(_param0 string, _param1 string, _param2 string) error {
	ret := _m.ctrl.Call(_m, "RemoveMemberRole", _param0, _param1, _param2)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockClientRecorder) RemoveMemberRole(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "RemoveMemberRole", arg0, arg1, arg2)
}

func (_m *MockClient) UpdateMember(_param0 string, _param1 string, _param2 []string) error {
	ret := _m.ctrl.Call(_m, "UpdateMember", _param0, _param1, _param2)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockClientRecorder) UpdateMember(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "UpdateMember", arg0, arg1, arg2)
}

//</editor-fold>

type mockError struct {
	message string
}

func (me *mockError) Error() string {
	return me.message
}

func TestUpdateRoles(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockClient := NewMockClient(mockCtrl)

	roleMap := &roleMapImpl{guildID: "1234567890", client: mockClient}

	gomock.InOrder(
		mockClient.EXPECT().GetAllRoles("1234567890").Return([]*discordgo.Role{
			{
				Name: "TEST ROLE 1",
				ID:   "0123456789",
			},
			{
				Name: "TEST ROLE 2",
				ID:   "0234567890",
			},
			{
				Name: "TEST ROLE 3",
				ID:   "0345678901",
			},
		}, nil),
	)

	err := roleMap.UpdateRoles()

	if err != nil {
		t.Fatalf("Received an error when none was expected: (%s)", err)
	}

	if len(roleMap.rolesByName) == 0 {
		t.Fatal("Expected more than zero roles")
	}

	if roleMap.rolesByName["TEST ROLE 1"] == nil {
		t.Fatal("Role 1 was not properly put into the map")
	}

	if roleMap.rolesByName["TEST ROLE 2"] == nil {
		t.Fatal("Role 2 was not properly put into the map")
	}

	if roleMap.rolesByName["TEST ROLE 3"] == nil {
		t.Fatal("Role 3 was not properly put into the map")
	}

	if roleMap.rolesByName["TEST ROLE 1"].ID != "0123456789" {
		t.Fatalf("Expected id for role 1: (0123456789) but received: (%s)", roleMap.rolesByName["TEST ROLE 1"].ID)
	}

	if roleMap.rolesByName["TEST ROLE 2"].ID != "0234567890" {
		t.Fatalf("Expected id for role 2: (0234567890) but received: (%s)", roleMap.rolesByName["TEST ROLE 1"].ID)
	}

	if roleMap.rolesByName["TEST ROLE 3"].ID != "0345678901" {
		t.Fatalf("Expected id for role 3: (0345678901) but received: (%s)", roleMap.rolesByName["TEST ROLE 1"].ID)
	}
}

func TestUpdateRolesWithError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockClient := NewMockClient(mockCtrl)

	roleMap := &roleMapImpl{guildID: "1234567890", client: mockClient}

	gomock.InOrder(
		mockClient.EXPECT().GetAllRoles("1234567890").Return(nil, &mockError{"OUCH!"}),
	)

	err := roleMap.UpdateRoles()

	if err == nil || err.Error() != "OUCH!" {
		t.Fatalf("Received nil or the wrong string, expected (OUCH!) but received: (%s)", err.Error())
	}
}

func TestGetRoles(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockClient := NewMockClient(mockCtrl)

	roleMap := &roleMapImpl{guildID: "1234567890", client: mockClient}

	gomock.InOrder(
		mockClient.EXPECT().GetAllRoles("1234567890").Return([]*discordgo.Role{
			{
				Name: "TEST ROLE 1",
				ID:   "0123456789",
			},
			{
				Name: "TEST ROLE 2",
				ID:   "0234567890",
			},
			{
				Name: "TEST ROLE 3",
				ID:   "0345678901",
			},
		}, nil),
	)

	err := roleMap.UpdateRoles()

	if err != nil {
		t.Fatalf("Received an error when none was expected: (%s)", err)
	}

	roles := roleMap.GetRoles()

	if len(roles) == 0 {
		t.Fatal("Expected more than zero roles")
	}

	if roles["TEST ROLE 1"] == nil {
		t.Fatal("Role 1 was not properly put into the map")
	}

	if roles["TEST ROLE 2"] == nil {
		t.Fatal("Role 2 was not properly put into the map")
	}

	if roles["TEST ROLE 3"] == nil {
		t.Fatal("Role 3 was not properly put into the map")
	}

	if roles["TEST ROLE 1"].ID != "0123456789" {
		t.Fatalf("Expected id for role 1: (0123456789) but received: (%s)", roles["TEST ROLE 1"].ID)
	}

	if roles["TEST ROLE 2"].ID != "0234567890" {
		t.Fatalf("Expected id for role 2: (0234567890) but received: (%s)", roles["TEST ROLE 1"].ID)
	}

	if roles["TEST ROLE 3"].ID != "0345678901" {
		t.Fatalf("Expected id for role 3: (0345678901) but received: (%s)", roles["TEST ROLE 1"].ID)
	}
}

func TestGetRoleId(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockClient := NewMockClient(mockCtrl)

	roleMap := &roleMapImpl{guildID: "1234567890", client: mockClient}

	gomock.InOrder(
		mockClient.EXPECT().GetAllRoles("1234567890").Return([]*discordgo.Role{
			{
				Name: "TEST ROLE 1",
				ID:   "0123456789",
			},
			{
				Name: "TEST ROLE 2",
				ID:   "0234567890",
			},
			{
				Name: "TEST ROLE 3",
				ID:   "0345678901",
			},
		}, nil),
	)

	err := roleMap.UpdateRoles()

	if err != nil {
		t.Fatalf("Received an error when none was expected: (%s)", err)
	}

	roleId := roleMap.GetRoleId("TEST ROLE 1")

	if len(roleId) == 0 {
		t.Fatal("Expected something as the role id but got 0 length string")
	}

	if roleId != "0123456789" {
		t.Fatalf("Expected role id: (%s) but received: (%s)", "0123456789", roleId)
	}
}

func TestGetRoleName(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockClient := NewMockClient(mockCtrl)

	roleMap := &roleMapImpl{guildID: "1234567890", client: mockClient}

	gomock.InOrder(
		mockClient.EXPECT().GetAllRoles("1234567890").Return([]*discordgo.Role{
			{
				Name: "TEST ROLE 1",
				ID:   "0123456789",
			},
			{
				Name: "TEST ROLE 2",
				ID:   "0234567890",
			},
			{
				Name: "TEST ROLE 3",
				ID:   "0345678901",
			},
		}, nil),
	)

	err := roleMap.UpdateRoles()

	if err != nil {
		t.Fatalf("Received an error when none was expected: (%s)", err)
	}

	roleName := roleMap.GetRoleName("0123456789")

	if roleName != "TEST ROLE 1" {
		t.Fatalf("Expected: (TEST ROLE 1) but recieved: (%s)", roleName)
	}
}

func TestGetRoleNameWithError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockClient := NewMockClient(mockCtrl)

	roleMap := &roleMapImpl{guildID: "1234567890", client: mockClient}

	gomock.InOrder(
		mockClient.EXPECT().GetAllRoles("1234567890").Return([]*discordgo.Role{
			{
				Name: "TEST ROLE 1",
				ID:   "0123456789",
			},
			{
				Name: "TEST ROLE 2",
				ID:   "0234567890",
			},
			{
				Name: "TEST ROLE 3",
				ID:   "0345678901",
			},
		}, nil),
	)

	err := roleMap.UpdateRoles()

	if err != nil {
		t.Fatalf("Received an error when none was expected: (%s)", err)
	}

	roleName := roleMap.GetRoleName("01234567890")

	if roleName != "" {
		t.Fatalf("Expected: (\"\") but recieved: (%s)", roleName)
	}
}

func TestGetRoleIdForNonRole(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockClient := NewMockClient(mockCtrl)

	roleMap := &roleMapImpl{guildID: "1234567890", client: mockClient}

	gomock.InOrder(
		mockClient.EXPECT().GetAllRoles("1234567890").Return([]*discordgo.Role{
			{
				Name: "TEST ROLE 1",
				ID:   "0123456789",
			},
			{
				Name: "TEST ROLE 2",
				ID:   "0234567890",
			},
			{
				Name: "TEST ROLE 3",
				ID:   "0345678901",
			},
		}, nil),
	)

	err := roleMap.UpdateRoles()

	if err != nil {
		t.Fatalf("Received an error when none was expected: (%s)", err)
	}

	roleId := roleMap.GetRoleId("DERP")

	if len(roleId) != 0 {
		t.Fatal("Received something when nothing was expected")
	}
}

func TestNewRoleMap(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockClient := NewMockClient(mockCtrl)

	roleMap := NewRoleMap("1234567890", mockClient).(*roleMapImpl)

	if roleMap.client != mockClient {
		t.Fatalf("Expected clients to match but they don't, original: (%+v) and result: (%+v)", mockClient, roleMap.client)
	}

	if roleMap.guildID != "1234567890" {
		t.Fatalf("Expected guild id's to match but they don't, origin: (1234567890) and result: (%s)", roleMap.guildID)
	}
}
