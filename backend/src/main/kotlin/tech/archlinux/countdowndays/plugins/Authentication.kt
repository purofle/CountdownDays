package tech.archlinux.countdowndays.plugins

import io.ktor.server.application.*
import io.ktor.server.auth.*
import kotlinx.serialization.Serializable
import net.mamoe.yamlkt.Yaml
import tech.archlinux.countdowndays.Config
import kotlin.io.path.Path
import kotlin.io.path.isReadable
import kotlin.io.path.readText



fun Application.configureAuthentication() {

    val authToken = Config.config.bot.auth

    install(Authentication) {
        basic("auth-bot") {
            realm = "Access to the '/' path"
            validate { credentials ->
                if (credentials.name == "bot" && credentials.password == authToken) {
                    UserIdPrincipal(credentials.name)
                } else {
                    null
                }
            }
        }
    }
}