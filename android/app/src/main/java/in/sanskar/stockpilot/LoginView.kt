package in.sanskar.stockpilot

import android.content.Context
import android.text.InputType
import android.view.Gravity
import android.view.View
import android.view.inputmethod.EditorInfo
import android.widget.ImageView
import android.widget.LinearLayout
import android.widget.ScrollView
import android.widget.TextView

class LoginView(context: Context) : ScrollView(context) {
    var onSubmit: ((serverUrl: String, email: String, password: String) -> Unit)? = null

    private val serverInput = context.stockInput(context.getString(R.string.api_url_label)).apply {
        inputType = InputType.TYPE_CLASS_TEXT or InputType.TYPE_TEXT_VARIATION_URI
        imeOptions = EditorInfo.IME_ACTION_NEXT
        setSingleLine(true)
    }
    private val emailInput = context.stockInput(context.getString(R.string.email_label)).apply {
        inputType = InputType.TYPE_CLASS_TEXT or InputType.TYPE_TEXT_VARIATION_EMAIL_ADDRESS
        imeOptions = EditorInfo.IME_ACTION_NEXT
        setSingleLine(true)
        setAutofillHints(View.AUTOFILL_HINT_EMAIL_ADDRESS)
    }
    private val passwordInput = context.stockInput(context.getString(R.string.password_label)).apply {
        inputType = InputType.TYPE_CLASS_TEXT or InputType.TYPE_TEXT_VARIATION_PASSWORD
        imeOptions = EditorInfo.IME_ACTION_DONE
        setSingleLine(true)
        setAutofillHints(View.AUTOFILL_HINT_PASSWORD)
    }
    private val errorText: TextView = context.stockText("", 14f, R.color.stockpilot_danger).apply {
        visibility = View.GONE
        accessibilityLiveRegion = View.ACCESSIBILITY_LIVE_REGION_POLITE
    }
    private val submitButton = context.stockButton(context.getString(R.string.sign_in_action), primary = true)
    private val progressText = context.stockText(context.getString(R.string.loading), 14f, R.color.stockpilot_muted).apply {
        gravity = Gravity.CENTER
        visibility = View.GONE
        accessibilityLiveRegion = View.ACCESSIBILITY_LIVE_REGION_POLITE
    }

    init {
        isFillViewport = true
        setBackgroundColor(context.getColor(R.color.stockpilot_background))

        val content = LinearLayout(context).apply {
            orientation = LinearLayout.VERTICAL
            gravity = Gravity.CENTER_HORIZONTAL
            setPadding(context.dp(24), context.dp(36), context.dp(24), context.dp(36))
        }
        addView(content, LayoutParams(LayoutParams.MATCH_PARENT, LayoutParams.WRAP_CONTENT))

        content.addView(ImageView(context).apply {
            setImageResource(R.drawable.ic_stockpilot)
            accessibleDescription(context.getString(R.string.app_name))
        }, LinearLayout.LayoutParams(context.dp(72), context.dp(72)).apply {
            bottomMargin = context.dp(24)
        })

        content.addView(context.stockText(context.getString(R.string.sign_in_title), 28f).apply {
            gravity = Gravity.CENTER
            setTypeface(typeface, android.graphics.Typeface.BOLD)
        }, fullWidth().apply { bottomMargin = context.dp(10) })

        content.addView(context.stockText(context.getString(R.string.sign_in_subtitle), 15f, R.color.stockpilot_muted).apply {
            gravity = Gravity.CENTER
        }, fullWidth().apply { bottomMargin = context.dp(28) })

        content.addView(label(context.getString(R.string.api_url_label)), fullWidth().apply { bottomMargin = context.dp(8) })
        content.addView(serverInput, fullWidth().apply { bottomMargin = context.dp(18) })
        content.addView(label(context.getString(R.string.email_label)), fullWidth().apply { bottomMargin = context.dp(8) })
        content.addView(emailInput, fullWidth().apply { bottomMargin = context.dp(18) })
        content.addView(label(context.getString(R.string.password_label)), fullWidth().apply { bottomMargin = context.dp(8) })
        content.addView(passwordInput, fullWidth().apply { bottomMargin = context.dp(18) })
        content.addView(errorText, fullWidth().apply { bottomMargin = context.dp(14) })
        content.addView(submitButton, fullWidth().apply { bottomMargin = context.dp(12) })
        content.addView(progressText, fullWidth())

        submitButton.setOnClickListener { submit() }
        passwordInput.setOnEditorActionListener { _, actionId, _ ->
            if (actionId == EditorInfo.IME_ACTION_DONE) {
                submit()
                true
            } else {
                false
            }
        }
    }

    fun setServerUrl(value: String) {
        if (serverInput.text.toString() != value) {
            serverInput.setText(value)
        }
    }

    fun setLoading(loading: Boolean) {
        serverInput.isEnabled = !loading
        emailInput.isEnabled = !loading
        passwordInput.isEnabled = !loading
        submitButton.isEnabled = !loading
        progressText.visibility = if (loading) View.VISIBLE else View.GONE
    }

    fun showError(message: String?) {
        errorText.text = message.orEmpty()
        errorText.visibility = if (message.isNullOrBlank()) View.GONE else View.VISIBLE
        if (!message.isNullOrBlank()) {
            errorText.announceForAccessibility(message)
        }
    }

    private fun submit() {
        showError(null)
        onSubmit?.invoke(
            serverInput.text.toString(),
            emailInput.text.toString(),
            passwordInput.text.toString(),
        )
    }

    private fun label(value: String): TextView = context.stockText(value, 14f).apply {
        setTypeface(typeface, android.graphics.Typeface.BOLD)
    }

    private fun fullWidth() = LinearLayout.LayoutParams(
        LinearLayout.LayoutParams.MATCH_PARENT,
        LinearLayout.LayoutParams.WRAP_CONTENT,
    )
}
