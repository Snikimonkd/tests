package usecase

type usecase struct {
	createUserRepository CreateUserRepository
}

func New(createUserRepository CreateUserRepository) usecase {
	return usecase{
		createUserRepository: createUserRepository,
	}
}
