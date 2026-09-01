package `in`.sanskar.stockpilot

import org.junit.Assert.assertEquals
import org.junit.Assert.assertThrows
import org.junit.Test

class ApiClientTest {
    @Test
    fun normalizeBaseUrl_trimsWhitespaceAndTrailingSlash() {
        assertEquals(
            "https://inventory.example.com",
            ApiClient.normalizeBaseUrl("  https://inventory.example.com/  "),
        )
    }

    @Test
    fun normalizeBaseUrl_acceptsEmulatorHttpEndpoint() {
        assertEquals(
            "http://10.0.2.2:8080",
            ApiClient.normalizeBaseUrl("http://10.0.2.2:8080"),
        )
    }

    @Test
    fun normalizeBaseUrl_rejectsUnsupportedSchemes() {
        assertThrows(IllegalArgumentException::class.java) {
            ApiClient.normalizeBaseUrl("ftp://inventory.example.com")
        }
    }

    @Test
    fun normalizeBaseUrl_rejectsMissingHost() {
        assertThrows(IllegalArgumentException::class.java) {
            ApiClient.normalizeBaseUrl("https://")
        }
    }

    @Test
    fun normalizeBaseUrl_rejectsEmbeddedCredentials() {
        assertThrows(IllegalArgumentException::class.java) {
            ApiClient.normalizeBaseUrl("https://user:secret@inventory.example.com")
        }
    }

    @Test
    fun normalizeBaseUrl_rejectsPath() {
        assertThrows(IllegalArgumentException::class.java) {
            ApiClient.normalizeBaseUrl("https://inventory.example.com/app")
        }
    }

    @Test
    fun normalizeBaseUrl_rejectsQuery() {
        assertThrows(IllegalArgumentException::class.java) {
            ApiClient.normalizeBaseUrl("https://inventory.example.com?tenant=one")
        }
    }

    @Test
    fun normalizeBaseUrl_rejectsFragment() {
        assertThrows(IllegalArgumentException::class.java) {
            ApiClient.normalizeBaseUrl("https://inventory.example.com#dashboard")
        }
    }

    @Test
    fun normalizeBaseUrl_rejectsOutOfRangePort() {
        assertThrows(IllegalArgumentException::class.java) {
            ApiClient.normalizeBaseUrl("https://inventory.example.com:65536")
        }
    }
}
