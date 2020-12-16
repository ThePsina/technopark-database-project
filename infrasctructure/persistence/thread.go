package persistence

import (
	"database/sql"
	"github.com/jackc/pgx"
	"strconv"
	"tech-db-project/domain/entity"
	"time"
)

type ThreadDB struct {
	db *pgx.ConnPool
}

func NewThreadDB(db *pgx.ConnPool) *ThreadDB {
	return &ThreadDB{db: db}
}

func (threadDB *ThreadDB) InsertIntoForumUsers(forum, nickname string) error {
	tx, err := threadDB.db.Begin()
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

	var buffer string
	err = tx.QueryRow("get_forum_user", forum, nickname).Scan(&buffer)
	if err != nil {
		_, err = threadDB.db.Exec("forum_users_insert_into", forum, nickname)
		if err != nil {
			return err
		}
	}

	return nil
}

func (threadDB *ThreadDB) InsertInto(thread *entity.Thread) error {
	tx, err := threadDB.db.Begin()
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

	slug := &sql.NullString{}
	if thread.Slug != "" {
		slug.String = thread.Slug
		slug.Valid = true
	}

	created := &sql.NullString{}
	if thread.Created != "" {
		created.String = thread.Created
		created.Valid = true
	}

	row := tx.QueryRow("thread_insert_into", thread.Author, created, thread.Forum, thread.Message, thread.Title, slug)
	if err := row.Scan(&thread.Id); err != nil {
		return err
	}

	if err = threadDB.InsertIntoForumUsers(thread.Forum, thread.Author); err != nil {
		return err
	}

	return nil
}

func (threadDB *ThreadDB) GetBySlug(thread *entity.Thread) error {
	tx, err := threadDB.db.Begin()
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

	row := tx.QueryRow("thread_get_by_slug", thread.Slug)

	created := sql.NullTime{}
	slug := sql.NullString{}

	err = row.Scan(
		&thread.Id,
		&thread.Title,
		&thread.Message,
		&created,
		&slug,
		&thread.Author,
		&thread.Forum,
		&thread.Votes,
	)

	if err != nil {
		return err
	}

	if created.Valid {
		thread.Created = created.Time.Format(time.RFC3339Nano)
	}

	if slug.Valid {
		thread.Slug = slug.String
	} else {
		thread.Slug = ""
	}

	return nil
}

func (threadDB *ThreadDB) GetById(thread *entity.Thread) error {
	tx, err := threadDB.db.Begin()
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

	row := tx.QueryRow("thread_get_by_id", thread.Id)

	created := sql.NullTime{}
	slug := sql.NullString{}

	err = row.Scan(
		&thread.Id,
		&thread.Title,
		&thread.Message,
		&created,
		&slug,
		&thread.Author,
		&thread.Forum,
		&thread.Votes,
	)
	if err != nil {
		return err
	}

	if created.Valid {
		thread.Created = created.Time.Format(time.RFC3339Nano)
	}

	if slug.Valid {
		thread.Slug = slug.String
	} else {
		thread.Slug = ""
	}

	return nil
}

func (threadDB *ThreadDB) GetBySlugOrId(thread *entity.Thread) error {
	tx, err := threadDB.db.Begin()
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

	Id, err := strconv.ParseInt(thread.Slug, 10, 64)
	if err == nil {
		thread.Id = Id
		thread.Slug = ""
	}

	if thread.Slug != "" {
		err = threadDB.GetBySlug(thread)
	} else {
		err = threadDB.GetById(thread)
	}

	if err != nil {
		return err
	}

	return nil
}

func (threadDB *ThreadDB) InsertIntoVotes(thread *entity.Thread, vote *entity.Vote) error {
	tx, err := threadDB.db.Begin()
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

	voteNum := int32(0)
	err = tx.QueryRow("votes_get_info", vote.Nickname, vote.Thread).Scan(&voteNum)

	if voteNum == 0 {
		err = tx.QueryRow("votes_insert_into", vote.Vote, vote.Nickname, vote.Thread).Scan(&vote.Thread)
		thread.Votes += int64(vote.Vote)
	} else {
		if voteNum != vote.Vote {
			err = tx.QueryRow("votes_update", vote.Vote, vote.Nickname, vote.Thread).Scan(&vote.Thread)
			thread.Votes += 2 * int64(vote.Vote)
		}
	}

	if err != nil {
		return err
	}

	return nil
}

func (threadDB *ThreadDB) Update(thread *entity.Thread) error {
	tx, err := threadDB.db.Begin()
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

	slug := sql.NullString{}
	created := sql.NullTime{}
	votes := sql.NullInt64{}

	switch true {
	case thread.Message == "" && thread.Title == "":
		err = threadDB.GetBySlugOrId(thread)
	case thread.Message != "" && thread.Title == "":
		err = tx.QueryRow("thread_update_message",
			thread.Message,
			thread.Slug,
		).Scan(
			&thread.Id,
			&thread.Title,
			&thread.Message,
			&created,
			&slug,
			&thread.Author,
			&thread.Forum,
			&votes,
		)
	case thread.Message == "" && thread.Title != "":
		err = tx.QueryRow("thread_update_title",
			thread.Title,
			thread.Slug,
		).Scan(
			&thread.Id,
			&thread.Title,
			&thread.Message,
			&created,
			&slug,
			&thread.Author,
			&thread.Forum,
			&votes,
		)
	case thread.Message != "" && thread.Title != "":
		err = tx.QueryRow("thread_update_all",
			thread.Message,
			thread.Title,
			thread.Slug,
		).Scan(
			&thread.Id,
			&thread.Title,
			&thread.Message,
			&created,
			&slug,
			&thread.Author,
			&thread.Forum,
			&votes,
		)
	}
	if err != nil {
		return err
	}

	if created.Valid {
		thread.Created = created.Time.Format(time.RFC3339Nano)
	}

	if slug.Valid {
		thread.Slug = slug.String
	}

	if votes.Valid {
		thread.Votes = votes.Int64
	}

	return nil
}

func (threadDB *ThreadDB) GetPosts(thread *entity.Thread, desc, sort, limit, since string) (entity.Posts, error) {
	tx, err := threadDB.db.Begin()
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

	posts := make([]entity.Post, 0)
	var rows *pgx.Rows

	if sort == "tree" {
		switch true {
		case desc != "true" && since == "" && limit == "":
			rows, err = tx.Query("thread_posts_tree_asc", thread.Id)

		case desc == "true" && since == "" && limit == "":
			rows, err = tx.Query("thread_posts_tree_desc", thread.Id)

		case desc != "true" && since != "" && limit == "":
			rows, err = tx.Query("thread_posts_tree_asc_with_since", thread.Id, since)

		case desc == "true" && since != "" && limit == "":
			rows, err = tx.Query("thread_posts_tree_desc_with_since", thread.Id, since)

		case desc != "true" && since == "" && limit != "":
			rows, err = tx.Query("thread_posts_tree_asc_with_limit", thread.Id, limit)

		case desc == "true" && since == "" && limit != "":
			rows, err = tx.Query("thread_posts_tree_desc_with_limit", thread.Id, limit)

		case desc != "true" && since != "" && limit != "":
			rows, err = tx.Query("thread_posts_tree_asc_with_since_with_limit", thread.Id, since, limit)

		case desc == "true" && since != "" && limit != "":
			rows, err = tx.Query("thread_posts_tree_desc_with_since_with_limit", thread.Id, since, limit)
		}
	} else if sort == "parent_tree" {
		switch true {
		case desc != "true" && since == "" && limit == "":
			rows, err = tx.Query("thread_posts_parent_asc", thread.Id)

		case desc == "true" && since == "" && limit == "":
			rows, err = tx.Query("thread_posts_parent_desc", thread.Id)

		case desc != "true" && since != "" && limit == "":
			rows, err = tx.Query("thread_posts_parent_asc_with_since", thread.Id, since)

		case desc == "true" && since != "" && limit == "":
			rows, err = tx.Query("thread_posts_parent_desc_with_since", thread.Id, since)

		case desc != "true" && since == "" && limit != "":
			rows, err = tx.Query("thread_posts_parent_asc_with_limit", thread.Id, limit)

		case desc == "true" && since == "" && limit != "":
			rows, err = tx.Query("thread_posts_parent_desc_with_limit", thread.Id, limit)

		case desc != "true" && since != "" && limit != "":
			rows, err = tx.Query("thread_posts_parent_asc_with_since_with_limit", thread.Id, since, limit)

		case desc == "true" && since != "" && limit != "":
			rows, err = tx.Query("thread_posts_parent_desc_with_since_with_limit", thread.Id, since, limit)
		}
	} else {
		switch true {
		case desc != "true" && since == "" && limit == "":
			rows, err = tx.Query("thread_post_flat_asc", thread.Id)

		case desc == "true" && since == "" && limit == "":
			rows, err = tx.Query("thread_post_flat_desc", thread.Id)

		case desc != "true" && since != "" && limit == "":
			rows, err = tx.Query("thread_post_flat_asc_with_since", thread.Id, since)

		case desc == "true" && since != "" && limit == "":
			rows, err = tx.Query("thread_post_flat_desc_with_since", thread.Id, since)

		case desc != "true" && since == "" && limit != "":
			rows, err = tx.Query("thread_post_flat_asc_with_limit", thread.Id, limit)

		case desc == "true" && since == "" && limit != "":
			rows, err = tx.Query("thread_post_flat_desc_with_limit", thread.Id, limit)

		case desc != "true" && since != "" && limit != "":
			rows, err = tx.Query("thread_post_flat_asc_with_since_with_limit", thread.Id, since, limit)

		case desc == "true" && since != "" && limit != "":
			rows, err = tx.Query("thread_post_flat_desc_with_since_with_limit", thread.Id, since, limit)
		}
	}

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		created := sql.NullTime{}
		p := entity.Post{}

		err := rows.Scan(&p.Id, &p.Author, &created, &p.Forum, &p.IsEdited, &p.Message, &p.Parent, &p.Thread)
		if err != nil {
			return nil, err
		}

		if created.Valid {
			p.Created = created.Time.Format(time.RFC3339Nano)
		}

		posts = append(posts, p)
	}
	rows.Close()

	return posts, nil
}

func (threadDB *ThreadDB) Prepare() error {
	_, err := threadDB.db.Prepare("thread_insert_into",
		"insert into thread (usr, created, forum, message, title, slug) values ($1, $2, $3, $4, $5, $6)"+
			"on conflict do nothing "+
			"returning id",
	)
	if err != nil {
		return err
	}

	_, err = threadDB.db.Prepare(
		"get_forum_user",
		"select nickname from forum_users where forum = $1 and nickname = $2 ")
	if err != nil {
		return err
	}

	_, err = threadDB.db.Prepare("forum_users_insert_into","insert into forum_users (forum, nickname) " +
		"values ($1, $2) on conflict do nothing")
	if err != nil {
		return err
	}

	_, err = threadDB.db.Prepare("thread_get_by_slug",
		"select t.id, t.title, t.message, t.created, t.slug, t.usr, t.forum, t.votes from thread t where t.slug = $1",
	)
	if err != nil {
		return err
	}

	_, err = threadDB.db.Prepare("thread_get_by_id",
		"select t.id, t.title, t.message, t.created, t.slug, t.usr, t.forum, t.votes from thread t where t.id = $1 ",
	)
	if err != nil {
		return err
	}

	_, err = threadDB.db.Prepare("votes_insert_into",
		"insert into vote (vote, usr, thread) VALUES ($1 , $2, $3) returning thread",
	)
	if err != nil {
		return err
	}

	_, err = threadDB.db.Prepare("votes_update",
		"update vote set vote = $1 where usr = $2 and thread = $3 returning thread",
	)
	if err != nil {
		return err
	}

	_, err = threadDB.db.Prepare("votes_get_info","select vote from vote where usr = $1 and thread = $2 ")
	if err != nil {
		return err
	}

	_, err = threadDB.db.Prepare("thread_update_all",
		"update thread set message = $1, title = $2 "+
			"where id::citext = $3 or slug = $3 "+
			"returning id, title, message, created, slug, usr, forum, votes",
	)
	if err != nil {
		return err
	}

	_, err = threadDB.db.Prepare("thread_update_message",
		"update thread set message = $1 "+
			"where id::citext = $2 or slug = $2 "+
			"returning id, title, message, created, slug, usr, forum, votes",
	)
	if err != nil {
		return err
	}

	_, err = threadDB.db.Prepare("thread_update_title",
		"update thread set title = $1 "+
			"where id::citext = $2 or slug = $2 "+
			"returning id, title, message, created, slug, usr, forum, votes",
	)
	if err != nil {
		return err
	}

	_, err = threadDB.db.Prepare("thread_posts_tree_asc",
		"select p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread "+
			"from post p where p.thread = $1 order by p.path ",
	)
	if err != nil {
		return err
	}

	_, err = threadDB.db.Prepare("thread_posts_tree_desc",
		"select p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread "+
			"from post p where p.thread = $1 order by p.path DESC ",
	)
	if err != nil {
		return err
	}

	_, err = threadDB.db.Prepare("thread_posts_tree_asc_with_limit",
		"select p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread "+
			"from post p where p.thread = $1 order by p.path limit $2 ",
	)
	if err != nil {
		return err
	}

	_, err = threadDB.db.Prepare("thread_posts_tree_desc_with_limit",
		"select p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread "+
			"from post p where p.thread = $1 order by p.path DESC limit $2 ",
	)
	if err != nil {
		return err
	}

	_, err = threadDB.db.Prepare("thread_posts_tree_asc_with_since",
		"select p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread "+
			"from post p "+
			"where p.thread = $1 and p.path::bigint[] > (select path from post where id = $2 )::bigint[] "+
			"order by p.path ",
	)
	if err != nil {
		return err
	}

	_, err = threadDB.db.Prepare("thread_posts_tree_desc_with_since",
		"select p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread "+
			"from post p "+
			"where p.thread = $1 AND p.path::bigint[] < (select path from post where id = $2 )::bigint[] "+
			"order by p.path DESC ",
	)
	if err != nil {
		return err
	}

	_, err = threadDB.db.Prepare("thread_posts_tree_asc_with_since_with_limit",
		"select p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread "+
			"from post p "+
			"where p.thread = $1 and p.path::bigint[] > (select path from post where id = $2 )::bigint[] "+
			"order by p.path "+
			"limit $3",
	)
	if err != nil {
		return err
	}

	_, err = threadDB.db.Prepare("thread_posts_tree_desc_with_since_with_limit",
		"select p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread "+
			"from post p "+
			"where p.thread = $1 and p.path::bigint[] < (select path from post where id = $2 )::bigint[] "+
			"order by p.path desc "+
			"limit $3",
	)
	if err != nil {
		return err
	}

	_, err = threadDB.db.Prepare("thread_posts_parent_asc",
		"select p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread FROM "+
			"(select * from post p2 where p2.thread = $1 and p2.parent = 0 order by p2.path) "+
			"as prt "+
			"join post p on prt.path[1] = p.path[1] "+
			"order by p.path[1], p.path ",
	)
	if err != nil {
		return err
	}

	_, err = threadDB.db.Prepare("thread_posts_parent_desc",
		"select p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread FROM "+
			"(select * from post p2 where p2.thread = $1 and p2.parent = 0 order by p2.path DESC) "+
			"as prt "+
			"join post p on prt.path[1] = p.path[1] "+
			"order by p.path[1] desc, p.path ",
	)
	if err != nil {
		return err
	}

	_, err = threadDB.db.Prepare("thread_posts_parent_asc_with_limit",
		"select p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread FROM "+
			"(select * from post p2 where p2.thread = $1 and p2.parent = 0 order by p2.path limit $2) "+
			"as prt "+
			"join post p on prt.path[1] = p.path[1] "+
			"order by p.path[1], p.path ",
	)
	if err != nil {
		return err
	}

	_, err = threadDB.db.Prepare("thread_posts_parent_desc_with_limit",
		"select p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread FROM "+
			"(select * from post p2 where p2.thread = $1 and p2.parent = 0 order by p2.path DESC limit $2) "+
			"as prt "+
			"join post p on prt.path[1] = p.path[1] "+
			"order by p.path[1] desc, p.path ",
	)
	if err != nil {
		return err
	}

	_, err = threadDB.db.Prepare("thread_posts_parent_asc_with_since",
		"select p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread FROM "+
			"(select * from post p2 where p2.thread = $1 and p2.parent = 0 " +
			"and p2.path[1] > (select path[1] from post where id = $2 ) order by p2.path) "+
			"as prt "+
			"join post p on prt.path[1] = p.path[1] "+
			"order by p.path[1], p.path ",
	)
	if err != nil {
		return err
	}

	_, err = threadDB.db.Prepare("thread_posts_parent_desc_with_since",
		"select p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread FROM "+
			"(select * from post p2 where p2.thread = $1 and p2.parent = 0 "+
			"and p2.path[1] < (select path[1] from post where id = $2 ) order by p2.path desc) "+
			"as prt "+
			"join post p on prt.path[1] = p.path[1] "+
			"order by p.path[1] desc , p.path ",
	)
	if err != nil {
		return err
	}

	_, err = threadDB.db.Prepare("thread_posts_parent_asc_with_since_with_limit",
		"select p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread FROM "+
			"(select * from post p2 where p2.thread = $1 and p2.parent = 0 "+
			"and p2.path[1] > (select path[1] from post where id = $2 ) order by p2.path limit $3) "+
			"as prt "+
			"join post p on prt.path[1] = p.path[1] "+
			"order by p.path[1], p.path ",
	)
	if err != nil {
		return err
	}

	_, err = threadDB.db.Prepare("thread_posts_parent_desc_with_since_with_limit",
		"select p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread FROM "+
			"(select * from post p2 where p2.thread = $1 and p2.parent = 0 "+
			"and p2.path[1] < (select path[1] from post where id = $2 ) order by p2.path desc limit $3) "+
			"as prt "+
			"join post p on prt.path[1] = p.path[1] "+
			"order by p.path[1] desc, p.path ",
	)
	if err != nil {
		return err
	}

	_, err = threadDB.db.Prepare("thread_post_flat_asc",
		"select p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread "+
			"from post p where p.thread = $1 order by p.created, p.id",
	)
	if err != nil {
		return err
	}

	_, err = threadDB.db.Prepare("thread_post_flat_desc",
		"select p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread "+
			"from post p where p.thread = $1 order by p.created desc, p.id desc ",
	)
	if err != nil {
		return err
	}

	_, err = threadDB.db.Prepare("thread_post_flat_asc_with_limit",
		"select p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread "+
			"from post p where p.thread = $1 order by p.created, p.id limit $2 ",
	)
	if err != nil {
		return err
	}

	_, err = threadDB.db.Prepare("thread_post_flat_desc_with_limit",
		"select p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread "+
			"from post p where p.thread = $1 order by p.created desc , p.id desc limit $2",
	)
	if err != nil {
		return err
	}

	_, err = threadDB.db.Prepare("thread_post_flat_asc_with_since",
		"select p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread "+
			"from post p where p.thread = $1 and p.id > $2 order by p.created, p.id",
	)
	if err != nil {
		return err
	}

	_, err = threadDB.db.Prepare("thread_post_flat_desc_with_since",
		"select p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread "+
			"from post p where p.thread = $1 and p.id < $2 order by p.created desc, p.id desc ",
	)
	if err != nil {
		return err
	}

	_, err = threadDB.db.Prepare("thread_post_flat_asc_with_since_with_limit",
		"select p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread "+
			"from post p where p.thread = $1 and p.id > $2 order by p.created, p.id limit $3 ",
	)
	if err != nil {
		return err
	}

	_, err = threadDB.db.Prepare("thread_post_flat_desc_with_since_with_limit",
		"select p.id, p.usr, p.created, p.forum, p.isEdited, p.message, p.parent, p.thread "+
			"from post p where p.thread = $1 AND p.id < $2 order by p.created desc , p.id desc limit $3",
	)
	if err != nil {
		return err
	}

	return nil
}
