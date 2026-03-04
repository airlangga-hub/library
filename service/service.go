package service

type service struct {
	Repo      Repository
	JWTSecret []byte
}

func NewService(repo Repository, jwtSecret []byte) *service {
	return &service{Repo: repo, JWTSecret: jwtSecret}
}