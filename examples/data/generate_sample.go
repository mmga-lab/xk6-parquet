// +build ignore

package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/parquet-go/parquet-go"
)

// User represents a sample user record
type User struct {
	ID           int64  `parquet:"id"`
	Username     string `parquet:"username"`
	Email        string `parquet:"email"`
	Name         string `parquet:"name"`
	Age          int32  `parquet:"age"`
	Subscription string `parquet:"subscription"`
	Active       bool   `parquet:"active"`
	CreatedAt    string `parquet:"created_at"`
	Balance      float64 `parquet:"balance"`
	Country      string  `parquet:"country"`
}

func generateUsers(count int) []User {
	rand.Seed(42)
	users := make([]User, count)

	subscriptions := []string{"free", "premium", "enterprise"}
	countries := []string{"US", "UK", "DE", "FR", "JP", "CN", "IN", "BR"}

	for i := 0; i < count; i++ {
		users[i] = User{
			ID:           int64(i + 1),
			Username:     fmt.Sprintf("user%d", i+1),
			Email:        fmt.Sprintf("user%d@example.com", i+1),
			Name:         fmt.Sprintf("User %d", i+1),
			Age:          int32(18 + rand.Intn(62)),
			Subscription: subscriptions[rand.Intn(len(subscriptions))],
			Active:       rand.Float32() > 0.3,
			CreatedAt:    time.Now().AddDate(0, 0, -rand.Intn(365)).Format(time.RFC3339),
			Balance:      rand.Float64() * 10000,
			Country:      countries[rand.Intn(len(countries))],
		}
	}

	return users
}

func writeParquetFile(filename string, users []User) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := parquet.NewGenericWriter[User](file, parquet.Compression(&parquet.Snappy))
	defer writer.Close()

	_, err = writer.Write(users)
	if err != nil {
		return fmt.Errorf("failed to write parquet: %w", err)
	}

	return nil
}

func main() {
	fmt.Println("Generating sample Parquet files...")

	files := []struct {
		name  string
		count int
	}{
		{"sample.parquet", 1000},
		{"medium.parquet", 10000},
		{"large.parquet", 100000},
	}

	for _, f := range files {
		fmt.Printf("Creating %s (%d records)...\n", f.name, f.count)
		users := generateUsers(f.count)

		if err := writeParquetFile(f.name, users); err != nil {
			fmt.Printf("Error creating %s: %v\n", f.name, err)
			continue
		}

		stat, _ := os.Stat(f.name)
		sizeMB := float64(stat.Size()) / (1024 * 1024)
		fmt.Printf("✓ Created %s: %d records, %.2f MB\n", f.name, f.count, sizeMB)
	}

	fmt.Println("\n✓ All sample files generated successfully!")
	fmt.Println("\nTo generate the files, run:")
	fmt.Println("  go run generate_sample.go")
}
