package params_query_string

import "fmt"

func Generate(table string, params map[string]any) (string, []any) {
	query := fmt.Sprintf(`SELECT * FROM %s WHERE `, table)

	values := make([]any, 0, len(params))

	it := 1
	for key, value := range params {
		if it == 1 {
			query += fmt.Sprintf(`%s = $%d `, key, it)
		} else {
			query += fmt.Sprintf(`AND %s = $%d `, key, it)
		}
		values = append(values, value)
		it++
	}

	return query, values
}
