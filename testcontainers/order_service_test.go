package testcontainers

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql" // mysql driver import needed
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestOrderService(t *testing.T) {
	ctx := context.Background()

	// Start a MySQL container
	request := testcontainers.ContainerRequest{
		Image:        "mysql:5.7.43",
		ExposedPorts: []string{"3306/tcp"},
		Env: map[string]string{
			"MYSQL_ROOT_PASSWORD": "password",
		},
		WaitingFor: wait.ForLog("port: 3306  MySQL Community Server"),
	}
	mysqlContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: request,
		Started:          true,
	})

	assert.NoError(t, err)
	defer mysqlContainer.Terminate(ctx)

	// Get the MySQL container's host and port
	host, err := mysqlContainer.Host(ctx)
	assert.NoError(t, err)
	port, err := mysqlContainer.MappedPort(ctx, "3306")
	assert.NoError(t, err)

	// Create a database connection
	dsn := fmt.Sprintf("root:password@tcp(%s:%s)/", host, port.Port())
	db, err := sql.Open("mysql", dsn)
	assert.NoError(t, err)
	defer db.Close()

	// Create database and table
	createDatabaseAndTable(db)

	// Initialize and use the order service
	orderService := NewOrderService(db)

	// Add an order
	err = orderService.AddOrder(100, 2)
	assert.NoError(t, err)

	// Get the added order
	productID, quantity, err := orderService.GetOrder(1)
	assert.NoError(t, err)
	assert.Equal(t, 100, productID)
	assert.Equal(t, 2, quantity)
}

func createDatabaseAndTable(db *sql.DB) {
	db.Exec(`CREATE DATABASE IF NOT EXISTS orders;`)
	db.Exec(`USE orders;`)
	db.Exec(`CREATE TABLE IF NOT EXISTS orders ( id INT AUTO_INCREMENT PRIMARY KEY, product_id INT, quantity INT);`)
}
