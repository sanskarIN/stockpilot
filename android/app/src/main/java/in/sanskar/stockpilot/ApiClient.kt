package in.sanskar.stockpilot

import org.json.JSONObject
import java.io.IOException
import java.net.HttpCookie
import java.net.HttpURLConnection
import java.net.URL

class ApiClient(private val sessionStore: SecureSessionStore) {
    fun login(baseUrl: String, email: String, password: String): User {
        val body = JSONObject()
            .put("email", email.trim())
            .put("password", password)
            .toString()
        val response = request(
            baseUrl = baseUrl,
            method = "POST",
            path = "/api/v1/auth/login",
            body = body,
            authenticated = false,
        )
        val token = response.setCookies
            .asSequence()
            .flatMap { runCatching { HttpCookie.parse(it).asSequence() }.getOrDefault(emptySequence()) }
            .firstOrNull { it.name == SESSION_COOKIE_NAME }
            ?.value
            ?.takeIf { it.isNotBlank() }
            ?: throw ApiException(500, "The server did not establish a session.")
        sessionStore.save(token)
        return User.fromJson(JSONObject(response.body).getJSONObject("user"))
    }

    fun currentUser(baseUrl: String): User =
        User.fromJson(JSONObject(request(baseUrl, "GET", "/api/v1/auth/me").body))

    fun dashboard(baseUrl: String): DashboardSnapshot {
        val user = currentUser(baseUrl)
        val products = JSONObject(request(baseUrl, "GET", "/api/v1/products?active=true&limit=50").body)
            .optJSONArray("items")
            ?.mapObjects(Product::fromJson)
            .orEmpty()
        val lowStock = JSONObject(request(baseUrl, "GET", "/api/v1/inventory/low-stock?limit=50").body)
            .optJSONArray("items")
            ?.mapObjects(StockBalance::fromJson)
            .orEmpty()
        val orders = JSONObject(request(baseUrl, "GET", "/api/v1/orders?limit=50").body)
            .optJSONArray("items")
            ?.mapObjects(PurchaseOrder::fromJson)
            .orEmpty()
        return DashboardSnapshot(user, products, lowStock, orders)
    }

    fun logout(baseUrl: String) {
        try {
            request(baseUrl, "POST", "/api/v1/auth/logout", body = null)
        } finally {
            sessionStore.clear()
        }
    }

    fun hasSession(): Boolean = sessionStore.token() != null

    private fun request(
        baseUrl: String,
        method: String,
        path: String,
        body: String? = null,
        authenticated: Boolean = true,
    ): Response {
        val normalizedBaseUrl = normalizeBaseUrl(baseUrl)
        val connection = (URL(normalizedBaseUrl + path).openConnection() as HttpURLConnection).apply {
            requestMethod = method
            connectTimeout = CONNECT_TIMEOUT_MS
            readTimeout = READ_TIMEOUT_MS
            useCaches = false
            instanceFollowRedirects = false
            setRequestProperty("Accept", "application/json")
            setRequestProperty("User-Agent", "StockPilot-Android/${BuildConfig.VERSION_NAME}")
        }

        if (authenticated) {
            val token = sessionStore.token() ?: throw AuthenticationRequiredException()
            connection.setRequestProperty("Cookie", "$SESSION_COOKIE_NAME=$token")
        }
        if (method != "GET" && method != "HEAD") {
            connection.setRequestProperty("X-StockPilot-CSRF", "1")
        }
        if (body != null) {
            connection.doOutput = true
            connection.setRequestProperty("Content-Type", "application/json; charset=utf-8")
            connection.outputStream.bufferedWriter(Charsets.UTF_8).use { it.write(body) }
        }

        return try {
            val status = connection.responseCode
            val responseBody = when {
                status == HttpURLConnection.HTTP_NO_CONTENT -> ""
                status in 200..299 -> connection.inputStream?.bufferedReader(Charsets.UTF_8)?.use { it.readText() }.orEmpty()
                else -> connection.errorStream?.bufferedReader(Charsets.UTF_8)?.use { it.readText() }.orEmpty()
            }
            val setCookies = connection.headerFields
                .filterKeys { it?.equals("Set-Cookie", ignoreCase = true) == true }
                .values
                .flatten()

            if (status !in 200..299) {
                if (status == HttpURLConnection.HTTP_UNAUTHORIZED) {
                    sessionStore.clear()
                }
                val message = runCatching { JSONObject(responseBody).optString("error") }
                    .getOrNull()
                    ?.takeIf { it.isNotBlank() }
                    ?: "Request failed with HTTP $status."
                throw ApiException(status, message)
            }
            Response(responseBody, setCookies)
        } catch (error: ApiException) {
            throw error
        } catch (error: IOException) {
            throw ApiException(0, error.message ?: "Unable to reach the StockPilot server.", error)
        } finally {
            connection.disconnect()
        }
    }

    companion object {
        private const val SESSION_COOKIE_NAME = "stockpilot_session"
        private const val CONNECT_TIMEOUT_MS = 10_000
        private const val READ_TIMEOUT_MS = 15_000

        fun normalizeBaseUrl(value: String): String {
            val normalized = value.trim().trimEnd('/')
            val url = runCatching { URL(normalized) }
                .getOrElse { throw IllegalArgumentException("Enter a valid http:// or https:// server URL.") }
            require(url.protocol == "http" || url.protocol == "https") {
                "Server URL must use http:// or https://."
            }
            require(url.host.isNotBlank()) { "Server URL must include a host." }
            return normalized
        }
    }

    private data class Response(val body: String, val setCookies: List<String>)
}

open class ApiException(
    val statusCode: Int,
    override val message: String,
    cause: Throwable? = null,
) : IOException(message, cause)

class AuthenticationRequiredException : ApiException(401, "Sign in to continue.")
