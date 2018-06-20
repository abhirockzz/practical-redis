package com.wordpress.simplydistributed.pracredis.tweetstats.model;

import java.util.Arrays;
import java.util.List;
import java.util.Map;
import javax.xml.bind.annotation.XmlAccessType;
import javax.xml.bind.annotation.XmlAccessorType;
import javax.xml.bind.annotation.XmlRootElement;

@XmlRootElement
@XmlAccessorType(XmlAccessType.FIELD)
public class Tweet {

    private String tweeter;
    private String tweet;
    private String created_date;
    private String tweet_id;
    private List<String> keywords;

    public Tweet() {
    }

    public Tweet(String tweet_id, String tweeter, String text, String created_date, List<String> keywords) {
        this.tweet_id = tweet_id;
        this.tweeter = tweeter;
        this.tweet = text;
        this.created_date = created_date;
        this.keywords = keywords;
    }

    public String getTweeter() {
        return tweeter;
    }

    public String getTweet() {
        return tweet;
    }

    public String getTweet_id() {
        return tweet_id;
    }

    public String getCreated_date() {
        return created_date;
    }

    public List<String> getKeywords() {
        return keywords;
    }

    @Override
    public String toString() {
        return "Tweet{" + "tweeter=" + tweeter + ", tweet=" + tweet + ", created_date=" + created_date + ", tweet_id=" + tweet_id + ", hashtags=" + keywords + '}';
    }

    public static Tweet fromMap(Map<String,String> tweetFromRedis) {
        Tweet tweet = new Tweet(tweetFromRedis.get("tweet_id"),
                                tweetFromRedis.get("tweeter"),
                                tweetFromRedis.get("text"),
                                tweetFromRedis.get("created_date"),
                                Arrays.asList(tweetFromRedis.get("terms").split(",")));
        
        return tweet;
    }

}
