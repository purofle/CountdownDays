package tech.archlinux.countdowndays.database

import org.jetbrains.exposed.dao.IntEntity
import org.jetbrains.exposed.dao.IntEntityClass
import org.jetbrains.exposed.dao.id.EntityID
import org.jetbrains.exposed.dao.id.IdTable
import org.jetbrains.exposed.sql.Column

object Users: IdTable<Int>() {
    override val id: Column<EntityID<Int>> = integer("telegram_id").entityId().uniqueIndex()
    override val primaryKey = PrimaryKey(id)

    val username = varchar("username", 50)
    val name = varchar("name", 50)
    val password = varchar("password", 50).nullable()
}

class User(id: EntityID<Int>): IntEntity(id) {
    companion object : IntEntityClass<User>(Users)

    var username by Users.username
    var name by Users.name
    var password by Users.password
}