package tech.archlinux.countdowndays.response

import kotlinx.datetime.LocalDate
import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable

@Serializable
data class RecordResponse(
    val id: Int,
    val name: String,
    val date: LocalDate,
    val description: String?,
    @SerialName("show_anniversary") val showAnniversary: Boolean,
    val owner: UserResponse
)
