package command

import (
	proto "github.com/abaeve/auth-srv/proto"
	botprot "github.com/micro/micro/bot/proto"
	gomock "github.com/golang/mock/gomock"
	client "github.com/micro/go-micro/client"
	context "golang.org/x/net/context"
	"testing"
	"github.com/bwmarrin/discordgo"
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

func (_m *MockClient) UpdateMember(_param0 string, _param1 string, _param2 []string) error {
	ret := _m.ctrl.Call(_m, "UpdateMember", _param0, _param1, _param2)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockClientRecorder) UpdateMember(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "UpdateMember", arg0, arg1, arg2)
}

//</editor-fold>

func TestBotExec(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockClient := NewMockClient(mockCtrl)
	mockAuthSvc := NewMockUserAuthenticationClient(mockCtrl)
	defer mockCtrl.Finish()

	expectedCharName := "Test Char Name"

	gomock.InOrder(
		mockAuthSvc.EXPECT().Confirm(
			context.Background(),
			&proto.AuthConfirmRequest{
				UserId:             "u123456",
				AuthenticationCode: "1234567890",
			},
			gomock.Any(),
		).Return(
			&proto.AuthConfirmResponse{
				Success:       true,
				Roles:         []string{"ROLE1", "ROLE2"},
				CharacterName: expectedCharName,
			},
			nil,
		),
		mockClient.EXPECT().UpdateMember("g123456", "u123456", []string{"ROLE1", "ROLE2"}),
	)

	cmd := Command{guildID: "g123456", client: mockClient, authSvc: mockAuthSvc, name: "test"}

	response := botprot.ExecResponse{}
	err := cmd.Exec(context.Background(), &botprot.ExecRequest{Sender: "g123456:u123456", Args: []string{"auth", "1234567890"}}, &response)

	if err != nil {
		t.Fatal("Received an error when one wasn't expected")
	}

	if string(response.Result) != "@u123456 *Success*: "+expectedCharName+" has been successfully authed" {
		t.Fatalf("Result string, expected (%s) received (%s)",
			"@u123456 *Success*: "+expectedCharName+" has been successfully authed",
			string(response.Result))
	}

	if len(response.Error) != 0 {
		t.Fatal("Bot set the error in the response when it shouldn't have")
	}
}

func TestInvalidCommandExecution(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockClient := NewMockClient(mockCtrl)
	mockAuthSvc := NewMockUserAuthenticationClient(mockCtrl)
	defer mockCtrl.Finish()

	gomock.InOrder(
		mockAuthSvc.EXPECT().Confirm(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
		mockClient.EXPECT().UpdateMember(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
	)

	cmd := Command{guildID: "g123456", client: mockClient, authSvc: mockAuthSvc, name: "test"}

	var response botprot.ExecResponse

	err := cmd.Exec(context.Background(), &botprot.ExecRequest{Sender: "g123456:u123456", Args: []string{"auth"}}, &response)

	if err == nil {
		t.Fatal("Expected an error but received nil")
	}

	expectedResponseError := "@u123456 I did not understand your command."

	if len(response.Error) == 0 || response.Error != expectedResponseError {
		t.Fatalf("Response error: (%s) did not match expected: (%s)", response.Error, expectedResponseError)
	}
}

func TestErrorFromAuthSvc(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockClient := NewMockClient(mockCtrl)
	mockAuthSvc := NewMockUserAuthenticationClient(mockCtrl)
	defer mockCtrl.Finish()

	authError := "There was an issue over there ->"

	gomock.InOrder(
		mockAuthSvc.EXPECT().Confirm(
			context.Background(),
			&proto.AuthConfirmRequest{
				UserId:             "u123456",
				AuthenticationCode: "1234567890",
			},
			gomock.Any(),
		).Return(
			nil,
			&botError{message: authError},
		),
		mockClient.EXPECT().UpdateMember("g123456", "u123456", []string{"ROLE1", "ROLE2"}).Times(0),
	)

	cmd := Command{guildID: "g123456", client: mockClient, authSvc: mockAuthSvc, name: "test"}

	response := botprot.ExecResponse{}
	err := cmd.Exec(context.Background(), &botprot.ExecRequest{Sender: "g123456:u123456", Args: []string{"auth", "1234567890"}}, &response)

	if err == nil {
		t.Fatal("Expected an error but received nil")
	}

	expectedAuthError := "Received an error from the auth service: " + authError

	if err.Error() != expectedAuthError {
		t.Fatalf("Error text: (%s) did not match expected: (%s)", err.Error(), expectedAuthError)
	}

	expectedUserError := "@u123456 I had an issue authing your request, please reauth or contact your administrator."

	if response.Error != expectedUserError {
		t.Fatalf("User error: (%s) did not match expected: (%s)", response.Error, authError)
	}
}

func TestErrorFromDiscord(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockClient := NewMockClient(mockCtrl)
	mockAuthSvc := NewMockUserAuthenticationClient(mockCtrl)
	defer mockCtrl.Finish()

	discordError := "There was an issue over there ->"

	gomock.InOrder(
		mockAuthSvc.EXPECT().Confirm(
			context.Background(),
			&proto.AuthConfirmRequest{
				UserId:             "u123456",
				AuthenticationCode: "1234567890",
			},
			gomock.Any(),
		).Return(
			&proto.AuthConfirmResponse{
				Success:       true,
				Roles:         []string{"ROLE1", "ROLE2"},
				CharacterName: "Test Char",
			},
			nil,
		),
		mockClient.EXPECT().UpdateMember("g123456", "u123456", []string{"ROLE1", "ROLE2"}).Return(&botError{message: discordError}),
	)

	cmd := Command{guildID: "g123456", client: mockClient, authSvc: mockAuthSvc, name: "test"}

	response := botprot.ExecResponse{}
	err := cmd.Exec(context.Background(), &botprot.ExecRequest{Sender: "g123456:u123456", Args: []string{"auth", "1234567890"}}, &response)

	if err == nil {
		t.Fatal("Expected an error but received nil")
	}

	expectedDiscordError := "Received (" + discordError + ") from the chat service."

	if err.Error() != expectedDiscordError {
		t.Fatalf("Error text: (%s) did not match expected: (%s)", err.Error(), expectedDiscordError)
	}

	expectedUserError := "@u123456 I had an issue talking to the chat service, please try again later."

	if response.Error != expectedUserError {
		t.Fatalf("User error: (%s) did not match expected: (%s)", response.Error, discordError)
	}
}

func TestHelp(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockClient := NewMockClient(mockCtrl)
	mockAuthSvc := NewMockUserAuthenticationClient(mockCtrl)
	defer mockCtrl.Finish()

	cmd := Command{guildID: "g123456", client: mockClient, authSvc: mockAuthSvc, name: "test"}

	response := botprot.HelpResponse{}

	err := cmd.Help(context.Background(), &botprot.HelpRequest{}, &response)

	if err != nil {
		t.Fatal("Received an error when none was expected")
	}

	expectedHelpString := "Authenticate your chat user id and link it to the character used to create the given token."

	if response.Description != expectedHelpString {
		t.Fatalf("Response description: (%s) does not match expected: (%s)", response.Description, expectedHelpString)
	}

	if response.Usage != "test" {
		t.Fatalf("Response name: (%s) does not equal expected: (test)", response.Usage)
	}
}
