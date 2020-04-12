package crawler

import (
	"fmt"
	"strconv"

	"database/sql"

	"flash/internal/app"

	_ "github.com/go-sql-driver/mysql"
)

// Config is exported.
type Config struct {
	Username  string `env:"DBC_USERNAME" seed:"profx"`
	Password  string `env:"DBC_PASSWORD" seed:"profx"`
	ReadHost  string `env:"DBC_READ_HOST" seed:"localhost"`
	ReadPort  string `env:"DBC_READ_PORT" seed:"3306"`
	WriteHost string `env:"DBC_WRITE_HOST" seed:"localhost"`
	WritePort string `env:"DBC_WRITE_PORT" seed:"3306"`
	Schema    string `env:"DBC_SCHEMA" seed:"profx"`

	ReadMaxConn   string `env:"DBC_READ_MAX_CONN" seed:"50"`
	ReadIdleConn  string `env:"DBC_READ_IDLE_CONN" seed:"10"`
	WriteMaxConn  string `env:"DBC_WRITE_MAX_CONN" seed:"50"`
	WriteIdleConn string `env:"DBC_WRITE_IDLE_CONN" seed:"10"`
}

const connectionString = "%v:%v@tcp(%v:%v)/%v?multiStatements=true&parseTime=true"

func CreateConnectionString(username string, password string, host string, port string, schema string) string {
	return fmt.Sprintf(connectionString, username, password, host, port, schema)
}

// DB is exported.
type DB struct {
	Read  *sql.DB
	Write *sql.DB
}

// New is exported.
func (c *Config) New() (*DB, error) {
	cs := CreateConnectionString(c.Username, c.Password, c.ReadHost, c.ReadPort, c.Schema)
	read, err := sql.Open(
		"mysql",
		cs,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to open connection to Write database with config = '%v' : %w", c, err)
	}

	err = read.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping Write database with config = '%v' : %w", c, err)
	}

	var maxConn, idleConn int
	maxConn, _ = strconv.Atoi(c.ReadMaxConn)
	idleConn, _ = strconv.Atoi(c.ReadIdleConn)

	read.SetMaxOpenConns(maxConn)
	read.SetMaxIdleConns(idleConn)

	write, err := sql.Open(
		"mysql",
		CreateConnectionString(c.Username, c.Password, c.WriteHost, c.WritePort, c.Schema),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to open connection to Read database with config = '%v' : %w", c, err)
	}

	err = write.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping master Read database with config = '%v' : %v", c, err)
	}

	maxConn, _ = strconv.Atoi(c.WriteMaxConn)
	idleConn, _ = strconv.Atoi(c.ReadIdleConn)

	write.SetMaxOpenConns(maxConn)
	write.SetMaxIdleConns(idleConn)

	db := DB{
		Read:  read,
		Write: write,
	}

	return &db, nil
}

func (db *DB) Ping() (err error) {
	err = db.Read.Ping()
	if err == nil {
		err = db.Write.Ping()
	}
	return
}

// HasLink is exported.
func (db *DB) HasLink(hash string) (*bool, error) {
	var res bool
	err := db.Read.QueryRow("SELECT EXISTS (SELECT * FROM links WHERE hash = ? LIMIT 1)", hash).Scan(&res)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &res, nil
}

// IsArticle is exported.
func (db *DB) IsArticle(hash string) (*bool, error) {
	var res bool

	var query = `
		SELECT EXISTS (
			SELECT * FROM resources WHERE hash = ? AND kind = 'article'
		)
	`
	err := db.Read.QueryRow(query, hash).Scan(&res)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &res, nil
}

// SaveLink is exported.
func (db *DB) SaveLink(link *app.Link) error {

	query := "INSERT INTO links (hash, url, source, from_url) VALUES (?, ?, ?, ?)"

	res, err := db.Write.Exec(query, link.Hash, link.URL, link.Source, link.FromURL)
	if err != nil {
		return fmt.Errorf("can't save link %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil || rows == 0 {
		return fmt.Errorf("unsuccessful link insert %w", err)
	}

	return nil
}

// GetLink is exported.
func (db *DB) GetLink(hash string) (*app.Link, error) {
	query := `SELECT * FROM links WHERE hash = ?`

	row, err := db.Read.Query(query, hash)
	if err != nil {
		return nil, fmt.Errorf("failed to get link by hash %s: %w", hash, err)
	}
	defer row.Close()

	link := app.Link{}
	row.Next()
	err = row.Scan(&link.Hash, &link.URL, &link.Source, &link.FromURL, &link.Created)
	if err != nil {
		return nil, fmt.Errorf("failed to scan link row: %w", err)
	}

	return &link, nil
}

// UpdateLink is exported.
func (db *DB) UpdateLink(link *app.Link) error {
	query := `
		UPDATE links
		SET source = ?, from_url = ?
		WHERE hash = ?
	`

	res, err := db.Write.Exec(query, link.Source, link.FromURL, link.Hash)
	if err != nil {
		return fmt.Errorf("failed update link: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("expected to have more than 0 affected rows")
	}

	return nil
}

// SaveResource is exported.
func (db *DB) SaveResource(link *app.Resource) error {
	query := "INSERT INTO resources (hash, url, kind, source, from_url, title, description) VALUES (?, ?, ?, ?, ?, ?, ?)"

	res, err := db.Write.Exec(query, link.Hash, link.URL, link.Kind, link.Source, link.FromURL, link.Title, link.Description)
	if err != nil {
		return fmt.Errorf("can't save resource %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil || rows == 0 {
		return fmt.Errorf("unsuccessful resource insert %w", err)
	}

	return nil
}

// SaveLog is exported.
func (db *DB) SaveLog(log *app.SysLog) error {
	query := "INSERT INTO sys_logs (message) VALUES (?)"

	res, err := db.Write.Exec(query, log.Message)
	if err != nil {
		return fmt.Errorf("can't save log %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil || rows == 0 {
		return fmt.Errorf("unsuccessful log insert %w", err)
	}

	return nil
}

func (db *DB) LoadSources() ([]app.Source, error) {
	query := `SELECT * FROM sources`

	rows, err := db.Read.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query sources: %w", err)
	}
	defer rows.Close()

	sources := []app.Source{}
	for rows.Next() {
		source := app.Source{}

		err := rows.Scan(&source.Name, &source.URL, &source.Kind, &source.Created)

		if err != nil {
			return nil, fmt.Errorf("faild to scan source row: %w", err)
		}
		sources = append(sources, source)
	}

	return sources, nil
}

func (db *DB) LoadRules() ([]app.Rule, error) {
	query := `SELECT * FROM rules`

	rows, err := db.Read.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query source rules: %w", err)
	}
	defer rows.Close()

	sources := []app.Rule{}
	for rows.Next() {
		source := app.Rule{}

		err := rows.Scan(&source.Source, &source.Type, &source.Rule, &source.Created)

		if err != nil {
			return nil, fmt.Errorf("faild to scan source rule row: %w", err)
		}
		sources = append(sources, source)
	}

	return sources, nil
}

func (db *DB) LoadResources() ([]app.Resource, error) {
	query := `SELECT source, url, created FROM resources ORDER BY created DESC`

	rows, err := db.Read.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query resources: %w", err)
	}
	defer rows.Close()

	resources := []app.Resource{}
	for rows.Next() {
		resource := app.Resource{}

		err := rows.Scan(&resource.Source, &resource.URL, &resource.Created)

		if err != nil {
			return nil, fmt.Errorf("faild to scan resource row: %w", err)
		}
		resources = append(resources, resource)
	}

	return resources, nil
}
