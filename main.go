package main

import (
	"database/sql"
	"embed"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/samber/lo"
)

//go:embed frontend/dist/*
var staticFiles embed.FS

//go:embed frontend/dist/index.html
var indexHTML embed.FS

func returnIndexHTML(c *gin.Context) {
	file, err := staticFiles.Open("frontend/dist/index.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "index.html not found")
		return
	}
	defer file.Close()
	c.DataFromReader(http.StatusOK, -1, "text/html", file, nil)
}

type TableColumn struct {
	Name          string      `json:"name"`
	Type          string      `json:"type"`
	NotNull       int         `json:"not_null"`
	PrimaryKey    int         `json:"primary_key"`
	ForeignKey    int         `json:"foreign_key"`
	AutoIncrement int         `json:"auto_increment"`
	DefaultValue  interface{} `json:"default_value"`
	Unique        int         `json:"unique"`
	References    string      `json:"references"`
	// Cid       int         `json:"cid"`
	// Name      string      `json:"name"`
	// Type      string      `json:"type"`
	// NotNull   int         `json:"not_null"`
	// DfltValue interface{} `json:"value"`
	// Pk        int         `json:"pk"`
}
type TableWithColumns struct {
	Name    string        `json:"name"`
	Columns []TableColumn `json:"columns"`
	// RAW     string        `json:"raw"`
}
type TableWithColumnsAndRAW struct {
	Name    string        `json:"name"`
	Columns []TableColumn `json:"columns"`
	RAW     string        `json:"raw"`
}
type TableWithColumnsAndASpecialColumn struct {
	Name     string                               `json:"name"`
	Columns  []TableColumn                        `json:"columns"`
	Column   TableColumn                          `json:"special_column"`
	RAW      string                               `json:"raw"`
	Children []*TableWithColumnsAndASpecialColumn `json:"children"`
}
type ForeignKey struct {
	ReferenceKey string
	TableName    string
	PrimaryKey   string
}

func getTableForeignKeys(db *sql.DB, tableName string) (map[string]ForeignKey, error) {
	query := fmt.Sprintf("PRAGMA foreign_key_list(%s);", tableName)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	keys := make(map[string]ForeignKey)
	for rows.Next() {
		var foreign ForeignKey
		var a1 string
		var a2 string
		var a3 string
		var a4 string
		var a5 string
		if err := rows.Scan(&a1, &a2, &foreign.TableName, &foreign.ReferenceKey, &foreign.PrimaryKey, &a3, &a4, &a5); err != nil {
			return nil, err
		}
		keys[foreign.ReferenceKey] = foreign
	}
	return keys, nil
}
func getTableColumns(db *sql.DB, tableName string) ([]TableColumn, error) {
	foreignKeys, err := getTableForeignKeys(db, tableName)
	if err != nil {
		return nil, err
	}
	query := fmt.Sprintf("PRAGMA table_info(%s);", tableName)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []TableColumn
	for rows.Next() {
		var col TableColumn
		var a1 string
		if err := rows.Scan(&a1, &col.Name, &col.Type, &col.NotNull, &col.DefaultValue, &col.PrimaryKey); err != nil {
			return nil, err
		}
		col.Type = strings.ToLower(col.Type)
		if col.Type == "real" {
			col.Type = "integer"
		}
		foreign, has := foreignKeys[col.Name]
		if has {
			col.References = foreign.TableName
		}
		columns = append(columns, col)
	}
	return columns, nil
}
func buildSQLColumn(column TableColumn, tableName string) (string, []string) {
	stmt := fmt.Sprintf("`%v` %v", column.Name, strings.ToUpper(column.Type))
	var constraint []string
	if column.PrimaryKey == 1 {
		stmt += " PRIMARY KEY"
	}
	if column.AutoIncrement == 1 {
		stmt += " AUTO INCREMENT"
	}
	if column.NotNull == 1 {
		stmt += " NOT NULL"
	}
	if column.DefaultValue != nil {
		var v string
		if column.Type == "text" {
			v = fmt.Sprintf("'%v'", column.DefaultValue)
		}
		if column.Type == "integer" {
			v = fmt.Sprintf("%v", column.DefaultValue)
		}
		if column.Type == "datetime" {
			if v2, ok := column.DefaultValue.(time.Time); ok {
				formattedTime := v2.Format("2006-01-02 15:04:05")
				v = fmt.Sprintf("'%v'", formattedTime)
			}
		}
		stmt += fmt.Sprintf(" DEFAULT %v", v)
	}
	if column.Unique == 1 {
		constraint = append(constraint, fmt.Sprintf("CONSTRAINT `uni_%v_%v` UNIQUE (`%v`)", tableName, column.Name, column.Name))
	}
	// fmt.Println(column.Name)
	// fmt.Println("check has foreign", column.References)
	if column.References != "" {
		// todo 这里 id 是 reference 表的主键，这里固定写死了 id
		constraint = append(constraint, fmt.Sprintf("CONSTRAINT `fk_%v_%v` FOREIGN KEY (`%v`) REFERENCES `%v`(`id`)", tableName, column.References, column.Name, column.References))
	}
	return stmt, constraint
}
func buildSQLCreateTableWithColumns(name string, columns []TableColumn) string {
	var columnStmts []string
	var constraintsParts []string
	for _, col := range columns {
		columnStmt, constraints := buildSQLColumn(col, name)
		columnStmts = append(columnStmts, columnStmt)
		constraintsParts = append(constraintsParts, constraints...)
	}
	fmt.Println(columnStmts, constraintsParts)
	stmts := joinStrWithSEQ(append(columnStmts, constraintsParts...), ", ")
	stmt := fmt.Sprintf("CREATE TABLE %v (%v);", name, stmts)
	return stmt
}

// 获取关联表
func getRelatedTables(db *sql.DB, tableName string) ([]string, error) {
	rows, err := db.Query(fmt.Sprintf("PRAGMA foreign_key_list(%s);", tableName))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	relatedTables := []string{}
	for rows.Next() {
		var id, seq int
		var table, from, to, onUpdate, onDelete, match string
		if err := rows.Scan(&id, &seq, &table, &from, &to, &onUpdate, &onDelete, &match); err != nil {
			return nil, err
		}
		relatedTables = append(relatedTables, table)
	}
	return relatedTables, nil
}

func generateSQLOfTable(db *sql.DB, id int, relatedTables []*TableWithColumnsAndASpecialColumn, processedTables []string) string {
	// relatedTables := findRelatedTables(tables, tableName)
	sql := ""
	for _, table := range relatedTables {
		existing := lo.Contains(processedTables, table.Name)
		if existing {
			continue
		}
		find := fmt.Sprintf(`SELECT * FROM %v WHERE %v = ?`, table.Name, table.Column.Name)
		rows, err := db.Query(find, id)
		if err != nil {
			fmt.Println("There is error", err.Error())
			return sql
		}
		var unique_ids []int
		defer rows.Close()
		for rows.Next() {
			values := make([]interface{}, len(table.Columns))
			dest := make([]interface{}, len(table.Columns))
			for i := range dest {
				dest[i] = &values[i]
			}
			if err := rows.Scan(dest...); err != nil {
				fmt.Println("There is error2", err.Error())
				return sql
			}
			var columns2 []string
			for _, col := range table.Columns {
				columns2 = append(columns2, col.Name)
			}
			insertStmt := fmt.Sprintf("INSERT OR REPLACE INTO %v (%v) VALUES (%v);\n",
				table.Name,
				joinColumnNames(columns2, ", "),
				formatInsert(values, table.Columns),
			)
			sql += insertStmt
			// callback(insertStmt)
			if v, ok := values[0].(int64); ok {
				unique_ids = append(unique_ids, int(v))
				// sql += generateSQLOfTable(db, int(v), table.Name, relatedTables, callback)
			}
		}
		processedTables = append(processedTables, table.Name)
		if len(table.Children) == 0 {
			continue
		}
		for _, v := range unique_ids {
			sql += generateSQLOfTable(db, v, table.Children, processedTables)
		}
	}
	return sql
}

func findRelatedTables(tables []TableWithColumns, tableName string) ([]*TableWithColumnsAndASpecialColumn, error) {
	var relatedTables []*TableWithColumnsAndASpecialColumn
	for _, table := range tables {
		for _, column := range table.Columns {
			if column.References == tableName {
				relatedTables = append(relatedTables, &TableWithColumnsAndASpecialColumn{
					Name:     table.Name,
					Columns:  table.Columns,
					Column:   column,
					Children: []*TableWithColumnsAndASpecialColumn{},
				})
			}
		}
	}
	// for _, table := range tables {
	// 	if table.Name == tableName {
	// 		return &TableWithColumnsAndASpecialColumn{
	// 			Name:     table.Name,
	// 			Columns:  table.Columns,
	// 			Column:   TableColumn{},
	// 			Children: relatedTables,
	// 		}, nil
	// 	}
	// }
	// return nil, errors.New("没有找到匹配的 table")
	return relatedTables, nil
}
func findRelatedTablesRec(tableName string, tables []TableWithColumns) (*TableWithColumnsAndASpecialColumn, error) {
	relatedTables, err := findRelatedTables(tables, tableName)
	if err != nil {
		return nil, err
	}
	var t *TableWithColumnsAndASpecialColumn
	for _, table := range tables {
		if table.Name == tableName {
			t = &TableWithColumnsAndASpecialColumn{
				Name:     table.Name,
				Columns:  table.Columns,
				Column:   TableColumn{},
				Children: relatedTables,
			}
		}
	}
	if t == nil {
		return nil, errors.New("没有找到匹配的 table")
	}
	t.Children = relatedTables
	// fmt.Printf("%v 's related tables is\n", t.Name)
	// for i, related := range relatedTables {
	// 	fmt.Printf("%v、%v\n", i, related.Name)
	// }
	for i, related := range relatedTables {
		sub, err := findRelatedTablesRec(related.Name, tables)
		if err != nil {
			continue
		}
		relatedTables[i] = &TableWithColumnsAndASpecialColumn{
			Name:     related.Name,
			Columns:  related.Columns,
			Column:   related.Column,
			Children: sub.Children,
		}
	}
	// fmt.Printf("\n")
	return t, nil
}

func generateSQL(db *sql.DB, tableName string, id int) (string, error) {
	var insertStmt string

	tables, err := FetchTables(db)
	// fmt.Printf("after return tables count is %v \n", len(tables))
	if err != nil {
		return insertStmt, err
	}
	var columns []TableColumn
	for _, t := range tables {
		if t.Name == tableName {
			columns = t.Columns
			break
		}
	}
	if columns == nil {
		return "", errors.New("没有找到匹配的 table")
	}
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = ?", tableName)
	row := db.QueryRow(query, id)
	values := make([]interface{}, len(columns))
	dest := make([]interface{}, len(columns))
	for i := range dest {
		dest[i] = &values[i]
	}
	if err := row.Scan(dest...); err != nil {
		return "", err
	}
	var columns2 []string
	for _, col := range columns {
		columns2 = append(columns2, col.Name)
	}
	insertStmt = fmt.Sprintf("INSERT OR REPLACE INTO %v (%v) VALUES (%v);\n",
		tableName,
		joinColumnNames(columns2, ", "),
		formatInsert(values, columns),
	)
	var simpleTables []TableWithColumns
	for _, t := range tables {
		simpleTables = append(simpleTables, TableWithColumns{
			Name:    t.Name,
			Columns: t.Columns,
		})
	}
	tree, err := findRelatedTablesRec(tableName, simpleTables)
	if err != nil {
		return "", err
	}
	processedTables := []string{}
	// fmt.Printf("the table tree is %v", tree.Name)
	if v, ok := values[0].(int64); ok {
		insertStmt += generateSQLOfTable(db, int(v), tree.Children, processedTables)
	}
	return insertStmt, nil
}

func joinColumnNames(values []string, seq string) string {
	result := ""
	count := len(values)
	for i, v := range values {
		result = result + fmt.Sprintf("`%v`", v)
		if i != count-1 {
			result += seq
		}
	}
	return result
}
func joinStrWithSEQ(values []string, seq string) string {
	result := ""
	count := len(values)
	for i, v := range values {
		result = result + fmt.Sprintf("%v", v)
		if i != count-1 {
			result += seq
		}
	}
	return result
}

func formatInsert(values []interface{}, columns []TableColumn) string {
	stmt := ""
	count := len(values)
	for i, col := range columns {
		// fmt.Println("format insert", col.Name, col.Type)
		v := values[i]
		if v == nil {
			stmt += "NULL"
			if i < count-1 {
				stmt += ", "
			}
			continue
		}
		if col.Type == "text" {
			vv := strings.ReplaceAll(fmt.Sprintf("%v", v), "'", "''")
			stmt += fmt.Sprintf("'%v'", vv)
			if i < count-1 {
				stmt += ", "
			}
			continue
		}
		if col.Type == "integer" {
			stmt += fmt.Sprintf("%v", v)
			if i < count-1 {
				stmt += ", "
			}
			continue
		}
		if col.Type == "datetime" {
			if v2, ok := v.(time.Time); ok {
				formattedTime := v2.Format("2006-01-02 15:04:05")
				stmt += fmt.Sprintf("'%v'", formattedTime)
			}
			if i < count-1 {
				stmt += ", "
			}
			continue
		}
	}
	return stmt
}

func FetchTables(db *sql.DB) ([]TableWithColumnsAndRAW, error) {
	query := "SELECT name, sql FROM sqlite_master WHERE type='table';"
	rows, err := db.Query(query)
	if err != nil {
		log.Fatalf("Failed to fetch table names: %v", err)
	}
	defer rows.Close()
	// var tables []pkg.Table
	var tables []TableWithColumnsAndRAW
	for rows.Next() {
		var name string
		var profile string
		if err := rows.Scan(&name, &profile); err != nil {
			return tables, err
		}
		// table := pkg.ParseTableCreateStatement(name, profile)
		columns, err := getTableColumns(db, name)
		if err != nil {
			// fmt.Printf("find table %v happen error %v,  \n", name, err.Error())
			return tables, err
		}
		// if table.Columns == nil {
		// 	fmt.Println(name, profile)
		// }
		// if table.Columns != nil {
		// }
		// fmt.Printf("before append table %v \n", name)
		tables = append(tables, TableWithColumnsAndRAW{
			Name:    name,
			Columns: columns,
			RAW:     profile,
		})
	}
	if err := rows.Err(); err != nil {
		return tables, err
	}
	// fmt.Printf("before return tables count is %v \n", len(tables))
	return tables, nil
}
func executeSelectQuery(db *sql.DB, sqlStatement string) ([][]interface{}, error) {
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return [][]interface{}{}, err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return [][]interface{}{}, err
	}
	var result [][]interface{}
	for rows.Next() {
		pointers := make([]interface{}, len(columns))
		values := make([]interface{}, len(columns))
		for i := range values {
			pointers[i] = &values[i]
		}
		if err := rows.Scan(pointers...); err != nil {
			return result, err
		}
		result = append(result, values)
	}
	if err := rows.Err(); err != nil {
		return result, err
	}
	return result, nil
}

var env string

func main() {
	var url string
	var port string

	flag.Parse()
	args := flag.Args()
	if len(args) > 0 {
		url = args[0]
	}
	for i, arg := range args {
		if (arg == "--port" || arg == "-port") && i+1 < len(args) {
			portValue := args[i+1]
			port = portValue
		}
	}
	if url == "" {
		fmt.Println("请指定数据库文件路径，例如 sqliteweb ./test.db")
		return
	}
	if port == "" {
		port = "8000"
	}
	if env == "" {
		env = "dev"
	}
	_, err := os.Stat(url)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("文件不存在")
			return
		}
		fmt.Println(err.Error())
		return
	}
	db, err := sql.Open("sqlite3", url)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	if env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	r.GET("/", returnIndexHTML)
	staticFS, _ := fs.Sub(staticFiles, "frontend/dist/assets")
	r.NoRoute(returnIndexHTML)
	r.StaticFS("/assets", http.FS(staticFS))

	apiV1 := r.Group("/api/v1")
	route := apiV1.Group("/database")
	{
		route.POST("/tables", func(c *gin.Context) {
			tables, err := FetchTables(db)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"code":  1,
					"error": err.Error(),
					"data":  nil,
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"code":  0,
				"error": "",
				"data":  tables,
			})
		})

		route.POST("/exec", func(c *gin.Context) {
			var body struct {
				Query string `json:"query"`
			}
			if err := c.Bind(&body); err != nil {
				c.JSON(http.StatusOK, gin.H{
					"code": 1,
					"msg":  "无效的请求体",
					"data": nil,
				})
				return
			}
			queryType := strings.ToUpper(strings.Split(body.Query, " ")[0])
			if queryType == "SELECT" {
				rows, err := executeSelectQuery(db, body.Query)
				if err != nil {
					c.JSON(http.StatusOK, gin.H{
						"code": 1,
						"msg":  err.Error(),
						"data": nil,
					})
					return
				}
				c.JSON(http.StatusOK, gin.H{
					"code": 0,
					"msg":  "",
					"data": rows,
				})
			} else {
				result, err := db.Exec(body.Query)
				if err != nil {
					c.JSON(http.StatusOK, gin.H{
						"code": 1,
						"msg":  err.Error(),
						"data": nil,
					})
					return
				}
				rowsAffected, err := result.RowsAffected()
				if err != nil {
					c.JSON(http.StatusOK, gin.H{
						"code": 1,
						"msg":  err.Error(),
						"data": nil,
					})
					return
				}
				c.JSON(http.StatusOK, gin.H{
					"code": 0,
					"msg":  strconv.Itoa(int(rowsAffected)) + "条记录被影响",
					"data": nil,
				})
			}
		})

		route.POST("/export", func(c *gin.Context) {
			var body struct {
				TableName string      `json:"table_name"`
				Id        interface{} `json:"id"`
			}
			if err := c.Bind(&body); err != nil {
				c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "无效的请求体", "data": nil})
				return
			}
			if body.TableName == "" {
				c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "缺少TableName", "data": nil})
				return
			}
			if body.Id == nil {
				c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "缺少记录id", "data": nil})
				return
			}
			v, ok := body.Id.(float64)
			if !ok {
				c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "id 类型错误", "data": nil})
				return
			}
			sql, err := generateSQL(db, body.TableName, int(v))
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"code": 1, "msg": err.Error(), "data": nil})
				return
			}
			sqlFile, err := os.Create(fmt.Sprintf("%v_%v.sql", body.TableName, body.Id))
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"code": 1, "msg": err.Error(), "data": nil})
				return
			}
			sqlFile.WriteString(sql)
			c.JSON(http.StatusOK, gin.H{
				"code": 0,
				"msg":  "",
				"data": sql,
			})
		})
		route.POST("/import", func(c *gin.Context) {
			file, err := c.FormFile("file")
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "Failed to get file", "data": nil})
				return
			}
			tempFilePath := "./" + file.Filename
			if err := c.SaveUploadedFile(file, tempFilePath); err != nil {
				c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "Failed to save file", "data": nil})
				return
			}
			defer os.Remove(tempFilePath)
			sqlBytes, err := os.ReadFile(tempFilePath)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "Failed to read file", "data": nil})
				return
			}
			sqlContent := string(sqlBytes)
			sqlStatements := strings.Split(sqlContent, ";\n")
			type ExecError struct {
				Err string `json:"err"`
				SQL string `json:"sql"`
			}
			var errors []ExecError
			for _, stmt := range sqlStatements {
				stmt = strings.TrimSpace(stmt) // 去除前后空白
				if stmt == "" {
					continue // 如果是空语句则跳过
				}
				_, err = db.Exec(fmt.Sprintf("%v;", stmt))
				if err != nil {
					fmt.Println("Exec sql failed", stmt)
					fmt.Println(err.Error())
					errors = append(errors, ExecError{
						Err: err.Error(),
						SQL: stmt,
					})
					continue
				}
			}
			if len(errors) != 0 {
				msg, err := json.Marshal(errors)
				if err != nil {
					c.JSON(http.StatusOK, gin.H{"code": 1, "msg": err.Error(), "data": nil})
					return
				}
				c.JSON(http.StatusOK, gin.H{"code": 1, "msg": msg, "data": nil})
				return
			}
			c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "SQL executed successfully", "data": gin.H{
				"msg": "共" + strconv.Itoa(len(sqlStatements)) + "条sql",
			}})
		})

		// 新增列 ALTER TABLE `words` ADD `text` text NOT NULL
		// CREATE UNIQUE INDEX `idx_line_subtitle` ON `paragraphs`(`line`,`subtitle_id`)

		route.POST("/remove_column", func(c *gin.Context) {
			var body struct {
				Table      TableWithColumns `json:"table"`
				ColumnName string           `json:"column_name"`
			}
			if err := c.Bind(&body); err != nil {
				c.JSON(http.StatusOK, gin.H{"code": 1, "msg": err.Error(), "data": nil})
				return
			}
			if body.Table.Name == "" {
				c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "缺少 table", "data": nil})
				return
			}
			if body.ColumnName == "" {
				c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "缺少 column name", "data": nil})
				return
			}
			// 开始一个事务
			tx, err := db.Begin()
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"code": 1, "msg": err.Error(), "data": nil})
				return
			}
			defer func() {
				fmt.Println("invoke defer func() {")
				if err != nil {
					fmt.Println("发生错误，回滚事务:", err)
					tx.Rollback()
					c.JSON(http.StatusOK, gin.H{"code": 1, "msg": err.Error(), "data": nil})
					return
				}
				// 提交事务
				err = tx.Commit()
				if err != nil {
					c.JSON(http.StatusOK, gin.H{"code": 1, "msg": err.Error(), "data": nil})
					return
				}
				c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "", "data": nil})
			}()
			var columnsRemaining []TableColumn
			var columnsRemaining2 []string
			for _, col := range body.Table.Columns {
				if col.Name != body.ColumnName {
					columnsRemaining = append(columnsRemaining, col)
					columnsRemaining2 = append(columnsRemaining2, col.Name)
				}
			}
			tmpTableName := "new_" + body.Table.Name
			stmt1 := buildSQLCreateTableWithColumns(tmpTableName, columnsRemaining)
			// "CREATE TABLE new_employees (id INTEGER PRIMARY KEY, name TEXT, salary REAL);"
			_, err = tx.Exec(stmt1)
			if err != nil {
				return
			}
			// 2. 将旧表的数据插入到新表中
			colStr := joinColumnNames(columnsRemaining2, ", ")
			stmt2 := fmt.Sprintf("INSERT INTO %v (%v) SELECT %v FROM %v;", tmpTableName, colStr, colStr, body.Table.Name)
			fmt.Println(stmt2)
			_, err = tx.Exec(stmt2)
			if err != nil {
				return
			}
			// 3. 删除旧表
			stmt3 := fmt.Sprintf("DROP TABLE %v;", body.Table.Name)
			fmt.Println(stmt3)
			_, err = tx.Exec(stmt3)
			if err != nil {
				return
			}
			// 4. 将新表重命名为旧表的名称
			stmt4 := fmt.Sprintf("ALTER TABLE %v RENAME TO %v;", tmpTableName, body.Table.Name)
			fmt.Println(stmt4)
			_, err = tx.Exec(stmt4)
			if err != nil {
				return
			}
			fmt.Println("成功删除列并更新表结构！")
		})
	}
	host := "127.0.0.1"
	if env == "prod" {
		host = "0.0.0.0"
	}
	web := host + ":" + port
	fmt.Printf("Server is running at: %v", web)
	r.Run(web)
}
