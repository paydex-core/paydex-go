package paydextoml

import "log"

// ExampleGetTOML gets the paydex.toml file for coins.asia
func ExampleClient_GetPaydexToml() {
	_, err := DefaultClient.GetPaydexToml("coins.asia")
	if err != nil {
		log.Fatal(err)
	}
}
