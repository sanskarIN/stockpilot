# StockPilot Feature Matrix

| Capability | Backend | Web | Android | Browser companion | Tests/CI |
| --- | --- | --- | --- | --- | --- |
| Authentication and sessions | Implemented | Implemented | Implemented | Public health/launcher | Covered |
| RBAC | Implemented | Role-aware UI | Server-authoritative | N/A | Covered |
| Product catalog | Implemented | Dashboard/read flow; management pending | Read flow foundation | N/A | Covered at domain/API level |
| Inventory balances | Implemented | Dashboard | Dashboard | N/A | Covered |
| Inventory movements/transfers | Implemented | Guided UI pending | Workflow UI pending | N/A | Backend coverage |
| Purchase orders | Implemented | Workflow UI pending | Workflow UI pending | N/A | Backend coverage |
| Lot/expiry tracking | Implemented | Receiving UI pending | Receiving UI pending | N/A | Domain/API coverage |
| Barcode lookup | Implemented | Typed API available | API foundation | Public companion integration pending | Backend/API coverage |
| Reorder recommendations | Implemented | Dashboard | Client workflow pending | N/A | Unit + integration coverage |
| Inventory valuation | Implemented | Dashboard | Client report pending | N/A | Unit + integration coverage |
| Audit reads | Implemented | Viewer pending | Viewer pending | N/A | Backend coverage |
| Audit writes | Planned | Planned | Planned | N/A | Planned |
| CSV import/export | Planned | Planned | Planned | N/A | Planned |
| Advanced analytics | Planned | Planned | Planned | N/A | Planned |
| Release/operations | Health, migrations, backup tooling | Build | Lint/test/build gate | Validation gate | CI + CodeQL |
