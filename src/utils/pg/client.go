// Package pg provides API to interact with a Postgres server.
package pg

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/tricorder/src/utils/log"
)

// Client wraps a connection to a Postgres database.
type Client struct {
	url string

	// Use a connection pool to manage access to multiple database connections from multiple goroutines.
	pool *pgxpool.Pool
}

func NewClient(url string) *Client {
	return &Client{
		url: url,
	}
}

func (c *Client) Connect() error {
	log.Infof("Connecting to Postgres Database at %s ...", c.url)
	pool, err := pgxpool.New(context.Background(), c.url)
	if err != nil {
		return fmt.Errorf(
			"while connecting PG at %s, failed to create pgxpool, error: %v",
			c.url,
			err,
		)
	}
	c.pool = pool
	return nil
}

func buildCreateTableSQL(schema *Schema) (string, error) {
	if len(schema.Name) == 0 {
		return "", fmt.Errorf("while building SQL for creating table, table name is empty")
	}
	cols := make([]string, 0, len(schema.Columns))
	for _, col := range schema.Columns {
		colDef, err := DefineColumn(col)
		if err != nil {
			return "", fmt.Errorf(
				"while creating SQL for creating table, failed to define column, error: %v",
				err,
			)
		}
		cols = append(cols, colDef)
	}
	sql := fmt.Sprintf(
		`CREATE TABLE IF NOT EXISTS %s ( %s );`,
		schema.Name,
		strings.Join(cols, ","),
	)
	return sql, nil
}

func (client *Client) CreateTable(schema *Schema) error {
	sql, err := buildCreateTableSQL(schema)
	if err != nil {
		return fmt.Errorf(
			"while creating table '%s', failed to build SQL, error: %v",
			schema.Name,
			err,
		)
	}
	_, err = client.pool.Exec(context.Background(), sql)
	if err != nil {
		return fmt.Errorf(
			"while creating table '%s', failed to execute SQL, error: %v",
			schema.Name,
			err,
		)
	}
	return nil
}

func (client *Client) CreateHTTPRequestTable() error {
	createSQL := `CREATE TABLE IF NOT EXISTS http (
	time TIMESTAMP WITH TIME ZONE NOT NULL,
	id TEXT NOT NULL,
	method TEXT,
	proto TEXT,
	url TEXT,
	header TEXT,
	body TEXT,
	PRIMARY KEY(time, id)
	);`
	_, err := client.pool.Exec(context.Background(), createSQL)
	return err
}

func (client *Client) JSON() *Json {
	return &Json{pool: client.pool}
}

func (client *Client) Clean(table string) error {
	_, err := client.pool.Exec(context.Background(), fmt.Sprintf("TRUNCATE TABLE %s", table))
	return err
}

func (client *Client) WriteHTTPRequest(req *http.Request) error {
	sql := `
	INSERT INTO http (id, method, proto, url, header, body, time)
	VALUES ($1, $2, $3, $4, $5, $6, $7);
	`
	id := req.Header.Get("Request-Id")
	b, err := json.Marshal(req.Header)
	if err != nil {
		return err
	}
	header := string(b)

	body := []byte{}
	if req.Body != nil {
		body, err = io.ReadAll(req.Body)
		if err != nil {
			return err
		}
	}

	_, err = client.pool.Exec(
		context.Background(),
		sql,
		id,
		req.Method,
		req.Proto,
		req.URL,
		header,
		string(body),
		time.Now(),
	)
	return err
}

func (client *Client) Close() {
	client.pool.Close()
}

// Returns a string in the form of '$1, $2, ... ${count}'.
func placeHolder(count int) string {
	res := make([]string, 0, count)
	for i := 1; i <= count; i = i + 1 {
		res = append(res, "$"+strconv.Itoa(i))
	}
	return strings.Join(res, ", ")
}

func colNames(schema *Schema) string {
	colNames := make([]string, 0, len(schema.Columns))
	for _, col := range schema.Columns {
		colNames = append(colNames, col.Name)
	}
	return strings.Join(colNames, ", ")
}

// WriteRecord writes a slice of values in string format, according to the table schema.
func (client *Client) WriteRecord(record []interface{}, schema *Schema) error {
	if len(record) != len(schema.Columns) {
		return fmt.Errorf(
			"while writing record, the record's field count differs from the schema's column count, "+
				"%d vs %d",
			len(record),
			len(schema.Columns),
		)
	}
	const writeRecordSQLTmpl = `INSERT INTO %s (%s) VALUES (%s)`
	sql := fmt.Sprintf(
		writeRecordSQLTmpl,
		schema.Name,
		colNames(schema),
		placeHolder(len(schema.Columns)),
	)
	_, err := client.pool.Exec(context.Background(), sql, record...)
	return err
}

// Query returns the value of the sql query statement, or error if failed.
func (client *Client) Query(sql string) ([][]interface{}, error) {
	rows, err := client.pool.Query(context.Background(), sql)
	if err != nil {
		return nil, fmt.Errorf(
			"while querying '%s', failed to execute the statement, error: %v",
			sql,
			err,
		)
	}
	// This closes the connection objects associated with this Row.
	// This row object has a associated connection to continuously pull data from the remote database.
	// This is critical, otherwise client.pool.Close() will bock indefinitely, as it waits for all connections
	// to be closed before closing itself.
	defer rows.Close()

	res := make([][]interface{}, 0)

	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return nil, fmt.Errorf(
				"while querying '%s', failed to retrieve values from rows, error: %v",
				sql,
				err,
			)
		}
		res = append(res, values)
	}
	return res, nil
}

// Json wrapper the corresponding operations(Get|List|Upsert|Delete) to specified struct
type Json struct {
	pool *pgxpool.Pool
}

// Get returns an object with the given clause
// object MUST be a pointer
func (j *Json) Get(table, object interface{}, clause ...string) error {
	row := j.pool.QueryRow(context.Background(), fmt.Sprintf("SELECT data FROM %s %s", table, strings.Join(clause, " ")))
	if err := row.Scan(object); err != nil {
		return err
	}
	return nil
}

// List returns objects into result
// result MUST be []*T pointer, e.g. &([]*T)
func (j *Json) List(table string, result interface{}, clause ...string) error {
	sql := fmt.Sprintf("SELECT data FROM %s %s", table, strings.Join(clause, " "))
	rows, err := j.pool.Query(context.Background(), sql)
	if err != nil {
		return fmt.Errorf("while listing objects on table '%s', failed to query with sql '%s', error: %v", table, sql, err)
	}
	defer rows.Close()
	slicev := reflect.ValueOf(result).Elem()
	for rows.Next() {
		// This element is an instance pointer of T from the given 'result' parameter
		elemt := reflect.New(slicev.Type().Elem().Elem())
		if err := rows.Scan(elemt.Interface()); err != nil {
			return err
		}
		slicev = reflect.Append(slicev, elemt)
	}
	reflect.ValueOf(result).Elem().Set(slicev)
	return nil
}

// Upsert = insert if not exist, otherwise update it, idPath is optional, [->'metadata'->>'uid'] by default
func (j *Json) Upsert(table, uid string, data []byte, idPath ...string) error {
	var (
		result = ""
		ctx    = context.Background()
	)

	row := j.pool.QueryRow(ctx, fmt.Sprintf("SELECT data FROM %s WHERE %s=$1", table, pgPath(idPath)), uid)
	if err := row.Scan(&result); err == pgx.ErrNoRows {
		_, err := j.pool.Exec(ctx, fmt.Sprintf("INSERT INTO %s (data) VALUES ($1)", table), data)
		return err
	}

	_, err := j.pool.Exec(ctx, fmt.Sprintf("UPDATE %s SET data=$1 WHERE %s=$2", table, pgPath(idPath)), data, uid)
	return err
}

func (j *Json) Delete(table, uid string, idPath ...string) error {
	_, err := j.pool.Exec(context.Background(), fmt.Sprintf("DELETE FROM %s WHERE %s=$1", table, pgPath(idPath)), uid)
	return err
}

func (client *Client) CheckTableExist(tableName string) error {
	sql := fmt.Sprintf(
		`select count(*) as c from %s ;`,
		tableName,
	)
	_, err := client.pool.Exec(context.Background(), sql)
	if err != nil {
		return fmt.Errorf(
			"while check table '%s' exist, failed to execute SQL, error: %v",
			tableName,
			err,
		)
	}
	return nil
}
