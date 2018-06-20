package com.wordpress.simplydistributed.pracredis.tweetingestion;

import com.fasterxml.jackson.databind.ObjectMapper;

import com.wordpress.simplydistributed.pracredis.tweetingestion.lcm.ServiceLifecycleManager;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.HashSet;
import java.util.List;
import java.util.Set;
import java.util.logging.Level;
import java.util.logging.Logger;
import java.util.stream.Collectors;

import redis.clients.jedis.Jedis;

import twitter4j.StallWarning;
import twitter4j.Status;
import twitter4j.StatusDeletionNotice;
import twitter4j.StatusListener;

public class TwitterStreamListener implements StatusListener {

    static final ObjectMapper MAPPER = new ObjectMapper();
    static final String TWEETS_LIST = "tweets";

    @Override
    public void onStatus(Status status) {

        List<String> matchedTerms = new ArrayList<>(getMatchedTerms(status.getText()));

        if (!status.isPossiblySensitive() && !matchedTerms.isEmpty()) {
            //System.out.println("Tweet text\n" + status.getText());
            System.out.println("Matched terms  - " + matchedTerms);

            TweetInfo tweet = new TweetInfo(status.getUser().getScreenName(),
                    status.getText(),
                    status.getCreatedAt(),
                    String.valueOf(status.getId()),
                    matchedTerms);

            System.out.println(tweet);

            String jsonTweet = null;
            Jedis redis = null;
            try {
                jsonTweet = MAPPER.writeValueAsString(tweet);
                System.out.println("json tweet " + jsonTweet);

                redis = RedisConnectionPool.getInstance().getResource();
                redis.lpush(TWEETS_LIST, jsonTweet);

                System.out.println("pushed tweet to Redis List");

            } catch (Exception ex) {
                Logger.getLogger(TwitterStreamListener.class.getName()).log(Level.SEVERE, null, ex);
            } finally {
                RedisConnectionPool.getInstance().returnResource(redis);
            }

        }

    }

    private Set<String> getMatchedTerms(String tweetText) {
        List<String> tweetWords = Arrays.asList(tweetText.split(" "));
        List<String> matchedTerms
                = tweetWords.stream()
                        .map((word) -> word.toLowerCase())
                        .filter((word) -> Arrays.asList(ServiceLifecycleManager.TRACKED_TERMS.split(",")).contains(word))
                        .collect(Collectors.toList());

        return new HashSet<>(matchedTerms);
    }

    @Override
    public void onDeletionNotice(StatusDeletionNotice statusDeletionNotice) {
    }

    @Override
    public void onTrackLimitationNotice(int numberOfLimitedStatuses) {
    }

    @Override
    public void onScrubGeo(long userId, long upToStatusId) {
    }

    @Override
    public void onStallWarning(StallWarning warning) {
    }

    @Override
    public void onException(Exception ex) {
        ex.printStackTrace();
    }
}
