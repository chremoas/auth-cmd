package mocks

//go:generate mockgen -source=../discord/discord.go -package=mocks -destination=discord_mocks.go Client
//go:generate mockgen -source=../discord/role_map.go -package=mocks -destination=role_map_mocks.go RoleMap
//go:generate mockgen -source=../command/command.go -package=mocks -destination=client_factory_mocks.go ClientFactory
