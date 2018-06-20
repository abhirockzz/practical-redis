package com.wordpress.simplydistributed.pracredis.tweetingestion;

import java.text.SimpleDateFormat;
import java.util.Date;
import java.util.List;

public class TweetInfo {

    private String tweetID;
    private String tweeter;
    private String tweet;
    private List<String> terms;
    private String createdDate;
    

    public TweetInfo() {
    }

    public TweetInfo(String tweeter, String tweet, Date created, String tweet_id, List<String> terms) {
        this.tweeter = tweeter;
        this.tweet = tweet;
        this.createdDate = new SimpleDateFormat("dd-MM-yyyy").format(created);
        this.tweetID = tweet_id;
        this.terms = terms;
    }

    public String getTweetID() {
        return tweetID;
    }

    public void setTweetID(String tweetID) {
        this.tweetID = tweetID;
    }

    public String getTweeter() {
        return tweeter;
    }

    public void setTweeter(String tweeter) {
        this.tweeter = tweeter;
    }

    public String getTweet() {
        return tweet;
    }

    public void setTweet(String tweet) {
        this.tweet = tweet;
    }

    public List<String> getTerms() {
        return terms;
    }

    public void setTerms(List<String> terms) {
        this.terms = terms;
    }

    public String getCreatedDate() {
        return createdDate;
    }

    public void setCreatedDate(String createdDate) {
        this.createdDate = createdDate;
    }

    @Override
    public String toString() {
        return "TweetInfo{" + "tweetID=" + tweetID + ", tweeter=" + tweeter + ", tweet=" + tweet + ", terms=" + terms + ", createdDate=" + createdDate + '}';
    }


    
}
