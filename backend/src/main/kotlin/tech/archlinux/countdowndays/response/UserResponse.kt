package tech.archlinux.countdowndays.response

import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable

@Serializable
data class UserResponse(
    val id: Int,
    @SerialName("telegram_id") val telegramId: Long,
    val username: String,
    val name: String
)