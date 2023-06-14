package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"

	//"github.com/jinzhu/gorm"
	"github.com/prae014/pokemon/graph"
	"github.com/prae014/pokemon/graph/model"

	//"gorm.io/driver/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const defaultPort = "8080"

//var db *gorm.DB

func initDB() *gorm.DB {
	//dsn := "host=localhost user=postgres password=admin dbname=test port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	dsn := "postgresql://postgres:admin@pokedex_postgres:5432/test"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	//db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&model.Pokemon{}, &model.PokemonType{}, &model.PokemonAbility{})
	poke1 := model.Pokemon{
		Name:        "Pikachu",
		Description: "Pikachu is an Electric type Pokemon introduced in Generation 1",
		Category:    "Mouse Pokemon",
		Abilities:   []*model.PokemonAbility{{Name: "Static"}, {Name: "Lightning Rod"}},
		Type:        []*model.PokemonType{{Name: "Electric"}},
	}

	poke2 := model.Pokemon{
		Name:        "Dedene",
		Description: "Dedene is an Electric/Fairy Pokemon introduced in Generation 6",
		Category:    "Antenna",
		Abilities:   []*model.PokemonAbility{{Name: "Cheek Pouch"}, {Name: "Pickup"}},
		Type:        []*model.PokemonType{{Name: "Fairy"}, {Name: "Electric"}},
	}

	poke3 := model.Pokemon{
		Name:        "Lucario",
		Description: "Lucario is a Fighting/Steel Pokemon introduced in Generation 4",
		Category:    "Aura",
		Abilities:   []*model.PokemonAbility{{Name: "Inner Focus"}, {Name: "Steadfast"}},
		Type:        []*model.PokemonType{{Name: "Fighting"}, {Name: "Steel"}},
	}

	db.Create(&poke1)
	db.Create(&poke2)
	db.Create(&poke3)

	fmt.Println(poke1)
	fmt.Println(poke2)
	fmt.Println(poke3)

	return db
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	db := initDB()

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(
		graph.Config{
			Resolvers: &graph.Resolver{
				DB: db,
			},
		}),
	)

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
