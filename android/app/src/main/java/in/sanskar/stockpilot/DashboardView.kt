package in.sanskar.stockpilot

import android.content.Context
import android.graphics.Typeface
import android.view.Gravity
import android.view.View
import android.widget.LinearLayout
import android.widget.ScrollView
import android.widget.TextView
import java.text.NumberFormat
import java.util.Currency

class DashboardView(context: Context) : ScrollView(context) {
    var onRefresh: (() -> Unit)? = null
    var onLogout: (() -> Unit)? = null

    private val content = LinearLayout(context).apply {
        orientation = LinearLayout.VERTICAL
        setPadding(context.dp(20), context.dp(24), context.dp(20), context.dp(32))
    }
    private val subtitle = context.stockText("", 14f, R.color.stockpilot_muted)
    private val statusText = context.stockText("", 14f, R.color.stockpilot_muted).apply {
        visibility = View.GONE
        accessibilityLiveRegion = View.ACCESSIBILITY_LIVE_REGION_POLITE
    }
    private val refreshButton = context.stockButton(context.getString(R.string.refresh_action), primary = true)
    private val logoutButton = context.stockButton(context.getString(R.string.sign_out_action), primary = false)
    private val dynamicContent = LinearLayout(context).apply {
        orientation = LinearLayout.VERTICAL
    }

    init {
        isFillViewport = true
        setBackgroundColor(context.getColor(R.color.stockpilot_background))
        addView(content, LayoutParams(LayoutParams.MATCH_PARENT, LayoutParams.WRAP_CONTENT))

        content.addView(context.stockText(context.getString(R.string.app_name), 28f).apply {
            setTypeface(typeface, Typeface.BOLD)
        }, fullWidth().apply { bottomMargin = context.dp(4) })
        content.addView(subtitle, fullWidth().apply { bottomMargin = context.dp(16) })

        val actions = LinearLayout(context).apply {
            orientation = LinearLayout.HORIZONTAL
            gravity = Gravity.CENTER_VERTICAL
            addView(refreshButton, LinearLayout.LayoutParams(0, LinearLayout.LayoutParams.WRAP_CONTENT, 1f).apply {
                marginEnd = context.dp(6)
            })
            addView(logoutButton, LinearLayout.LayoutParams(0, LinearLayout.LayoutParams.WRAP_CONTENT, 1f).apply {
                marginStart = context.dp(6)
            })
        }
        content.addView(actions, fullWidth().apply { bottomMargin = context.dp(12) })
        content.addView(statusText, fullWidth().apply { bottomMargin = context.dp(12) })
        content.addView(dynamicContent, fullWidth())

        refreshButton.setOnClickListener { onRefresh?.invoke() }
        logoutButton.setOnClickListener { onLogout?.invoke() }
    }

    fun render(snapshot: DashboardSnapshot) {
        subtitle.text = "${snapshot.user.displayName} · ${snapshot.user.role.replaceFirstChar { it.uppercase() }}"
        dynamicContent.removeAllViews()

        val metrics = LinearLayout(context).apply {
            orientation = LinearLayout.HORIZONTAL
            addView(metricCard(context.getString(R.string.products_title), snapshot.products.size), weightedCard(end = 6))
            addView(metricCard(context.getString(R.string.low_stock_title), snapshot.lowStock.size), weightedCard(start = 6, end = 6))
            addView(metricCard(context.getString(R.string.orders_title), snapshot.orders.size), weightedCard(start = 6))
        }
        dynamicContent.addView(metrics, fullWidth().apply { bottomMargin = context.dp(16) })

        dynamicContent.addView(
            sectionCard(
                title = context.getString(R.string.products_title),
                emptyMessage = "No active products yet.",
                rows = snapshot.products.take(MAX_ROWS).map { product ->
                    row(
                        primary = product.name.ifBlank { product.sku },
                        secondary = "${product.sku} · ${formatMoney(product.unitCostMinor, product.currency)} / ${product.unit}",
                    )
                },
            ),
            fullWidth().apply { bottomMargin = context.dp(14) },
        )

        val productNames = snapshot.products.associate { it.id to it.name }
        dynamicContent.addView(
            sectionCard(
                title = context.getString(R.string.low_stock_title),
                emptyMessage = "No low-stock balances.",
                rows = snapshot.lowStock.take(MAX_ROWS).map { balance ->
                    val productName = productNames[balance.productId].orEmpty().ifBlank { balance.productId }
                    row(
                        primary = productName,
                        secondary = "${balance.quantity} on hand · location ${balance.locationId}",
                        danger = true,
                    )
                },
            ),
            fullWidth().apply { bottomMargin = context.dp(14) },
        )

        dynamicContent.addView(
            sectionCard(
                title = context.getString(R.string.orders_title),
                emptyMessage = "No purchase orders yet.",
                rows = snapshot.orders.take(MAX_ROWS).map { order ->
                    row(
                        primary = order.number,
                        secondary = "${order.status.replace('_', ' ')} · ${order.currency}",
                    )
                },
            ),
            fullWidth().apply { bottomMargin = context.dp(14) },
        )

        dynamicContent.addView(
            sectionCard(
                title = context.getString(R.string.account_title),
                emptyMessage = "",
                rows = listOf(
                    row(snapshot.user.email, "Role: ${snapshot.user.role}"),
                    row("Server session", if (snapshot.user.active) "Active" else "Account disabled"),
                    row(context.getString(R.string.app_credit), "StockPilot Android ${BuildConfig.VERSION_NAME}"),
                ),
            ),
            fullWidth(),
        )
    }

    fun setLoading(loading: Boolean) {
        refreshButton.isEnabled = !loading
        logoutButton.isEnabled = !loading
        if (loading) {
            statusText.setTextColor(context.getColor(R.color.stockpilot_muted))
            statusText.text = context.getString(R.string.loading)
            statusText.visibility = View.VISIBLE
        } else if (statusText.text == context.getString(R.string.loading)) {
            statusText.visibility = View.GONE
        }
    }

    fun showError(message: String?) {
        if (message.isNullOrBlank()) {
            statusText.visibility = View.GONE
            return
        }
        statusText.setTextColor(context.getColor(R.color.stockpilot_danger))
        statusText.text = message
        statusText.visibility = View.VISIBLE
        statusText.announceForAccessibility(message)
    }

    private fun metricCard(label: String, value: Int): View = LinearLayout(context).apply {
        orientation = LinearLayout.VERTICAL
        gravity = Gravity.CENTER
        setPadding(context.dp(8), context.dp(14), context.dp(8), context.dp(14))
        background = context.roundedBackground(
            fillColor = context.getColor(R.color.stockpilot_surface),
            strokeColor = context.getColor(R.color.stockpilot_border),
            radiusDp = 14,
        )
        addView(context.stockText(value.toString(), 24f).apply {
            gravity = Gravity.CENTER
            setTypeface(typeface, Typeface.BOLD)
        })
        addView(context.stockText(label, 12f, R.color.stockpilot_muted).apply { gravity = Gravity.CENTER })
    }

    private fun sectionCard(title: String, emptyMessage: String, rows: List<View>): View = LinearLayout(context).apply {
        orientation = LinearLayout.VERTICAL
        setPadding(context.dp(16), context.dp(16), context.dp(16), context.dp(16))
        background = context.roundedBackground(
            fillColor = context.getColor(R.color.stockpilot_surface),
            strokeColor = context.getColor(R.color.stockpilot_border),
            radiusDp = 16,
        )
        addView(context.stockText(title, 18f).apply { setTypeface(typeface, Typeface.BOLD) }, fullWidth().apply {
            bottomMargin = context.dp(12)
        })
        if (rows.isEmpty()) {
            addView(context.stockText(emptyMessage, 14f, R.color.stockpilot_muted))
        } else {
            rows.forEachIndexed { index, view ->
                addView(view, fullWidth())
                if (index != rows.lastIndex) {
                    addView(View(context).apply { setBackgroundColor(context.getColor(R.color.stockpilot_border)) }, fullWidth().apply {
                        height = context.dp(1)
                        topMargin = context.dp(10)
                        bottomMargin = context.dp(10)
                    })
                }
            }
        }
    }

    private fun row(primary: String, secondary: String, danger: Boolean = false): View = LinearLayout(context).apply {
        orientation = LinearLayout.VERTICAL
        addView(context.stockText(primary, 15f, if (danger) R.color.stockpilot_danger else R.color.stockpilot_text).apply {
            setTypeface(typeface, Typeface.BOLD)
        })
        addView(context.stockText(secondary, 13f, R.color.stockpilot_muted), fullWidth().apply {
            topMargin = context.dp(3)
        })
    }

    private fun weightedCard(start: Int = 0, end: Int = 0) = LinearLayout.LayoutParams(
        0,
        LinearLayout.LayoutParams.WRAP_CONTENT,
        1f,
    ).apply {
        marginStart = context.dp(start)
        marginEnd = context.dp(end)
    }

    private fun fullWidth() = LinearLayout.LayoutParams(
        LinearLayout.LayoutParams.MATCH_PARENT,
        LinearLayout.LayoutParams.WRAP_CONTENT,
    )

    private fun formatMoney(minor: Long, currencyCode: String): String = runCatching {
        NumberFormat.getCurrencyInstance().apply {
            currency = Currency.getInstance(currencyCode)
        }.format(minor / 100.0)
    }.getOrElse { "$currencyCode ${minor / 100.0}" }

    companion object {
        private const val MAX_ROWS = 8
    }
}
