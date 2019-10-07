package main

import (
	"github.com/graphql-go/graphql"
	"log"
	"mockrift/pkg/models"
)

func getSchema() graphql.Schema {
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
					Type: graphql.String,
				},
			},
		},
	)

	var requestType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Request",
			Fields: graphql.Fields{
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
					Type: graphql.String,
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

	apps := []*models.App{
		{
			Slug: "blah",
			Name: "Blah",
		},
	}

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
