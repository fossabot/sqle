package sqle

import (
	"database/sql"
	"fmt"
)

func (s Sqle) Select(query string, args []interface{}, dest []interface{}) error {
	_, err := s.SelectExists(query, args, dest)
	return err
}

func (s Sqle) SelectExists(query string, args []interface{}, dest []interface{}) (exists bool, err error) {
	tx, err := s.db.Begin()
	if err != nil {
		return false, err
	}

	exists, err = s.SelectExistsTx(tx, query, args, dest)
	if err != nil {
		rlbErr := tx.Rollback()
		if rlbErr == nil {
			return false, err
		}
		return false, fmt.Errorf("sqle.SelectExists: multiple errors occured: %q followed by tx.Rollback error: %q", err.Error(), rlbErr.Error())
	}

	err = tx.Commit()
	if err != nil {
		return false, err
	}

	return
}

func (s Sqle) SelectExistsTx(tx *sql.Tx, query string, args []interface{}, dest []interface{}) (exists bool, err error) {
	if len(dest) == 0 {
		return false, fmt.Errorf("sqle.SelectExistsTx: no dest to scan to")
	}

	stmt, err := tx.Prepare(query)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	if len(args) > 0 {
		err = stmt.QueryRow(args...).Scan(dest...)
	} else if len(args) == 0 {
		err = stmt.QueryRow().Scan(dest...)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (s Sqle) SelectRange(query string, args []interface{}, dest []interface{}, handleRow func()) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	err = s.SelectRangeTx(tx, query, args, dest, handleRow)
	if err != nil {
		rlbErr := tx.Rollback()
		if rlbErr == nil {
			return err
		}

		return fmt.Errorf("sqle.SelectRange: multiple errors occured: %q followed by tx.Rollback error: %q", err.Error(), rlbErr.Error())
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (s Sqle) SelectRangeTx(tx *sql.Tx, query string, args []interface{}, dest []interface{}, handleRow func()) error {
	if len(dest) == 0 {
		return fmt.Errorf("sqle.SelectRangeTx: no dest to scan to")
	}

	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	rows, err := stmt.Query(args...)
	if err != nil {
		return err
	}

	for rows.Next() {
		err = rows.Scan(dest...)
		if err != nil {
			return err
		}
		handleRow()
	}
	if rows.Err() != nil {
		return err
	}

	return nil
}
