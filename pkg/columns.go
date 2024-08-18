package pkg

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

type Column struct {
	Name          string `json:"name"`
	Type          string `json:"type"`
	NotNull       bool   `json:"not_null"`
	PrimaryKey    bool   `json:"primary_key"`
	ForeignKey    bool   `json:"foreign_key"`
	AutoIncrement bool   `json:"auto_increment"`
	// Default       int    `json:"default"`
	Unique     bool   `json:"unique"`
	References string `json:"references"`
}
type PendingColumn struct {
	Column
	Index int `json:"index"`
}

type Table struct {
	Name    string   `json:"name"`
	Raw     string   `json:"raw"`
	Columns []Column `json:"columns"`
}

func ParseTableCreateStatement(name string, content string) Table {
	var result Table
	// Extract table name
	// parts := strings.Split(content, "|")
	// if len(parts) != 2 {
	// 	return result
	// }
	// result.Name = parts[0]
	result.Name = name
	result.Raw = content

	regexp1 := regexp.MustCompile("CREATE TABLE [`\"]([a-zA-Z0-9_-]{1,})[`\"] {1,}\\(([.\\s\\S]{1,})\\)")
	matches := regexp1.FindStringSubmatch(content)

	if len(matches) < 2 {
		return result
	}
	tableName := matches[1]
	columnStr := matches[2]

	result.Name = tableName

	columnStrParts := strings.Split(columnStr, ",")
	// fmt.Println(columnStrParts)

	columnNameReg := regexp.MustCompile("^[`\"]([a-zA-Z0-9_-]{1,})[`\"]")
	columnTypeReg := regexp.MustCompile("(TEXT|CLOB|INT|BIGINT|INTEGER|REAL|DOUBLE|FLOAT|DATETIME|NONE|BLOB|text|clob|int|bigint|integer|real|double|float|datetime|none|blob|varchar\\([0-9]{1,}\\))")
	columnPrimaryKeyReg := regexp.MustCompile("PRIMARY KEY")
	columnReg4 := regexp.MustCompile("AUTOINCREMENT")
	columnReg5 := regexp.MustCompile("NOT NULL")

	columnsMap := make(map[string]PendingColumn)
	// columns := make([]Column, len(columnStrParts))

	for i, columnStr := range columnStrParts {
		var c PendingColumn
		c.Index = i
		c.ForeignKey = false
		c.Unique = false
		columnStr = strings.TrimSpace(columnStr)
		if strings.Contains(columnStr, "CONSTRAINT") {
			// fmt.Println("process column str", columnStr)
			re1 := regexp.MustCompile("UNIQUE \\([`\"]([a-zA-Z0-9_-]{1,})[`\"]\\)")
			// 约束某个字段唯一
			// CONSTRAINT `uni_media_episode_profiles_unique_id` UNIQUE (`unique_id`)
			m1 := re1.FindStringSubmatch(columnStr)
			// fmt.Println("is unique", m1)
			if len(m1) == 2 {
				col, has := columnsMap[m1[1]]
				if has {
					col.Unique = true
					columnsMap[col.Name] = col
				}
				continue
			}
			// 外键
			// CONSTRAINT `fk_media_profiles_episodes` FOREIGN KEY (`media_profile_id`) REFERENCES `media_profiles`(`id`),
			re2 := regexp.MustCompile("FOREIGN KEY \\([`\"]([a-zA-Z0-9_-]{1,})[`\"]\\) REFERENCES [`\"]([a-zA-Z0-9_-]{1,})[`\"] {0,1}\\([`\"]([a-zA-Z0-9_-]{1,})[`\"]\\)")
			m2 := re2.FindStringSubmatch(columnStr)
			// fmt.Println("has reference table", m2, len(m2))
			if len(m2) == 4 {
				col, has := columnsMap[m2[1]]
				// fmt.Printf("has reference table 2 %v", col)
				if has {
					col.ForeignKey = true
					col.References = m2[2]
					columnsMap[col.Name] = col
				}
				// _, has = columnsMap[m2[2]]
				// if !has {
				// 	c.Name = m2[2]
				// 	c.Type = "table"
				// 	columnsMap[c.Name] = c
				// }
				continue
			}
			continue
		}
		// 获取字段名
		matches1 := columnNameReg.FindStringSubmatch(columnStr)
		if len(matches1) < 2 {
			fmt.Println("matches1", matches1, columnStr)
			continue
		}
		c.Name = matches1[1]
		// 获取类型
		name := matches1[0]
		columnStr = strings.TrimSpace(strings.Replace(columnStr, name, "", 1))
		matches2 := columnTypeReg.FindStringSubmatch(columnStr)
		if len(matches2) < 2 {
			fmt.Println("matches2", matches2, columnStr)
			continue
		}
		varcharReg := regexp.MustCompile("varchar")
		c.Type = strings.ToLower(matches2[1])
		varcharMatched := varcharReg.FindStringSubmatch(matches2[1])
		if len(varcharMatched) == 1 {
			c.Type = "text"
		}
		// 是否为主键
		columnStr = strings.TrimSpace(strings.Replace(columnStr, matches2[0], "", 1))
		matches3 := columnPrimaryKeyReg.FindStringSubmatch(columnStr)
		// fmt.Println("matches3", matches3)
		c.PrimaryKey = false
		if len(matches3) == 1 {
			c.PrimaryKey = true
			columnStr = strings.TrimSpace(strings.Replace(columnStr, matches3[0], "", 1))
		}
		// 是否自增
		matches4 := columnReg4.FindStringSubmatch(columnStr)
		c.AutoIncrement = false
		if len(matches4) == 1 {
			c.AutoIncrement = true
			columnStr = strings.TrimSpace(strings.Replace(columnStr, matches4[0], "", 1))
		}
		// 不能为空
		matches5 := columnReg5.FindStringSubmatch(columnStr)
		c.NotNull = false
		if len(matches5) == 1 {
			c.NotNull = true
		}
		// columns[i] = c
		// fmt.Println("before columnsMap[c.Name]")
		columnsMap[c.Name] = c
	}
	// fmt.Println(columnsMap)
	columns := make([]PendingColumn, 0, len(columnsMap))
	for _, value := range columnsMap {
		columns = append(columns, value)
	}
	sort.Slice(columns, func(i, j int) bool {
		if columns[i].ForeignKey && !columns[j].ForeignKey {
			return false
		}
		if !columns[i].ForeignKey && columns[j].ForeignKey {
			return true
		}
		return columns[i].Index < columns[j].Index
	})
	// result.Columns = columns
	// fmt.Println(columns)
	var c []Column
	for _, value := range columns {
		c = append(c, Column{
			Name:          value.Name,
			Type:          value.Type,
			NotNull:       value.NotNull,
			PrimaryKey:    value.PrimaryKey,
			ForeignKey:    value.ForeignKey,
			AutoIncrement: value.AutoIncrement,
			Unique:        value.Unique,
			References:    value.References,
		})
	}
	result.Columns = c

	// for _, colDef := range columnDefs {
	// 	colName := colDef[1]
	// 	colType := colDef[2]
	// 	colConstraints := strings.TrimSpace(colDef[3])

	// 	column := Column{
	// 		Name:       colName,
	// 		Type:       colType,
	// 		NotNull:    strings.Contains(strings.ToUpper(colConstraints), "NOT NULL"),
	// 		PrimaryKey: strings.Contains(strings.ToUpper(colConstraints), "PRIMARY KEY"),
	// 		ForeignKey: false,
	// 		References: "",
	// 	}

	// 	columns = append(columns, column)
	// }

	// Check for foreign keys
	// fkRegex := regexp.MustCompile(`(?i)FOREIGN KEY\s*\((\w+)\)\s*REFERENCES\s*(\w+)\s*\((\w+)\)`)
	// fkMatches := fkRegex.FindAllStringSubmatch(content, -1)
	// for _, fkMatch := range fkMatches {
	// 	fkColumn := fkMatch[1]
	// 	fkRefTable := fkMatch[2]
	// 	fkRefColumn := fkMatch[3]

	// 	// Find the column and mark it as a foreign key
	// 	for i := range columns {
	// 		if columns[i].Name == fkColumn {
	// 			columns[i].ForeignKey = true
	// 			columns[i].References = fmt.Sprintf("%s(%s)", fkRefTable, fkRefColumn)
	// 			break
	// 		}
	// 	}
	// }

	return result
}
