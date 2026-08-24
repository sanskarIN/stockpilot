package in.sanskar.stockpilot

import android.app.Activity
import android.os.Build
import android.os.Bundle
import android.os.Handler
import android.os.Looper
import android.view.View
import android.view.WindowInsets
import java.util.concurrent.Executors

class MainActivity : Activity() {
    private val executor = Executors.newSingleThreadExecutor()
    private val mainHandler = Handler(Looper.getMainLooper())

    private lateinit var preferences: AppPreferences
    private lateinit var sessionStore: SecureSessionStore
    private lateinit var apiClient: ApiClient
    private var loginView: LoginView? = null
    private var dashboardView: DashboardView? = null

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        preferences = AppPreferences(this)
        sessionStore = SecureSessionStore(this)
        apiClient = ApiClient(sessionStore)

        if (apiClient.hasSession() && preferences.serverUrl().isNotBlank()) {
            showDashboard()
            refreshDashboard()
        } else {
            showLogin()
        }
    }

    override fun onDestroy() {
        executor.shutdownNow()
        super.onDestroy()
    }

    private fun showLogin(error: String? = null) {
        val view = LoginView(this).apply {
            setServerUrl(preferences.serverUrl())
            onSubmit = ::signIn
            showError(error)
        }
        loginView = view
        dashboardView = null
        installContent(view)
    }

    private fun showDashboard() {
        val view = DashboardView(this).apply {
            onRefresh = ::refreshDashboard
            onLogout = ::signOut
        }
        dashboardView = view
        loginView = null
        installContent(view)
    }

    private fun signIn(rawServerUrl: String, email: String, password: String) {
        val serverUrl = try {
            ApiClient.normalizeBaseUrl(rawServerUrl)
        } catch (error: IllegalArgumentException) {
            loginView?.showError(error.message)
            return
        }
        if (!BuildConfig.DEBUG && !serverUrl.startsWith("https://", ignoreCase = true)) {
            loginView?.showError("Release builds require an HTTPS StockPilot server.")
            return
        }
        if (email.isBlank() || password.isBlank()) {
            loginView?.showError("Email and password are required.")
            return
        }

        preferences.setServerUrl(serverUrl)
        loginView?.setLoading(true)
        executor.execute {
            runCatching { apiClient.login(serverUrl, email, password) }
                .onSuccess {
                    postToUi {
                        showDashboard()
                        refreshDashboard()
                    }
                }
                .onFailure { error ->
                    postToUi {
                        loginView?.setLoading(false)
                        loginView?.showError(error.userMessage())
                    }
                }
        }
    }

    private fun refreshDashboard() {
        val serverUrl = preferences.serverUrl()
        if (serverUrl.isBlank()) {
            sessionStore.clear()
            showLogin("Enter the StockPilot server URL to continue.")
            return
        }

        dashboardView?.showError(null)
        dashboardView?.setLoading(true)
        executor.execute {
            runCatching { apiClient.dashboard(serverUrl) }
                .onSuccess { snapshot ->
                    postToUi {
                        dashboardView?.render(snapshot)
                        dashboardView?.setLoading(false)
                    }
                }
                .onFailure { error ->
                    postToUi {
                        if (error is AuthenticationRequiredException || (error is ApiException && error.statusCode == 401)) {
                            showLogin("Your session expired. Sign in again.")
                        } else {
                            dashboardView?.setLoading(false)
                            dashboardView?.showError(error.userMessage())
                        }
                    }
                }
        }
    }

    private fun signOut() {
        val serverUrl = preferences.serverUrl()
        dashboardView?.setLoading(true)
        executor.execute {
            val failure = runCatching { apiClient.logout(serverUrl) }.exceptionOrNull()
            postToUi {
                showLogin(
                    failure?.let { "Signed out on this device. The server could not be reached to close the remote session." },
                )
            }
        }
    }

    private fun installContent(view: View) {
        applySystemBarInsets(view)
        setContentView(view)
    }

    private fun applySystemBarInsets(view: View) {
        view.setOnApplyWindowInsetsListener { target, insets ->
            if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.R) {
                val bars = insets.getInsets(WindowInsets.Type.systemBars())
                target.setPadding(bars.left, bars.top, bars.right, bars.bottom)
            } else {
                @Suppress("DEPRECATION")
                target.setPadding(
                    insets.systemWindowInsetLeft,
                    insets.systemWindowInsetTop,
                    insets.systemWindowInsetRight,
                    insets.systemWindowInsetBottom,
                )
            }
            insets
        }
        view.requestApplyInsets()
    }

    private fun postToUi(block: () -> Unit) {
        mainHandler.post {
            if (!isFinishing && !isDestroyed) {
                block()
            }
        }
    }

    private fun Throwable.userMessage(): String = when (this) {
        is ApiException -> message
        else -> message ?: "Something went wrong. Please try again."
    }
}
