package tech.archlinux.countdowndays

import kotlinx.serialization.Serializable
import net.mamoe.yamlkt.Yaml
import kotlin.io.path.Path
import kotlin.io.path.isReadable
import kotlin.io.path.readText

object Config {
    @Serializable
    data class YamlConfig(
        val bot: Bot,
        val database: DatabaseConfig
    ) {
        @Serializable
        data class Bot(
            val auth: String
        )

        @Serializable
        data class DatabaseConfig(
            val url: String,
            val driver: String,
            val user: String,
            val password: String
        )
    }

    val config by lazy {
        val configText =
            Path("config.yaml").takeIf { it.isReadable() }?.readText() ?: throw Exception("config.yaml not found")

        Yaml().decodeFromString(YamlConfig.serializer(), configText)
    }
}