// Command seed populates the production database with the demo dataset
// shipped in ../../seed.sql (relative to this source file). It is
// intended to be invoked once after the first deployment:
//
//	DATABASE_URL=... go run ./cmd/seed
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	databaseURL := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if databaseURL == "" {
		log.Fatalf("DATABASE_URL is required")
	}

	sqlPath := os.Getenv("SEED_SQL_PATH")
	if sqlPath == "" {
		// Look for seed.sql in: $SEED_DIR, current dir, parent dir,
		// then ../../seed.sql (project root, relative to this source).
		candidates := []string{
			"seed.sql",
			"../seed.sql",
			"../../seed.sql",
		}
		if dir := strings.TrimSpace(os.Getenv("SEED_DIR")); dir != "" {
			candidates = []string{filepath.Join(dir, "seed.sql")}
		}
		for _, c := range candidates {
			if _, err := os.Stat(c); err == nil {
				sqlPath = c
				break
			}
		}
		if sqlPath == "" {
			log.Fatalf("could not locate seed.sql — set SEED_SQL_PATH or SEED_DIR")
		}
	}

	rawSQL, err := os.ReadFile(sqlPath)
	if err != nil {
		log.Fatalf("read seed file %q: %v", sqlPath, err)
	}

	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("open postgres: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("get sql handle: %v", err)
	}
	defer func() {
		if cerr := sqlDB.Close(); cerr != nil {
			log.Printf("close database: %v", cerr)
		}
	}()

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("ping postgres: %v", err)
	}

	if _, err := sqlDB.Exec(string(rawSQL)); err != nil {
		// Detect the common "tables don't exist yet" case so the
		// message is actionable instead of a raw SQL error.
		if isMissingRelation(err) {
			log.Fatalf("seed failed: tables do not exist yet — start the backend once so AutoMigrate runs, then re-run this command")
		}
		log.Fatalf("execute seed: %v", err)
	}

	fmt.Println("seed: ok")
}

func isMissingRelation(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "does not exist") &&
		(strings.Contains(msg, "relation") || strings.Contains(msg, "table"))
}
