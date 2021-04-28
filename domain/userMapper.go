package domain

func UserMapper(user *User) *UserDto {
	userDto := new(UserDto)
	userDto.Id = user.Id
	userDto.Email = user.Email
	userDto.Username = user.Username
	userDto.ProfilePictureUrl = user.ProfilePictureUrl
	userDto.CurrentTagLine = user.CurrentTagLine
	userDto.UnlockedTagLine = user.UnlockedTagLine
	userDto.CurrentBadgeUrl = user.CurrentBadgeUrl
	userDto.ProfileIsViewable = user.ProfileIsViewable
	userDto.UnlockedBadgesUrls = user.UnlockedBadgesUrls
	userDto.AcceptMessages = user.AcceptMessages

	return userDto
}

func UserDtoMapper(dto UserDto) *User {
	user := new(User)
	user.Id = dto.Id
	user.Email = dto.Email
	user.Username = dto.Username
	user.ProfilePictureUrl = dto.ProfilePictureUrl
	user.CurrentTagLine = dto.CurrentTagLine
	user.UnlockedTagLine = dto.UnlockedTagLine
	user.CurrentBadgeUrl = dto.CurrentBadgeUrl
	user.ProfileIsViewable = dto.ProfileIsViewable
	user.UnlockedBadgesUrls = dto.UnlockedBadgesUrls
	user.AcceptMessages = dto.AcceptMessages
	user.IsVerified = dto.IsVerified

	return user
}