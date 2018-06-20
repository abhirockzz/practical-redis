package com.wordpress.simplydistributed.pracredis.tweetingestion.lcm;

import com.wordpress.simplydistributed.pracredis.tweetingestion.TwitterStreamListener;

import java.util.concurrent.atomic.AtomicBoolean;
import java.util.logging.Logger;

import twitter4j.FilterQuery;
import twitter4j.TwitterStream;
import twitter4j.TwitterStreamFactory;
import twitter4j.conf.ConfigurationBuilder;

public final class ServiceLifecycleManager {

    private static final Logger LOGGER = Logger.getLogger(ServiceLifecycleManager.class.getName());
    private static ServiceLifecycleManager INSTANCE = null;
    private final AtomicBoolean RUNNING = new AtomicBoolean(false);
    private final TwitterStream twitterStream;
    private final FilterQuery query;
    public static final String TRACKED_TERMS = System.getenv().getOrDefault("TWITTER_TRACKED_TERMS", "redis,java,golang,nosql,database");

    private ServiceLifecycleManager() {

        String _consumerKey = System.getenv().getOrDefault("TWITTER_CONSUMER_KEY", "GWHkZvPZV78MvKfEKO3wMAgHB");
        String _consumerSecret = System.getenv().getOrDefault("TWITTER_CONSUMER_SECRET", "pfferqmVGgKu0fGOEYPQTbuPjZVImFqrXToqZAWY6c5SgeywNo");
        String _accessToken = System.getenv().getOrDefault("TWITTER_ACCESS_TOKEN", "565672040-is2Bwmxvqv6ynBTp2yUJntJdvPK2vE99Uf4a8Ouu");
        String _accessTokenSecret = System.getenv().getOrDefault("TWITTER_ACCESS_TOKEN_SECRET", "uBssoJk7gvlLgUVYWRcfSQLyxinOosEvt0FWrm3ngK8nw");

        ConfigurationBuilder configurationBuilder = new ConfigurationBuilder();
        configurationBuilder.setOAuthConsumerKey(_consumerKey)
                .setOAuthConsumerSecret(_consumerSecret)
                .setOAuthAccessToken(_accessToken)
                .setOAuthAccessTokenSecret(_accessTokenSecret);

        twitterStream = new TwitterStreamFactory(configurationBuilder.build()).getInstance();
        twitterStream.addListener(new TwitterStreamListener());

        query = new FilterQuery();
        query.track(TRACKED_TERMS.split(","));
    }

    public static ServiceLifecycleManager getInstance() {
        if (INSTANCE == null) {
            INSTANCE = new ServiceLifecycleManager();
        }
        return INSTANCE;
    }

    public void start() throws Exception {
        if (RUNNING.get()) {
            throw new IllegalStateException("Service is already running");
        }
        twitterStream.filter(query);

        LOGGER.info("Started Tweets Producer thread");
        RUNNING.set(true);
    }

    public void stop() throws Exception {
        if (!RUNNING.get()) {
            throw new IllegalStateException("Service is NOT running. Cannot stop");
        }
        twitterStream.shutdown();
        LOGGER.info("Stopped Tweet Producer thread");
        RUNNING.set(false);
    }

}
