package com.wordpress.simplydistributed.pracredis.api;

import com.wordpress.simplydistributed.pracredis.newsapp.RedisConnectionPool;

import javax.ws.rs.Consumes;
import javax.ws.rs.GET;
import javax.ws.rs.POST;
import javax.ws.rs.Path;
import javax.ws.rs.PathParam;
import javax.ws.rs.Produces;
import javax.ws.rs.core.MediaType;
import javax.ws.rs.core.Response;

import redis.clients.jedis.Jedis;

@Path("news/{newsid}/upvotes")
public class UpvotesResource {

    final static String REDIS_NEWS_UPVOTES_SORTED_SET = "news-upvotes";

    @POST
    @Consumes(MediaType.TEXT_PLAIN)
    public Response upvoteNewsItem(@PathParam("newsid") String newsID) {
        Response response = null;
        Jedis redis = null;
        try {
            System.out.println("upvoting news " + newsID);

            redis = RedisConnectionPool.getInstance().getResource();
            System.out.println("obtained conn from pool");

            redis.zincrby(REDIS_NEWS_UPVOTES_SORTED_SET, 1, newsID);
            System.out.println("upvoted successfully");            

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

    @GET
    @Produces(MediaType.TEXT_PLAIN)
    public Response getUpvotesForNewsItem(@PathParam("newsid") String newsID) {
        Response response = null;
        Jedis redis = null;

        try {
            redis = RedisConnectionPool.getInstance().getResource();
            System.out.println("Getting upvotes for news " + newsID);

            Double upvotes = redis.zscore(REDIS_NEWS_UPVOTES_SORTED_SET, newsID);
            
            if (upvotes == null){
                upvotes = 0.0; //leaving upvotes to null will not add it to the JSON output sent to client
            }
            
            System.out.println("News item "+ newsID + " has "+ upvotes + " upvotes");
            
            response = Response.ok(String.valueOf(upvotes.intValue())).build();
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
