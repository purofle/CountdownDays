package tech.archlinux.countdowndays.plugins

import io.ktor.http.*
import io.ktor.server.application.*
import io.ktor.server.auth.*
import io.ktor.server.request.*
import io.ktor.server.response.*
import io.ktor.server.routing.*
import org.jetbrains.exposed.sql.transactions.transaction
import tech.archlinux.countdowndays.database.User
import tech.archlinux.countdowndays.database.dbQuery
import tech.archlinux.countdowndays.request.AddUserRequest
import tech.archlinux.countdowndays.response.UserResponse

fun Application.configureRouting() {
    routing {

        get("/user/{id}") {
            val id = call.parameters["id"] ?: return@get call.respondText(
                "Missing or malformed id",
                status = HttpStatusCode.BadRequest
            )
            val telegramId = runCatching {
                id.toInt()
            }.getOrNull() ?: return@get call.respondText("Invalid id", status = HttpStatusCode.BadRequest)

            dbQuery {
                User.findById(telegramId)
            }?.let {
                call.respond(
                    UserResponse(
                        id = it.id.value,
                        username = it.username,
                        name = it.name
                    )
                )
                return@get
            }

            call.respondText("$telegramId not in database", status = HttpStatusCode.NotFound)
        }

        authenticate("auth-bot") {
            post("/user") {
                val userRequest = call.receive<AddUserRequest>()

                dbQuery {
                    User.new(userRequest.telegramId) {
                        username = userRequest.username
                        name = userRequest.name
                    }
                }
            }
        }
    }
}
