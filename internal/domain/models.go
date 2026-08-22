package domain

import "time"

type Role string

const (
	RoleOwner   Role = "owner"
	RoleAdmin   Role = "admin"
	RoleManager Role = "manager"
	RoleViewer  Role = "viewer"
)

type User struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Role         Role      `json:"role"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"createdAt"`
}

type Category struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

type Supplier struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email,omitempty"`
	Phone     string    `json:"phone,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}

type Warehouse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	Address   string    `json:"address,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}

type Location struct {
	ID          string    `json:"id"`
	WarehouseID string    `json:"warehouseId"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	CreatedAt   time.Time `json:"createdAt"`
}

type Product struct {
	ID           string    `json:"id"`
	SKU          string    `json:"sku"`
	Name         string    `json:"name"`
	Barcode      string    `json:"barcode,omitempty"`
	CategoryID   string    `json:"categoryId,omitempty"`
	SupplierID   string    `json:"supplierId,omitempty"`
	Unit         string    `json:"unit"`
	ReorderPoint float64   `json:"reorderPoint"`
	CostCents    int64     `json:"costCents"`
	PriceCents   int64     `json:"priceCents"`
	TrackLots    bool      `json:"trackLots"`
	TrackExpiry  bool      `json:"trackExpiry"`
	Active       bool      `json:"active"`
	CurrentStock float64   `json:"currentStock"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type InventoryBalance struct {
	ID            string     `json:"id"`
	ProductID     string     `json:"productId"`
	SKU           string     `json:"sku"`
	ProductName   string     `json:"productName"`
	LocationID    string     `json:"locationId"`
	LocationName  string     `json:"locationName"`
	WarehouseID   string     `json:"warehouseId"`
	Warehouse     string     `json:"warehouse"`
	LotNumber     string     `json:"lotNumber,omitempty"`
	ExpiryDate    *time.Time `json:"expiryDate,omitempty"`
	Quantity      float64    `json:"quantity"`
	UnitCostCents int64      `json:"unitCostCents"`
	UpdatedAt     time.Time  `json:"updatedAt"`
}

type MovementType string

const (
	MovementIn          MovementType = "in"
	MovementOut         MovementType = "out"
	MovementAdjustment  MovementType = "adjustment"
	MovementTransferIn  MovementType = "transfer_in"
	MovementTransferOut MovementType = "transfer_out"
	MovementReceipt     MovementType = "receipt"
)

type StockMovement struct {
	ID            string       `json:"id"`
	ProductID     string       `json:"productId"`
	SKU           string       `json:"sku,omitempty"`
	ProductName   string       `json:"productName,omitempty"`
	LocationID    string       `json:"locationId"`
	LocationName  string       `json:"locationName,omitempty"`
	WarehouseName string       `json:"warehouseName,omitempty"`
	Type          MovementType `json:"type"`
	Quantity      float64      `json:"quantity"`
	LotNumber     string       `json:"lotNumber,omitempty"`
	ExpiryDate    *time.Time   `json:"expiryDate,omitempty"`
	Reference     string       `json:"reference,omitempty"`
	Note          string       `json:"note,omitempty"`
	ActorUserID   string       `json:"actorUserId"`
	CreatedAt     time.Time    `json:"createdAt"`
}

type PurchaseOrder struct {
	ID            string              `json:"id"`
	Number        string              `json:"number"`
	SupplierID    string              `json:"supplierId"`
	SupplierName  string              `json:"supplierName,omitempty"`
	WarehouseID   string              `json:"warehouseId"`
	WarehouseName string              `json:"warehouseName,omitempty"`
	Status        string              `json:"status"`
	ExpectedAt    *time.Time          `json:"expectedAt,omitempty"`
	Notes         string              `json:"notes,omitempty"`
	Items         []PurchaseOrderItem `json:"items"`
	CreatedBy     string              `json:"createdBy"`
	CreatedAt     time.Time           `json:"createdAt"`
	ReceivedAt    *time.Time          `json:"receivedAt,omitempty"`
}

type PurchaseOrderItem struct {
	ID               string  `json:"id"`
	PurchaseOrderID  string  `json:"purchaseOrderId"`
	ProductID        string  `json:"productId"`
	SKU              string  `json:"sku,omitempty"`
	ProductName      string  `json:"productName,omitempty"`
	QuantityOrdered  float64 `json:"quantityOrdered"`
	QuantityReceived float64 `json:"quantityReceived"`
	UnitCostCents    int64   `json:"unitCostCents"`
}

type LowStockItem struct {
	ProductID    string  `json:"productId"`
	SKU          string  `json:"sku"`
	Name         string  `json:"name"`
	CurrentStock float64 `json:"currentStock"`
	ReorderPoint float64 `json:"reorderPoint"`
	SuggestedQty float64 `json:"suggestedQty"`
}

type ValuationRow struct {
	ProductID  string  `json:"productId"`
	SKU        string  `json:"sku"`
	Name       string  `json:"name"`
	Quantity   float64 `json:"quantity"`
	ValueCents int64   `json:"valueCents"`
}

type ValuationReport struct {
	Method          string         `json:"method"`
	TotalValueCents int64          `json:"totalValueCents"`
	Rows            []ValuationRow `json:"rows"`
	GeneratedAt     time.Time      `json:"generatedAt"`
}

type Dashboard struct {
	ProductCount        int64   `json:"productCount"`
	WarehouseCount      int64   `json:"warehouseCount"`
	LowStockCount       int64   `json:"lowStockCount"`
	OpenPurchaseOrders  int64   `json:"openPurchaseOrders"`
	StockUnits          float64 `json:"stockUnits"`
	InventoryValueCents int64   `json:"inventoryValueCents"`
	MovementsLast30Days int64   `json:"movementsLast30Days"`
	ExpiringLots30Days  int64   `json:"expiringLots30Days"`
}

type AgingBucket struct {
	Label      string  `json:"label"`
	Quantity   float64 `json:"quantity"`
	ValueCents int64   `json:"valueCents"`
}

type AgingReport struct {
	Basis       string        `json:"basis"`
	GeneratedAt time.Time     `json:"generatedAt"`
	Buckets     []AgingBucket `json:"buckets"`
}

type AuditEvent struct {
	ID          int64     `json:"id"`
	ActorUserID string    `json:"actorUserId,omitempty"`
	Action      string    `json:"action"`
	EntityType  string    `json:"entityType"`
	EntityID    string    `json:"entityId,omitempty"`
	Metadata    string    `json:"metadata,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
}
