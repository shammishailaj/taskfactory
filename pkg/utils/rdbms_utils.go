package utils

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type RDBMS struct {
	Conn    *sql.DB
	ConnErr error
	Conf    *RDBMSConfig
}

func (d *RDBMS) Connect() error {
	d.Conf = new(RDBMSConfig)
	d.Conf.DriverName = os.Getenv("DB_DRIVER")
	d.Conf.HostName = os.Getenv("DB_HOST")
	d.Conf.Port, _ = strconv.ParseInt(os.Getenv("DB_PORT"), 10, 64)
	d.Conf.UserName = os.Getenv("DB_USER")
	d.Conf.Password = os.Getenv("DB_PASS")
	d.Conf.Database = os.Getenv("DB_NAME")
	d.Conf.MaxIdleConns, _ = strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONN"))
	d.Conf.MaxOpenConns, _ = strconv.Atoi(os.Getenv("DB_MAX_OPEN_CONN"))
	d.Conf.MaxLifeTimeConn, _ = strconv.Atoi(os.Getenv("DB_MAX_LIFE_TIME_CONN"))
	d.Conn, d.ConnErr = d.Conf.Connect()
	return d.ConnErr
}

func (d *RDBMS) ConnectM2() error {
	d.Conf = new(RDBMSConfig)
	d.Conf.DriverName = os.Getenv("M2_DB_DRIVER")
	d.Conf.HostName = os.Getenv("M2_DB_HOST")
	d.Conf.Port, _ = strconv.ParseInt(os.Getenv("M2_DB_PORT"), 10, 64)
	d.Conf.UserName = os.Getenv("M2_DB_USER")
	d.Conf.Password = os.Getenv("M2_DB_PASS")
	d.Conf.Database = os.Getenv("M2_DB_NAME")
	d.Conf.MaxIdleConns, _ = strconv.Atoi(os.Getenv("M2_DB_MAX_IDLE_CONN"))
	d.Conf.MaxOpenConns, _ = strconv.Atoi(os.Getenv("M2_DB_MAX_OPEN_CONN"))
	d.Conf.MaxLifeTimeConn, _ = strconv.Atoi(os.Getenv("M2_DB_MAX_LIFE_TIME_CONN"))
	d.Conn, d.ConnErr = d.Conf.Connect()
	return d.ConnErr
}

func (d *RDBMS) Close() error {
	return d.Conn.Close()
}

/**
 * @param type string $dbObj
 * @param type string $tablename
 * @param type string $columns
 * @param type string $condition
 * @param type integer $limit
 * @param type integer $offset
 * @param type string $order_by
 * @return array
 */
func (d *RDBMS) SelectRows(tableName, columns, condition, limit, offset, orderBy string) *sql.Rows {
	sql := "SELECT " + columns + " FROM `" + tableName + "`"
	if len(condition) > 0 {
		sql += " WHERE " + condition
	}
	if len(orderBy) > 0 {
		sql += " ORDER BY " + orderBy
	}
	if len(offset) > 0 && len(limit) > 0 {
		sql += " LIMIT " + offset + ", " + limit
	} else if len(limit) > 0 {
		sql += " LIMIT 0," + limit
	}

	log.Printf("SQLEXEC:: %s", sql)
	rows, err := d.Conn.Query(sql)

	if err != nil {
		log.Printf("ERROR executing SQL:\n%s \nDetails: %s", sql, err.Error())
		return nil
	}

	return rows
}

/**
 * @param type string $dbObj
 * @param type string $tablename
 * @param type string insertDataArray
 * @return type int
 */
func (d *RDBMS) InsertRow(table string, insertDataArray map[string]string) int {

	sql := "Insert into " + table

	value := []interface{}{}

	v := reflect.ValueOf(insertDataArray)
	typeOfAlertTable := v.Type()

	if v.Len() > 0 {
		sql += " SET "
		for i := 0; i < v.NumField(); i++ {
			//fmt.Printf("Field: %s\tValue: %v\n", typeOfS.Field(i).Name, v.Field(i).Interface())
			columnName := typeOfAlertTable.Field(i).Name
			columnValue := v.Field(i).Interface()
			sql += columnName + "=?,"
			value = append(value, columnValue)
		}
		sql = strings.TrimRight(sql, ",")
	}

	log.Printf("Insert Query: %#v", sql)

	//if len(insertDataArray)>0 {
	//    sql+=" SET "
	//    for key,val := range insertDataArray{
	//        sql += key+"=?,"
	//        value = append(value, val)
	//    }
	//    sql=strings.TrimRight(sql,",")
	//}

	stmt, err := d.Conn.Prepare(sql)
	if err != nil {
		log.Println("Cannot prepare DB statement", err)
	}

	res, err := stmt.Exec(value...)
	if err != nil {
		log.Println("Cannot run insert statement", err)
	}

	id, lastInsertIdErr := res.LastInsertId()
	if lastInsertIdErr != nil {
		log.Println(" lastInsertIdErr = %#v\n", lastInsertIdErr)
	}
	return int(id)
}

/**
 * @param type string $dbObj
 * @param type string $tablename
 * @param type string updateDataArray
 * @param type string optionDataArray
 * @return type rows
 */
func (d *RDBMS) UpdateRows(table string, updateDataArray map[string]string, optionDataArray map[string]string) *sql.Rows {
	sql := "UPDATE `" + table + "` SET "

	for key, val := range updateDataArray {
		sql += "`" + key + "` = '" + val + "',"
	}
	sql = strings.TrimRight(sql, ",")

	if len(optionDataArray) > 0 {
		sql += " WHERE "
		for k, v := range optionDataArray {
			sql += "`" + k + "` = '" + v + "' AND "
		}
	}

	sql = strings.TrimRight(sql, "AND ")

	log.Printf("SQLEXEC:: %s", sql)
	rows, err := d.Conn.Query(sql)

	if err != nil {
		log.Printf("ERROR updating row via SQL:\n%s\nDetails: %s\n", sql, err.Error())
	}
	return rows
}

func (d *RDBMS) UpdateTable(table string, updateDataArray map[string]string, optionDataArray map[string]string) (sql.Result, error) {
	sql := "UPDATE `" + table + "` SET "

	for key, val := range updateDataArray {
		sql += "`" + key + "` = '" + val + "',"
	}
	sql = strings.TrimRight(sql, ",")

	if len(optionDataArray) > 0 {
		sql += " WHERE "
		for k, v := range optionDataArray {
			sql += "`" + k + "` = '" + v + "' AND "
		}
	}

	sql = strings.TrimRight(sql, "AND ")

	log.Printf("SQLEXEC:: %s", sql)
	stmt, stmtErr := d.Conn.Prepare(sql)
	if stmtErr != nil {
		log.Infof("Cannot prepare DB statement. %s", stmtErr.Error())
		return nil, stmtErr
	}

	return stmt.Exec()
}

func (d *RDBMS) DropTable(tableName string) error {
	var dropTableErr error
	qry := "DROP TABLE `" + tableName + "`"

	stmt, stmtErr := d.Conn.Prepare(qry)
	dropTableErr = stmtErr
	if stmtErr != nil {
		log.Printf("Cannot prepare DB statement. %s", stmtErr.Error())
	} else {
		log.Printf("Successfully prepared DB statement")
		res, err := stmt.Exec()
		dropTableErr = err
		if err != nil {
			log.Printf("Cannot execute DROP query. %s. %s", err.Error(), qry)
		} else {
			log.Printf("Successfully dropped table: %s", tableName)
			aRows, aRowsErr := res.RowsAffected()
			dropTableErr = aRowsErr
			if aRowsErr != nil {
				log.Printf("Error getting affected rows. %s", aRowsErr.Error())
			} else {
				log.Printf("Successfully dropped. Affected Rows = %d", aRows)
			}
		}
	}
	return dropTableErr
}

