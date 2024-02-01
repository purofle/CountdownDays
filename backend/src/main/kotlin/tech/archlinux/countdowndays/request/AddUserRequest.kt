package tech.archlinux.countdowndays.request

data class AddUserRequest(
    val telegramId: Int,
    val username: String,
    val name: String
)