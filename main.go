package main

import (
	"database/sql"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
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
	RAW     string        `json:"raw"`
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
		foreign, has := foreignKeys[col.Name]
		if has {
			col.References = foreign.TableName
		}
		columns = append(columns, col)
	}
	return columns, nil
}

func FetchTables(db *sql.DB) ([]TableWithColumns, error) {
	query := "SELECT name, sql FROM sqlite_master WHERE type='table';"
	rows, err := db.Query(query)
	if err != nil {
		log.Fatalf("Failed to fetch table names: %v", err)
	}
	defer rows.Close()
	// var tables []pkg.Table
	var tables []TableWithColumns
	for rows.Next() {
		var name string
		var profile string
		if err := rows.Scan(&name, &profile); err != nil {
			return tables, err
		}
		// table := pkg.ParseTableCreateStatement(name, profile)
		columns, err := getTableColumns(db, name)
		if err != nil {
			return tables, err
		}
		// if table.Columns == nil {
		// 	fmt.Println(name, profile)
		// }
		// if table.Columns != nil {
		// }
		tables = append(tables, TableWithColumns{
			Name:    name,
			Columns: columns,
			RAW:     profile,
		})
	}
	if err := rows.Err(); err != nil {
		return tables, err
	}
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
					"code":  1,
					"error": "无效的请求体",
					"data":  nil,
				})
				return
			}
			queryType := strings.ToUpper(strings.Split(body.Query, " ")[0])
			if queryType == "SELECT" {
				rows, err := executeSelectQuery(db, body.Query)
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
					"data":  rows,
				})
			} else {
				result, err := db.Exec(body.Query)
				if err != nil {
					c.JSON(http.StatusOK, gin.H{
						"code":  1,
						"error": err.Error(),
						"data":  nil,
					})
					return
				}
				rowsAffected, err := result.RowsAffected()
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
					"error": strconv.Itoa(int(rowsAffected)) + "条记录被影响",
					"data":  nil,
				})
			}
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
