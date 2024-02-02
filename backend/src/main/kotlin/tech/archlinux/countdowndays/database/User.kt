package tech.archlinux.countdowndays.database

import org.jetbrains.exposed.dao.IntEntity
import org.jetbrains.exposed.dao.IntEntityClass
import org.jetbrains.exposed.dao.id.EntityID
import org.jetbrains.exposed.dao.id.IntIdTable
import tech.archlinux.countdowndays.response.UserResponse

object Users: IntIdTable() {
    val telegramId = long("telegram_id").uniqueIndex()
    val username = varchar("username", 50)
    val name = text("name")
    val password = varchar("password", 50).nullable()
}

class User(id: EntityID<Int>): IntEntity(id) {
    companion object : IntEntityClass<User>(Users)

    var telegramId by Users.telegramId
    var username by Users.username
    var name by Users.name
    var password by Users.password
}

fun User.Companion.findByTelegramId(telegramId: Long): User? {
    return User.find { Users.telegramId eq telegramId }.firstOrNull()
}

fun User.toResponse() = UserResponse(
        id = this.id.value,
        telegramId = this.telegramId,
        username = this.username,
        name = this.name
    )