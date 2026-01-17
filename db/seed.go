package main

import (
	"fmt"

	"github.com/crazydw4rf/oil-bank-backend/internal/services/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
)

const defaultPassword = "password123"

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed the database with sample data",
	Long:  `Populate the database with sample data including users, addresses, and transactions.`,
	RunE:  runSeed,
}

func runSeed(cmd *cobra.Command, args []string) error {
	cfg, err := config.InitConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Create database connection
	db, err := sqlx.Connect("postgres", cfg.DATABASE_URL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	fmt.Println("Connected to database successfully!")
	fmt.Printf("Default password for all users: %s\n", defaultPassword)
	fmt.Println("Starting seeding process...")

	// Seed data
	if err := seedData(db); err != nil {
		return fmt.Errorf("failed to seed data: %w", err)
	}

	fmt.Println("Seeding completed successfully!")
	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func seedData(db *sqlx.DB) error {
	// Start transaction
	tx, err := db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Seed Addresses
	fmt.Println("Seeding addresses...")
	addressIDs, err := seedAddresses(tx)
	if err != nil {
		return fmt.Errorf("failed to seed addresses: %w", err)
	}

	// Seed Users (Sellers, Collectors, Companies)
	fmt.Println("Seeding users...")
	if err := seedUsers(tx, addressIDs); err != nil {
		return fmt.Errorf("failed to seed users: %w", err)
	}

	// Get specific user IDs by type
	sellerIDs, collectorIDs, companyIDs, err := getUserIDsByType(tx)
	if err != nil {
		return fmt.Errorf("failed to get user IDs by type: %w", err)
	}

	// Seed Sell Transactions
	fmt.Println("Seeding sell transactions...")
	if err := seedSellTransactions(tx, sellerIDs, collectorIDs); err != nil {
		return fmt.Errorf("failed to seed sell transactions: %w", err)
	}

	// Seed Distribute Transactions
	fmt.Println("Seeding distribute transactions...")
	if err := seedDistributeTransactions(tx, collectorIDs, companyIDs); err != nil {
		return fmt.Errorf("failed to seed distribute transactions: %w", err)
	}

	// Seed Oil inventory (automatically created by triggers, just verify)
	fmt.Println("Oil inventory automatically created by triggers...")
	if err := verifyOilInventory(tx, collectorIDs); err != nil {
		return fmt.Errorf("failed to verify oil inventory: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func seedAddresses(tx *sqlx.Tx) ([]int64, error) {
	addresses := []map[string]any{
		{"street_address": "Jl. Merdeka No. 123", "city": "Jakarta Pusat", "regency": "Jakarta Pusat", "province": "DKI Jakarta"},
		{"street_address": "Jl. Sudirman No. 456", "city": "Jakarta Selatan", "regency": "Jakarta Selatan", "province": "DKI Jakarta"},
		{"street_address": "Jl. Gatot Subroto No. 789", "city": "Bandung", "regency": "Bandung", "province": "Jawa Barat"},
		{"street_address": "Jl. Diponegoro No. 321", "city": "Semarang", "regency": "Semarang", "province": "Jawa Tengah"},
		{"street_address": "Jl. Ahmad Yani No. 654", "city": "Surabaya", "regency": "Surabaya", "province": "Jawa Timur"},
		{"street_address": "Jl. Pemuda No. 111", "city": "Medan", "regency": "Medan", "province": "Sumatera Utara"},
		{"street_address": "Jl. Asia Afrika No. 222", "city": "Bandung", "regency": "Bandung", "province": "Jawa Barat"},
		{"street_address": "Jl. Thamrin No. 333", "city": "Jakarta Pusat", "regency": "Jakarta Pusat", "province": "DKI Jakarta"},
		{"street_address": "Jl. Malioboro No. 444", "city": "Yogyakarta", "regency": "Yogyakarta", "province": "DI Yogyakarta"},
		{"street_address": "Jl. Veteran No. 555", "city": "Surabaya", "regency": "Surabaya", "province": "Jawa Timur"},
	}

	query := `
		INSERT INTO "Address" (street_address, city, regency, province)
		VALUES (:street_address, :city, :regency, :province)
	`

	_, err := tx.NamedExec(query, addresses)
	if err != nil {
		return nil, err
	}

	// Retrieve inserted IDs
	rows, err := tx.Queryx(`SELECT id FROM "Address" ORDER BY created_at`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var addressIDs []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		addressIDs = append(addressIDs, id)
	}

	return addressIDs, nil
}

func seedUsers(tx *sqlx.Tx, addressIDs []int64) error {
	// Hash the default password once
	hashedPassword, err := hashPassword(defaultPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	users := []map[string]any{
		{"address_id": addressIDs[0], "username": "john", "email": "john.seller@example.com", "password_hash": hashedPassword, "user_type": "SELLER"},
		{"address_id": addressIDs[1], "username": "jane", "email": "jane.seller@example.com", "password_hash": hashedPassword, "user_type": "SELLER"},
		{"address_id": addressIDs[2], "username": "bob", "email": "bob.seller@example.com", "password_hash": hashedPassword, "user_type": "SELLER"},
		{"address_id": addressIDs[3], "username": "alice", "email": "alice.collector@example.com", "password_hash": hashedPassword, "user_type": "COLLECTOR"},
		{"address_id": addressIDs[4], "username": "charlie", "email": "charlie.collector@example.com", "password_hash": hashedPassword, "user_type": "COLLECTOR"},
		{"address_id": addressIDs[5], "username": "david", "email": "david.collector@example.com", "password_hash": hashedPassword, "user_type": "COLLECTOR"},
		{"address_id": addressIDs[6], "username": "techcorp", "email": "contact@techcorp.com", "password_hash": hashedPassword, "user_type": "COMPANY"},
		{"address_id": addressIDs[7], "username": "greenoil", "email": "info@greenoil.com", "password_hash": hashedPassword, "user_type": "COMPANY"},
		{"address_id": addressIDs[8], "username": "emma", "email": "emma.seller@example.com", "password_hash": hashedPassword, "user_type": "SELLER"},
		{"address_id": addressIDs[9], "username": "frank", "email": "frank.collector@example.com", "password_hash": hashedPassword, "user_type": "COLLECTOR"},
	}

	query := `
		INSERT INTO "User" (address_id, username, email, password_hash, user_type)
		VALUES (:address_id, :username, :email, :password_hash, :user_type)
	`

	_, err = tx.NamedExec(query, users)
	return err
}

func getAllUserIDs(tx *sqlx.Tx) ([]int64, error) {
	rows, err := tx.Queryx(`SELECT id FROM "User" ORDER BY created_at`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userIDs []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		userIDs = append(userIDs, id)
	}

	return userIDs, nil
}

func getUserIDsByType(tx *sqlx.Tx) (sellerIDs []int64, collectorIDs []int64, companyIDs []int64, err error) {
	// Get Seller IDs
	sellerQuery := `SELECT id FROM "Seller" ORDER BY id`
	rows, err := tx.Queryx(sellerQuery)
	if err != nil {
		return nil, nil, nil, err
	}
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			rows.Close()
			return nil, nil, nil, err
		}
		sellerIDs = append(sellerIDs, id)
	}
	rows.Close()

	// Get Collector IDs
	collectorQuery := `SELECT id FROM "Collector" ORDER BY id`
	rows, err = tx.Queryx(collectorQuery)
	if err != nil {
		return nil, nil, nil, err
	}
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			rows.Close()
			return nil, nil, nil, err
		}
		collectorIDs = append(collectorIDs, id)
	}
	rows.Close()

	// Get Company IDs
	companyQuery := `SELECT id FROM "Company" ORDER BY id`
	rows, err = tx.Queryx(companyQuery)
	if err != nil {
		return nil, nil, nil, err
	}
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			rows.Close()
			return nil, nil, nil, err
		}
		companyIDs = append(companyIDs, id)
	}
	rows.Close()

	return sellerIDs, collectorIDs, companyIDs, nil
}

func seedSellTransactions(tx *sqlx.Tx, sellerIDs, collectorIDs []int64) error {
	if len(sellerIDs) == 0 || len(collectorIDs) == 0 {
		return fmt.Errorf("no sellers or collectors found")
	}

	transactions := []map[string]any{
		// Collector 0 will receive: 100 + 75.25 + 150 = 325.25L
		{"seller_id": sellerIDs[0], "collector_id": collectorIDs[0], "volume": 100.00, "price": 5000.00},
		{"seller_id": sellerIDs[1], "collector_id": collectorIDs[0], "volume": 75.25, "price": 5050.00},
		{"seller_id": sellerIDs[2], "collector_id": collectorIDs[0], "volume": 150.00, "price": 4900.00},

		// Collector 1 will receive: 50.50 + 125.50 + 100 = 276L
		{"seller_id": sellerIDs[0], "collector_id": collectorIDs[1], "volume": 50.50, "price": 5100.00},
		{"seller_id": sellerIDs[2], "collector_id": collectorIDs[1], "volume": 125.50, "price": 4900.00},
		{"seller_id": sellerIDs[3], "collector_id": collectorIDs[1], "volume": 100.00, "price": 5150.00},

		// Collector 2 will receive: 100 + 90 + 200 = 390L
		{"seller_id": sellerIDs[1], "collector_id": collectorIDs[2], "volume": 100.00, "price": 5200.00},
		{"seller_id": sellerIDs[3], "collector_id": collectorIDs[2], "volume": 90.00, "price": 5150.00},
		{"seller_id": sellerIDs[0], "collector_id": collectorIDs[2], "volume": 200.00, "price": 5000.00},

		// Collector 3 will receive: 60 + 80 + 120 = 260L
		{"seller_id": sellerIDs[0], "collector_id": collectorIDs[3], "volume": 60.00, "price": 5000.00},
		{"seller_id": sellerIDs[1], "collector_id": collectorIDs[3], "volume": 80.00, "price": 5100.00},
		{"seller_id": sellerIDs[2], "collector_id": collectorIDs[3], "volume": 120.00, "price": 4950.00},
	}

	query := `
		INSERT INTO "SellTransaction" (seller_id, collector_id, volume, price)
		VALUES (:seller_id, :collector_id, :volume, :price)
	`

	_, err := tx.NamedExec(query, transactions)
	return err
}

func seedDistributeTransactions(tx *sqlx.Tx, collectorIDs, companyIDs []int64) error {
	if len(collectorIDs) == 0 || len(companyIDs) == 0 {
		return fmt.Errorf("no collectors or companies found")
	}

	transactions := []map[string]any{
		// Collector 0 has 325.25L, selling: 150 + 100 = 250L (leaving 75.25L)
		{"collector_id": collectorIDs[0], "company_id": companyIDs[0], "volume": 150.00, "price": 5500.00},
		{"collector_id": collectorIDs[0], "company_id": companyIDs[1], "volume": 100.00, "price": 5600.00},

		// Collector 1 has 276L, selling: 120 + 80 = 200L (leaving 76L)
		{"collector_id": collectorIDs[1], "company_id": companyIDs[0], "volume": 120.00, "price": 5550.00},
		{"collector_id": collectorIDs[1], "company_id": companyIDs[1], "volume": 80.00, "price": 5500.00},

		// Collector 2 has 390L, selling: 200 + 100 = 300L (leaving 90L)
		{"collector_id": collectorIDs[2], "company_id": companyIDs[1], "volume": 200.00, "price": 5700.00},
		{"collector_id": collectorIDs[2], "company_id": companyIDs[0], "volume": 100.00, "price": 5650.00},

		// Collector 3 has 260L, selling: 150 + 50 = 200L (leaving 60L)
		{"collector_id": collectorIDs[3], "company_id": companyIDs[0], "volume": 150.00, "price": 5650.00},
		{"collector_id": collectorIDs[3], "company_id": companyIDs[1], "volume": 50.00, "price": 5600.00},
	}

	query := `
		INSERT INTO "DistributeTransaction" (collector_id, company_id, volume, price)
		VALUES (:collector_id, :company_id, :volume, :price)
	`

	_, err := tx.NamedExec(query, transactions)
	return err
}

func verifyOilInventory(tx *sqlx.Tx, collectorIDs []int64) error {
	// Oil inventory is automatically created and updated by triggers
	// SellTransaction trigger adds volume to collector's inventory
	// DistributeTransaction trigger subtracts volume from collector's inventory

	// Verify that oil records exist for all collectors
	for _, collectorID := range collectorIDs {
		var exists bool
		err := tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM "Oil" WHERE collector_id = $1)`, collectorID).Scan(&exists)
		if err != nil {
			return fmt.Errorf("failed to check oil inventory for collector %d: %w", collectorID, err)
		}
		if !exists {
			return fmt.Errorf("oil inventory not found for collector %d", collectorID)
		}
	}

	// Display current inventory
	rows, err := tx.Queryx(`
		SELECT o.collector_id, c.collector_name, o.total_volume
		FROM "Oil" o
		JOIN "Collector" c ON o.collector_id = c.id
		ORDER BY o.collector_id
	`)
	if err != nil {
		return fmt.Errorf("failed to query oil inventory: %w", err)
	}
	defer rows.Close()

	fmt.Println("\nCurrent Oil Inventory:")
	fmt.Println("=====================")
	for rows.Next() {
		var collectorID int64
		var collectorName string
		var totalVolume float64
		if err := rows.Scan(&collectorID, &collectorName, &totalVolume); err != nil {
			return err
		}
		fmt.Printf("Collector #%d (%s): %.2f liters\n", collectorID, collectorName, totalVolume)
	}
	fmt.Println("=====================")

	return nil
}
