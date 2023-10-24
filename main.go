package main

import (
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/gocolly/colly"
    "sync"
    "time"
    "net/http"
)

type CrawlData struct {
    link         string
    LastCrawled  time.Time
    IsPayingUser bool
}

var dataMutex sync.Mutex
var crawlDatas = make(map[string]CrawlData)
var maxRetryAttempts = 3

// Define separate queues for paying and non-paying customers
var payingQueue = make(chan string, 5)  // Up to 5 concurrent crawlers for paying customers
var nonPayingQueue = make(chan string, 2)  // Up to 2 concurrent crawlers for non-paying customers

func main() {
    r := gin.Default()

    r.LoadHTMLGlob("./index.html")

    r.POST("/crawl", func(c *gin.Context) {
        url := c.DefaultQuery("url", "")
        customerType := c.DefaultQuery("customerType", "")

        result, err := crawlPage(url, customerType)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.String(http.StatusOK, result)
    })

    r.Run(":8080")
}

func crawlURL(url string, isPaying bool) {
    if isPaying {
        payingQueue <- url  // Add the URL to the paying customer queue
    } else {
        nonPayingQueue <- url  // Add the URL to the non-paying customer queue
    }
    go crawlWorker(url, isPaying)
}

func crawlWorker(url string, isPaying bool) {
    c := colly.NewCollector()

    c.OnRequest(func(r *colly.Request) {
        fmt.Println("Visiting: ", r.URL)
    })

    c.OnError(func(_ *colly.Response, err error) {
        fmt.Println("Something went wrong: ", err)
        // Implement retry logic here
        if isPaying {
            payingQueue <- url
        } else {
            nonPayingQueue <- url
        }
    })

    c.OnResponse(func(r *colly.Response) {
        fmt.Println("Page visited: ", r.Request.URL)
    })

    c.OnHTML("a[href]", func(e *colly.HTMLElement) {
        crawlData := CrawlData{}
        crawlData.link = e.ChildAttr("a", "href")
        crawlDatas[crawlData.link] = crawlData
        fmt.Println("Scraped!")
    })

    c.OnScraped(func(r *colly.Response) {
        // Add your storage logic here.
        fmt.Println(r.Request.URL, "Scraped!")
    })

    c.Visit(url)
}

func crawlPage(url, customerType string) (string, error) {
    // Check if the URL has been crawled in the last 60 minutes
    dataMutex.Lock()
    crawlData, exists := crawlDatas[url]
    dataMutex.Unlock()

    if exists && time.Since(crawlData.LastCrawled) < time.Minute*60 {
        return fmt.Sprintf("URL: %s was crawled %v minutes ago. Content: %s", url, time.Since(crawlData.LastCrawled).Minutes(), crawlData.link), nil
    }

    // Implement your crawling logic here
    var isPaying = customerType == "paying"

    // Check the priority queue based on customer type
    var queue chan string
    if isPaying {
        queue = payingQueue
    } else {
        queue = nonPayingQueue
    }

    // Check if the URL is already in the queue to avoid duplicate crawling
    inQueue := false
    for {
        select {
        case item := <-queue:
            if item == url {
                inQueue = true
            }
            // You can process other items in the queue if needed.
        default:
            // The queue is empty or the item was not found.
            break
        }
    }

    if !inQueue {
        queue <- url
    }

    return fmt.Sprintf("Crawling URL: %s\nCustomer Type: %s", url, customerType), nil
}
