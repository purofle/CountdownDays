package tech.archlinux.countdowndays

import io.ktor.server.application.*
import tech.archlinux.countdowndays.database.DatabaseManager
import tech.archlinux.countdowndays.plugins.*

fun main(args: Array<String>) {
    io.ktor.server.netty.EngineMain.main(args)
}

fun Application.module() {

    DatabaseManager()

    configureHTTP()
    configureMonitoring()
    configureSerialization()
    configureAuthentication()
    configureRouting()
}
