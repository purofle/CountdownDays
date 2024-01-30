package tech.archlinux.countdowndays.plugins

import io.ktor.server.application.*
import io.ktor.server.auth.*
import io.ktor.server.request.*
import io.ktor.server.response.*
import io.ktor.server.routing.*

fun Application.configureRouting() {
    routing {
        authenticate("auth-bot") {
            post("/add") {
                call.receive()
            }
            get("/all_countdown") {

            }
            get("/countdown/{id}") {

            }
        }
    }
}
