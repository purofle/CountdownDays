package tech.archlinux.countdowndays

import io.ktor.server.application.*
import tech.archlinux.countdowndays.plugins.*

fun main(args: Array<String>) {
    io.ktor.server.netty.EngineMain.main(args)
}

fun Application.module() {
    configureHTTP()
    configureMonitoring()
    configureSerialization()
    configureRouting()
}
