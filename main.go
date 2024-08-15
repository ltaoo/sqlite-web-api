package main

import (
	"database/sql"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
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
	Cid       int         `json:"cid"`
	Name      string      `json:"name"`
	Type      string      `json:"type"`
	NotNull   int         `json:"not_null"`
	DfltValue interface{} `json:"value"`
	Pk        int         `json:"pk"`
}
type TableWithColumns struct {
	Name    string        `json:"name"`
	Columns []TableColumn `json:"columns"`
}

func getTableColumns(db *sql.DB, tableName string) ([]TableColumn, error) {
	query := fmt.Sprintf("PRAGMA table_info(%s);", tableName)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []TableColumn
	for rows.Next() {
		var col TableColumn
		if err := rows.Scan(&col.Cid, &col.Name, &col.Type, &col.NotNull, &col.DfltValue, &col.Pk); err != nil {
			return nil, err
		}
		columns = append(columns, col)
	}

	return columns, nil
}

func FetchTables(db *sql.DB) ([]TableWithColumns, error) {
	query := "SELECT name FROM sqlite_master WHERE type='table';"
	rows, err := db.Query(query)
	if err != nil {
		log.Fatalf("Failed to fetch table names: %v", err)
	}
	defer rows.Close()
	var tables []TableWithColumns
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return tables, err
		}
		columns, err := getTableColumns(db, tableName)
		if err != nil {
			return tables, err
		}
		table := TableWithColumns{
			Name:    tableName,
			Columns: columns,
		}
		tables = append(tables, table)
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

func main() {
	var url string

	flag.StringVar(&url, "url", "", "The Sqlite3 db file")
	flag.Parse()

	if url == "" {
		fmt.Println("Error: the media_id is required")
		flag.Usage()
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

	// 启动服务器
	r.Run("127.0.0.1:8000")
}
