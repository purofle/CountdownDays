package tech.archlinux.countdowndays.response

import kotlinx.serialization.Serializable

@Serializable
data class UserResponse(
    val id: Int,
    val telegramId: Long,
    val username: String,
    val name: String
)