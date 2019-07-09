package scalar

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	uuid "github.com/satori/go.uuid"
)

type SpecialDate struct {
	time.Time
}

func NewSpecialDate(v time.Time) *SpecialDate {
	return &SpecialDate{v}
}

func (sd *SpecialDate) UnmarshalJSON(input []byte) error {
	strInput := string(input)
	strInput = strings.Trim(strInput, `"`)
	newTime, err := time.Parse("2006-01-02T15:04:05", strInput)
	if err != nil {
		return err
	}

	sd.Time = newTime
	return nil
}

var SpecialDateScalar = graphql.NewScalar(graphql.ScalarConfig{
	Name:        "SpecialDate",
	Description: "The `SpecialDate` scalar type represents an Time Object.",
	// Serialize serializes `SpecialDate` to string.
	Serialize: func(value interface{}) interface{} {
		switch value := value.(type) {
		case time.Time:
			return value.Format(time.RFC3339)
		case SpecialDate:
			return value.Format(time.RFC3339)
		case string:
			return value
		case int:
			t := time.Unix(int64(value), 0)
			return t.Format(time.RFC3339)
		case float64:
			t := time.Unix(int64(value), 0)
			return t.Format(time.RFC3339)
		default:
			return fmt.Errorf("wrong type")
		}
	},
	// ParseValue parses GraphQL variables from `string` to `SpecialDate`.
	ParseValue: func(value interface{}) interface{} {
		switch value := value.(type) {
		case time.Time:
			return value
		case SpecialDate:
			return value
		case string:
			t, err := time.Parse(time.RFC3339, value)
			if err != nil {
				return err
			} else {
				return t
			}
		case int:
			t := time.Unix(int64(value), 0)
			return t
		case float64:
			t := time.Unix(int64(value), 0)
			return t
		default:
			return fmt.Errorf("wrong type")
		}
	},
	// ParseLiteral parses GraphQL AST value to `SpecialDate`.
	ParseLiteral: func(valueAST ast.Value) interface{} {
		switch valueAST := valueAST.(type) {
		case *ast.StringValue:
			t, err := time.Parse(time.RFC3339, valueAST.Value)
			if err != nil {
				return err
			} else {
				return t
			}
		default:
			return fmt.Errorf("wrong type")
		}
	},
})

var NullStringScalar = graphql.NewScalar(graphql.ScalarConfig{
	Name:        "NullString",
	Description: "The `NullString` scalar type represents an String Object.",
	// Serialize serializes `NullString` to string.
	Serialize: func(value interface{}) interface{} {
		switch value := value.(type) {
		case sql.NullString:
			if value.Valid {
				return value.String
			} else {
				return nil
			}
		case string:
			return value
		case uuid.NullUUID:
			if value.Valid {
				return value.UUID.String()
			} else {
				return nil
			}
		default:
			return fmt.Errorf("wrong type")
		}
	},
	// ParseValue parses GraphQL variables from `string` to `NullString`.
	ParseValue: func(value interface{}) interface{} {
		switch value := value.(type) {
		case string:
			return sql.NullString{Valid: true, String: value}
		case sql.NullString:
			return value
		case uuid.NullUUID:
			return value.UUID.String()
		default:
			return fmt.Errorf("wrong type")
		}
	},
	// ParseLiteral parses GraphQL AST value to `NullString`.
	ParseLiteral: func(valueAST ast.Value) interface{} {
		switch valueAST := valueAST.(type) {
		case *ast.StringValue:
			return sql.NullString{Valid: true, String: valueAST.Value}
		default:
			return fmt.Errorf("wrong type")
		}
	},
})
