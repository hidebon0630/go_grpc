package main

import (
	"go_grpc/article/client"
	"go_grpc/graph"
	"go_grpc/graph/generated"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// articleClientを生成
	articleClient, err := client.NewClient("localhost:50051")
	if err != nil {
		articleClient.Close()
		log.Fatalf("Failed to create article client: %v\n", err)
	}

	// GraphQLサーバーに先程のResolverを実装
	srv := handler.NewDefaultServer(
		generated.NewExecutableSchema(
			generated.Config{
				Resolvers: &graph.Resolver{
					ArticleClient: articleClient,
				}}))

	// GraphQL playgroundのエンドポイント
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))

	// 実装したクエリが実行可能なGraphQLサーバーのエンドポイント
	http.Handle("/query", srv)

	// GraphQLサーバーを起動
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
