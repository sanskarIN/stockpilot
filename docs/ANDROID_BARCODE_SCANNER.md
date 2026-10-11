# Android Barcode/QR Scanning

StockPilot Android uses Google Code Scanner for barcode and QR capture.

## Flow

1. The dashboard starts the scanner with supported retail barcode formats and QR codes.
2. A successful scan returns the raw code value to the app.
3. The value is resolved through the authenticated exact-barcode API.
4. The matching product is displayed together with current positive lot/location inventory.
5. Scanner cancellation or failure leaves the existing dashboard state unchanged.

## Security and privacy

The app does not request the Android camera permission for this scanner integration. The scan is delegated to Google Play services, and the StockPilot app receives the resulting code value only. Release server communication continues to require HTTPS.

## Dependency

`com.google.android.gms:play-services-code-scanner:16.1.0`

The scanner module is declared for install-time preparation through the application manifest.
