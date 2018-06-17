package com.wordpress.simplydistributed.pracredis.model;

import java.util.List;
import javax.xml.bind.annotation.XmlAccessType;
import javax.xml.bind.annotation.XmlAccessorType;
import javax.xml.bind.annotation.XmlRootElement;

@XmlRootElement
@XmlAccessorType(XmlAccessType.FIELD)
public class Comments {
    
    private String newsID;
    private List<String> comments;

    public Comments() {
    }

    public Comments(String newsID, List<String> comments) {
        this.newsID = newsID;
        this.comments = comments;
    }

    public String getNewsID() {
        return newsID;
    }

    public void setNewsID(String newsID) {
        this.newsID = newsID;
    }

    public List<String> getComments() {
        return comments;
    }

    public void setComments(List<String> comments) {
        this.comments = comments;
    }

    @Override
    public String toString() {
        return "Comments{" + "newsID=" + newsID + ", comments=" + comments + '}';
    }
    
    
}
