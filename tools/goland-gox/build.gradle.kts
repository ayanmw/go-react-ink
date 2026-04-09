plugins {
    id("java")
    id("org.jetbrains.kotlin.jvm") version "1.9.22"
    id("org.jetbrains.intellij") version "1.17.2"
}

group = "com.ayanmw"
version = "0.1.0"

repositories {
    mavenCentral()
    maven { url = uri("https://cache-redirector.jetbrains.com/intellij-dependencies") }
}

// Configure Gradle IntelliJ Plugin
intellij {
    pluginName.set("GoX - JSX for Go")
    version.set("2024.3.6")
    type.set("GO") // GoLand

    // Plugin Dependencies
    plugins.set(listOf("org.jetbrains.plugins.go"))
}

kotlin {
    jvmToolchain(21)
}

tasks {
    // Set the JVM compatibility versions
    withType<JavaCompile> {
        sourceCompatibility = "21"
        targetCompatibility = "21"
        options.encoding = "UTF-8"
    }

    withType<org.jetbrains.kotlin.gradle.tasks.KotlinCompile> {
        kotlinOptions {
            jvmTarget = "21"
            freeCompilerArgs = listOf("-Xjsr305=strict")
        }
    }

    patchPluginXml {
        sinceBuild.set("243")
        untilBuild.set("243.*")
    }

    // 签名配置 - 仅发布时需要
    // signPlugin {
    //     certificateChain.set(System.getenv("CERTIFICATE_CHAIN"))
    //     privateKey.set(System.getenv("PRIVATE_KEY"))
    //     privateKeyPassword.set(System.getenv("PRIVATE_KEY_PASSWORD"))
    // }

    // 发布配置 - 仅发布时需要
    // publishPlugin {
    //     token.set(System.getenv("PUBLISH_TOKEN"))
    // }
}

dependencies {
    implementation("org.jetbrains.kotlin:kotlin-stdlib-jdk8")
}