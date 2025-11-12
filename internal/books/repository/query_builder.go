package repository

import (
	"fmt"
	"strings"

	"github.com/bkiran6398/library/internal/books/domain"
)

// buildListQuery constructs a SQL query and arguments for listing books based on the filter.
func buildListQuery(filter domain.ListFilter) (query string, args []any) {
	whereConditions := buildWhereConditions(filter)
	queryArgs := buildQueryArguments(filter)

	baseQuery := `
SELECT id, title, author, isbn, published_year, copies_total, copies_available, created_at, updated_at
FROM books`

	if len(whereConditions) > 0 {
		baseQuery += " WHERE " + strings.Join(whereConditions, " AND ")
	}

	baseQuery += " ORDER BY created_at DESC"

	if filter.Limit > 0 {
		baseQuery += fmt.Sprintf(" LIMIT %d", filter.Limit)
	}

	if filter.Offset > 0 {
		baseQuery += fmt.Sprintf(" OFFSET %d", filter.Offset)
	}

	return baseQuery, queryArgs
}

// buildWhereConditions builds WHERE clause conditions based on filter criteria.
func buildWhereConditions(filter domain.ListFilter) []string {
	var conditions []string
	placeholderIndex := 1

	if filter.Title != nil && *filter.Title != "" {
		conditions = append(conditions, fmt.Sprintf("title ILIKE $%d", placeholderIndex))
		placeholderIndex++
	}

	if filter.Author != nil && *filter.Author != "" {
		conditions = append(conditions, fmt.Sprintf("author ILIKE $%d", placeholderIndex))
		placeholderIndex++
	}

	if filter.ISBN != nil && *filter.ISBN != "" {
		conditions = append(conditions, fmt.Sprintf("isbn = $%d", placeholderIndex))
		placeholderIndex++
	}

	return conditions
}

// buildQueryArguments builds the argument list for the SQL query based on filter criteria.
func buildQueryArguments(filter domain.ListFilter) []any {
	var queryArguments []any

	if filter.Title != nil && *filter.Title != "" {
		queryArguments = append(queryArguments, "%"+*filter.Title+"%")
	}

	if filter.Author != nil && *filter.Author != "" {
		queryArguments = append(queryArguments, "%"+*filter.Author+"%")
	}

	if filter.ISBN != nil && *filter.ISBN != "" {
		queryArguments = append(queryArguments, *filter.ISBN)
	}

	return queryArguments
}
