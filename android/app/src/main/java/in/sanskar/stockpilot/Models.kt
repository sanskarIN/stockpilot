package `in`.sanskar.stockpilot

import org.json.JSONArray
import org.json.JSONObject

data class User(
    val id: String,
    val email: String,
    val displayName: String,
    val role: String,
    val active: Boolean,
    val createdAt: String,
    val updatedAt: String,
    val lastLoginAt: String?,
) {
    companion object {
        fun fromJson(json: JSONObject) = User(
            id = json.optString("id"),
            email = json.optString("email"),
            displayName = json.optString("displayName"),
            role = json.optString("role"),
            active = json.optBoolean("active", true),
            createdAt = json.optString("createdAt"),
            updatedAt = json.optString("updatedAt"),
            lastLoginAt = json.nullableString("lastLoginAt"),
        )
    }
}

data class Product(
    val id: String,
    val sku: String,
    val name: String,
    val unit: String,
    val unitCostMinor: Long,
    val currency: String,
    val reorderPoint: Long,
    val active: Boolean,
) {
    companion object {
        fun fromJson(json: JSONObject) = Product(
            id = json.optString("id"),
            sku = json.optString("sku"),
            name = json.optString("name"),
            unit = json.optString("unit"),
            unitCostMinor = json.optLong("unitCostMinor"),
            currency = json.optString("currency", "INR"),
            reorderPoint = json.optLong("reorderPoint"),
            active = json.optBoolean("active", true),
        )
    }
}

data class StockBalance(
    val productId: String,
    val locationId: String,
    val lotId: String?,
    val quantity: Long,
    val updatedAt: String,
) {
    companion object {
        fun fromJson(json: JSONObject) = StockBalance(
            productId = json.optString("productId"),
            locationId = json.optString("locationId"),
            lotId = json.nullableString("lotId"),
            quantity = json.optLong("quantity"),
            updatedAt = json.optString("updatedAt"),
        )
    }
}

data class PurchaseOrder(
    val id: String,
    val number: String,
    val supplierId: String,
    val warehouseId: String,
    val status: String,
    val currency: String,
    val createdAt: String,
) {
    companion object {
        fun fromJson(json: JSONObject) = PurchaseOrder(
            id = json.optString("id"),
            number = json.optString("number"),
            supplierId = json.optString("supplierId"),
            warehouseId = json.optString("warehouseId"),
            status = json.optString("status"),
            currency = json.optString("currency", "INR"),
            createdAt = json.optString("createdAt"),
        )
    }
}

data class DashboardSnapshot(
    val user: User,
    val products: List<Product>,
    val lowStock: List<StockBalance>,
    val orders: List<PurchaseOrder>,
)

internal fun JSONObject.nullableString(name: String): String? =
    if (isNull(name)) null else optString(name).takeIf { it.isNotBlank() }

internal inline fun <T> JSONArray.mapObjects(transform: (JSONObject) -> T): List<T> =
    buildList(length()) {
        for (index in 0 until length()) {
            val item = optJSONObject(index) ?: continue
            add(transform(item))
        }
    }
