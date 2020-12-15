package persistence

import (
	"database/sql"
	"github.com/jackc/pgx"
	"tech-db-project/domain/entity"
	"time"
)

type ForumDB struct {
	db *pgx.ConnPool
}

func NewForumDB(db *pgx.ConnPool) *ForumDB {
	return &ForumDB{db: db}
}

func (forumDB *ForumDB) InsertInto(forum *entity.Forum) error {
	tx, err := forumDB.db.Begin()
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

	row := tx.QueryRow("forum_insert_into", forum.Slug, forum.Title, forum.User)

	var info string
	if err = row.Scan(&info); err != nil {
		return err
	}

	return nil
}

func (forumDB *ForumDB) GetBySlug(forum *entity.Forum) error {
	tx, err := forumDB.db.Begin()
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

	row := tx.QueryRow("forum_get_by_slug", forum.Slug)

	if err := row.Scan(&forum.Posts, &forum.Slug, &forum.Threads, &forum.Title, &forum.User); err != nil {
		return err
	}

	return nil
}

func (forumDB *ForumDB) GetThreads(forum *entity.Forum, desc, limit, since string) (entity.Threads, error) {
	tx, err := forumDB.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err == nil {
			_ = tx.Commit()
		} else {
			_ = tx.Rollback()
		}
	}()

	threads := make([]entity.Thread, 0)
	var rows *pgx.Rows

	if since == "" {
		if desc == "true" {
			since = "infinity"
		} else {
			since = "-infinity"
		}
	}
	if desc == "true" {
		if limit != "" {
			rows, err = tx.Query("forum_get_threads_desc_with_limit", forum.Slug, since, limit)
		} else {
			rows, err = tx.Query("forum_get_threads_desc", forum.Slug, since)
		}
	} else {
		if limit != "" {
			rows, err = tx.Query("forum_get_threads_asc_with_limit", forum.Slug, since, limit)
		} else {
			rows, err = tx.Query("forum_get_threads_asc", forum.Slug, since)
		}
	}

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		created := sql.NullTime{}
		slug := sql.NullString{}
		thread := entity.Thread{}
		votes := sql.NullInt64{}

		if err = rows.Scan(&thread.Id, &thread.Title, &thread.Message, &created, &slug, &thread.Author, &thread.Forum, &votes); err != nil {
			return nil, err
		}

		if slug.Valid {
			thread.Slug = slug.String
		}
		if votes.Valid {
			thread.Votes = votes.Int64
		}
		if created.Valid {
			thread.Created = created.Time.Format(time.RFC3339Nano)
		}

		threads = append(threads, thread)
	}
	rows.Close()

	return threads, nil
}

func (forumDB *ForumDB) GetUsers(forum *entity.Forum, desc, limit, since string) (entity.Users, error) {
	tx, err := forumDB.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err == nil {
			_ = tx.Commit()
		} else {
			_ = tx.Rollback()
		}
	}()

	users := make([]entity.User, 0)
	var rows *pgx.Rows

	switch true {
	case desc != "true" && since == "" && limit == "":
		rows, err = tx.Query("forum_get_users", forum.Slug)

	case desc == "true" && since == "" && limit == "":
		rows, err = tx.Query("forum_get_users_desc", forum.Slug)

	case desc != "true" && since != "" && limit == "":
		rows, err = tx.Query("forum_get_users_asc_with_since", forum.Slug, since)

	case desc == "true" && since != "" && limit == "":
		rows, err = tx.Query("forum_get_users_desc_with_since", forum.Slug, since)

	case desc != "true" && since == "" && limit != "":
		rows, err = tx.Query("forum_get_users_with_limit", forum.Slug, limit)

	case desc == "true" && since == "" && limit != "":
		rows, err = tx.Query("forum_get_users_desc_with_limit", forum.Slug, limit)

	case desc != "true" && since != "" && limit != "":
		rows, err = tx.Query("forum_get_users_asc_with_since_with_limit", forum.Slug, since, limit)

	case desc == "true" && since != "" && limit != "":
		rows, err = tx.Query("forum_get_users_desc_with_since_with_limit", forum.Slug, since, limit)

	}
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		user := entity.User{}
		if err := rows.Scan(&user.Email, &user.Fullname, &user.Nickname, &user.About); err != nil {
			rows.Close()
			return nil, err
		}
		users = append(users, user)
	}
	rows.Close()

	return users, nil
}

func (forumDB *ForumDB) Prepare() error {
	_, err := forumDB.db.Prepare("forum_insert_into",
		"INSERT INTO forum (slug, title, usr) "+
			"VALUES ($1, $2, $3) "+
			"RETURNING title",
	)
	if err != nil {
		return err
	}

	_, err = forumDB.db.Prepare("forum_get_by_slug",
		"SELECT f.posts, f.slug, f.threads,f.title, f.usr "+
			"FROM forum f "+
			"WHERE f.slug = $1 ",
	)
	if err != nil {
		return err
	}

	_, err = forumDB.db.Prepare("forum_get_threads_desc",
		"SELECT t.id, t.title, t.message, t.created, t.slug, t.usr, t.forum, t.votes FROM thread t "+
			"WHERE t.forum = $1 AND t.created <=  $2::timestamptz "+
			"ORDER BY t.created DESC ",
	)
	if err != nil {
		return err
	}

	_, err = forumDB.db.Prepare("forum_get_threads_desc_with_limit",
		"SELECT t.id, t.title, t.message, t.created, t.slug, t.usr, t.forum, t.votes FROM thread t "+
			"WHERE t.forum = $1 AND t.created <=  $2::timestamptz "+
			"ORDER BY t.created DESC "+
			"LIMIT $3",
	)
	if err != nil {
		return err
	}

	_, err = forumDB.db.Prepare("forum_get_threads_asc",
		"SELECT t.id, t.title, t.message, t.created, t.slug, t.usr, t.forum, t.votes FROM thread t "+
			"WHERE t.forum = $1 AND t.created >=  $2::timestamptz "+
			"ORDER BY t.created ",
	)
	if err != nil {
		return err
	}
	_, err = forumDB.db.Prepare("forum_get_threads_asc_with_limit",
		"SELECT t.id, t.title, t.message, t.created, t.slug, t.usr, t.forum, t.votes FROM thread t "+
			"WHERE t.forum = $1 AND t.created >=  $2::timestamptz "+
			"ORDER BY t.created "+
			"LIMIT $3 ",
	)
	if err != nil {
		return err
	}

	_, err = forumDB.db.Prepare("forum_get_users",
		"SELECT u.email, u.fullname, u.nickname, u.about "+
			"FROM forum_users "+
			"JOIN usr u on forum_users.nickname = u.nickname "+
			"WHERE forum = $1 "+
			"ORDER BY u.nickname ",
	)
	if err != nil {
		return err
	}

	_, err = forumDB.db.Prepare("forum_get_users_with_limit",
		"SELECT u.email, u.fullname, u.nickname, u.about "+
			"FROM forum_users "+
			"JOIN usr u on forum_users.nickname = u.nickname "+
			"WHERE forum = $1 "+
			"ORDER BY u.nickname "+
			"LIMIT $2 ",
	)
	if err != nil {
		return err
	}

	_, err = forumDB.db.Prepare("forum_get_users_desc",
		"SELECT u.email, u.fullname, u.nickname, u.about "+
			"FROM forum_users "+
			"JOIN usr u on forum_users.nickname = u.nickname "+
			"WHERE forum = $1 "+
			"ORDER BY u.nickname DESC ",
	)
	if err != nil {
		return err
	}

	_, err = forumDB.db.Prepare("forum_get_users_desc_with_limit",
		"SELECT u.email, u.fullname, u.nickname, u.about "+
			"FROM forum_users "+
			"JOIN usr u on forum_users.nickname = u.nickname "+
			"WHERE forum = $1 "+
			"ORDER BY u.nickname DESC "+
			"LIMIT $2 ",
	)
	if err != nil {
		return err
	}

	_, err = forumDB.db.Prepare("forum_get_users_desc_with_since_with_limit",
		"SELECT u.email, u.fullname, u.nickname, u.about "+
			"FROM forum_users "+
			"JOIN usr u on forum_users.nickname = u.nickname "+
			"WHERE forum = $1 AND u.nickname < $2 "+
			"ORDER BY u.nickname DESC "+
			"LIMIT $3 ",
	)
	if err != nil {
		return err
	}

	_, err = forumDB.db.Prepare("forum_get_users_desc_with_since",
		"SELECT u.email, u.fullname, u.nickname, u.about "+
			"FROM forum_users "+
			"JOIN usr u on forum_users.nickname = u.nickname "+
			"WHERE forum = $1 AND u.nickname < $2 "+
			"ORDER BY u.nickname DESC",
	)
	if err != nil {
		return err
	}

	_, err = forumDB.db.Prepare("forum_get_users_asc_with_since_with_limit",
		"SELECT u.email, u.fullname, u.nickname, u.about "+
			"FROM forum_users "+
			"JOIN usr u on forum_users.nickname = u.nickname "+
			"WHERE forum = $1 AND u.nickname > $2 "+
			"ORDER BY u.nickname "+
			"LIMIT $3 ",
	)
	if err != nil {
		return err
	}

	_, err = forumDB.db.Prepare("forum_get_users_asc_with_since",
		"SELECT u.email, u.fullname, u.nickname, u.about "+
			"FROM forum_users "+
			"JOIN usr u on forum_users.nickname = u.nickname "+
			"WHERE forum = $1 AND u.nickname > $2 "+
			"ORDER BY u.nickname ",
	)
	if err != nil {
		return err
	}

	return nil
}
