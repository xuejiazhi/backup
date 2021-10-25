package core

import (
	"bibt-SpeedSkat/backup/utils"
	"errors"
	"fmt"
	"github.com/spf13/cast"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io"
	"os"
	"strings"
	"time"
)

type PgSql struct {
	Schema     string
	TableName  string
	FieldList  []map[string]interface{}
	ResultList []map[string]interface{}
}

func PgSqldump(dsn, schema string) {
	//set default value
	if err := SetPgConn(dsn); err != nil {
		os.Exit(0)
	}

	schemaList := strings.Split(schema, "|")
	if len(schemaList) == 0 {
		return
	}

	//get tableList
	for _, sonSchema := range schemaList {
		tableNames := GetPgTableList(sonSchema)
		if len(tableNames) == 0 {
			continue
		}
		for _, t := range tableNames {
			//get field list
			fieldList := GetPgFieldList(t)
			if len(fieldList) == 0 {
				continue
			}

			if ok, _ := PathExists(BackTmpDir + "/" + sonSchema); !ok {
				Mkdir(BackTmpDir + "/" + sonSchema)
			}

			GetSqlBulk(PgSql{
				Schema:    sonSchema,
				TableName: fmt.Sprintf("%s.%s", sonSchema, t),
				FieldList: fieldList,
			})
		}
	}

}

func SetPgConn(dsn string) (err error) {
	if PgConn, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	}); err != nil {
		fmt.Println("PostgreSQL  Connect Is Faled! ")
		err = errors.New("PostgreSQL  Connect Is Faled! ")
	}

	return
}

func QueryTableColumns(queryStr string) (results []map[string]interface{}, err error) {
	err = PgConn.Raw(queryStr).Take(&results).Error
	return
}

func QueryTableNums(tableName string) (count int) {
	sqlStr := fmt.Sprintf("select count(*) from %s", tableName)
	PgConn.Raw(sqlStr).First(&count)
	return
}

func TypeTransForm(typename, fieldname string) (columnStr string) {
	switch strings.ToLower(typename) {
	case "timestamptz":
		columnStr = fmt.Sprintf("to_char(%s,'YYYY-MM-DD hh24:mi:ss') AS %s", fieldname, fieldname)
	case "timestamp":
		columnStr = fmt.Sprintf("to_char(%s,'YYYY-MM-DD hh24:mi:ss') AS %s", fieldname, fieldname)
	default:
		columnStr = fieldname
	}
	return
}

func GetSqlBulk(result PgSql) {
	//get count
	var err error
	count := QueryTableNums(result.TableName)
	pgCount := count/PgLimit + 1
	pgBar := utils.LocalBar{
		BarCount:    pgCount,
		Start:       0,
		Notice:      "Generate PgSQL " + result.TableName + ":",
		Graph:       "#",
		NoticeColor: 2,
		SleepTime:   1,
	}
	pgBar.GenBar()
	for i := 0; i < pgCount; i++ {
		//get result
		for true {
			if result.ResultList, err = QueryTableColumns(GetQuerySql(result.TableName, result.FieldList, i)); err == nil {
				break
			}
			fmt.Println("Get Result Is Error,", err.Error())
			time.Sleep(2 * time.Second)
		}

		sqlBulk := GetBulkSql(result)
		fileName := fmt.Sprintf("%s/%s/%s.sql",
			BackTmpDir,
			result.Schema,
			utils.If(i > 0, result.TableName+"_"+cast.ToString(i), result.TableName),
		)

		//create Zip file
		var f *os.File
		if ok, _ := PathExists(fileName); ok { //如果文件存在
			f, _ = os.OpenFile(fileName, os.O_APPEND, 0666) //打开文件
		} else {
			f, _ = os.Create(fileName) //创建文件
		}
		io.WriteString(f, sqlBulk) //写入文件(字符串)
		f.Close()
		//print bar
		pgBar.PrintBar()
	}
	pgBar.EndBar()
}

func GetPgTableList(schema string) []string {
	var tbNames []string
	sqlStr := fmt.Sprintf("select tablename from pg_tables where schemaname = '%s'", schema)
	PgConn.Raw(sqlStr).Take(&tbNames)
	return tbNames
}

func GetPgFieldList(tableName string) (results []map[string]interface{}) {
	sqlStr := fmt.Sprintf(`SELECT 
		a.attname AS field,
		t.typname as typename
		FROM pg_class c,
		pg_attribute a
	LEFT OUTER JOIN pg_description b ON a.attrelid=b.objoid AND a.attnum = b.objsubid,
		pg_type t
	WHERE c.relname = '%s'
	and a.attnum > 0
	and a.attrelid = c.oid
	and a.atttypid = t.oid
	ORDER BY a.attnum;`, tableName)
	PgConn.Raw(sqlStr).Take(&results)
	return
}

func GetQuerySql(tableName string, fields []map[string]interface{}, pageSize int) string {
	fieldList := []string{}
	for _, v := range fields {
		fieldList = append(fieldList, TypeTransForm(cast.ToString(v["typename"]), cast.ToString(v["field"])))
	}
	fieldStr := strings.Join(fieldList, ",")

	sqlStr := fmt.Sprintf("SELECT %s FROM %s  OFFSET %d LIMIT %d", fieldStr, tableName, pageSize*PgLimit, PgLimit)
	return sqlStr
}

func GetBulkSql(result PgSql) (sqlStr string) {
	if len(result.ResultList) == 0 {
		return
	}
	for _, v := range result.ResultList {
		Sql := GetPreInsertSql(result.TableName, result.FieldList, v)
		sqlStr += Sql + ";\n"
	}
	return
}

func GetPreInsertSql(tableName string, fieldList []map[string]interface{}, val map[string]interface{}) string {
	//set list
	colNameList := []string{}
	valueList := []string{}
	for _, v := range fieldList {
		columnName := cast.ToString(v["field"])
		colNameList = append(colNameList, columnName)
		if _, ok := val[columnName]; ok {
			valueList = append(valueList, fmt.Sprintf("'%s'", cast.ToString(val[columnName])))
		} else {
			return ""
		}
	}

	sqlStr := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		tableName,
		strings.Join(colNameList, ","),
		strings.Join(valueList, ","))

	return sqlStr
}
