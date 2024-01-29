package tech.archlinux.countdowndays.plugins

import io.ktor.server.application.*
import io.ktor.server.auth.*
import io.ktor.server.response.*
import io.ktor.server.routing.*

fun Application.configureRouting() {
    routing {
        authenticate("auth-bot") {
            get("/") {
                call.respondText("Hello World!")
            }
        }
    }
}
