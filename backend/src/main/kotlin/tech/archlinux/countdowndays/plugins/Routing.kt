package tech.archlinux.countdowndays.plugins

import io.ktor.http.*
import io.ktor.server.application.*
import io.ktor.server.auth.*
import io.ktor.server.request.*
import io.ktor.server.response.*
import io.ktor.server.routing.*
import org.jetbrains.exposed.sql.SqlExpressionBuilder.eq
import tech.archlinux.countdowndays.database.*
import tech.archlinux.countdowndays.request.AddUserRequest
import tech.archlinux.countdowndays.response.UserResponse

fun Application.configureRouting() {
    routing {

        // get user public information
        get("/user/{id}") {

            val telegramId = runCatching {
                call.parameters["id"]!!.toLong()
            }.getOrNull() ?: return@get call.respondText(
                "Missing or invalid id",
                status = HttpStatusCode.BadRequest
            )

            dbQuery {
                User.findByTelegramId(telegramId)
            }?.let {
                return@get call.respond(it.toResponse())
            }

            call.respondText("$telegramId not in database", status = HttpStatusCode.NotFound)
        }

        authenticate("auth-bot") {
            // add user
            post("/user") {
                val userRequest = call.receive<AddUserRequest>()

                val newUser = dbQuery {
                    User.new {
                        telegramId = userRequest.telegramId
                        username = userRequest.username
                        name = userRequest.name
                    }
                }

                return@post call.respond(
                    UserResponse(
                        id = newUser.id.value,
                        username = newUser.username,
                        name = newUser.name,
                        telegramId = newUser.telegramId
                    )
                )
            }
        }
    }
}
