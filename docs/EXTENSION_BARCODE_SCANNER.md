# Browser companion barcode scanner

The Manifest V3 companion can scan QR and common retail barcode formats from its popup when the browser exposes `BarcodeDetector` and camera access is available.

The extension stores only the configured StockPilot origin. After a successful scan it opens the StockPilot web origin with a `barcode` query parameter. The extension does not read, copy, or transmit StockPilot session cookies.

Because the destination is the normal StockPilot origin, the signed-in web application remains responsible for authenticated product/inventory resolution.

A manual value entry remains available when camera scanning is unavailable or permission is declined.
