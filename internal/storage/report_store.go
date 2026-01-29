package storage

import (
	"context"
	"database/sql"
	"errors"

	global_varables "github.com/sirUnchained/my-go-instagram/internal/global"
	"github.com/sirUnchained/my-go-instagram/internal/payloads"
	"github.com/sirUnchained/my-go-instagram/internal/storage/models"
)

type reportStore struct {
	db *sql.DB
}

func (rs *reportStore) Create(ctx context.Context, creatorId int64, reportP payloads.CreateReportPayload) error {
	query := `INSERT INTO reports (creator, post, comment, reason, content) VALUES ($1, $2, $3, $4);`

	_, err := rs.db.ExecContext(ctx, query, creatorId, reportP.PostID, reportP.CommentID, reportP.Reason, reportP.Content)

	return err
}

func (rs *reportStore) GetReports(ctx context.Context, limit, offset int64) ([]models.ReportModel, error) {
	query := `
	SELECT 
			r.id, 
			r.creator, 
			r.post, 
			r.comment, 
			r.reason, 
			r.content, 
			r.created_at,
			u.username,
			u.email
	FROM 	reports 	AS r 
	JOIN 	users 		AS u ON u.id = r.creator
	LIMIT $1 OFFSET $2;
	`

	var reports []models.ReportModel

	rows, err := rs.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var report models.ReportModel
		report.Creator = models.UserModel{}
		report.Post = &models.PostModel{}
		report.Comment = &models.CommentModel{}

		err := rows.Scan(
			report.ID,
			report.Creator.Id,
			report.Post.Id,
			report.Comment.ID,
			report.Reason,
			report.Content,
			report.CreatedAt,
			report.Creator.Username,
			report.Creator.Email,
		)

		if err != nil {
			return nil, err
		}

		reports = append(reports, report)
	}

	return reports, err
}

func (rs *reportStore) Delete(ctx context.Context, reportId int64) error {
	query := `DELETE FROM reports WHERE id = $1;`

	_, err := rs.db.ExecContext(ctx, query, reportId)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return global_varables.NOT_FOUND_ROW
		default:
			return err
		}
	}

	return nil
}
