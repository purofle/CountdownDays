package tech.archlinux.countdowndays.plugins

import io.ktor.server.application.*
import io.ktor.server.auth.*
import kotlinx.serialization.Serializable
import net.mamoe.yamlkt.Yaml
import kotlin.io.path.Path
import kotlin.io.path.isReadable
import kotlin.io.path.readText

@Serializable
data class Config(
    val bot: Bot
) {
    @Serializable
    data class Bot(
        val auth: String
    )
}

fun Application.configureAuthentication() {

    val configText =
        Path("config.yaml").takeIf { it.isReadable() }?.readText() ?: throw Exception("config.yaml not found")
    val authToken = Yaml().decodeFromString(Config.serializer(), configText).bot.auth

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