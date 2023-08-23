package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/DavidHernandez21/golang_argocd/internal/driver/argocd"
	"github.com/joho/godotenv"
)

func loadEnvFile(filepath string) error {

	if err := godotenv.Load(filepath); err != nil {
		return err
	}

	return nil
}

func main() {

	err := loadEnvFile(".env")

	if err != nil {
		panic(err)
	}

	const address = "127.0.0.1:8080"

	connection := argocd.Connection{
		Address: address,
		Token:   os.Getenv("ARGOCD_TOKEN"),
	}

	// fmt.Println(os.Getenv("ARGOCD_TOKEN"))

	client, err := argocd.NewClient(&connection)
	if err != nil {
		fmt.Println("error creating client")
		client.CloseConnections()
		panic(err)
	}

	defer client.CloseConnections()

	// createProject, err := client.CreateProject("foo")
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(createProject.UID)

	// err = client.AddDestination(createProject.Name, "server", "namespace", "name")
	// if err != nil {
	// 	panic(err)
	// }

	// getProject, err := client.GetProject("foo")
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(getProject.Namespace)

	// err = client.DeleteProject(getProject.Name)
	// if err != nil {
	// 	panic(err)
	// }
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	const applicationName = "guestbook"
	getApplication, err := client.GetApplication(ctx, applicationName)
	if err != nil {
		fmt.Println("error getting application")
		panic(err)
	}

	fmt.Printf("%#v\n", getApplication.Spec.Destination)
	fmt.Printf("%#v\n", getApplication.Spec.Source)

	clusters, err := client.GetClusters(ctx)
	if err != nil {
		fmt.Println("error getting clusters")
		client.CloseConnections()
		panic(err)
	}

	for _, cluster := range clusters {
		fmt.Println(cluster.Name)
	}
}
