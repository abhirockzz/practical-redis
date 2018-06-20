package com.wordpress.simplydistributed.pracredis.tweetingestion;

import com.wordpress.simplydistributed.pracredis.tweetingestion.lcm.ServiceLifecycleManager;
import com.wordpress.simplydistributed.pracredis.tweetingestion.lcm.ServiceManagerResource;

import java.io.IOException;
import java.net.URI;
import java.util.Optional;
import java.util.logging.Level;
import java.util.logging.Logger;
import javax.ws.rs.core.UriBuilder;

import org.glassfish.grizzly.http.server.HttpServer;
import org.glassfish.jersey.grizzly2.httpserver.GrizzlyHttpServerFactory;
import org.glassfish.jersey.server.ResourceConfig;

public class Bootstrap {

    private static final Logger LOGGER = Logger.getLogger(Bootstrap.class.getName());

    private static void bootstrap() throws IOException {

        String hostname = "0.0.0.0"; //bind all avaialble interfaces
        String port = Optional.ofNullable(System.getenv("PORT")).orElse("8080");

        URI baseUri = UriBuilder.fromUri("http://" + hostname + "/").port(Integer.parseInt(port)).build();

        ResourceConfig config = new ResourceConfig(ServiceManagerResource.class);

        HttpServer server = GrizzlyHttpServerFactory.createHttpServer(baseUri, config);
        LOGGER.log(Level.INFO, "Application accessible at {0}", baseUri.toString());

        //gracefully exit Grizzly services when app is shut down
        Runtime.getRuntime().addShutdownHook(new Thread(new Runnable() {
            @Override
            public void run() {
                LOGGER.log(Level.INFO, "Exiting......");
                try {
                    server.shutdownNow();
                    LOGGER.log(Level.INFO, "REST services stopped");
                    
                    ServiceLifecycleManager.getInstance().stop();
                    LOGGER.log(Level.INFO, "Twitter producer stopped");
                    
                } catch (Exception ex) {
                    //log & continue....
                    LOGGER.log(Level.SEVERE, ex, ex::getMessage);
                }

            }
        }));
        server.start();

    }


    public static void main(String[] args) throws Exception {

        bootstrap();

    }
}
