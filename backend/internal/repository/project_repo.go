package repository

import (
	"context"
	"fmt"

	"oa-nsdiy/backend/ent"
	"oa-nsdiy/backend/ent/project"
	"oa-nsdiy/backend/ent/projectmember"
	"oa-nsdiy/backend/internal/domain"
)

type ProjectRepository struct {
	client *ent.Client
}

func NewProjectRepository(client *ent.Client) *ProjectRepository {
	return &ProjectRepository{client: client}
}

// Type aliases for backward compatibility
type Project = domain.Project
type ProjectDetail = domain.ProjectDetail
type ProjectMember = domain.ProjectMember
type ProjectMemberDetail = domain.ProjectMemberDetail

func (r *ProjectRepository) GetByID(ctx context.Context, id int) (*Project, error) {
	e, err := r.client.Project.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return toProject(e), nil
}

func (r *ProjectRepository) GetByProjectNo(ctx context.Context, projectNo string) (*Project, error) {
	e, err := r.client.Project.Query().Where(project.ProjectNoEQ(projectNo)).First(ctx)
	if err != nil {
		return nil, err
	}
	return toProject(e), nil
}

func (r *ProjectRepository) GetDetailByID(ctx context.Context, id int) (*ProjectDetail, error) {
	e, err := r.client.Project.Query().Where(project.IDEQ(id)).WithOwner().First(ctx)
	if err != nil {
		return nil, err
	}
	detail := &ProjectDetail{Project: *toProject(e)}
	if e.Edges.Owner != nil {
		detail.OwnerNickname = stringPtr(e.Edges.Owner.Nickname)
	}
	return detail, nil
}

func (r *ProjectRepository) GenerateUniqueProjectNo(ctx context.Context) (string, error) {
	count, err := r.client.Project.Query().Count(ctx)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("PRJ%04d", count+1), nil
}

func (r *ProjectRepository) Create(ctx context.Context, p *Project) error {
	e, err := r.client.Project.Create().
		SetName(p.Name).
		SetProjectNo(p.ProjectNo).
		SetNillableDescription(p.Description).
		SetStatus(p.Status).
		SetPriority(p.Priority).
		SetNillableExpectedStartDate(p.ExpectedStartDate).
		SetNillableExpectedEndDate(p.ExpectedEndDate).
		SetNillableStartDate(p.StartDate).
		SetNillableEndDate(p.EndDate).
		SetOwnerID(p.OwnerID).
		Save(ctx)
	if err != nil {
		return err
	}
	p.ID = e.ID
	return nil
}

func (r *ProjectRepository) Update(ctx context.Context, p *Project) error {
	_, err := r.client.Project.UpdateOneID(p.ID).
		SetName(p.Name).
		SetNillableDescription(p.Description).
		SetStatus(p.Status).
		SetPriority(p.Priority).
		SetNillableExpectedStartDate(p.ExpectedStartDate).
		SetNillableExpectedEndDate(p.ExpectedEndDate).
		SetNillableStartDate(p.StartDate).
		SetNillableEndDate(p.EndDate).
		Save(ctx)
	return err
}

func (r *ProjectRepository) Delete(ctx context.Context, id int) error {
	return r.client.Project.DeleteOneID(id).Exec(ctx)
}

func (r *ProjectRepository) ListByUserID(ctx context.Context, userID int, keyword string, page, pageSize int) ([]*ProjectDetail, int64, error) {

	q := r.client.Project.Query().Where(
		project.Or(
			project.OwnerIDEQ(userID),
			project.HasMembersWith(projectmember.UserIDEQ(userID)),
		),
	)

	if keyword != "" {
		q.Where(
			project.Or(
				project.NameContains(keyword),
				project.ProjectNoContains(keyword),
			),
		)
	}

	count, err := q.Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	total := int64(count)

	offset := (page - 1) * pageSize
	entities, err := q.WithOwner().
		Order(ent.Desc(project.FieldID)).
		Limit(pageSize).
		Offset(offset).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	var projects []*ProjectDetail
	for _, e := range entities {
		detail := &ProjectDetail{Project: *toProject(e)}
		if e.Edges.Owner != nil {
			detail.OwnerNickname = stringPtr(e.Edges.Owner.Nickname)
		}
		projects = append(projects, detail)
	}

	return projects, total, nil
}

func (r *ProjectRepository) IsMember(ctx context.Context, projectID, userID int) (bool, error) {
	exist, err := r.client.ProjectMember.Query().
		Where(projectmember.ProjectIDEQ(projectID), projectmember.UserIDEQ(userID)).
		Exist(ctx)
	if err != nil {
		return false, err
	}
	return exist, nil
}

func (r *ProjectRepository) GetMember(ctx context.Context, projectID, userID int) (*ProjectMember, error) {
	m, err := r.client.ProjectMember.Query().
		Where(projectmember.ProjectIDEQ(projectID), projectmember.UserIDEQ(userID)).
		First(ctx)
	if err != nil {
		return nil, err
	}
	return toProjectMember(m), nil
}

func (r *ProjectRepository) AddMember(ctx context.Context, member *ProjectMember) error {
	e, err := r.client.ProjectMember.Create().
		SetProjectID(member.ProjectID).
		SetUserID(member.UserID).
		SetRole(member.Role).
		Save(ctx)
	if err != nil {
		return err
	}
	member.ID = e.ID
	return nil
}

func (r *ProjectRepository) UpdateMemberRole(ctx context.Context, projectID, userID int, role string) error {
	_, err := r.client.ProjectMember.Update().
		Where(projectmember.ProjectIDEQ(projectID), projectmember.UserIDEQ(userID)).
		SetRole(role).
		Save(ctx)
	return err
}

func (r *ProjectRepository) RemoveMember(ctx context.Context, projectID, userID int) error {
	_, err := r.client.ProjectMember.Delete().
		Where(projectmember.ProjectIDEQ(projectID), projectmember.UserIDEQ(userID)).
		Exec(ctx)
	return err
}

func (r *ProjectRepository) GetMembers(ctx context.Context, projectID int) ([]*ProjectMemberDetail, error) {
	members, err := r.client.ProjectMember.Query().
		Where(projectmember.ProjectIDEQ(projectID)).
		WithUser().
		Order(ent.Asc(projectmember.FieldJoinedAt)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	var result []*ProjectMemberDetail
	for _, m := range members {
		detail := &ProjectMemberDetail{
			ProjectMember: *toProjectMember(m),
		}
		if m.Edges.User != nil {
			detail.Username = m.Edges.User.Username
			detail.Nickname = stringPtr(m.Edges.User.Nickname)
			detail.Email = m.Edges.User.Email
			detail.IsActive = m.Edges.User.IsActive
		}
		result = append(result, detail)
	}

	return result, nil
}

// --- Converters ---

func toProject(e *ent.Project) *Project {
	return &Project{
		ID:                e.ID,
		Name:              e.Name,
		ProjectNo:         e.ProjectNo,
		Description:       stringPtr(e.Description),
		Status:            e.Status,
		Priority:          e.Priority,
		ExpectedStartDate: timePtr(e.ExpectedStartDate),
		ExpectedEndDate:   timePtr(e.ExpectedEndDate),
		StartDate:         timePtr(e.StartDate),
		EndDate:           timePtr(e.EndDate),
		OwnerID:           e.OwnerID,
		CreatedAt:         e.CreatedAt,
		UpdatedAt:         e.UpdatedAt,
	}
}

func toProjectMember(e *ent.ProjectMember) *ProjectMember {
	return &ProjectMember{
		ID:        e.ID,
		ProjectID: e.ProjectID,
		UserID:    e.UserID,
		Role:      e.Role,
		JoinedAt:  e.JoinedAt,
	}
}
