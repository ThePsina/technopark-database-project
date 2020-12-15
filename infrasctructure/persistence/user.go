package persistence

import (
	"database/sql"
	"github.com/jackc/pgx"
	"tech-db-project/domain/entity"
)

type UserDB struct {
	db *pgx.ConnPool
}

func NewUserRepo(db *pgx.ConnPool) *UserDB {
	return &UserDB{db: db}
}

func (userDB *UserDB) InsertInto(user *entity.User) error {
	tx, err := userDB.db.Begin()
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

	var info string
	about := &sql.NullString{}

	if user.About != "" {
		about.String = user.About
		about.Valid = true
	}

	err = tx.QueryRow("user_insert", user.Email, user.Fullname, user.Nickname, about).Scan(&info)
	if err != nil {
		return err
	}

	return nil
}

func (userDB *UserDB) GetByNickname(user *entity.User) error {
	tx, err := userDB.db.Begin()
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

	row := tx.QueryRow("user_get_by_nickname", user.Nickname)

	fullname := &sql.NullString{}
	about := &sql.NullString{}

	if err := row.Scan(&user.Email, fullname, &user.Nickname, about); err != nil {
		return err
	}

	if about.Valid {
		user.About = about.String
	}

	if fullname.Valid {
		user.Fullname = fullname.String
	}

	return nil
}

func (userDB *UserDB) GetByNicknameOrEmail(user *entity.User) (entity.Users, error) {
	tx, err := userDB.db.Begin()
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

	rows, err := tx.Query("user_get_by_nickname_or_email", user.Nickname, user.Email)
	if err != nil {
		return nil, err
	}

	users := make([]entity.User, 0)
	for rows.Next() {
		fullname := &sql.NullString{}
		about := &sql.NullString{}

		if err := rows.Scan(&user.Email, fullname, &user.Nickname, about); err != nil {
			return nil, err
		}

		if about.Valid {
			user.About = about.String
		}

		if fullname.Valid {
			user.Fullname = fullname.String
		}

		users = append(users, *user)
	}
	rows.Close()

	return users, nil
}

func (userDB *UserDB) Update(user *entity.User) error {
	tx, err := userDB.db.Begin()
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

	_, err = tx.Exec("user_update", user.Email, user.Nickname, user.Fullname, user.About)
	if err != nil {
		return err
	}

	return nil
}

func (userDB *UserDB) DeleteAll() error {
	tx, err := userDB.db.Begin()
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

	_, err = userDB.db.Exec("DELETE FROM usr")
	if err != nil {
		return err
	}

	return nil
}

func (userDB *UserDB) GetStatus(s *entity.Status) error {
	tx, err := userDB.db.Begin()
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

	rows, err := tx.Query(
		"SELECT count(*) FROM forum " +
			"UNION ALL " +
			"SELECT count(*) " +
			"FROM post " +
			"UNION ALL " +
			"SELECT count(*) FROM thread " +
			"UNION ALL " +
			"SELECT count(*) FROM usr",
	)
	if err != nil {
		return err
	}

	i := 0
	for rows.Next() {
		var err error
		switch i {
		case 0:
			err = rows.Scan(&s.Forum)
		case 1:
			err = rows.Scan(&s.Post)
		case 2:
			err = rows.Scan(&s.Thread)
		case 3:
			err = rows.Scan(&s.User)
		}
		if err != nil {
			rows.Close()
			return err
		}
		i++
	}
	rows.Close()

	return nil
}

func (userDB *UserDB) Prepare() error {
	_, err := userDB.db.Prepare("user_insert",
		"insert into usr (email, fullname, nickname, about) "+
			"values ($1, $2, $3, $4) on conflict do nothing returning email",
	)
	if err != nil {
		return err
	}

	_, err = userDB.db.Prepare("user_get_by_nickname",
		"select u.email, u.fullname, u.nickname, u.about "+
			"from usr u where nickname = $1 ",
	)
	if err != nil {
		return err
	}

	_, err = userDB.db.Prepare("user_get_by_nickname_or_email",
		"select u.email, u.fullname, u.nickname, u.about "+
			"from usr u where nickname = $1 or email = $2",
	)
	if err != nil {
		return err
	}

	_, err = userDB.db.Prepare("user_update",
		"update usr set email = $1, nickname = $2, fullname = $3, about = $4 where nickname = $2 returning email",
	)
	if err != nil {
		return err
	}

	return nil
}
