package main

import (
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"log"
)

func makeString(value interface{}) interface{} {
	if v, ok := value.(*[]byte); ok {
		if v == nil {
			return nil
		}
		return *v
	}
	return value
}

var Base64Type = graphql.NewScalar(graphql.ScalarConfig{
	Name:        "Base64Type",
	Description: "Base64Type representation of a string.",
	Serialize:   makeString,
	ParseValue:  makeString,
	ParseLiteral: func(valueAST ast.Value) interface{} {
		switch valueAST := valueAST.(type) {
		case *ast.StringValue:
			return valueAST.Value
		}
		return nil
	},
})

func (s *server) getSchema() graphql.Schema {
	apps := s.apps.GetAll()

	var headerType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Header",
			Fields: graphql.Fields{
				"name": &graphql.Field{
					Type: graphql.String,
				},
				"value": &graphql.Field{
					Type: graphql.NewList(graphql.String),
				},
			},
		},
	)

	var responseType = graphql.NewObject(
		graphql.ObjectConfig{
			Name:        "Response",
			Description: "",
			Fields: graphql.Fields{
				"id": &graphql.Field{
					Type: graphql.ID,
				},
				"active": &graphql.Field{
					Type: graphql.Boolean,
				},
				"status_code": &graphql.Field{
					Type: graphql.Int,
				},
				"header": &graphql.Field{
					Type: graphql.NewList(headerType),
				},
				"body": &graphql.Field{
					Type: Base64Type,
				},
			},
		},
	)

	var requestType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Request",
			Fields: graphql.Fields{
				"id": &graphql.Field{
					Type: graphql.ID,
				},
				"method": &graphql.Field{
					Type: graphql.String,
				},
				"url": &graphql.Field{
					Type: graphql.String,
				},
				"header": &graphql.Field{
					Type: graphql.NewList(headerType),
				},
				"body": &graphql.Field{
					Type: Base64Type,
				},
				"responses": &graphql.Field{
					Type: graphql.NewList(responseType),
				},
			},
		},
	)

	var appType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "App",
			Fields: graphql.Fields{
				"slug": &graphql.Field{
					Type: graphql.String,
				},
				"name": &graphql.Field{
					Type: graphql.String,
				},
				"requests": &graphql.Field{
					Type: graphql.NewList(requestType),
				},
			},
		},
	)

	fields := graphql.Fields{
		"app": &graphql.Field{
			Type:        appType,
			Description: "",
			Args: graphql.FieldConfigArgument{
				"slug": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				for _, app := range apps {
					if app.Slug == params.Args["slug"].(string) {
						return app, nil
					}
				}

				return nil, nil
			},
		},
		"apps": &graphql.Field{
			Type:        graphql.NewList(appType),
			Description: "",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return apps, nil
			},
		},
	}

	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	return schema
}
