package in.sanskar.stockpilot

import android.content.Context

class AppPreferences(context: Context) {
    private val preferences = context.getSharedPreferences(PREFERENCES_NAME, Context.MODE_PRIVATE)

    fun serverUrl(): String =
        preferences.getString(KEY_SERVER_URL, BuildConfig.DEFAULT_API_URL).orEmpty().trim().trimEnd('/')

    fun setServerUrl(value: String) {
        preferences.edit().putString(KEY_SERVER_URL, value.trim().trimEnd('/')).apply()
    }

    companion object {
        private const val PREFERENCES_NAME = "stockpilot_preferences"
        private const val KEY_SERVER_URL = "server_url"
    }
}
