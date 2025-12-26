package study

import (
	"context"
	"fmt"
	"time"

	db "local-lib/database"
	"study-service/internal/model"
	"study-service/internal/repository"
	"study-service/internal/repository/study/converter"
	repoModel "study-service/internal/repository/study/model"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

const (
	tableName = "study"

	idColumn                  = "id"
	nameColumn                = "name"
	statusColumn              = "status"
	ownerIDColumn             = "owner_id"
	imageIDColumn             = "image_id"
	predictedClassIDColumn    = "predicted_class_id"
	predictedClassScoreColumn = "predicted_class_score"
	errorDescriptionColumn    = "error_description"
	createdAtColumn           = "created_at"
	updatedAtColumn           = "updated_at"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.StudyRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, study *model.Study) (uuid.UUID, error) {
	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(
			idColumn,
			nameColumn,
			statusColumn,
			ownerIDColumn,
			imageIDColumn,
			predictedClassIDColumn,
			predictedClassScoreColumn,
			createdAtColumn,
		).
		Values(
			uuid.New(),
			study.Name,
			study.Status,
			study.OwnerID,
			study.ImageID,
			study.PredictedClassID,
			study.PredictedScore,
			time.Now().UTC(),
		).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("create study error: %w", err)
	}

	q := db.Query{
		Name:     "study_repository.Create",
		QueryRaw: query,
	}

	var id uuid.UUID
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&id)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("create study error: %w", err)
	}

	return id, nil
}

func (r *repo) Get(ctx context.Context, id uuid.UUID) (*model.Study, error) {
	builder := sq.Select(
		idColumn,
		nameColumn,
		statusColumn,
		ownerIDColumn,
		imageIDColumn,
		predictedClassIDColumn,
		predictedClassScoreColumn,
		errorDescriptionColumn,
		createdAtColumn,
		updatedAtColumn,
	).
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("get study error: %w", err)
	}

	q := db.Query{
		Name:     "study_repository.Get",
		QueryRaw: query,
	}

	var study repoModel.Study
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(
		&study.ID,
		&study.Name,
		&study.Status,
		&study.OwnerID,
		&study.ImageID,
		&study.PredictedClassID,
		&study.PredictedScore,
		&study.ErrorDescription,
		&study.CreatedAt,
		&study.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get study error: %w", err)
	}

	return converter.ToStudyFromRepo(&study), nil
}

func (r *repo) Update(ctx context.Context, s *model.Study) (*model.Study, error) {
	builder := sq.Update(tableName).
		PlaceholderFormat(sq.Dollar).
		Set(predictedClassIDColumn, s.PredictedClassID).
		Set(updatedAtColumn, time.Now().UTC()).
		Where(sq.Eq{idColumn: s.ID}).
		Suffix(`RETURNING 
					id,
					name,
					status,
					owner_id,
					image_id,
					predicted_class_id,
					predicted_class_score,
					error_description,
					created_at,
					updated_at`)

	if s.Name != nil {
		builder = builder.Set(nameColumn, s.Name)
	}
	if s.Status != "" {
		builder = builder.Set(statusColumn, s.Status)
	}
	if s.PredictedScore != nil {
		builder = builder.Set(predictedClassScoreColumn, s.PredictedScore)
	}
	if s.ErrorDescription != nil {
		builder = builder.Set(errorDescriptionColumn, s.ErrorDescription)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("update study error: %w", err)
	}

	q := db.Query{
		Name:     "study_repository.Update",
		QueryRaw: query,
	}

	var study repoModel.Study
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(
		&study.ID,
		&study.Name,
		&study.Status,
		&study.OwnerID,
		&study.ImageID,
		&study.PredictedClassID,
		&study.PredictedScore,
		&study.ErrorDescription,
		&study.CreatedAt,
		&study.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("update study error: %w", err)
	}
	return converter.ToStudyFromRepo(&study), nil
}

//type Filter struct {
//}
//
//func (r *repo) List(ctx context.Context, filter *Filter, limit, offset int) ([]*model.Study, error) {
//	builder := sq.Select(
//		idColumn,
//		nameColumn,
//		statusIDColumn,
//		ownerIDColumn,
//		imageIDColumn,
//		predictedClassIDColumn,
//		predictedClassScoreColumn,
//		createdAtColumn,
//		updatedAtColumn,
//	).
//		From(tableName).
//		PlaceholderFormat(sq.Dollar).
//		Limit(uint64(limit)).
//		Offset(uint64(offset))
//
//	if filter != nil {
//		if filter.StudyID != nil {
//			builder = builder.Where(sq.Eq{idColumn: *filter.StudyID})
//		}
//		if filter.UserID != nil {
//			builder = builder.Where(sq.Eq{ownerIDColumn: *filter.UserID})
//		}
//	}
//
//	query, args, err := builder.ToSql()
//	if err != nil {
//		return nil, fmt.Errorf("list studies error: %w", err)
//	}
//
//	q := db.Query{
//		Name:     "study_repository.List",
//		QueryRaw: query,
//	}
//
//	rows, err := r.db.DB().QueryContext(ctx, q, args...)
//	if err != nil {
//		return nil, fmt.Errorf("list studies error: %w", err)
//	}
//	defer rows.Close()
//
//	var studies []*model.Study
//	for rows.Next() {
//		var study model.Study
//		err := rows.Scan(
//			&study.ID,
//			&study.Name,
//			&study.StatusID,
//			&study.OwnerID,
//			&study.ImageID,
//			&study.PredictedClassID,
//			&study.PredictedScore,
//			&study.CreatedAt,
//			&study.UpdatedAt,
//		)
//		if err != nil {
//			return nil, fmt.Errorf("list studies scan error: %w", err)
//		}
//		studies = append(studies, &study)
//	}
//
//	if err = rows.Err(); err != nil {
//		return nil, fmt.Errorf("list studies rows error: %w", err)
//	}
//
//	return studies, nil
//}

func (r *repo) Delete(ctx context.Context, id uuid.UUID) error {
	builder := sq.Delete(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id})

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("delete study error: %w", err)
	}

	q := db.Query{
		Name:     "study_repository.Delete",
		QueryRaw: query,
	}

	result, err := r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("delete study error: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("study not found")
	}

	return nil
}
