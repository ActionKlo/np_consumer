package handlers

import (
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"go.uber.org/zap"
	"log"
	"np_consumer/internal/db"
	"np_consumer/internal/db/gen"
	"np_consumer/internal/models"
	"sync"
)

func CreateSomeCustomers(d *db.ServiceDB) {
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			newCustomer := &models.CustomerInfo{
				Address: &gen.Address{
					Country: gofakeit.Country(),
					Street:  gofakeit.Street(),
					City:    gofakeit.City(),
					ZipCode: gofakeit.Zip(),
				},
				Customer: &gen.Customer{
					Name:        gofakeit.Name(),
					LastName:    gofakeit.LastName(),
					Email:       gofakeit.Email(),
					PhoneNumber: gofakeit.Phone(),
				},
			}

			err := d.CreateCustomer(newCustomer)
			if err != nil {
				d.Logger.Error("failed to create new customer:", zap.Error(err))
			}
		}(i)
	}

	wg.Wait()
}

func GetALlCustomersWithAddress(d *db.ServiceDB) ([]models.CustomerInfo, error) {
	customers, err := d.GetAllCustomers()
	if err != nil {
		log.Fatal(err)
	}

	var customersInfo []models.CustomerInfo
	for _, c := range customers {
		customersInfo = append(customersInfo, models.CustomerInfo{
			Address: &gen.Address{
				AddressID: c.AddressID,
				Country:   c.Country,
				Street:    c.Street,
				City:      c.City,
				ZipCode:   c.ZipCode,
			},
			Customer: &gen.Customer{
				CustomerID:        c.CustomerID,
				CustomerAddressID: c.CustomerAddressID,
				Name:              c.Name,
				LastName:          c.LastName,
				Email:             c.Email,
				PhoneNumber:       c.PhoneNumber,
			},
		})
	}

	for _, customer := range customersInfo {
		fmt.Println(customer.Customer.Email, customer.Address.City)
	}

	return customersInfo, nil
}
