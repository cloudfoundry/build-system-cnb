# `build-system-buildpack`
The Cloud Foundry Build System Buildpack is a Cloud Native Buildpack V3 that enables the building of JVM applications from source.

This buildpack is designed to work in collaboration with other buildpacks.

## Detection
The detection phase passes if

* `<APPLICATION_ROOT>/build.gradle` exists
  * Contributes `gradle` to the build plan
  * Contributes `jvm-application` to the build plan
  * Contributes `openjdk-jdk` to the build plan

* `<APPLICATION_ROOT>/pom.xml` exists
  * Contributes `maven` to the build plan
  * Contributes `jvm-application` to the build plan
  * Contributes `openjdk-jdk` to the build plan

## Build
If the build plan contains

* `gradle`
  * Links Gradle cache to a cache layer
  * Runs `<APPLICATION_ROOT>/gradlew -x test build` if `<APPLICATION_ROOT>/gradlew` exists
  * Runs `<GRADLE_ROOT>/bin/gradle -x test build` otherwise

* `maven`
  * Links Maven cache to a cache layer
  * Runs `<APPLICATION_ROOT>/mvnw -Dmaven.test.skip=true package` if `<APPLICATION_ROOT>/mvnw` exists
  * Runs `<MAVEN_ROOT>/bin/mvn -Dmaven.test.skip=true package` otherwise

## License
This buildpack is released under version 2.0 of the [Apache License][a].

[a]: http://www.apache.org/licenses/LICENSE-2.0

