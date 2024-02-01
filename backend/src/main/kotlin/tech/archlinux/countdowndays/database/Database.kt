package tech.archlinux.countdowndays.database

import kotlinx.coroutines.Dispatchers
import org.jetbrains.exposed.sql.*
import org.jetbrains.exposed.sql.transactions.experimental.newSuspendedTransaction
import org.jetbrains.exposed.sql.transactions.transaction
import tech.archlinux.countdowndays.Config

class DatabaseManager {
    init {
        transaction(db) {
            addLogger(StdOutSqlLogger)

            SchemaUtils.createMissingTablesAndColumns(Users, Records)
        }
    }

    companion object {
        val db by lazy {
            val config = Config.config.database

            Database.connect(
                config.url,
                driver = config.driver,
                user = config.user,
                password = config.password
            )
        }
    }
}

suspend fun <T> dbQuery(block: () -> T): T =
    newSuspendedTransaction(Dispatchers.IO, DatabaseManager.db) {
        addLogger(StdOutSqlLogger)
        block()
    }