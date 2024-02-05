package tech.archlinux.countdowndays.plugins

import io.ktor.http.*
import io.ktor.server.application.*
import io.ktor.server.auth.*
import io.ktor.server.request.*
import io.ktor.server.response.*
import io.ktor.server.routing.*
import kotlinx.datetime.LocalDate
import org.jetbrains.exposed.sql.SortOrder
import tech.archlinux.countdowndays.database.*
import tech.archlinux.countdowndays.request.AddCountdownRequest
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
                val request = call.receive<AddUserRequest>()

                val newUser = dbQuery {
                    User.new {
                        telegramId = request.telegramId
                        username = request.username
                        name = request.name
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

            // add countdown
            post("/countdown") {
                val request = call.receive<AddCountdownRequest>()

                val user = dbQuery {
                    User.findByTelegramId(request.telegramId)
                } ?: return@post call.respondText(
                    "User not found",
                    status = HttpStatusCode.NotFound
                )

                // String to LocalDate
                // 1989-06-04
                val localDate = LocalDate.parse(request.date)

                val newCountdown = dbQuery {
                    Record.new {
                        name = request.name
                        date = localDate
                        owner = user
                    }
                }

                return@post call.respond(newCountdown.date)
            }

            // get all countdowns for a user
            get("/countdown/{id}/all") {
                val telegramId = runCatching {
                    call.parameters["id"]!!.toLong()
                }.getOrNull() ?: return@get call.respondText(
                    "Missing or invalid id",
                    status = HttpStatusCode.BadRequest
                )

                val user = dbQuery {
                    User.findByTelegramId(telegramId)
                } ?: return@get call.respondText(
                    "User not found",
                    status = HttpStatusCode.NotFound
                )

                val countdowns = dbQuery {
                    Record.find { Records.owner eq user.id }
                        .limit(20)
                        .orderBy(Records.date to SortOrder.ASC)
                        .map { it.toResponse() }
                }

                call.respond(countdowns)
            }

            delete("/countdown/{id}") {
                val countdownId = runCatching {
                    call.parameters["id"]!!.toInt()
                }.getOrNull() ?: return@delete call.respondText(
                    "Missing or invalid id",
                    status = HttpStatusCode.BadRequest
                )

                val countdown = dbQuery {
                    Record.findById(countdownId)
                }

                dbQuery {
                    countdown?.delete()
                } ?: return@delete call.respondText(
                    "Countdown not found",
                    status = HttpStatusCode.NotFound
                )

                call.respondText("Countdown deleted", status = HttpStatusCode.OK)
            }
        }
    }
}
