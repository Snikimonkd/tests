package handler

type handler struct {
	createUserUsecase CreateUserUsecase
}

func New(createUserUsecase CreateUserUsecase) handler {
	return handler{
		createUserUsecase: createUserUsecase,
	}
}
