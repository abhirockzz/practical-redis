package com.wordpress.simplydistributed.pracredis.newsapp.api;

import com.wordpress.simplydistributed.pracredis.newsapp.model.Comments;
import com.wordpress.simplydistributed.pracredis.newsapp.RedisConnectionPool;

import java.util.List;

import javax.ws.rs.Consumes;
import javax.ws.rs.GET;
import javax.ws.rs.POST;
import javax.ws.rs.Path;
import javax.ws.rs.PathParam;
import javax.ws.rs.Produces;
import javax.ws.rs.core.MediaType;
import javax.ws.rs.core.Response;

import redis.clients.jedis.Jedis;

@Path("news/{newsid}/comments")
public class CommentsResource {

    @POST
    @Consumes(MediaType.TEXT_PLAIN)
    public Response postCommentsForNewsItem(@PathParam("newsid") String newsID, String comment) {
        Response response = null;
        Jedis redis = null;
        try {
            System.out.println("posting comment " + comment + " for news " + newsID);

            redis = RedisConnectionPool.getInstance().getResource();
            System.out.println("obtained conn from pool");

            String commentListName = "news:" + newsID + ":comments";
            redis.lpush(commentListName, comment);
            System.out.println("comment stored news in LIST " + commentListName + " successfully");

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
    @Produces(MediaType.APPLICATION_JSON)
    public Response getCommentsForNewsItem(@PathParam("newsid") String newsID) {
        Response response = null;
        Jedis redis = null;

        try {
            redis = RedisConnectionPool.getInstance().getResource();

            String commentListName = "news:" + newsID + ":comments";
            System.out.println("Getting comments for news " + newsID);

            Long length = redis.llen(commentListName); //get length
            System.out.println("no. of comments " + length);

            List<String> comments = null;
            if (length > 0) {
                System.out.println("getting all comments from list " + commentListName);
                comments = redis.lrange(commentListName, 0, length);
            }
            
            response = Response.ok(new Comments(newsID, comments)).build();
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
