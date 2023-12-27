-- name: GetAllCustomers :many
SELECT * FROM customers JOIN addresses a ON a.address_id = customers.customer_address_id
ORDER BY customer_id;

-- name: CreateCustomerAddress :one
INSERT INTO addresses (
    country, street, city, zip_code
) VALUES (
    $1, $2, $3, $4
) RETURNING address_id;

-- name: CreateCustomer :one
INSERT INTO customers (
    customer_address_id, name, last_name, email, phone_number
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING customer_id;
