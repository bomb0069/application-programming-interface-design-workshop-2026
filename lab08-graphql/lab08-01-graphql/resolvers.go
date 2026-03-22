package main

import (
	"database/sql"
	"fmt"

	"github.com/graphql-go/graphql"
)

func resolveProducts(p graphql.ResolveParams) (interface{}, error) {
	query := "SELECT id, name, price, category FROM products WHERE 1=1"
	args := []interface{}{}
	idx := 1

	if category, ok := p.Args["category"].(string); ok {
		query += fmt.Sprintf(" AND category = $%d", idx)
		args = append(args, category)
		idx++
	}
	if minPrice, ok := p.Args["minPrice"].(float64); ok {
		query += fmt.Sprintf(" AND price >= $%d", idx)
		args = append(args, minPrice)
		idx++
	}
	if maxPrice, ok := p.Args["maxPrice"].(float64); ok {
		query += fmt.Sprintf(" AND price <= $%d", idx)
		args = append(args, maxPrice)
		idx++
	}

	query += " ORDER BY id"

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []Product{}
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Category); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func resolveProduct(p graphql.ResolveParams) (interface{}, error) {
	id := p.Args["id"].(int)
	var product Product
	err := db.QueryRow("SELECT id, name, price, category FROM products WHERE id = $1", id).
		Scan(&product.ID, &product.Name, &product.Price, &product.Category)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return product, nil
}

func resolveCategories(p graphql.ResolveParams) (interface{}, error) {
	rows, err := db.Query("SELECT category, COUNT(*) FROM products GROUP BY category ORDER BY category")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := []Category{}
	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.Name, &c.Count); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func resolveCreateProduct(p graphql.ResolveParams) (interface{}, error) {
	name := p.Args["name"].(string)
	price := p.Args["price"].(float64)
	category := p.Args["category"].(string)

	var product Product
	err := db.QueryRow(
		"INSERT INTO products (name, price, category) VALUES ($1, $2, $3) RETURNING id, name, price, category",
		name, price, category,
	).Scan(&product.ID, &product.Name, &product.Price, &product.Category)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func resolveUpdateProduct(p graphql.ResolveParams) (interface{}, error) {
	id := p.Args["id"].(int)

	var product Product
	err := db.QueryRow("SELECT id, name, price, category FROM products WHERE id = $1", id).
		Scan(&product.ID, &product.Name, &product.Price, &product.Category)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("product not found")
	}
	if err != nil {
		return nil, err
	}

	if name, ok := p.Args["name"].(string); ok {
		product.Name = name
	}
	if price, ok := p.Args["price"].(float64); ok {
		product.Price = price
	}
	if category, ok := p.Args["category"].(string); ok {
		product.Category = category
	}

	err = db.QueryRow(
		"UPDATE products SET name=$1, price=$2, category=$3 WHERE id=$4 RETURNING id, name, price, category",
		product.Name, product.Price, product.Category, id,
	).Scan(&product.ID, &product.Name, &product.Price, &product.Category)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func resolveDeleteProduct(p graphql.ResolveParams) (interface{}, error) {
	id := p.Args["id"].(int)
	result, err := db.Exec("DELETE FROM products WHERE id = $1", id)
	if err != nil {
		return false, err
	}
	rows, _ := result.RowsAffected()
	return rows > 0, nil
}
