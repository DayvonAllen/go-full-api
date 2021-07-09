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
	userDto.DisplayFollowerCount = user.DisplayFollowerCount
	userDto.FollowerCount = user.FollowerCount
	userDto.Following = user.Following
	userDto.Followers = user.Followers

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
	user.UnlockedBadgesUrls = dto.UnlockedBadgesUrls
	user.AcceptMessages = dto.AcceptMessages
	user.IsVerified = dto.IsVerified
	user.DisplayFollowerCount = dto.DisplayFollowerCount
	user.Followers = dto.Followers
	user.FollowerCount = dto.FollowerCount
	user.Following = dto.Following

	return user
}