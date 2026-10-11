package `in`.sanskar.stockpilot

import android.app.Activity
import com.google.mlkit.vision.barcode.common.Barcode
import com.google.mlkit.vision.codescanner.GmsBarcodeScanner
import com.google.mlkit.vision.codescanner.GmsBarcodeScannerOptions
import com.google.mlkit.vision.codescanner.GmsBarcodeScanning

class BarcodeScanCoordinator(activity: Activity) {
    private val scanner: GmsBarcodeScanner

    init {
        val options = GmsBarcodeScannerOptions.Builder()
            .setBarcodeFormats(
                Barcode.FORMAT_EAN_8,
                Barcode.FORMAT_EAN_13,
                Barcode.FORMAT_UPC_A,
                Barcode.FORMAT_UPC_E,
                Barcode.FORMAT_CODE_39,
                Barcode.FORMAT_CODE_128,
                Barcode.FORMAT_QR_CODE,
            )
            .enableAutoZoom()
            .build()
        scanner = GmsBarcodeScanning.getClient(activity, options)
    }

    fun start(
        onDetected: (String) -> Unit,
        onCancelled: () -> Unit,
        onFailure: (Exception) -> Unit,
    ) {
        scanner.startScan()
            .addOnSuccessListener { barcode ->
                barcode.rawValue?.trim()?.takeIf { it.isNotEmpty() }?.let(onDetected) ?: onCancelled()
            }
            .addOnCanceledListener(onCancelled)
            .addOnFailureListener(onFailure)
    }
}
