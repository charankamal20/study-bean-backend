package initializers

import (
	"log"
	// "os"
	// "path/filepath"
	"github.com/joho/godotenv"
)

// func LoadEnvVariables() {
// 	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	environmentPath := filepath.Join(dir, ".env")
// 	err = godotenv.Load(environmentPath)

// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

func LoadEnvVariables() {
    // dir, err := filepath.Abs(filepath.Dir(os.Args[0]))

    // if err != nil {
    //     log.Fatal(err)
    // }
    // fmt.Print(dir)

    // environmentPath := filepath.Join(dir, ".env")
    err := godotenv.Load(".env")

    if err != nil {
        log.Fatal(err)
    }
}