package connector

import (
	"database/sql"
	"fmt"
	"time"
)

func Query(query string, db *sql.DB) *QueryResponse {

	response := new(QueryResponse)

	if db != nil {
		rows, err := db.Query(query)
		if err != nil {
			response.Error = true
			response.ErrorMessage = err.Error()
			return response
		}
		defer rows.Close()

		columnNames, err := rows.Columns()
		if err != nil {
			response.Error = true
			response.ErrorMessage = err.Error()
			return response
		}

		rc := NewMapStringScan(columnNames)
		for rows.Next() {
			err := rc.Update(rows)
			if err != nil {
				response.Error = true
				response.ErrorMessage = err.Error()
				return response
			}
			cv := rc.Get()
			row := new(Row)

			for k, v := range cv {
				col := new(Column)
				col.Index = k

				for key, value := range v {
					col.Name = key
					col.Value = value
				}

				row.Columns = append(row.Columns, *col)
			}
			response.Rows = append(response.Rows, *row)

		}

	} else {
		response.Error = true
		response.ErrorMessage = "Received a nil object as a DB connection "
		return response
	}
	response.Timestamp = time.Now().UTC().Unix()
	return response
}

type mapStringScan struct {
	cp []interface{}

	row      map[int]map[string]string
	colCount int
	colNames []string
}

func NewMapStringScan(columnNames []string) *mapStringScan {
	lenCN := len(columnNames)
	s := &mapStringScan{
		cp:       make([]interface{}, lenCN),
		row:      make(map[int]map[string]string, lenCN),
		colCount: lenCN,
		colNames: columnNames,
	}
	for i := 0; i < lenCN; i++ {
		s.cp[i] = new(sql.RawBytes)
	}
	return s
}

func (s *mapStringScan) Update(rows *sql.Rows) error {
	if err := rows.Scan(s.cp...); err != nil {
		return err
	}

	for i := 0; i < s.colCount; i++ {
		if rb, ok := s.cp[i].(*sql.RawBytes); ok {

			s.row[i] = map[string]string{s.colNames[i]: string(*rb)}
			*rb = nil // reset pointer to discard current value to avoid a bug
		} else {
			return fmt.Errorf("Cannot convert index %d column %s to type *sql.RawBytes", i, s.colNames[i])
		}
	}
	return nil
}

func (s *mapStringScan) Get() map[int]map[string]string {
	return s.row
}
