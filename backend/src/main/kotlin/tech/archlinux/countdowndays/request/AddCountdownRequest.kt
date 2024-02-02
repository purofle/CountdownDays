package tech.archlinux.countdowndays.request

import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable

@Serializable
data class AddCountdownRequest(
    @SerialName("telegram_id") val telegramId: Long,
    val name: String,
    val date: String,
)