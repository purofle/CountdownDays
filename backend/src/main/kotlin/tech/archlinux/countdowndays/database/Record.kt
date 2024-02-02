package tech.archlinux.countdowndays.database

import org.jetbrains.exposed.dao.IntEntity
import org.jetbrains.exposed.dao.IntEntityClass
import org.jetbrains.exposed.dao.id.EntityID
import org.jetbrains.exposed.dao.id.IntIdTable
import org.jetbrains.exposed.sql.ReferenceOption
import org.jetbrains.exposed.sql.kotlin.datetime.date
import tech.archlinux.countdowndays.response.RecordResponse

object Records : IntIdTable() {
    var name = varchar("name", 50)
    var date = date("date")
    var description = varchar("description", 255).nullable().default(null)
    var showAnniversary = bool("show_anniversary").default(false)

    var owner = reference("owner", Users, onDelete = ReferenceOption.CASCADE)
}

class Record(id: EntityID<Int>) : IntEntity(id) {
    companion object : IntEntityClass<Record>(Records)

    var name by Records.name
    var date by Records.date
    var description by Records.description
    var showAnniversary by Records.showAnniversary

    var owner by User referencedOn Records.owner
}

fun Record.toResponse() = RecordResponse(
    id = id.value,
    name = name,
    date = date,
    description = description,
    showAnniversary = showAnniversary,
    owner = owner.toResponse()
)