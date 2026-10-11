package `in`.sanskar.stockpilot

import android.content.Context
import android.content.res.ColorStateList
import android.graphics.Color
import android.graphics.drawable.GradientDrawable
import android.util.TypedValue
import android.view.Gravity
import android.view.View
import android.widget.Button
import android.widget.EditText
import android.widget.TextView
import kotlin.math.roundToInt

internal fun Context.dp(value: Int): Int =
    TypedValue.applyDimension(TypedValue.COMPLEX_UNIT_DIP, value.toFloat(), resources.displayMetrics).roundToInt()

internal fun Context.stockText(
    text: CharSequence,
    sizeSp: Float = 16f,
    colorRes: Int = R.color.stockpilot_text,
): TextView = TextView(this).apply {
    this.text = text
    setTextSize(TypedValue.COMPLEX_UNIT_SP, sizeSp)
    setTextColor(getColor(colorRes))
    includeFontPadding = false
}

internal fun Context.stockButton(text: CharSequence, primary: Boolean): Button = Button(this).apply {
    this.text = text
    isAllCaps = false
    gravity = Gravity.CENTER
    setTextSize(TypedValue.COMPLEX_UNIT_SP, 15f)
    minHeight = dp(48)
    minimumHeight = dp(48)
    setPadding(dp(18), dp(10), dp(18), dp(10))
    if (primary) {
        setTextColor(Color.WHITE)
        backgroundTintList = ColorStateList.valueOf(getColor(R.color.stockpilot_accent))
    } else {
        setTextColor(getColor(R.color.stockpilot_text))
        backgroundTintList = ColorStateList.valueOf(getColor(R.color.stockpilot_surface_alt))
    }
}

internal fun Context.stockInput(hint: String): EditText = EditText(this).apply {
    this.hint = hint
    setTextColor(getColor(R.color.stockpilot_text))
    setHintTextColor(getColor(R.color.stockpilot_muted))
    setTextSize(TypedValue.COMPLEX_UNIT_SP, 16f)
    setPadding(dp(14), dp(12), dp(14), dp(12))
    minHeight = dp(52)
    background = roundedBackground(
        fillColor = getColor(R.color.stockpilot_surface),
        strokeColor = getColor(R.color.stockpilot_border),
        radiusDp = 12,
    )
}

internal fun Context.roundedBackground(
    fillColor: Int,
    strokeColor: Int? = null,
    radiusDp: Int = 16,
): GradientDrawable = GradientDrawable().apply {
    shape = GradientDrawable.RECTANGLE
    setColor(fillColor)
    cornerRadius = dp(radiusDp).toFloat()
    if (strokeColor != null) {
        setStroke(dp(1), strokeColor)
    }
}

internal fun View.accessibleDescription(value: String) {
    contentDescription = value
    importantForAccessibility = View.IMPORTANT_FOR_ACCESSIBILITY_YES
}
