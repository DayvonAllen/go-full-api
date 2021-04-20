package domain

func UserMapper(user *User) *UserDto {
	userDto := new(UserDto)
	userDto.Email = user.Email
	userDto.ProfilePictureUrl = user.ProfilePictureUrl
	userDto.CurrentTagLine = user.CurrentTagLine
	userDto.UnlockedTagLine = user.UnlockedTagLine
	userDto.CurrentBadgeUrl = user.CurrentBadgeUrl
	userDto.ProfileIsViewable = user.ProfileIsViewable
	userDto.UnlockedBadgesUrls = user.UnlockedBadgesUrls
	userDto.AcceptMessages = user.AcceptMessages

	return userDto
}
