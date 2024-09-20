package graphql

import (
	"write-stream-go/internal/models"

	"github.com/graphql-go/graphql"
)

var userType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id":        &graphql.Field{Type: graphql.String},
			"SocialID":  &graphql.Field{Type: graphql.String},
			"email":     &graphql.Field{Type: graphql.String},
			"name":      &graphql.Field{Type: graphql.String},
			"avatarUrl": &graphql.Field{Type: graphql.String},
			"createdAt": &graphql.Field{Type: graphql.DateTime},
			"updatedAt": &graphql.Field{Type: graphql.DateTime},
		},
	},
)

func NewSchema() (graphql.Schema, error) {
	return graphql.NewSchema(
		graphql.SchemaConfig{
			Query: graphql.NewObject(
				graphql.ObjectConfig{
					Name: "Query",
					Fields: graphql.Fields{
						"me": &graphql.Field{
							Type: userType,
							Resolve: func(p graphql.ResolveParams) (interface{}, error) {
								if user, ok := p.Context.Value("user").(*models.User); ok {
									return user, nil
								}
								return nil, nil
							},
						},
					},
				},
			),
		},
	)
}
