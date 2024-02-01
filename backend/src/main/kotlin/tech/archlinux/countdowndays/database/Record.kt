package tech.archlinux.countdowndays.database

import org.jetbrains.exposed.dao.IntEntity
import org.jetbrains.exposed.dao.IntEntityClass
import org.jetbrains.exposed.dao.id.EntityID
import org.jetbrains.exposed.dao.id.IntIdTable
import org.jetbrains.exposed.sql.ReferenceOption
import org.jetbrains.exposed.sql.kotlin.datetime.date

object Records : IntIdTable() {
    var name = varchar("name", 50)
    var date = date("date")
    var description = varchar("description", 255)
    var showAnniversary = bool("show_anniversary")

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