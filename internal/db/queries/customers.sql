-- name: GetAllCustomers :many
SELECT * FROM customers JOIN addresses a ON a.address_id = customers.customer_address_id
ORDER BY customer_id;

-- name: CreateCustomer :exec
INSERT INTO customers (
    customer_id, customer_address_id, name, last_name, email, phone_number
) VALUES (
    $1, $2, $3, $4, $5, $6
)ON CONFLICT DO NOTHING;