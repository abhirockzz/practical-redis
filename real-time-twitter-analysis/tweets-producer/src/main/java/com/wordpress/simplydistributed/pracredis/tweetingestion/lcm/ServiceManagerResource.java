package com.wordpress.simplydistributed.pracredis.tweetingestion.lcm;

import java.util.logging.Level;
import java.util.logging.Logger;
import javax.ws.rs.DELETE;
import javax.ws.rs.GET;
import javax.ws.rs.Path;
import javax.ws.rs.core.Response;

@Path("tweets/producer")
public class ServiceManagerResource {

    @GET
    public Response start() {
        Response r = null;
        try {
            ServiceLifecycleManager.getInstance().start();
            r = Response.ok("Tweets Producer started")
                .build();
        } catch (Exception ex) {
            Logger.getLogger(ServiceManagerResource.class.getName()).log(Level.SEVERE, null, ex);
            r = Response.serverError().build();
        }
        return r;
    }

    @DELETE
    public Response stop() {
        Response r = null;
        try {
            ServiceLifecycleManager.getInstance().stop();
            r = Response.ok("Tweets Producer stopped")
                .build();
        } catch (Exception ex) {
            Logger.getLogger(ServiceManagerResource.class.getName()).log(Level.SEVERE, null, ex);
            r = Response.serverError().build();
        }
        return r;
    }

}
