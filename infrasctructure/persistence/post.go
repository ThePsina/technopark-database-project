package persistence

import (
	"database/sql"
	"github.com/jackc/pgx"
	"github.com/sirupsen/logrus"
	"tech-db-project/domain/entity"
	"time"
)

type PostDB struct {
	db *pgx.ConnPool
}

func NewPostDB(db *pgx.ConnPool) *PostDB {
	return &PostDB {db: db}
}

func (postDB *PostDB) InsertInto(posts []*entity.Post, thread *entity.Thread) error {
	tx, err := postDB.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err == nil {
			_ = tx.Commit()
		} else {
			_ = tx.Rollback()
		}
	}()

	created := sql.NullTime{}
	for i, _ := range posts {
		posts[i].Thread = thread.Id
		posts[i].Forum = thread.Forum

		err = tx.QueryRow(
			"post_insert_into",
			posts[i].Author,
			posts[i].Message,
			posts[i].Parent,
			posts[i].Thread,
			posts[i].Forum).Scan(&posts[i].Id, &created)

		if err != nil {
			return err
		}

		if created.Valid {
			posts[i].Created = created.Time.Format(time.RFC3339Nano)
		}
	}

	if len(posts) > 0 {
		_, err := tx.Exec("forum_posts_update", len(posts), posts[0].Forum)
		if err != nil {
			logrus.Error("Error while update post count: " + err.Error())
		}
	}

	return nil
}

func (postDB *PostDB) GetById(post *entity.Post) error {
	tx, err := postDB.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err == nil {
			_ = tx.Commit()
		} else {
			_ = tx.Rollback()
		}
	}()

	created := sql.NullTime{}
	err = tx.QueryRow("post_get_by_id", post.Id).
		Scan(&post.Author, &created, &post.Forum, &post.IsEdited, &post.Message, &post.Parent, &post.Thread, &post.Path)
	if err != nil {
		return err
	}

	if created.Valid {
		post.Created = created.Time.Format(time.RFC3339Nano)
	}

	return nil
}

func (postDB *PostDB) Update(post *entity.Post) error {
	tx, err := postDB.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err == nil {
			_ = tx.Commit()
		} else {
			_ = tx.Rollback()
		}
	}()

	created := sql.NullTime{}
	err = tx.QueryRow("post_update", post.Message, post.Id).
		Scan(&post.Author, &created, &post.Forum, &post.IsEdited, &post.Message, &post.Parent, &post.Thread)
	if err != nil {
		return err
	}

	if created.Valid {
		post.Created = created.Time.Format(time.RFC3339Nano)
	}

	return nil
}

func (postDB *PostDB) Prepare() error {
	_, err := postDB.db.Prepare("post_insert_into",
		"insert into post (usr, message,  parent, thread, forum, created) "+
			"values ($1, $2, $3, $4, $5, current_timestamp) returning id, created",
	)
	if err != nil {
		return err
	}

	_, err = postDB.db.Prepare("post_get_by_id",
		"select p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread, p.path "+
			"from post p where p.id = $1",
	)
	if err != nil {
		return err
	}

	_, err = postDB.db.Prepare("post_update",
		"update post set message = $1, isEdited = true "+
			"where id = $2 returning usr, created, forum, isEdited, message, parent, thread",
	)
	if err != nil {
		return err
	}

	_, err = postDB.db.Prepare("forum_posts_update",
		"update forum  set posts = (posts + $1) where slug = $2",
	)
	if err != nil {
		return err
	}

	return nil
}
