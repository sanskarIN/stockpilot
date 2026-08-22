package domain

import (
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"
)

var (
	ErrNotFound          = errors.New("not found")
	ErrConflict          = errors.New("conflict")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrForbidden         = errors.New("forbidden")
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrBootstrapComplete = errors.New("bootstrap already completed")
	ErrInvalid           = errors.New("invalid input")
)

type UserCreate struct {
	Name         string
	Email        string
	PasswordHash string
	Role         Role
}

type ManagedUserInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     Role   `json:"role"`
}

func (in *ManagedUserInput) NormalizeAndValidate() error {
	in.Name = strings.TrimSpace(in.Name)
	in.Email = strings.ToLower(strings.TrimSpace(in.Email))
	if err := ValidateBootstrap(in.Name, in.Email, in.Password); err != nil {
		return err
	}
	switch in.Role {
	case RoleOwner, RoleAdmin, RoleManager, RoleViewer:
	default:
		return errors.New("role must be owner, admin, manager, or viewer")
	}
	return nil
}

type ProductInput struct {
	SKU          string  `json:"sku"`
	Name         string  `json:"name"`
	Barcode      string  `json:"barcode"`
	CategoryID   string  `json:"categoryId"`
	SupplierID   string  `json:"supplierId"`
	Unit         string  `json:"unit"`
	ReorderPoint float64 `json:"reorderPoint"`
	CostCents    int64   `json:"costCents"`
	PriceCents   int64   `json:"priceCents"`
	TrackLots    bool    `json:"trackLots"`
	TrackExpiry  bool    `json:"trackExpiry"`
	Active       *bool   `json:"active,omitempty"`
}

func (in *ProductInput) NormalizeAndValidate() error {
	in.SKU = strings.ToUpper(strings.TrimSpace(in.SKU))
	in.Name = strings.TrimSpace(in.Name)
	in.Barcode = strings.TrimSpace(in.Barcode)
	in.CategoryID = strings.TrimSpace(in.CategoryID)
	in.SupplierID = strings.TrimSpace(in.SupplierID)
	in.Unit = strings.ToLower(strings.TrimSpace(in.Unit))
	if in.Unit == "" {
		in.Unit = "unit"
	}
	if len(in.SKU) < 2 || len(in.SKU) > 64 {
		return errors.New("sku must be between 2 and 64 characters")
	}
	if len(in.Name) < 2 || len(in.Name) > 160 {
		return errors.New("name must be between 2 and 160 characters")
	}
	if len(in.Barcode) > 128 || len(in.Unit) > 24 {
		return errors.New("barcode or unit is too long")
	}
	if in.ReorderPoint < 0 || in.ReorderPoint > 1_000_000_000 {
		return errors.New("reorder point is out of range")
	}
	if in.CostCents < 0 || in.PriceCents < 0 {
		return errors.New("cost and price cannot be negative")
	}
	if in.TrackExpiry && !in.TrackLots {
		return errors.New("expiry tracking requires lot tracking")
	}
	return nil
}

type CategoryInput struct {
	Name string `json:"name"`
}
type SupplierInput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}
type WarehouseInput struct {
	Name    string `json:"name"`
	Code    string `json:"code"`
	Address string `json:"address"`
}
type LocationInput struct {
	WarehouseID string `json:"warehouseId"`
	Name        string `json:"name"`
	Code        string `json:"code"`
}

func (in *CategoryInput) NormalizeAndValidate() error {
	in.Name = strings.TrimSpace(in.Name)
	if len(in.Name) < 2 || len(in.Name) > 100 {
		return errors.New("category name must be between 2 and 100 characters")
	}
	return nil
}
func (in *SupplierInput) NormalizeAndValidate() error {
	in.Name, in.Email, in.Phone = strings.TrimSpace(in.Name), strings.ToLower(strings.TrimSpace(in.Email)), strings.TrimSpace(in.Phone)
	if len(in.Name) < 2 || len(in.Name) > 160 {
		return errors.New("supplier name must be between 2 and 160 characters")
	}
	if in.Email != "" {
		if _, err := mail.ParseAddress(in.Email); err != nil {
			return errors.New("supplier email is invalid")
		}
	}
	if len(in.Phone) > 40 {
		return errors.New("supplier phone is too long")
	}
	return nil
}
func (in *WarehouseInput) NormalizeAndValidate() error {
	in.Name, in.Code, in.Address = strings.TrimSpace(in.Name), strings.ToUpper(strings.TrimSpace(in.Code)), strings.TrimSpace(in.Address)
	if len(in.Name) < 2 || len(in.Name) > 120 {
		return errors.New("warehouse name must be between 2 and 120 characters")
	}
	if len(in.Code) < 2 || len(in.Code) > 32 {
		return errors.New("warehouse code must be between 2 and 32 characters")
	}
	if len(in.Address) > 500 {
		return errors.New("warehouse address is too long")
	}
	return nil
}
func (in *LocationInput) NormalizeAndValidate() error {
	in.WarehouseID, in.Name, in.Code = strings.TrimSpace(in.WarehouseID), strings.TrimSpace(in.Name), strings.ToUpper(strings.TrimSpace(in.Code))
	if in.WarehouseID == "" {
		return errors.New("warehouseId is required")
	}
	if len(in.Name) < 2 || len(in.Name) > 120 {
		return errors.New("location name must be between 2 and 120 characters")
	}
	if len(in.Code) < 1 || len(in.Code) > 32 {
		return errors.New("location code must be between 1 and 32 characters")
	}
	return nil
}

type MovementInput struct {
	ProductID  string       `json:"productId"`
	LocationID string       `json:"locationId"`
	Type       MovementType `json:"type"`
	Quantity   float64      `json:"quantity"`
	LotNumber  string       `json:"lotNumber"`
	ExpiryDate *time.Time   `json:"expiryDate"`
	Reference  string       `json:"reference"`
	Note       string       `json:"note"`
}

func (in *MovementInput) NormalizeAndValidate() error {
	in.ProductID, in.LocationID = strings.TrimSpace(in.ProductID), strings.TrimSpace(in.LocationID)
	in.LotNumber, in.Reference, in.Note = strings.TrimSpace(in.LotNumber), strings.TrimSpace(in.Reference), strings.TrimSpace(in.Note)
	if in.ProductID == "" || in.LocationID == "" {
		return errors.New("productId and locationId are required")
	}
	if in.Type != MovementIn && in.Type != MovementOut && in.Type != MovementAdjustment {
		return errors.New("type must be in, out, or adjustment")
	}
	if in.Quantity == 0 || in.Quantity < -1_000_000_000 || in.Quantity > 1_000_000_000 {
		return errors.New("quantity must be non-zero and within range")
	}
	if (in.Type == MovementIn || in.Type == MovementOut) && in.Quantity < 0 {
		return errors.New("in/out quantity must be positive")
	}
	if len(in.LotNumber) > 100 || len(in.Reference) > 120 || len(in.Note) > 1000 {
		return errors.New("movement text field is too long")
	}
	return nil
}

type TransferInput struct {
	ProductID      string     `json:"productId"`
	FromLocationID string     `json:"fromLocationId"`
	ToLocationID   string     `json:"toLocationId"`
	Quantity       float64    `json:"quantity"`
	LotNumber      string     `json:"lotNumber"`
	ExpiryDate     *time.Time `json:"expiryDate"`
	Reference      string     `json:"reference"`
	Note           string     `json:"note"`
}

func (in *TransferInput) NormalizeAndValidate() error {
	in.ProductID, in.FromLocationID, in.ToLocationID = strings.TrimSpace(in.ProductID), strings.TrimSpace(in.FromLocationID), strings.TrimSpace(in.ToLocationID)
	in.LotNumber, in.Reference, in.Note = strings.TrimSpace(in.LotNumber), strings.TrimSpace(in.Reference), strings.TrimSpace(in.Note)
	if in.ProductID == "" || in.FromLocationID == "" || in.ToLocationID == "" {
		return errors.New("productId, fromLocationId, and toLocationId are required")
	}
	if in.FromLocationID == in.ToLocationID {
		return errors.New("source and destination locations must differ")
	}
	if in.Quantity <= 0 || in.Quantity > 1_000_000_000 {
		return errors.New("quantity must be positive and within range")
	}
	return nil
}

type PurchaseOrderItemInput struct {
	ProductID       string  `json:"productId"`
	QuantityOrdered float64 `json:"quantityOrdered"`
	UnitCostCents   int64   `json:"unitCostCents"`
}
type PurchaseOrderInput struct {
	SupplierID  string                   `json:"supplierId"`
	WarehouseID string                   `json:"warehouseId"`
	ExpectedAt  *time.Time               `json:"expectedAt"`
	Notes       string                   `json:"notes"`
	Items       []PurchaseOrderItemInput `json:"items"`
}

func (in *PurchaseOrderInput) NormalizeAndValidate() error {
	in.SupplierID, in.WarehouseID, in.Notes = strings.TrimSpace(in.SupplierID), strings.TrimSpace(in.WarehouseID), strings.TrimSpace(in.Notes)
	if in.SupplierID == "" || in.WarehouseID == "" {
		return errors.New("supplierId and warehouseId are required")
	}
	if len(in.Items) == 0 || len(in.Items) > 500 {
		return errors.New("purchase order must contain between 1 and 500 items")
	}
	seen := map[string]bool{}
	for i := range in.Items {
		item := &in.Items[i]
		item.ProductID = strings.TrimSpace(item.ProductID)
		if item.ProductID == "" || item.QuantityOrdered <= 0 || item.QuantityOrdered > 1_000_000_000 || item.UnitCostCents < 0 {
			return fmt.Errorf("invalid purchase order item %d", i+1)
		}
		if seen[item.ProductID] {
			return errors.New("duplicate product in purchase order")
		}
		seen[item.ProductID] = true
	}
	if len(in.Notes) > 2000 {
		return errors.New("notes are too long")
	}
	return nil
}

type ReceiptLineInput struct {
	ProductID  string     `json:"productId"`
	LotNumber  string     `json:"lotNumber"`
	ExpiryDate *time.Time `json:"expiryDate"`
}

type ReceivePurchaseOrderInput struct {
	LocationID string             `json:"locationId"`
	LotNumber  string             `json:"lotNumber,omitempty"`
	ExpiryDate *time.Time         `json:"expiryDate,omitempty"`
	Lines      []ReceiptLineInput `json:"lines,omitempty"`
}

func (in *ReceivePurchaseOrderInput) NormalizeAndValidate() error {
	in.LocationID, in.LotNumber = strings.TrimSpace(in.LocationID), strings.TrimSpace(in.LotNumber)
	if in.LocationID == "" {
		return errors.New("locationId is required")
	}
	if len(in.Lines) > 500 {
		return errors.New("too many receipt lines")
	}
	seen := map[string]bool{}
	for i := range in.Lines {
		in.Lines[i].ProductID = strings.TrimSpace(in.Lines[i].ProductID)
		in.Lines[i].LotNumber = strings.TrimSpace(in.Lines[i].LotNumber)
		if in.Lines[i].ProductID == "" {
			return errors.New("receipt line productId is required")
		}
		if seen[in.Lines[i].ProductID] {
			return errors.New("duplicate receipt line productId")
		}
		seen[in.Lines[i].ProductID] = true
		if len(in.Lines[i].LotNumber) > 100 {
			return errors.New("receipt lot number is too long")
		}
	}
	return nil
}

func ValidateBootstrap(name, email, password string) error {
	if len(strings.TrimSpace(name)) < 2 || len(strings.TrimSpace(name)) > 120 {
		return errors.New("name must be between 2 and 120 characters")
	}
	if _, err := mail.ParseAddress(strings.TrimSpace(email)); err != nil {
		return errors.New("email is invalid")
	}
	if len(password) < 12 || len(password) > 256 {
		return errors.New("password must be between 12 and 256 characters")
	}
	return nil
}

func RoleAllows(role Role, minimum Role) bool {
	rank := map[Role]int{RoleViewer: 1, RoleManager: 2, RoleAdmin: 3, RoleOwner: 4}
	return rank[role] >= rank[minimum]
}
