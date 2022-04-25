package tms

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

//compile time interface implementation assertion
var _ Storer = (*PostgresStore)(nil)

//PostgresStore implements Store backend by Postgresql
type PostgresStore struct {
	db *sqlx.DB
}

//NewPostgresStore creates a new postgres store
func NewPosgresStore(dsn string) (*PostgresStore, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to connect: %w", err)
	}
	return &PostgresStore{
		db: db,
	}, nil
}

//Close closes the underlying db store
func (p *PostgresStore) Close() error {
	return p.db.Close()
}

//Create user creates a user
func (p *PostgresStore) CreateUser(ctx context.Context, firstName, lastName, email string) error {
	_, err := p.db.ExecContext(ctx, insertUserSQL, firstName, lastName, email)
	//TODO check for email uniqueness violation
	return err
}

//CreateTask creates a task
func (p *PostgresStore) CreateTask(ctx context.Context, description string, userID uint, startTime, endTime, reminderPeriod time.Time) error {
	_, err := p.db.ExecContext(ctx, insertTaskSQL, description, startTime, endTime, reminderPeriod, userID)
	return err
}

//DeleteUser deletes a user
func (p *PostgresStore) DeleteUser(ctx context.Context, id uint) error {
	_, err := p.db.ExecContext(ctx, deleteUserSQL, id)
	return err
}

//DeleteTask deletes a task from store
func (p *PostgresStore) DeleteTask(ctx context.Context, id uint) error {
	_, err := p.db.ExecContext(ctx, deleteTaskSQL, id)
	return err
}



//GetTasks gets tasks from store
func (p *PostgresStore) GetTasks(ctx context.Context, limit, offset uint) ([]Task, error) {
	var sqlTasks []sqlTask
	if err := p.db.SelectContext(ctx, &sqlTasks, selectTasksSQL, limit, offset); err != nil {
		return nil, err
	}
	tasks := make([]Task, 0, len(sqlTasks))
	for _, task := range sqlTasks {
		tasks = append(tasks, *task.toTask())
	}
	return tasks, nil
}
//GetTask gets task from store
func (p *PostgresStore) GetTask(ctx context.Context, id uint) (*Task, error) {
	var t sqlTask
	if err := p.db.GetContext(ctx, &t, selectTaskSQL, id); err != nil {
		return nil, err
	}
	return t.toTask(), nil
}
//GetUser gets a user from storage
func (p *PostgresStore) GetUser(ctx context.Context, id uint) (*User, error) {
	var user sqlUser
	if err := p.db.GetContext(ctx, &user, selectUserSQL, id); err != nil {
		return nil, err
	}
	return user.toUser(), nil
}

//GetUsers gets users from store
func (p *PostgresStore) GetUsers(ctx context.Context, limit, offset uint) ([]User, error) {
	var u []sqlUser
	if err := p.db.SelectContext(ctx, &u, selectUsersSQL, limit, offset); err != nil {
		return nil, err
	}
	users := make([]User, 0, len(u))
	for _, user := range u {
		users = append(users, *user.toUser())
	}
	return users, nil
}

//GetUserTasks gets tasks from store
func (p *PostgresStore) GetUserTasks(ctx context.Context, userID, limit, offset uint) ([]Task, error) {
	var sqlTasks []sqlTask
	if err := p.db.SelectContext(ctx, &sqlTasks, selectUserTasksSQL, userID, limit, offset); err != nil {
		return nil, err
	}
	tasks := make([]Task, 0, len(sqlTasks))
	for _, task := range sqlTasks {
		tasks = append(tasks, *task.toTask())
	}
	return tasks, nil
}

var (
	deleteTaskSQL = `
		delete from tasks where id = $1
	`
	deleteUserSQL = `
		delete from users where id = $1
	`
	insertTaskSQL = `
		insert into tasks (description, start_time, end_time, reminder_period, user_id)
		values ($1, $2::timestamp, $3::timestamp, $4::timestamp, $5)
	`
	insertUserSQL = `
		insert into users (firstname, lastname, email)
		values ($1, $2, $3)
	`
	selectTaskSQL = `
		select id, description, start_time, end_time, reminder_period, created_at from tasks where id = $1
	`
	selectTasksSQL = `
		select id, description, start_time, end_time, reminder_period, created_at from tasks limit $1 offset $2
	`
	selectUserSQL = `
		select id, firstname, lastname, email, created_at from users where id = $1
	`
	selectUsersSQL = `
		select id, firstname, lastname, email, created_at from users order by id asc limit $1 offset $2
	`
	selectUserTasksSQL = `
		select id, description, start_time, end_time, reminder_period, created_at from tasks where user_id = $1 limit $2 offset $3
	`
)

type sqlUser struct {
	ID        uint      `db:"id"`
	FirstName string    `db:"firstname"`
	LastName  string    `db:"lastname"`
	Email     string    `db:"email"`
	CreatedAt time.Time `db:"created_at"`
}

func (s sqlUser) toUser() *User {
	return &User{
		ID:        s.ID,
		FirstName: s.FirstName,
		LastName:  s.LastName,
		Email:     s.Email,
		CreatedAt: s.CreatedAt,
	}
}

type sqlTask struct {
	ID             uint      `db:"id"`
	Description    string    `db:"description"`
	StartTime      time.Time `db:"start_time"`
	EndTime        time.Time `db:"end_time"`
	ReminderPeriod time.Time `db:"reminder_period"`
	CreatedAt      time.Time `db:"created_at"`
}

func (s sqlTask) toTask() *Task {
	return &Task{
		ID:             s.ID,
		CreatedAt:      s.CreatedAt,
		StartTime:      s.StartTime,
		Duration:       s.EndTime.Sub(s.StartTime),
		ReminderPeriod: s.ReminderPeriod,
		Description:    s.Description,
	}
}
