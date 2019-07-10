package scalar

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/lib/pq"

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

var NullScalar = graphql.NewScalar(graphql.ScalarConfig{
	Name:        "NullScalar",
	Description: "The `NullScalar` scalar type converts null to nil.",
	// Serialize serializes `NullScalar` to its corresponding type from pq.
	Serialize: func(value interface{}) interface{} {
		switch value := value.(type) {
		case sql.NullString:
			if value.Valid {
				return value.String
			} else {
				return nil
			}
		case sql.NullInt64:
			if value.Valid {
				return value.Int64
			} else {
				return nil
			}
		case sql.NullBool:
			if value.Valid {
				return value.Bool
			} else {
				return nil
			}
		case sql.NullFloat64:
			if value.Valid {
				return value.Float64
			} else {
				return nil
			}
		case pq.NullTime:
			if value.Valid {
				return value.Time
			} else {
				return nil
			}
		case uuid.NullUUID:
			if value.Valid {
				return value.UUID.String()
			} else {
				return nil
			}
		case string:
			return value
		case float64:
			return value
		case bool:
			return value
		case uuid.UUID:
			return value
		case time.Time:
			return value
		default:
			return fmt.Errorf("wrong type")
		}
	},
	// ParseValue parses GraphQL variables from `string` to `NullScalar`.
	ParseValue: func(value interface{}) interface{} {
		switch value := value.(type) {
		case string:
			return sql.NullString{Valid: true, String: value}
		case bool:
			return sql.NullBool{Valid: true, Bool: value}
		case float64:
			return sql.NullFloat64{Valid: true, Float64: value}
		case time.Time:
			return pq.NullTime{Valid: true, Time: value}
		case sql.NullString:
			return value
		case sql.NullBool:
			return value
		case sql.NullFloat64:
			return value
		case pq.NullTime:
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
		case *ast.BooleanValue:
			return sql.NullBool{Valid: true, Bool: valueAST.Value}
		case *ast.FloatValue:
			i, err := strconv.ParseFloat(valueAST.Value, 64)
			if err != nil {
				return sql.NullFloat64{Valid: false, Float64: i}
			}
			return sql.NullFloat64{Valid: true, Float64: i}
		default:
			return fmt.Errorf("wrong type")
		}
	},
})
