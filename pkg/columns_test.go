package pkg_test

import (
	"fmt"
	"os"
	"reflect"
	"sqlite-web/pkg"
	"testing"
)

func TestParseSQLQuery(t *testing.T) {
	// t.Run("CREATE Table 1 with line break", func(c *testing.T) {
	// 	bytes, err := os.ReadFile("./sql_with_linebreak.txt")
	// 	if err != nil {
	// 		t.Errorf("Read file content failed, %v", err.Error())
	// 		return
	// 	}
	// 	expect := pkg.Table{
	// 		Name: "users",
	// 		Columns: []pkg.Column{
	// 			pkg.Column{
	// 				Index:         0,
	// 				Name:          "id",
	// 				Type:          "integer",
	// 				NotNull:       false,
	// 				PrimaryKey:    true,
	// 				AutoIncrement: true,
	// 				ForeignKey:    false,
	// 				Unique:        false,
	// 				References:    "",
	// 			},
	// 			pkg.Column{
	// 				Index:         1,
	// 				Name:          "created_at",
	// 				Type:          "datetime",
	// 				NotNull:       false,
	// 				PrimaryKey:    false,
	// 				AutoIncrement: false,
	// 				ForeignKey:    false,
	// 				Unique:        false,
	// 				References:    "",
	// 			},
	// 			pkg.Column{
	// 				Index:         2,
	// 				Name:          "name",
	// 				Type:          "varchar(100)",
	// 				NotNull:       true,
	// 				PrimaryKey:    false,
	// 				AutoIncrement: false,
	// 				ForeignKey:    false,
	// 				Unique:        true,
	// 				References:    "",
	// 			},
	// 			pkg.Column{
	// 				Index:         3,
	// 				Name:          "avatar",
	// 				Type:          "text",
	// 				NotNull:       false,
	// 				PrimaryKey:    false,
	// 				AutoIncrement: false,
	// 				ForeignKey:    false,
	// 				Unique:        false,
	// 				References:    "",
	// 			},
	// 			pkg.Column{
	// 				Index:         4,
	// 				Name:          "subtitle_id",
	// 				Type:          "integer",
	// 				NotNull:       false,
	// 				PrimaryKey:    false,
	// 				AutoIncrement: false,
	// 				ForeignKey:    true,
	// 				Unique:        false,
	// 				References:    "",
	// 			},
	// 			pkg.Column{
	// 				Index:         5,
	// 				Name:          "subtitles",
	// 				Type:          "table",
	// 				NotNull:       false,
	// 				PrimaryKey:    false,
	// 				AutoIncrement: false,
	// 				ForeignKey:    false,
	// 				Unique:        false,
	// 				References:    "",
	// 			},
	// 		},
	// 	}
	// 	content := string(bytes)
	// 	result := pkg.ParseTableCreateStatement("users", content)
	// 	fmt.Println("after pkg.ParseTableCreateStatement")
	// 	if result.Name != expect.Name {
	// 		t.Errorf("Test failed, The table name %v is different %v", result.Name, expect.Name)
	// 		return
	// 	}
	// 	if len(expect.Columns) != len(result.Columns) {
	// 		t.Errorf("Test failed, The columns len %v is different %v", len(result.Columns), len(expect.Columns))
	// 		return
	// 	}
	// 	for i, expectedColumn := range expect.Columns {
	// 		fmt.Println("------")
	// 		actuallyColumn := result.Columns[i]
	// 		expectedValue := reflect.ValueOf(expectedColumn)
	// 		actuallyValue := reflect.ValueOf(actuallyColumn)
	// 		typ := actuallyValue.Type()
	// 		for j := 0; j < actuallyValue.NumField(); j++ {
	// 			field := typ.Field(j).Name
	// 			actuallyV := fmt.Sprintf("%v", actuallyValue.Field(j).Interface())
	// 			expectedV := fmt.Sprintf("%v", expectedValue.Field(j).Interface())
	// 			fmt.Println(field, "'"+actuallyV+"'/'"+expectedV+"'", reflect.TypeOf(actuallyV))
	// 			if actuallyV != expectedV {
	// 				t.Errorf("%v | %v '%v' - '%v'", expectedColumn, field, actuallyV, expectedV)
	// 			}
	// 		}
	// 	}
	// })

	// t.Run("CREATE Table 1 with line break", func(c *testing.T) {
	// 	bytes, err := os.ReadFile("./test1.txt")
	// 	if err != nil {
	// 		t.Errorf("Read file content failed, %v", err.Error())
	// 		return
	// 	}
	// 	expect := pkg.Table{
	// 		Name: "word_in_paragraphs",
	// 		Columns: []pkg.Column{
	// 			pkg.Column{
	// 				Name:          "id",
	// 				Type:          "integer",
	// 				NotNull:       false,
	// 				PrimaryKey:    true,
	// 				AutoIncrement: true,
	// 				ForeignKey:    false,
	// 				Unique:        false,
	// 				References:    "",
	// 			},
	// 			pkg.Column{
	// 				Name:          "text",
	// 				Type:          "text",
	// 				NotNull:       true,
	// 				PrimaryKey:    false,
	// 				AutoIncrement: false,
	// 				ForeignKey:    false,
	// 				Unique:        false,
	// 				References:    "",
	// 			},
	// 			pkg.Column{
	// 				Name:          "sw",
	// 				Type:          "text",
	// 				NotNull:       false,
	// 				PrimaryKey:    false,
	// 				AutoIncrement: false,
	// 				ForeignKey:    false,
	// 				Unique:        false,
	// 				References:    "",
	// 			},
	// 			pkg.Column{
	// 				Name:          "phonetic",
	// 				Type:          "text",
	// 				NotNull:       false,
	// 				PrimaryKey:    false,
	// 				AutoIncrement: false,
	// 				ForeignKey:    false,
	// 				Unique:        false,
	// 				References:    "",
	// 			},
	// 			pkg.Column{
	// 				Name:          "definition",
	// 				Type:          "text",
	// 				NotNull:       false,
	// 				PrimaryKey:    false,
	// 				AutoIncrement: false,
	// 				ForeignKey:    false,
	// 				Unique:        false,
	// 				References:    "",
	// 			},
	// 			pkg.Column{
	// 				Name:          "translation",
	// 				Type:          "text",
	// 				NotNull:       false,
	// 				PrimaryKey:    false,
	// 				AutoIncrement: false,
	// 				ForeignKey:    false,
	// 				Unique:        false,
	// 				References:    "",
	// 			},
	// 			pkg.Column{
	// 				Name:          "pos",
	// 				Type:          "text",
	// 				NotNull:       false,
	// 				PrimaryKey:    false,
	// 				AutoIncrement: false,
	// 				ForeignKey:    false,
	// 				Unique:        false,
	// 				References:    "",
	// 			},
	// 			pkg.Column{
	// 				Name:          "collins",
	// 				Type:          "integer",
	// 				NotNull:       false,
	// 				PrimaryKey:    false,
	// 				AutoIncrement: false,
	// 				ForeignKey:    false,
	// 				Unique:        false,
	// 				References:    "",
	// 			},
	// 			pkg.Column{
	// 				Name:          "oxford",
	// 				Type:          "integer",
	// 				NotNull:       false,
	// 				PrimaryKey:    false,
	// 				AutoIncrement: false,
	// 				ForeignKey:    false,
	// 				Unique:        false,
	// 				References:    "",
	// 			},
	// 			pkg.Column{
	// 				Name:          "tag",
	// 				Type:          "text",
	// 				NotNull:       false,
	// 				PrimaryKey:    false,
	// 				AutoIncrement: false,
	// 				ForeignKey:    false,
	// 				Unique:        false,
	// 				References:    "",
	// 			},
	// 			pkg.Column{
	// 				Name:          "bnc",
	// 				Type:          "integer",
	// 				NotNull:       false,
	// 				PrimaryKey:    false,
	// 				AutoIncrement: false,
	// 				ForeignKey:    false,
	// 				Unique:        false,
	// 				References:    "",
	// 			},
	// 			pkg.Column{
	// 				Name:          "frq",
	// 				Type:          "integer",
	// 				NotNull:       false,
	// 				PrimaryKey:    false,
	// 				AutoIncrement: false,
	// 				ForeignKey:    false,
	// 				Unique:        false,
	// 				References:    "",
	// 			},
	// 			pkg.Column{
	// 				Name:          "exchange",
	// 				Type:          "text",
	// 				NotNull:       false,
	// 				PrimaryKey:    false,
	// 				AutoIncrement: false,
	// 				ForeignKey:    false,
	// 				Unique:        false,
	// 				References:    "",
	// 			},
	// 			pkg.Column{
	// 				Name:          "start",
	// 				Type:          "integer",
	// 				NotNull:       false,
	// 				PrimaryKey:    false,
	// 				AutoIncrement: false,
	// 				ForeignKey:    false,
	// 				Unique:        false,
	// 				References:    "",
	// 			},
	// 			pkg.Column{
	// 				Name:          "end",
	// 				Type:          "integer",
	// 				NotNull:       false,
	// 				PrimaryKey:    false,
	// 				AutoIncrement: false,
	// 				ForeignKey:    false,
	// 				Unique:        false,
	// 				References:    "",
	// 			},
	// 			pkg.Column{
	// 				Name:          "subtitle_id",
	// 				Type:          "integer",
	// 				NotNull:       false,
	// 				PrimaryKey:    false,
	// 				AutoIncrement: false,
	// 				ForeignKey:    true,
	// 				Unique:        false,
	// 				References:    "subtitles",
	// 			},
	// 			pkg.Column{
	// 				Name:          "paragraph_id",
	// 				Type:          "integer",
	// 				NotNull:       false,
	// 				PrimaryKey:    false,
	// 				AutoIncrement: false,
	// 				ForeignKey:    true,
	// 				Unique:        false,
	// 				References:    "paragraphs",
	// 			},
	// 			pkg.Column{
	// 				Name:          "profile_id",
	// 				Type:          "integer",
	// 				NotNull:       false,
	// 				PrimaryKey:    false,
	// 				AutoIncrement: false,
	// 				ForeignKey:    true,
	// 				Unique:        false,
	// 				References:    "words",
	// 			},
	// 		},
	// 	}
	// 	content := string(bytes)
	// 	result := pkg.ParseTableCreateStatement("word_in_paragraphs", content)
	// 	fmt.Println("after pkg.ParseTableCreateStatement")
	// 	if result.Name != expect.Name {
	// 		t.Errorf("Test failed, The table name %v is different %v", result.Name, expect.Name)
	// 		return
	// 	}
	// 	if len(expect.Columns) != len(result.Columns) {
	// 		t.Errorf("Test failed, The columns len %v is different %v", len(result.Columns), len(expect.Columns))
	// 		return
	// 	}
	// 	for i, expectedColumn := range expect.Columns {
	// 		fmt.Println("------")
	// 		actuallyColumn := result.Columns[i]
	// 		expectedValue := reflect.ValueOf(expectedColumn)
	// 		actuallyValue := reflect.ValueOf(actuallyColumn)
	// 		typ := actuallyValue.Type()
	// 		for j := 0; j < actuallyValue.NumField(); j++ {
	// 			field := typ.Field(j).Name
	// 			actuallyV := fmt.Sprintf("%v", actuallyValue.Field(j).Interface())
	// 			expectedV := fmt.Sprintf("%v", expectedValue.Field(j).Interface())
	// 			fmt.Println(field, "'"+actuallyV+"'/'"+expectedV+"'", reflect.TypeOf(actuallyV))
	// 			if actuallyV != expectedV {
	// 				t.Errorf("%v | %v '%v' - '%v'", expectedColumn, field, actuallyV, expectedV)
	// 			}
	// 		}
	// 	}
	// })

	// t.Run("CREATE Table 1 with line break", func(c *testing.T) {
	// 	bytes, err := os.ReadFile("./test2.txt")
	// 	if err != nil {
	// 		t.Errorf("Read file content failed, %v", err.Error())
	// 		return
	// 	}
	// 	expect := pkg.Table{
	// 		Name: "student",
	// 		Columns: []pkg.Column{
	// 			pkg.Column{
	// 				Name:          "id",
	// 				Type:          "integer",
	// 				NotNull:       true,
	// 				PrimaryKey:    true,
	// 				AutoIncrement: true,
	// 				ForeignKey:    false,
	// 				Unique:        false,
	// 				References:    "",
	// 			},
	// 			pkg.Column{
	// 				Name:          "created",
	// 				Type:          "datetime",
	// 				NotNull:       true,
	// 				PrimaryKey:    false,
	// 				AutoIncrement: false,
	// 				ForeignKey:    false,
	// 				Unique:        false,
	// 				References:    "",
	// 			},
	// 			pkg.Column{
	// 				Name:          "updated",
	// 				Type:          "datetime",
	// 				NotNull:       true,
	// 				PrimaryKey:    false,
	// 				AutoIncrement: false,
	// 				ForeignKey:    false,
	// 				Unique:        false,
	// 				References:    "",
	// 			},
	// 			pkg.Column{
	// 				Name:          "name",
	// 				Type:          "text",
	// 				NotNull:       true,
	// 				PrimaryKey:    false,
	// 				AutoIncrement: false,
	// 				ForeignKey:    false,
	// 				Unique:        false,
	// 				References:    "",
	// 			},
	// 			pkg.Column{
	// 				Name:          "age",
	// 				Type:          "integer",
	// 				NotNull:       true,
	// 				PrimaryKey:    false,
	// 				AutoIncrement: false,
	// 				ForeignKey:    false,
	// 				Unique:        false,
	// 				References:    "",
	// 			},
	// 			pkg.Column{
	// 				Name:          "height",
	// 				Type:          "integer",
	// 				NotNull:       true,
	// 				PrimaryKey:    false,
	// 				AutoIncrement: false,
	// 				ForeignKey:    false,
	// 				Unique:        false,
	// 				References:    "",
	// 			},
	// 			pkg.Column{
	// 				Name:          "birth_at",
	// 				Type:          "datetime",
	// 				NotNull:       true,
	// 				PrimaryKey:    false,
	// 				AutoIncrement: false,
	// 				ForeignKey:    false,
	// 				Unique:        false,
	// 				References:    "",
	// 			},
	// 			pkg.Column{
	// 				Name:          "hobby",
	// 				Type:          "text",
	// 				NotNull:       false,
	// 				PrimaryKey:    false,
	// 				AutoIncrement: false,
	// 				ForeignKey:    false,
	// 				Unique:        false,
	// 				References:    "",
	// 			},
	// 			pkg.Column{
	// 				Name:          "grade_id",
	// 				Type:          "integer",
	// 				NotNull:       true,
	// 				PrimaryKey:    false,
	// 				AutoIncrement: false,
	// 				ForeignKey:    true,
	// 				Unique:        false,
	// 				References:    "grade_class",
	// 			},
	// 			pkg.Column{
	// 				Name:          "friend_id",
	// 				Type:          "integer",
	// 				NotNull:       false,
	// 				PrimaryKey:    false,
	// 				AutoIncrement: false,
	// 				ForeignKey:    true,
	// 				Unique:        false,
	// 				References:    "student",
	// 			},
	// 		},
	// 	}
	// 	content := string(bytes)
	// 	result := pkg.ParseTableCreateStatement("word_in_paragraphs", content)
	// 	fmt.Println("after pkg.ParseTableCreateStatement")
	// 	if result.Name != expect.Name {
	// 		t.Errorf("Test failed, The table name %v is different %v", result.Name, expect.Name)
	// 		return
	// 	}
	// 	if len(expect.Columns) != len(result.Columns) {
	// 		t.Errorf("Test failed, The columns len %v is different %v", len(result.Columns), len(expect.Columns))
	// 		return
	// 	}
	// 	for i, expectedColumn := range expect.Columns {
	// 		fmt.Println("------")
	// 		actuallyColumn := result.Columns[i]
	// 		expectedValue := reflect.ValueOf(expectedColumn)
	// 		actuallyValue := reflect.ValueOf(actuallyColumn)
	// 		typ := actuallyValue.Type()
	// 		for j := 0; j < actuallyValue.NumField(); j++ {
	// 			field := typ.Field(j).Name
	// 			actuallyV := fmt.Sprintf("%v", actuallyValue.Field(j).Interface())
	// 			expectedV := fmt.Sprintf("%v", expectedValue.Field(j).Interface())
	// 			fmt.Println(field, "'"+actuallyV+"'/'"+expectedV+"'", reflect.TypeOf(actuallyV))
	// 			if actuallyV != expectedV {
	// 				t.Errorf("%v | %v '%v' - '%v'", expectedColumn, field, actuallyV, expectedV)
	// 			}
	// 		}
	// 	}
	// })

	t.Run("CREATE Table 1 with line break", func(c *testing.T) {
		bytes, err := os.ReadFile("./test3.txt")
		if err != nil {
			t.Errorf("Read file content failed, %v", err.Error())
			return
		}
		expect := pkg.Table{
			Name: "PlayHistoryV2",
			Columns: []pkg.Column{
				pkg.Column{
					Name:          "id",
					Type:          "text",
					NotNull:       true,
					PrimaryKey:    true,
					AutoIncrement: false,
					ForeignKey:    false,
					Unique:        false,
					References:    "",
				},
				pkg.Column{
					Name:          "created",
					Type:          "datetime",
					NotNull:       true,
					PrimaryKey:    false,
					AutoIncrement: false,
					ForeignKey:    false,
					Unique:        false,
					References:    "",
				},
				pkg.Column{
					Name:          "updated",
					Type:          "datetime",
					NotNull:       true,
					PrimaryKey:    false,
					AutoIncrement: false,
					ForeignKey:    false,
					Unique:        false,
					References:    "",
				},
				pkg.Column{
					Name:          "text",
					Type:          "text",
					NotNull:       true,
					PrimaryKey:    false,
					AutoIncrement: false,
					ForeignKey:    false,
					Unique:        false,
					References:    "",
				},
				pkg.Column{
					Name:          "duration",
					Type:          "real",
					NotNull:       true,
					PrimaryKey:    false,
					AutoIncrement: false,
					ForeignKey:    false,
					Unique:        false,
					References:    "",
				},
				pkg.Column{
					Name:          "current_time",
					Type:          "real",
					NotNull:       true,
					PrimaryKey:    false,
					AutoIncrement: false,
					ForeignKey:    false,
					Unique:        false,
					References:    "",
				},
				pkg.Column{
					Name:          "thumbnail_path",
					Type:          "text",
					NotNull:       false,
					PrimaryKey:    false,
					AutoIncrement: false,
					ForeignKey:    false,
					Unique:        false,
					References:    "",
				},
				pkg.Column{
					Name:          "file_id",
					Type:          "text",
					NotNull:       false,
					PrimaryKey:    false,
					AutoIncrement: false,
					ForeignKey:    false,
					Unique:        false,
					References:    "",
				},
				pkg.Column{
					Name:          "media_id",
					Type:          "text",
					NotNull:       true,
					PrimaryKey:    false,
					AutoIncrement: false,
					ForeignKey:    true,
					Unique:        false,
					References:    "Media",
				},
				pkg.Column{
					Name:          "media_source_id",
					Type:          "text",
					NotNull:       true,
					PrimaryKey:    false,
					AutoIncrement: false,
					ForeignKey:    true,
					Unique:        false,
					References:    "MediaSource",
				},
				pkg.Column{
					Name:          "member_id",
					Type:          "text",
					NotNull:       true,
					PrimaryKey:    false,
					AutoIncrement: false,
					ForeignKey:    true,
					Unique:        false,
					References:    "Member",
				},
			},
		}
		content := string(bytes)
		result := pkg.ParseTableCreateStatement("word_in_paragraphs", content)
		fmt.Println("after pkg.ParseTableCreateStatement")
		if result.Name != expect.Name {
			t.Errorf("Test failed, The table name %v is different %v", result.Name, expect.Name)
			return
		}
		if len(expect.Columns) != len(result.Columns) {
			t.Errorf("Test failed, The columns len %v is different %v", len(result.Columns), len(expect.Columns))
			return
		}
		for i, expectedColumn := range expect.Columns {
			fmt.Println("------")
			actuallyColumn := result.Columns[i]
			expectedValue := reflect.ValueOf(expectedColumn)
			actuallyValue := reflect.ValueOf(actuallyColumn)
			typ := actuallyValue.Type()
			for j := 0; j < actuallyValue.NumField(); j++ {
				field := typ.Field(j).Name
				actuallyV := fmt.Sprintf("%v", actuallyValue.Field(j).Interface())
				expectedV := fmt.Sprintf("%v", expectedValue.Field(j).Interface())
				fmt.Println(field, "'"+actuallyV+"'/'"+expectedV+"'", reflect.TypeOf(actuallyV))
				if actuallyV != expectedV {
					t.Errorf("%v | %v '%v' - '%v'", expectedColumn, field, actuallyV, expectedV)
				}
			}
		}
	})
}
