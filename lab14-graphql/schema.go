package main

import (
	"github.com/graphql-go/graphql"
)

var productType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Product",
	Fields: graphql.Fields{
		"id":       &graphql.Field{Type: graphql.Int},
		"name":     &graphql.Field{Type: graphql.String},
		"price":    &graphql.Field{Type: graphql.Float},
		"category": &graphql.Field{Type: graphql.String},
	},
})

var categoryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Category",
	Fields: graphql.Fields{
		"name":  &graphql.Field{Type: graphql.String},
		"count": &graphql.Field{Type: graphql.Int},
	},
})

func buildQueryType() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"products": &graphql.Field{
				Type:        graphql.NewList(productType),
				Description: "List all products",
				Args: graphql.FieldConfigArgument{
					"category": &graphql.ArgumentConfig{
						Type:        graphql.String,
						Description: "Filter by category",
					},
					"minPrice": &graphql.ArgumentConfig{
						Type:        graphql.Float,
						Description: "Minimum price filter",
					},
					"maxPrice": &graphql.ArgumentConfig{
						Type:        graphql.Float,
						Description: "Maximum price filter",
					},
				},
				Resolve: resolveProducts,
			},
			"product": &graphql.Field{
				Type:        productType,
				Description: "Get a single product by ID",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: resolveProduct,
			},
			"categories": &graphql.Field{
				Type:        graphql.NewList(categoryType),
				Description: "List categories with product counts",
				Resolve:     resolveCategories,
			},
		},
	})
}

func buildMutationType() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"createProduct": &graphql.Field{
				Type:        productType,
				Description: "Create a new product",
				Args: graphql.FieldConfigArgument{
					"name": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"price": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Float),
					},
					"category": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: resolveCreateProduct,
			},
			"updateProduct": &graphql.Field{
				Type:        productType,
				Description: "Update an existing product",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
					"name":     &graphql.ArgumentConfig{Type: graphql.String},
					"price":    &graphql.ArgumentConfig{Type: graphql.Float},
					"category": &graphql.ArgumentConfig{Type: graphql.String},
				},
				Resolve: resolveUpdateProduct,
			},
			"deleteProduct": &graphql.Field{
				Type:        graphql.Boolean,
				Description: "Delete a product by ID",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: resolveDeleteProduct,
			},
		},
	})
}
