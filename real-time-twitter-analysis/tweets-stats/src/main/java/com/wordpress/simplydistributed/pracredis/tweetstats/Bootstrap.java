package com.wordpress.simplydistributed.pracredis.tweetstats;

import com.wordpress.simplydistributed.pracredis.tweetstats.api.TweetStatsResource;
import java.io.IOException;
import java.net.URI;
import java.util.Optional;

import javax.ws.rs.core.UriBuilder;
import org.glassfish.grizzly.http.server.HttpServer;
import org.glassfish.jersey.grizzly2.httpserver.GrizzlyHttpServerFactory;
import org.glassfish.jersey.moxy.json.MoxyJsonFeature;
import org.glassfish.jersey.server.ResourceConfig;

public class Bootstrap {

    public static void main(String[] args) throws IOException, InterruptedException {

        String host = "0.0.0.0"; //bind to ALL interfaces
        String port = Optional.ofNullable(System.getenv("PORT")).orElse("8080");

        URI baseUri = UriBuilder.fromUri("http://" + host + "/").port(Integer.parseInt(port)).build();

        ResourceConfig config = new ResourceConfig(TweetStatsResource.class)
                                        .register(MoxyJsonFeature.class);

        HttpServer server = GrizzlyHttpServerFactory.createHttpServer(baseUri, config);

        //gracefully exit Grizzly when app is shut down
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
