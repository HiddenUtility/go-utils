package main

import (
	"fmt"
	"github.com/username/go-utils/pkg/jsonutils"
	"github.com/username/go-utils/pkg/jsonutils/models"
	"log"
)

func main() {
	user := models.User{
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   30,
		Address: models.Address{
			Street: "123 Main St",
			City:   "New York",
		},
		Settings: models.Settings{
			NotificationsEnabled: true,
			Theme:                "dark",
		},
	}

	// Serialize to JSON
	jsonStr, err := jsonutils.ToJSON(user)
	if err != nil {
		log.Fatalf("failed to serialize: %v", err)
	}
	fmt.Println("Serialized JSON:")
	fmt.Println(jsonStr)

	// Deserialize from JSON
	decodedUser, err := jsonutils.FromJSON[models.User](jsonStr)
	if err != nil {
		log.Fatalf("failed to deserialize: %v", err)
	}
	fmt.Printf("\nDeserialized User: %+v\n", decodedUser)

	// Pretty Serialize
	prettyJSON, err := jsonutils.ToJSONPretty(user)
	if err != nil {
		log.Fatalf("failed to serialize pretty: %v", err)
	}
	fmt.Println("\nPretty JSON:")
	fmt.Println(prettyJSON)
}
