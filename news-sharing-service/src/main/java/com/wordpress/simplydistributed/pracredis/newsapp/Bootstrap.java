package com.wordpress.simplydistributed.pracredis.newsapp;

import com.wordpress.simplydistributed.pracredis.api.UpvotesResource;
import com.wordpress.simplydistributed.pracredis.api.NewsItemResource;
import com.wordpress.simplydistributed.pracredis.api.CommentsResource;
import com.wordpress.simplydistributed.pracredis.api.NewUserRegistrationResource;

import java.io.IOException;
import java.net.URI;
import java.util.Optional;

import javax.ws.rs.core.UriBuilder;
import org.glassfish.grizzly.http.server.HttpServer;
import org.glassfish.jersey.grizzly2.httpserver.GrizzlyHttpServerFactory;
import org.glassfish.jersey.moxy.json.MoxyJsonFeature;
import org.glassfish.jersey.server.ResourceConfig;

public class Bootstrap {

    public static final String BASE_URI = "http://0.0.0.0:8080/";

    public static void main(String[] args) throws IOException, InterruptedException {

        String host = "0.0.0.0"; //bind to ALL interfaces
        String port = Optional.ofNullable(System.getenv("PORT")).orElse("8080");

        URI baseUri = UriBuilder.fromUri("http://" + host + "/").port(Integer.parseInt(port)).build();

        ResourceConfig config = new ResourceConfig(NewsItemResource.class,CommentsResource.class,
                                                    UpvotesResource.class,NewUserRegistrationResource.class)
                                        .register(MoxyJsonFeature.class);

        HttpServer server = GrizzlyHttpServerFactory.createHttpServer(baseUri, config);

        //gracefully exit
        Runtime.getRuntime().addShutdownHook(new Thread(new Runnable() {
            @Override
            public void run() {
                server.shutdownNow();
                System.out.println("Grizzly server shut down");
                RedisConnectionPool.getInstance().close();
            }
        }));
        server.start();
    }

}
