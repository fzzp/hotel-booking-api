package db

// Repository 是由各个模块的CURD功能组成，整合一起，提供给service层使用
type Repository struct {
	UserRepo    UserRepo
	SessionRepo SessionRepo
}

// NewRepository 实例化，返回所有模块db操作
func NewRepository(qb Queryable) *Repository {
	return &Repository{
		UserRepo:    NewUserRepo(qb),
		SessionRepo: NewSessionRepo(qb),
	}
}
