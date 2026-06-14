package repository

import "oa-nsdiy/backend/ent"

type Repositories struct {
	User         *UserRepository
	Role         *RoleRepository
	Permission   *PermissionRepository
	NewsGroup    *NewsGroupRepository
	News         *NewsRepository
	Article      *ArticleRepository
	Project      *ProjectRepository
	Task         *TaskRepository
	MediaAccount *MediaAccountRepository
	MediaContent *MediaContentRepository
	File         *FileRepository
	ApiToken     *ApiTokenRepository
}

func NewRepositories(client *ent.Client) *Repositories {
	return &Repositories{
		User:         NewUserRepository(client),
		Role:         NewRoleRepository(client),
		Permission:   NewPermissionRepository(client),
		NewsGroup:    NewNewsGroupRepository(client),
		News:         NewNewsRepository(client),
		Article:      NewArticleRepository(client),
		Project:      NewProjectRepository(client),
		Task:         NewTaskRepository(client),
		MediaAccount: NewMediaAccountRepository(client),
		MediaContent: NewMediaContentRepository(client),
		File:         NewFileRepository(client),
		ApiToken:     NewApiTokenRepository(client),
	}
}
