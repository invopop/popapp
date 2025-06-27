// Package domain contains the domain logic for the application.
package domain

// ABOUT: Domain is where the business logic resides.
// It defines the structs that represent the core concepts of the application.
// These structs connect the interfaces with the models.
// The best way to work is to create one struct per conceptand then initialize them when creating a new domain.Setup.
// Each struct should be implemented in its own file to separate concerns and improve maintainability.
//
// In addition, each method should have it's accompanying structs for requests. It may also have structs for responses.
// For example, an Invoices struct, which handles the business logic for invoices, has a method called CreateInvoice.
// This method should have a corresponding struct for the request.
// In this case, the method may return an Invoice from the model or an InvoiceResponse struct.

type Setup struct{}

func New() *Setup {
	return &Setup{}
}
