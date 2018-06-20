FROM maven:3.5-jdk-8-slim
RUN mkdir /app
WORKDIR /app
COPY src /app/src
COPY pom.xml /app
RUN mvn clean install
CMD ["java","-jar", "target/tweets-stats.jar"]