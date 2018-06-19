package com.wordpress.simplydistributed.pracredis.newsapp.model;

import javax.xml.bind.annotation.XmlAccessType;
import javax.xml.bind.annotation.XmlAccessorType;
import javax.xml.bind.annotation.XmlRootElement;

@XmlRootElement
@XmlAccessorType(XmlAccessType.FIELD)
public class NewsItem {

    private String newsID;
    private String url;
    private String title;
    private String submittedBy;
    private String numUpvotes;
    private String numComments;

    public NewsItem() {
    }

    public NewsItem(String newsID, String url, String title, String submittedBy, String numUpvotes, String numComments) {
        this.newsID = newsID;
        this.url = url;
        this.title = title;
        this.submittedBy = submittedBy;
        this.numUpvotes = numUpvotes;
        this.numComments = numComments;
    }

    public String getNewsID() {
        return newsID;
    }

    public void setNewsID(String newsID) {
        this.newsID = newsID;
    }

    public String getUrl() {
        return url;
    }

    public void setUrl(String url) {
        this.url = url;
    }

    public String getTitle() {
        return title;
    }

    public void setTitle(String title) {
        this.title = title;
    }

    public String getSubmittedBy() {
        return submittedBy;
    }

    public void setSubmittedBy(String submittedBy) {
        this.submittedBy = submittedBy;
    }

    public String getNumUpvotes() {
        return numUpvotes;
    }

    public void setNumUpvotes(String numUpvotes) {
        this.numUpvotes = numUpvotes;
    }

    public String getNumComments() {
        return numComments;
    }

    public void setNumComments(String numComments) {
        this.numComments = numComments;
    }

    @Override
    public String toString() {
        return "NewsItem{" + "newsID=" + newsID + ", url=" + url + ", title=" + title + ", submittedBy=" + submittedBy + ", numUpvotes=" + numUpvotes + ", numComments=" + numComments + '}';
    }

}
