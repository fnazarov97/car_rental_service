package postgres

import (
	"car_rental/genprotos/article"
	"car_rental/genprotos/author"
	"errors"
	"time"
)

// AddArticle ...
func (p Postgres) AddArticle(id string, req *article.AddArticleReq) error {
	Id := &author.Id{
		Id: req.AuthorId,
	}
	_, err := p.GetAuthorByID(Id)
	if err != nil {
		return err
	}

	if req.Content == nil {
		req.Content = &article.AddArticleReq_Post{}
	}
	_, err = p.DB.Exec(`Insert into article(id, title, body, author_id, created_at) 
						 VALUES($1, $2, $3, $4, now())
						`, id, req.Content.Title, req.Content.Body, req.AuthorId)
	if err != nil {
		return err
	}
	return nil
}

// GetArticleByID ...
func (p Postgres) GetArticleByID(id string) (*article.GetArticleByIdRes, error) {
	res := &article.GetArticleByIdRes{
		Content: &article.GetArticleByIdRes_Post{},
		Authori: &article.GetArticleByIdRes_Author{},
	}
	var deletedAt *time.Time
	var updatedAt, authorUpdatedAt *string
	err := p.DB.QueryRow(`SELECT 
		ar.id,
		ar.title,
		ar.body,
		ar.created_at,
		ar.updated_at,
		ar.deleted_at,
		au.id,
		au.fullname,
		au.created_at,
		au.updated_at
    FROM article AS ar JOIN author AS au ON ar.author_id = au.id WHERE ar.id = $1`, id).Scan(
		&res.Id,
		&res.Content.Title,
		&res.Content.Body,
		&res.CreatedAt,
		&updatedAt,
		&deletedAt,
		&res.Authori.Id,
		&res.Authori.Fullname,
		&res.Authori.CreatedAt,
		&authorUpdatedAt,
	)
	if err != nil {
		return res, err
	}

	if updatedAt != nil {
		res.UpdatedAt = *updatedAt
	}

	if authorUpdatedAt != nil {
		res.Authori.UpdatedAt = *authorUpdatedAt
	}

	if deletedAt != nil {
		return res, errors.New("article not found")
	}

	return res, err
}

// GetArticleList ...
func (p Postgres) GetArticleList(offset, limit int, search string) (*article.GetArticleListRes, error) {
	resp := &article.GetArticleListRes{
		Articles: make([]*article.AddArticleRes, 0),
	}

	rows, err := p.DB.Queryx(`SELECT
	id,
	title,
	body,
	author_id,
	created_at,
	updated_at
	FROM article WHERE deleted_at IS NULL AND ((title ILIKE '%' || $1 || '%') OR (body ILIKE '%' || $1 || '%'))
	LIMIT $2
	OFFSET $3
	`, search, limit, offset)
	if err != nil {
		return resp, err
	}

	for rows.Next() {
		a := &article.AddArticleRes{
			Content: &article.AddArticleRes_Post{},
		}

		var updatedAt *string

		err := rows.Scan(
			&a.Id,
			&a.Content.Title,
			&a.Content.Body,
			&a.AuthorId,
			&a.CreatedAt,
			&updatedAt,
		)
		if err != nil {
			return resp, err
		}

		if updatedAt != nil {
			a.UpdatedAt = *updatedAt
		}

		resp.Articles = append(resp.Articles, a)
	}

	return resp, err
}

// UpdateArticle ...
func (p Postgres) UpdateArticle(entity *article.UpdateArticleReq) error {
	if entity.Content == nil {
		entity.Content = &article.UpdateArticleReq_Post{}
	}
	res, err := p.DB.NamedExec("UPDATE article  SET title=:t, body=:b, updated_at=now() WHERE deleted_at IS NULL AND id=:id", map[string]interface{}{
		"id": entity.Id,
		"t":  entity.Content.Title,
		"b":  entity.Content.Body,
	})
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if n > 0 {
		return nil
	}

	return errors.New("article not found")
}

// DeleteArticle ...
func (p Postgres) DeleteArticle(id string) error {
	res, err := p.DB.Exec("UPDATE article SET deleted_at=now() WHERE id=$1 AND deleted_at IS NULL", id)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if n > 0 {
		return nil
	}

	return errors.New("article not found")
}
