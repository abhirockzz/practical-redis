package com.wordpress.simplydistributed.pracredis.newsapp.api;

import com.wordpress.simplydistributed.pracredis.newsapp.RedisConnectionPool;

import javax.ws.rs.Consumes;
import javax.ws.rs.POST;
import javax.ws.rs.Path;
import javax.ws.rs.core.MediaType;
import javax.ws.rs.core.Response;

import redis.clients.jedis.Jedis;

@Path("users")
public class NewUserRegistrationResource {

    final static String REDIS_USERS_SET = "users";

    @POST
    @Consumes(MediaType.TEXT_PLAIN)
    public Response register(String username) {
        Response response = null;
        Jedis redis = null;
        try {
            System.out.println("Registering new user " + username);
            redis = RedisConnectionPool.getInstance().getResource();
 
            Long numAdded = redis.sadd(REDIS_USERS_SET, username);
            
            //SADD returns number of elemnts added. if we get 0, it means that the user elready eixsts
            if(numAdded == 0){
                System.out.println("User " + username + " already exists!");
                return Response.status(Response.Status.CONFLICT).build();
            }
            System.out.println("Added new user to SET");

            response = Response.noContent().build();
        } catch (Exception e) {
            System.out.println("Error - " + e.getMessage());
            response = Response.serverError().entity(e.getMessage()).build();
        } finally {
            RedisConnectionPool.getInstance().returnResource(redis);
            System.out.println("returned to pool");
        }
        return response;
    }

}
