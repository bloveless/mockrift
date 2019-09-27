package main

import (
	"github.com/graphql-go/graphql"
	"log"
	"mockrift/pkg/models"
)

func getSchema() graphql.Schema {
	var appType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "App",
			Fields: graphql.Fields{
				"name": &graphql.Field{
					Type: graphql.String,
				},
				"alias": &graphql.Field{
					Type: graphql.String,
				},
			},
		},
	)

	apps := []*models.App{
		{
			Name:  "Blah",
			Alias: "blah",
		},
	}

	fields := graphql.Fields{
		"app": &graphql.Field{
			Type:        appType,
			Description: "",
			Args: graphql.FieldConfigArgument{
				"alias": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				for _, app := range apps {
					if app.Alias == params.Args["alias"].(string) {
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
