package tech.archlinux.countdowndays.response

import kotlinx.datetime.LocalDate
import kotlinx.serialization.Serializable

@Serializable
data class RecordResponse(
    val id: Int,
    val name: String,
    val date: LocalDate,
    val description: String?,
    val showAnniversary: Boolean,
    val owner: UserResponse
)
