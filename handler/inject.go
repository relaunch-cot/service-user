package handler

import "github.com/relaunch-cot/service-user/repositories"

type Handlers struct {
	User IUserHandler
}

func (h *Handlers) Inject(repositories *repositories.Repositories) {
	h.User = NewUserHandler(repositories)
}
