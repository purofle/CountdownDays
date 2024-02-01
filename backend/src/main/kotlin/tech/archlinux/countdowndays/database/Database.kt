package tech.archlinux.countdowndays.database

import org.jetbrains.exposed.sql.*
import org.jetbrains.exposed.sql.transactions.transaction
import tech.archlinux.countdowndays.Config

class DatabaseManager {
    init {

        val config = Config.config.database

        Database.connect(
            config.url,
            driver = config.driver,
            user = config.user,
            password = config.password
        )

        transaction {
            addLogger(StdOutSqlLogger)

            SchemaUtils.createMissingTablesAndColumns(Records)
        }
    }
}