package tech.archlinux.countdowndays.request

import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable

@Serializable
data class AddUserRequest(
    @SerialName("telegram_id") val telegramId: Long,
    val username: String,
    val name: String
)