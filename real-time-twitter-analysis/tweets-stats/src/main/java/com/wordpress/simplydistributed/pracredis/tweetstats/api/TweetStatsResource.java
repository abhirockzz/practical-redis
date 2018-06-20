package com.wordpress.simplydistributed.pracredis.tweetstats.api;

import com.wordpress.simplydistributed.pracredis.tweetstats.RedisConnectionPool;
import com.wordpress.simplydistributed.pracredis.tweetstats.model.Tweets;
import com.wordpress.simplydistributed.pracredis.tweetstats.model.Tweet;

import java.text.ParseException;
import java.text.SimpleDateFormat;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Date;
import java.util.List;
import java.util.Set;
import java.util.UUID;
import java.util.stream.Collectors;

import javax.ws.rs.GET;
import javax.ws.rs.Path;
import javax.ws.rs.PathParam;
import javax.ws.rs.Produces;
import javax.ws.rs.QueryParam;
import javax.ws.rs.core.MediaType;
import javax.ws.rs.core.Response;

import redis.clients.jedis.Jedis;

@Path("tweets")
public class TweetStatsResource {

    @GET
    @Produces(MediaType.APPLICATION_JSON)
    public Response tweetsWithKeyword(@QueryParam("keywords") String keywords, @QueryParam("op") String op) {
        Response response = null;
        Jedis redis = null;
        String finalSetName = null;
        try {

            System.out.println("searching for tweets with keywords " + keywords);

            redis = RedisConnectionPool.getInstance().getResource();
            System.out.println("obtained conn from pool");

            System.out.println("Operation " + op);

            List<String> keywordsList = Arrays.asList(keywords.split(","));
            List<String> hashNameList = keywordsList.stream()
                    .map((keyword) -> {
                        return "keyword_tweets:" + keyword;
                    })
                    .collect(Collectors.toList());

            String[] hashNameArray = hashNameList.toArray(new String[hashNameList.size()]);

            if (op == null) {
                finalSetName = hashNameArray[0];
            } else if (op.equalsIgnoreCase("OR")) {
                finalSetName = UUID.randomUUID().toString();
                redis.sunionstore(finalSetName, hashNameArray);

            } else if (op.equalsIgnoreCase("AND")) {
                finalSetName = UUID.randomUUID().toString();
                redis.sinterstore(finalSetName, hashNameArray);

            }

            //System.out.println("searching SET " + finalSetName);
            Set<String> tweetIDs = redis.smembers(finalSetName);

            System.out.println("found " + tweetIDs.size() + " tweets with keywords " + hashNameList);
            List<Tweet> tweets = new ArrayList<>();

            for (String tweetID : tweetIDs) {
                String hashName = "tweet:" + tweetID;
                Tweet tweet = Tweet.fromMap(redis.hgetAll(hashName));
                tweets.add(tweet);
            }

            response = Response.ok(new Tweets(tweets)).build();
        } catch (Exception e) {
            System.out.println("Error - " + e.getMessage());
            response = Response.serverError().entity(e.getMessage()).build();
        } finally {
            if (op != null) {
                redis.del(finalSetName);
                //System.out.println("SET " + finalSetName + " deleted");
            }
            RedisConnectionPool.getInstance().returnResource(redis);
            System.out.println("returned to pool");
        }
        return response;
    }

    @GET
    @Path("{date}")
    @Produces(MediaType.APPLICATION_JSON)
    public Response tweetsWithKeywordOnDate(@QueryParam("keywords") String keywords, @QueryParam("op") String op, @PathParam("date") String on) {
        Response response = null;
        Jedis redis = null;
        String finalSetName = null;
        try {

            System.out.println("searching for tweets with keywords " + keywords + " on " + on);

            redis = RedisConnectionPool.getInstance().getResource();
            System.out.println("obtained conn from pool");

            System.out.println("Operation " + op);

            List<String> keywordsList = Arrays.asList(keywords.split(","));
            List<String> hashNameList = keywordsList.stream()
                    .map((keyword) -> {
                        return "keyword_tweets:" + keyword + ":" + on;
                    })
                    .collect(Collectors.toList());

            String[] hashNameArray = hashNameList.toArray(new String[hashNameList.size()]);

            if (op == null) {
                finalSetName = hashNameArray[0];
            } else if (op.equalsIgnoreCase("OR")) {
                finalSetName = UUID.randomUUID().toString();
                redis.sunionstore(finalSetName, hashNameArray);

            } else if (op.equalsIgnoreCase("AND")) {
                finalSetName = UUID.randomUUID().toString();
                redis.sinterstore(finalSetName, hashNameArray);

            }

            //System.out.println("searching SET " + finalSetName);
            Set<String> tweetIDs = redis.smembers(finalSetName);

            System.out.println("found " + tweetIDs.size() + " tweets with keywords " + hashNameList);
            List<Tweet> tweets = new ArrayList<>();

            for (String tweetID : tweetIDs) {
                String hashName = "tweet:" + tweetID;
                Tweet tweet = Tweet.fromMap(redis.hgetAll(hashName));
                tweets.add(tweet);
            }

            response = Response.ok(new Tweets(tweets)).build();
        } catch (Exception e) {
            System.out.println("Error - " + e.getMessage());
            response = Response.serverError().entity(e.getMessage()).build();
        } finally {
            if (op != null) {
                redis.del(finalSetName);
                //System.out.println("SET " + finalSetName + " deleted");
            }
            RedisConnectionPool.getInstance().returnResource(redis);
            System.out.println("returned to pool");
        }
        return response;
    }

}
