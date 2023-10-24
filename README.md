# Web Crawling Application

This is a web crawling application built using Go and the Gin web framework. The application allows users to enter a URL, specify their customer type, and initiate real-time web crawling. It also provides a basic queuing system for managing concurrent crawling for different customer types and checks for recently crawled pages for optimization.

## URL and Query Parameters

- **URL**: This is where you enter the URL you want to crawl. The URL parameter is used to specify the web page you want to crawl.

- **Customer Type**: You can specify your customer type (e.g., paying or non-paying) using the dropdown menu. The `customerType` parameter helps the application prioritize crawling for different customer types.

  **Example URL**:
http://localhost:8080/crawl?url=https://example.com&customerType=paying


In this example, the URL is "https://example.com," and the customer type is "paying."

## Dependencies

The following Go packages and libraries are used in this project:

- [Gin](https://github.com/gin-gonic/gin) - The web framework for building the application.
- [Gocolly](https://github.com/gocolly/colly) - A package for web scraping and crawling.
