package com.wordpress.simplydistributed.pracredis.model;

import java.util.Collections;
import java.util.HashMap;
import java.util.Map;
import javax.xml.bind.annotation.XmlAccessType;
import javax.xml.bind.annotation.XmlAccessorType;
import javax.xml.bind.annotation.XmlRootElement;

@XmlRootElement
@XmlAccessorType(XmlAccessType.FIELD)
public class NewsItemSubmission {
    private String title;
    private String url;
    private String submittedBy;

    public NewsItemSubmission() {
    }
    
    public NewsItemSubmission(String title, String url) {
        this.title = title;
        this.url = url;
    }

    public String getTitle() {
        return title;
    }

    public void setTitle(String title) {
        this.title = title;
    }

    public String getUrl() {
        return url;
    }

    public void setUrl(String url) {
        this.url = url;
    }

    public String getSubmittedBy() {
        return submittedBy;
    }

    public void setSubmittedBy(String submittedBy) {
        this.submittedBy = submittedBy;
    }
    
    
    
    public Map<String,String> toMap(){
        Map<String,String> details = new HashMap<>();
        details.put("title", title);
        details.put("url", url);
        details.put("submittedBy", submittedBy);
        
        Map<String,String> ImmutableDetails = Collections.unmodifiableMap(details);
        return ImmutableDetails;
    }

    @Override
    public String toString() {
        return "NewsItemSubmission{" + "title=" + title + ", url=" + url + ", submittedBy=" + submittedBy + '}';
    }
  
    
}
