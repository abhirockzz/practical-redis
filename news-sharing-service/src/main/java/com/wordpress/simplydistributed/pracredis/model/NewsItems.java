package com.wordpress.simplydistributed.pracredis.model;

import java.util.ArrayList;
import java.util.List;
import javax.xml.bind.annotation.XmlAccessType;
import javax.xml.bind.annotation.XmlAccessorType;
import javax.xml.bind.annotation.XmlRootElement;

@XmlRootElement
@XmlAccessorType(XmlAccessType.FIELD)
public class NewsItems {
    private final List<NewsItem> newsItems;

    public NewsItems() {
        this.newsItems = new ArrayList<>();
    }
    
    public void add(NewsItem newsItem){
        this.newsItems.add(newsItem);
    }
    
    public int size(){
        return this.newsItems.size();
    }
}
