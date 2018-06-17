package com.wordpress.simplydistributed.pracredis.api;

import com.wordpress.simplydistributed.pracredis.newsapp.RedisConnectionPool;
import com.wordpress.simplydistributed.pracredis.model.NewsItems;
import com.wordpress.simplydistributed.pracredis.model.NewsItemSubmission;
import com.wordpress.simplydistributed.pracredis.model.NewsItem;

import java.util.Map;
import java.util.Set;

import javax.ws.rs.Consumes;
import javax.ws.rs.GET;
import javax.ws.rs.POST;
import javax.ws.rs.Path;
import javax.ws.rs.PathParam;
import javax.ws.rs.Produces;
import javax.ws.rs.core.MediaType;
import javax.ws.rs.core.Response;

import redis.clients.jedis.Jedis;


@Path("news")
public class NewsItemResource {

    final static String REDIS_NEWS_HASH_PREFIX = "news";
    final static String REDIS_NEWS_ID_COUNTER = "news-id-counter";
    final static String REDIS_NEWS_UPVOTES_SORTED_SET = "news-upvotes";


    @POST
    @Consumes(MediaType.APPLICATION_JSON)
    public Response postNewsItem(NewsItemSubmission news) {
        Response response = null;
        Jedis redis = null;
        try {
            System.out.println("trying to store news " + news);
            redis = RedisConnectionPool.getInstance().getResource();
            System.out.println("obtained conn from pool");
            
            //generate unique news ID
            Long newsID = redis.incr(REDIS_NEWS_ID_COUNTER);
            String newsHash = REDIS_NEWS_HASH_PREFIX + ":" + newsID;

            //store details in hash
            redis.hmset(newsHash, news.toMap());
            System.out.println("stored news in HASH successfully");

            //add to upvotes sorted set
            redis.zadd(REDIS_NEWS_UPVOTES_SORTED_SET, 0, String.valueOf(newsID));
            System.out.println("added to upvotes sorted set successfully");

            response = Response.ok(newsID).build();
        } catch (Exception e) {
            System.out.println("Error - " + e.getMessage());
            response = Response.serverError().entity(e.getMessage()).build();
        } finally {
            RedisConnectionPool.getInstance().returnResource(redis);
            System.out.println("returned to pool");
        }
        return response;
    }

    @Path("{id}")
    @GET
    @Produces(MediaType.APPLICATION_JSON)
    public Response getNewsItem(@PathParam("id") String newsID) {
        Response response = null;
        NewsItem news = null;

        try {
            news = getDetailsForNewsItem(newsID);
            if (news == null) {
                response = Response.noContent().build();
            } else {
                response = Response.ok(getDetailsForNewsItem(newsID)).build();
            }

        } catch (Exception e) {
            System.out.println("Error - " + e.getMessage());
            response = Response.serverError().entity(e.getMessage()).build();
        }

        return response;
    }

    //ALL news items sorted by upvotes (desc)
    @GET
    @Produces(MediaType.APPLICATION_JSON)
    public Response getNewsItems() {
        Response response = null;
        Jedis redis = null;

        try {
            System.out.println("Getting all news items sorted by upvotes...");
            redis = RedisConnectionPool.getInstance().getResource();

            Long numItems = redis.zcard(REDIS_NEWS_UPVOTES_SORTED_SET); //we need sorted set size to find limit for zrevrange
            Set<String> newsIDs = redis.zrevrange(REDIS_NEWS_UPVOTES_SORTED_SET, 0, numItems - 1);

            NewsItems allNewsItems = new NewsItems();

            for (String newsID : newsIDs) {
                NewsItem newsItem = getDetailsForNewsItem(newsID);
                allNewsItems.add(newsItem);
            }

            System.out.println("returning " + allNewsItems.size() + " news items");
            response = Response.ok(allNewsItems).build();
        } catch (Exception e) {
            System.out.println("Error - " + e.getMessage());
            response = Response.serverError().entity(e.getMessage()).build();
        } finally {
            RedisConnectionPool.getInstance().returnResource(redis);
            System.out.println("returned to pool");
        }

        return response;
    }

    private NewsItem getDetailsForNewsItem(String newsID) {
        NewsItem news = null;
        Jedis redis = null;

        try {
            String newsHash = REDIS_NEWS_HASH_PREFIX + ":" + newsID;
            System.out.println("Searching for news " + newsHash);

            redis = RedisConnectionPool.getInstance().getResource();

            Map<String, String> basicNewsDetails = redis.hgetAll(newsHash); //get basic details
            System.out.println("basic details " + basicNewsDetails);

            if (basicNewsDetails.isEmpty()) {
                System.out.println("no info for news ID " + newsID);
                return null;
            }
            Double upvotes = redis.zscore(REDIS_NEWS_UPVOTES_SORTED_SET, newsID); //get upvotes
            //get comments
            String commentListName = "news:" + newsID + ":comments";
            System.out.println("Getting comments for news " + newsID);
            Long length = redis.llen(commentListName); //get length
            System.out.println("no. of comments " + length);
            
            news = new NewsItem(newsID, basicNewsDetails.get("url"),
                    basicNewsDetails.get("title"),
                    basicNewsDetails.get("submittedBy"),
                    String.valueOf(upvotes.intValue()),
                    String.valueOf(length));

            System.out.println("Got news " + news);

        } catch (Exception e) {
            System.out.println("Error - " + e.getMessage());
        } finally {
            RedisConnectionPool.getInstance().returnResource(redis);
            System.out.println("returned to pool");
        }

        return news;
    }

}
