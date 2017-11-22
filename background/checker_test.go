package background

import (
	"github.com/chremoas/auth-cmd/mocks"
	uauthmocks "github.com/chremoas/auth-srv/mocks"
	proto "github.com/chremoas/auth-srv/proto"
	"github.com/bwmarrin/discordgo"
	gomock "github.com/golang/mock/gomock"
	context "golang.org/x/net/context"
	"testing"
	"time"
)

type mockError struct {
	message string
}

func (me *mockError) Error() string {
	return me.message
}

func TestPollHappyPathNoRoles(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockClient := mocks.NewMockClient(mockCtrl)
	mockAuthSvc := uauthmocks.NewMockUserAuthenticationClient(mockCtrl)
	mockFactory := mocks.NewMockClientFactory(mockCtrl)
	mockRoleMap := mocks.NewMockRoleMap(mockCtrl)
	defer mockCtrl.Finish()

	checker := NewChecker("g1234567890", mockClient, mockFactory, mockRoleMap, time.Millisecond*500)

	mockClient.EXPECT().GetUser("@me").Return(&discordgo.User{ID: "12345678901"}, nil).Times(1)
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

//BGN Error Path Tests
//From above to the END Error Path Tests I'm merely doing stuff to get code coverage up.  These unit tests will need
//revisited once a better error handling path is chosen for full validation
//TODO: Revisit once error handling path is chosen (such as talking in a specific channel or something)
func TestPollWithErrorAtUpdateRoles(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockClient := mocks.NewMockClient(mockCtrl)
	mockAuthSvc := uauthmocks.NewMockUserAuthenticationClient(mockCtrl)
	mockFactory := mocks.NewMockClientFactory(mockCtrl)
	mockRoleMap := mocks.NewMockRoleMap(mockCtrl)
	defer mockCtrl.Finish()

	checker := NewChecker("g1234567890", mockClient, mockFactory, mockRoleMap, time.Millisecond*500)

	mockClient.EXPECT().GetUser("@me").Return(&discordgo.User{ID: "12345678901"}, nil).Times(1)
	mockRoleMap.EXPECT().UpdateRoles().Return(&mockError{message: "ERROR!"}).Times(1)
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
	).Times(0)
	mockFactory.EXPECT().NewClient().Return(mockAuthSvc).Times(0)
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
	).Times(0)
	mockRoleMap.EXPECT().GetRoleId("ROLE1").Return("r1234567890").Times(0)
	mockRoleMap.EXPECT().GetRoleId("ROLE2").Return("r2234567890").Times(0)
	mockClient.EXPECT().UpdateMember("g1234567890", "u1234567890", []string{"r1234567890", "r2234567890"}).Times(0)

	checker.Start()

	time.Sleep(time.Millisecond * 520)

	checker.Stop()
}

func TestPollWithErrorAtGetAllMembers(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockClient := mocks.NewMockClient(mockCtrl)
	mockAuthSvc := uauthmocks.NewMockUserAuthenticationClient(mockCtrl)
	mockFactory := mocks.NewMockClientFactory(mockCtrl)
	mockRoleMap := mocks.NewMockRoleMap(mockCtrl)
	defer mockCtrl.Finish()

	checker := NewChecker("g1234567890", mockClient, mockFactory, mockRoleMap, time.Millisecond*500)

	mockClient.EXPECT().GetUser("@me").Return(&discordgo.User{ID: "12345678901"}, nil).Times(1)
	mockRoleMap.EXPECT().UpdateRoles().Return(nil).Times(1)
	mockClient.EXPECT().GetAllMembers("g1234567890", "", 1000).Return(
		nil,
		&mockError{message: "ERROR!"},
	).Times(1)
	mockFactory.EXPECT().NewClient().Return(mockAuthSvc).Times(0)
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
	).Times(0)
	mockRoleMap.EXPECT().GetRoleId("ROLE1").Return("r1234567890").Times(0)
	mockRoleMap.EXPECT().GetRoleId("ROLE2").Return("r2234567890").Times(0)
	mockClient.EXPECT().UpdateMember("g1234567890", "u1234567890", []string{"r1234567890", "r2234567890"}).Times(0)

	checker.Start()

	time.Sleep(time.Millisecond * 520)

	checker.Stop()
}

func TestPollWithErrorAtAuthClientGetRoles(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockClient := mocks.NewMockClient(mockCtrl)
	mockAuthSvc := uauthmocks.NewMockUserAuthenticationClient(mockCtrl)
	mockFactory := mocks.NewMockClientFactory(mockCtrl)
	mockRoleMap := mocks.NewMockRoleMap(mockCtrl)
	defer mockCtrl.Finish()

	checker := NewChecker("g1234567890", mockClient, mockFactory, mockRoleMap, time.Millisecond*500)

	mockClient.EXPECT().GetUser("@me").Return(&discordgo.User{ID: "12345678901"}, nil).Times(1)
	mockRoleMap.EXPECT().UpdateRoles().Return(nil).Times(1)
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
	).Times(1)
	mockFactory.EXPECT().NewClient().Return(mockAuthSvc).Times(1)
	mockAuthSvc.EXPECT().GetRoles(
		context.Background(),
		&proto.GetRolesRequest{UserId: "u1234567890"},
	).Return(
		nil,
		&mockError{message: "ERROR!"},
	).Times(1)
	mockRoleMap.EXPECT().GetRoleId("ROLE1").Return("r1234567890").Times(0)
	mockRoleMap.EXPECT().GetRoleId("ROLE2").Return("r2234567890").Times(0)
	mockClient.EXPECT().UpdateMember("g1234567890", "u1234567890", []string{"r1234567890", "r2234567890"}).Times(0)

	checker.Start()

	//Sleep for just one tick and then a tiny bit more
	time.Sleep(time.Millisecond * 520)

	checker.Stop()
}

func TestPollWithErrorAtClientUpdateMember(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockClient := mocks.NewMockClient(mockCtrl)
	mockAuthSvc := uauthmocks.NewMockUserAuthenticationClient(mockCtrl)
	mockFactory := mocks.NewMockClientFactory(mockCtrl)
	mockRoleMap := mocks.NewMockRoleMap(mockCtrl)
	defer mockCtrl.Finish()

	checker := NewChecker("g1234567890", mockClient, mockFactory, mockRoleMap, time.Millisecond*500)

	mockClient.EXPECT().GetUser("@me").Return(&discordgo.User{ID: "12345678901"}, nil).Times(1)
	mockRoleMap.EXPECT().UpdateRoles().Return(nil).Times(1)
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
	).Times(1)
	mockFactory.EXPECT().NewClient().Return(mockAuthSvc).Times(1)
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
	).Times(1)
	mockRoleMap.EXPECT().GetRoleId("ROLE1").Return("r1234567890").Times(1)
	mockRoleMap.EXPECT().GetRoleId("ROLE2").Return("r2234567890").Times(1)
	mockClient.EXPECT().UpdateMember("g1234567890", "u1234567890", []string{"r1234567890", "r2234567890"}).Return(&mockError{message: "ERROR!"}).Times(1)

	checker.Start()

	time.Sleep(time.Millisecond * 520)

	checker.Stop()
}

func TestPollWithErrorAtClientRemoveMemberRole(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockClient := mocks.NewMockClient(mockCtrl)
	mockAuthSvc := uauthmocks.NewMockUserAuthenticationClient(mockCtrl)
	mockFactory := mocks.NewMockClientFactory(mockCtrl)
	mockRoleMap := mocks.NewMockRoleMap(mockCtrl)
	defer mockCtrl.Finish()

	checker := NewChecker("g1234567890", mockClient, mockFactory, mockRoleMap, time.Millisecond*500)

	mockClient.EXPECT().GetUser("@me").Return(&discordgo.User{ID: "12345678901"}, nil).Times(1)
	mockRoleMap.EXPECT().UpdateRoles().Return(nil).Times(1)
	mockClient.EXPECT().GetAllMembers("g1234567890", "", 1000).Return(
		[]*discordgo.Member{
			{
				User: &discordgo.User{
					ID:       "u1234567890",
					Username: "Test User 1",
				},
				Roles: []string{"r2234567890"},
			},
		},
		nil,
	).Times(1)
	mockFactory.EXPECT().NewClient().Return(mockAuthSvc).Times(1)
	mockAuthSvc.EXPECT().GetRoles(
		context.Background(),
		&proto.GetRolesRequest{UserId: "u1234567890"},
	).Return(
		&proto.AuthConfirmResponse{
			Success:       true,
			CharacterName: "Test Character Name 1",
			Roles:         []string{},
		},
		nil,
	).Times(1)
	mockClient.EXPECT().RemoveMemberRole("g1234567890", "u1234567890", "r2234567890").Return(&mockError{message: "ERROR!"}).Times(1)

	checker.Start()

	time.Sleep(time.Millisecond * 520)

	checker.Stop()
}

//END Error Path Tests

func TestPollRemoveAllRoles(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockClient := mocks.NewMockClient(mockCtrl)
	mockAuthSvc := uauthmocks.NewMockUserAuthenticationClient(mockCtrl)
	mockFactory := mocks.NewMockClientFactory(mockCtrl)
	mockRoleMap := mocks.NewMockRoleMap(mockCtrl)
	defer mockCtrl.Finish()

	checker := NewChecker("g1234567890", mockClient, mockFactory, mockRoleMap, time.Millisecond*500)

	mockClient.EXPECT().GetUser("@me").Return(&discordgo.User{ID: "12345678901"}, nil).Times(1)
	mockRoleMap.EXPECT().UpdateRoles().Return(nil).Times(10)
	mockClient.EXPECT().GetAllMembers("g1234567890", "", 1000).Return(
		[]*discordgo.Member{
			{
				User: &discordgo.User{
					ID:       "u1234567890",
					Username: "Test User 1",
				},
				Roles: []string{
					"r1234567890",
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
			Roles:         []string{},
		},
		nil,
	).Times(10)
	mockClient.EXPECT().RemoveMemberRole("g1234567890", "u1234567890", "r1234567890").Return(nil).Times(10)
	mockClient.EXPECT().UpdateMember("g1234567890", "u1234567890", gomock.Any()).Times(0)

	checker.Start()

	//Sleep a little longer than 10 ticks so we get all the calls we want to happen
	time.Sleep(time.Millisecond * 501 * 10)

	checker.Stop()
}

func TestPollRemoveARole(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockClient := mocks.NewMockClient(mockCtrl)
	mockAuthSvc := uauthmocks.NewMockUserAuthenticationClient(mockCtrl)
	mockFactory := mocks.NewMockClientFactory(mockCtrl)
	mockRoleMap := mocks.NewMockRoleMap(mockCtrl)
	defer mockCtrl.Finish()

	checker := NewChecker("g1234567890", mockClient, mockFactory, mockRoleMap, time.Millisecond*500)

	mockClient.EXPECT().GetUser("@me").Return(&discordgo.User{ID: "12345678901"}, nil).Times(1)
	mockRoleMap.EXPECT().UpdateRoles().Return(nil).Times(1)
	mockClient.EXPECT().GetAllMembers("g1234567890", "", 1000).Return(
		[]*discordgo.Member{
			{
				User: &discordgo.User{
					ID:       "u1234567890",
					Username: "Test User 1",
				},
				Roles: []string{
					"r1234567890",
					"r2234567890",
				},
			},
		},
		nil,
	).Times(1)
	mockFactory.EXPECT().NewClient().Return(mockAuthSvc).Times(1)
	mockAuthSvc.EXPECT().GetRoles(
		context.Background(),
		&proto.GetRolesRequest{UserId: "u1234567890"},
	).Return(
		&proto.AuthConfirmResponse{
			Success:       true,
			CharacterName: "Test Character Name 1",
			Roles:         []string{"ROLE2"},
		},
		nil,
	).Times(1)
	mockRoleMap.EXPECT().GetRoleId("ROLE2").Return("r2234567890")
	mockRoleMap.EXPECT().GetRoleName("r1234567890").Return("ROLE1")
	mockRoleMap.EXPECT().GetRoleName("r2234567890").Return("ROLE2")
	mockClient.EXPECT().RemoveMemberRole("g1234567890", "u1234567890", "r1234567890").Return(nil).Times(0)
	mockClient.EXPECT().UpdateMember("g1234567890", "u1234567890", []string{"r2234567890"}).Times(1)

	checker.Start()

	//Sleep for just one tick and then a tiny bit more
	time.Sleep(time.Millisecond * 520)

	checker.Stop()
}

func TestDontMessWithRolesIfTheyAllExist(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockClient := mocks.NewMockClient(mockCtrl)
	mockAuthSvc := uauthmocks.NewMockUserAuthenticationClient(mockCtrl)
	mockFactory := mocks.NewMockClientFactory(mockCtrl)
	mockRoleMap := mocks.NewMockRoleMap(mockCtrl)
	defer mockCtrl.Finish()

	checker := NewChecker("g1234567890", mockClient, mockFactory, mockRoleMap, time.Millisecond*500)

	mockClient.EXPECT().GetUser("@me").Return(&discordgo.User{ID: "12345678901"}, nil).Times(1)
	mockRoleMap.EXPECT().UpdateRoles().Return(nil).Times(10)
	mockClient.EXPECT().GetAllMembers("g1234567890", "", 1000).Return(
		[]*discordgo.Member{
			{
				User: &discordgo.User{
					ID:       "u1234567890",
					Username: "Test User 1",
				},
				Roles: []string{
					"r1234567890",
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
			Roles:         []string{"ROLE1"},
		},
		nil,
	).Times(10)
	mockRoleMap.EXPECT().GetRoleId("ROLE1").Return("r1234567890").Times(10)
	mockRoleMap.EXPECT().GetRoleName("r1234567890").Return("ROLE1").Times(10)
	mockClient.EXPECT().RemoveMemberRole("g1234567890", "u1234567890", "r1234567890").Return(nil).Times(0)

	checker.Start()

	//Sleep a little longer than 10 ticks so we get all the calls we want to happen
	time.Sleep(time.Millisecond * 501 * 10)

	checker.Stop()
}

func TestAddRolesThatAreMissing(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockClient := mocks.NewMockClient(mockCtrl)
	mockAuthSvc := uauthmocks.NewMockUserAuthenticationClient(mockCtrl)
	mockFactory := mocks.NewMockClientFactory(mockCtrl)
	mockRoleMap := mocks.NewMockRoleMap(mockCtrl)
	defer mockCtrl.Finish()

	checker := NewChecker("g1234567890", mockClient, mockFactory, mockRoleMap, time.Millisecond*500)

	mockClient.EXPECT().GetUser("@me").Return(&discordgo.User{ID: "12345678901"}, nil).Times(1)
	mockRoleMap.EXPECT().UpdateRoles().Return(nil).Times(10)
	mockClient.EXPECT().GetAllMembers("g1234567890", "", 1000).Return(
		[]*discordgo.Member{
			{
				User: &discordgo.User{
					ID:       "u1234567890",
					Username: "Test User 1",
				},
				Roles: []string{
					"r1234567890",
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
	mockRoleMap.EXPECT().GetRoleName("r1234567890").Return("ROLE1").Times(10)
	mockClient.EXPECT().RemoveMemberRole("g1234567890", "u1234567890", "r1234567890").Return(nil).Times(0)
	mockClient.EXPECT().UpdateMember("g1234567890", "u1234567890", []string{"r1234567890", "r2234567890"}).Times(10)

	checker.Start()

	//Sleep a little longer than 10 ticks so we get all the calls we want to happen
	time.Sleep(time.Millisecond * 501 * 10)

	checker.Stop()
}
