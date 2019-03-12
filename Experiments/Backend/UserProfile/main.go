package main

// I'm trying to combine the features from GORM and graphql-go here

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"net/http"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/graphql-go/graphql"
)

type userProfile struct {
	gorm.Model
	Name string `json:"name"`
	Contact string `json:"contact"`
}

var db *gorm.DB

/*
   Create User object type with fields "id" and "name" by using GraphQLObjectTypeConfig:
       - Name: name of object type
       - Fields: a map of fields by using GraphQLFields
   Setup type of field use GraphQLFieldConfig
*/
var userProfileType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"contact": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

/*
   Create Query object type with fields "user" has type [userType] by using GraphQLObjectTypeConfig:
       - Name: name of object type
       - Fields: a map of fields by using GraphQLFields
   Setup type of field use GraphQLFieldConfig to define:
       - Type: type of field
       - Args: arguments to query with current field
       - Resolve: function to query data using params from [Args] and return value with current type
*/
var userProfileQuery = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"user": &graphql.Field{
				Type: userProfileType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					var profile userProfile
					db.First(&profile, p.Args["id"].(int))
					return profile, nil
				},
			},
		},
	})

var userProfileSchema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query: userProfileQuery,
	},
)

func executeUserProfileQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("wrong result, unexpected errors: %v", result.Errors)
	}
	return result
}

func main() {
	db, err := gorm.Open("sqlite3", "sqlite/user_profile.db")
	if err != nil {
		panic("failed to connect database" + err.Error())
	}
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&userProfile{})

	db.Create(&userProfile{Name: "billy flex", Contact: "big@flex.com"})

	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		result := executeUserProfileQuery(r.URL.Query().Get("query"), userProfileSchema)
		json.NewEncoder(w).Encode(result)
	})

	fmt.Println("Now server is running on port 8080")
	fmt.Println("Test with Get      : curl -g 'http://localhost:8080/graphql?query={user(id:\"1\"){name}}'")
	http.ListenAndServe(":8080", nil)
}