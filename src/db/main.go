package main

import (
	"bytes"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"go/format"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	connectionString = "host=localhost user=gebetamap password=gebetamap dbname=openstreetmap port=5432 sslmode=disable"
	schema           = "public"
	path             = "./src/db/generated/gorm"
)

var modelTemplate = template.Must(template.New("").Parse(`
package models 
{{.Imports}}
type {{.StructName}} struct {
    {{- range .Fields }}
       {{ .FieldName }} {{ .FieldType }} ` + "`gorm:\"{{if .IsPrimaryKey}}primaryKey;{{end}}column:{{ .ColumnName }}\" json:\"{{ .ColumnName }}\"`" + `
    {{- end }}

    {{- /* Belongs To relationships */}}   
    {{- range .Relationships}}
       {{.FieldName}} {{.FieldType}} ` + "`{{.Tag}}`" + `
    {{- end}}

    {{- /* Has Many relationships */}}
    {{- range .HasMany}}
       {{.FieldName}} {{.FieldType}} ` + "`{{.Tag}}`" + `
    {{- end}}
}

// TableName sets the insert table name for this struct type
func (m *{{.StructName}}) TableName() string {
    return "{{.TableName}}"
}
`))

type FieldData struct {
	FieldName    string
	FieldType    string
	ColumnName   string
	IsPrimaryKey bool
}
type TemplateData struct {
	StructName    string
	TableName     string
	Fields        []FieldData
	Relationships []RelationshipData
	HasMany       []HasManyData
	Imports       template.HTML
}

type RelationshipData struct {
	FieldName string
	FieldType string
	Tag       template.HTML
}

type HasManyData struct {
	FieldName string
	FieldType string
	Tag       template.HTML
}

type ColumnInfo struct {
	ColumnName string
	DataType   string
	IsNullable string
}
type ForeignKeyInfo struct {
	TableName         string
	ColumnName        string
	ForeignTableName  string
	ForeignColumnName string
}

func getForeignKeys(db *sql.DB) (map[string][]ForeignKeyInfo, error) {
	query := `
SELECT
    cl1.relname AS table_name,
    att2.attname AS column_name,
    cl2.relname AS foreign_table_name,
    att1.attname AS foreign_column_name
FROM
    pg_constraint AS con
    JOIN pg_class AS cl1 ON con.conrelid = cl1.oid
    JOIN pg_class AS cl2 ON con.confrelid = cl2.oid
    JOIN pg_attribute AS att1 ON att1.attrelid = con.confrelid AND att1.attnum = con.confkey[1]
    JOIN pg_attribute AS att2 ON att2.attrelid = con.conrelid AND att2.attnum = con.conkey[1]
WHERE
    con.contype = 'f'`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	fkMap := make(map[string][]ForeignKeyInfo)
	for rows.Next() {
		fkInfo := ForeignKeyInfo{}

		err := rows.Scan(&fkInfo.TableName, &fkInfo.ColumnName, &fkInfo.ForeignTableName, &fkInfo.ForeignColumnName)
		if err != nil {
			return nil, err
		}
		fkMap[fkInfo.TableName] = append(fkMap[fkInfo.TableName], fkInfo)
	}
	return fkMap, nil
}

func getColumns(db *sql.DB) (map[string][]ColumnInfo, error) {
	query := `SELECT table_name , column_name , data_type , is_nullable FROM information_schema.columns WHERE table_schema = $1;`

	results, err := db.Query(query, schema)

	if err != nil {
		return nil, err
	}
	defer results.Close()
	tables := make(map[string][]ColumnInfo)
	for results.Next() {
		columnInfo := ColumnInfo{}
		var tableName string
		err := results.Scan(&tableName, &columnInfo.ColumnName, &columnInfo.DataType, &columnInfo.IsNullable)
		if err != nil {
			return nil, err
		}
		tables[tableName] = append(tables[tableName], columnInfo)
	}
	return tables, nil
}

func getPrimaryKeys(db *sql.DB) (map[string]string, error) {
	query := `
SELECT
    tc.table_name, 
    kcu.column_name
FROM 
    information_schema.table_constraints AS tc 
    JOIN information_schema.key_column_usage AS kcu
      ON tc.constraint_name = kcu.constraint_name
      AND tc.table_schema = kcu.table_schema
WHERE 
    tc.constraint_type = 'PRIMARY KEY' 
    AND tc.table_schema = $1;
`
	rows, err := db.Query(query, schema)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pkMap := make(map[string]string)
	for rows.Next() {
		var tableName, columnName string
		if err := rows.Scan(&tableName, &columnName); err != nil {
			return nil, err
		}
		pkMap[tableName] = columnName
	}
	return pkMap, nil
}

func main() {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	tables, err := getColumns(db)
	if err != nil {
		log.Fatal(err)
	}
	foreignKeys, err := getForeignKeys(db)
	if err != nil {
		log.Fatal(err)
	}
	primaryKeys, err := getPrimaryKeys(db)
	if err != nil {
		log.Fatal(err)
	}

	inverseRelationships := make(map[string][]ForeignKeyInfo)
	for _, fks := range foreignKeys {
		for _, fk := range fks {
			inverseRelationships[fk.ForeignTableName] = append(inverseRelationships[fk.ForeignTableName], fk)
		}
	}
	if err := os.MkdirAll(path, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}
	for tableName, columns := range tables {
		data := TemplateData{}
		data.TableName = tableName
		data.StructName = toPascalCase(tableName)
		shouldIncludeTime := false
		pkColumn := primaryKeys[tableName]

		for _, column := range columns {

			goType := mapPsqlTypeToGo(column.DataType, column.IsNullable)
			if goType == "time.Time" {
				shouldIncludeTime = true
			}
			data.Fields = append(data.Fields, FieldData{
				FieldName:    toPascalCase(column.ColumnName),
				FieldType:    goType,
				ColumnName:   column.ColumnName,
				IsPrimaryKey: column.ColumnName == pkColumn,
			})
		}

		// Populate "belongs to" relationships
		if fks, ok := foreignKeys[tableName]; ok {
			for _, fk := range fks {
				baseName := strings.TrimSuffix(fk.ColumnName, "_id")
				baseName = strings.TrimSuffix(baseName, "_code")
				baseName = strings.TrimSuffix(baseName, "_by")
				relFieldName := toPascalCase(baseName)

				relFieldType := toPascalCase(fk.ForeignTableName)
				gormTag := fmt.Sprintf(`gorm:"foreignKey:%s;references:%s"`, toPascalCase(fk.ColumnName), toPascalCase(fk.ForeignColumnName))

				data.Relationships = append(data.Relationships, RelationshipData{
					FieldName: relFieldName,
					FieldType: relFieldType,
					Tag:       template.HTML(gormTag),
				})
			}
		}

		if invFks, ok := inverseRelationships[tableName]; ok {
			groupedByTable := make(map[string][]ForeignKeyInfo)
			for _, fk := range invFks {
				groupedByTable[fk.TableName] = append(groupedByTable[fk.TableName], fk)
			}

			var sortedTableNames []string
			for name := range groupedByTable {
				sortedTableNames = append(sortedTableNames, name)
			}
			sort.Strings(sortedTableNames)

			for _, referencingTable := range sortedTableNames {
				fkGroup := groupedByTable[referencingTable]
				if len(fkGroup) == 1 {
					fk := fkGroup[0]
					relFieldName := toPascalCase(fk.TableName)
					relFieldType := "[]" + toPascalCase(fk.TableName)
					gormTag := fmt.Sprintf(`gorm:"foreignKey:%v"`, toPascalCase(fk.ColumnName))
					data.HasMany = append(data.HasMany, HasManyData{
						FieldName: relFieldName,
						FieldType: relFieldType,
						Tag:       template.HTML(gormTag),
					})
				} else {
					for _, fk := range fkGroup {
						baseName := strings.TrimSuffix(fk.ColumnName, "_id")
						relFieldName := toPascalCase(fk.TableName) + "As" + toPascalCase(baseName)
						relFieldType := "[]" + toPascalCase(fk.TableName)
						gormTag := fmt.Sprintf(`gorm:"foreignKey:%v"`, toPascalCase(fk.ColumnName))
						data.HasMany = append(data.HasMany, HasManyData{
							FieldName: relFieldName,
							FieldType: relFieldType,
							Tag:       template.HTML(gormTag),
						})
					}
				}
			}
		}

		if shouldIncludeTime {
			data.Imports = template.HTML(`import "time"`)
		}
		var buf bytes.Buffer
		err := modelTemplate.Execute(&buf, data)
		if err != nil {
			log.Fatal(err)
		}
		formattedCode, err := format.Source(buf.Bytes())
		if err != nil {
			log.Fatalf("Failed to format code for table %s: %v\nGenerated Code:\n%s", tableName, err, buf.String())
		}

		filePath := filepath.Join(path, tableName+".go")
		if err := os.WriteFile(filePath, formattedCode, 0644); err != nil {
			log.Fatalf("Failed to write file for table %s: %v", tableName, err)
		}
		fmt.Printf("Generated model for table: %s -> %s.go\n", tableName, tableName)
	}
}

func toPascalCase(s string) string {
	if strings.ToLower(s) == "id" {
		return "ID"
	}
	s = strings.ReplaceAll(s, "_", " ")
	s = strings.ReplaceAll(s, "-", " ")
	s = strings.Title(s)
	return strings.ReplaceAll(s, " ", "")
}

func mapPsqlTypeToGo(psqlType, isNullable string) string {
	nullable := isNullable == "YES"

	switch psqlType {
	case "bigint":
		if nullable {
			return "*int64"
		}
		return "int64"
	case "integer", "smallint":
		if nullable {
			return "*int"
		}
		return "int"
	case "character varying", "text", "character", "uuid":
		if nullable {
			return "*string"
		}
		return "string"
	case "boolean":
		if nullable {
			return "*bool"
		}
		return "bool"
	case "timestamp without time zone", "timestamp with time zone", "date", "time without time zone":
		if nullable {
			return "*time.Time"
		}
		return "time.Time"
	case "numeric", "decimal", "double precision":
		if nullable {
			return "*float64"
		}
		return "float64"
	case "bytea":
		return "[]byte"
	case "json", "jsonb":
		if nullable {
			return "*string"
		}
		return "string"
	default:
		if nullable {
			return "*string"
		}
		return "string"
	}
}
