package users

import "math"

type UserService interface {
	ReadUser(filter UserFilter, page, size int) (UserList, error)
	GetProfile(uid string) (*ProfileView, error)
}

type UserServiceImpl struct {
	UserRepository UserRepository
}

func ProvideUserServiceImpl(userRepository UserRepository) *UserServiceImpl {
	return &UserServiceImpl{
		UserRepository: userRepository,
	}
}

func (s *UserServiceImpl) GetProfile(uid string) (*ProfileView, error) {
	return s.UserRepository.GetProfile(uid)
}

func (s *UserServiceImpl) ReadUser(filter UserFilter, page, size int) (UserList, error) {
	users, err := s.UserRepository.GetData(filter, page, size)
	if err != nil {
		return UserList{}, err
	}

	totalData, err := s.UserRepository.CountTotalData(filter)
	if err != nil {
		return UserList{}, err
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(totalData) / float64(size)))

	// Populate the response
	var nextPage *int
	if page < totalPages {
		nextPageValue := page + 1
		nextPage = &nextPageValue
	}

	var previousPage *int
	if page > 1 {
		previousPageValue := page - 1
		previousPage = &previousPageValue
	}

	response := UserList{
		Data:         users,
		TotalData:    totalData,
		TotalPages:   totalPages,
		CurrentPage:  page,
		NextPage:     nextPage,
		PreviousPage: previousPage,
	}

	return response, nil
}
