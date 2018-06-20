package com.wordpress.simplydistributed.pracredis.tweetingestion;

import redis.clients.jedis.Jedis;
import redis.clients.jedis.JedisPool;
import redis.clients.jedis.JedisPoolConfig;

public final class RedisConnectionPool {

    private JedisPool pool = null;

    private RedisConnectionPool() {
        JedisPoolConfig config = new JedisPoolConfig();
        config.setTestOnBorrow(true);
        config.setMaxWaitMillis(5000);
        config.setMaxTotal(15);

        String redisHost = System.getenv().getOrDefault("REDIS_HOST", "192.168.99.100");
        String redisPort = System.getenv().getOrDefault("REDIS_PORT", "6379");
        
        pool = new JedisPool(config, redisHost, Integer.valueOf(redisPort), 10000);
        System.out.println("Jedis Pool initialized");
    }

    private final static RedisConnectionPool INSTANCE = new RedisConnectionPool();
    
    public static RedisConnectionPool getInstance(){
        return INSTANCE;
    }
    
    public Jedis getResource(){
        return pool.getResource();
    }
    
    public void returnResource(Jedis resource){
        pool.returnResource(resource);
    }
    
    public void close(){
        pool.close();
        System.out.println("Redis Connection pool shut down");
    }
}
