package main

import (
	"errors"
	"flag"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/DavidHernandez21/gingonic/getserver"
	"github.com/DavidHernandez21/gingonic/servershutdown"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

type Invoice struct {
	InvoiceId   int    `json:"invoiceId"`
	CustomerId  int    `json:"customerId" binding:"required,gte=0"`
	Price       int    `json:"price" binding:"required,gte=0"`
	Description string `json:"description" binding:"required"`
}

type PrintJob struct {
	JobId     int    `json:"jobId"`
	InvoiceId int    `json:"invoiceId"`
	Format    string `json:"format"`
}

func handleWriteStringError(err error, c *gin.Context) {

	log.Printf("InvoiceGenerator: %s", err.Error())
	c.JSON(500, gin.H{"error": "Unable to connect to PrinterService"})
}

func createPrintJob(invoiceId int, URL string) error {
	client := resty.New()
	var p PrintJob
	// Call PrinterService via RESTful interface
	_, err := client.R().
		SetBody(PrintJob{Format: "A4", InvoiceId: invoiceId}).
		SetResult(&p).
		Post(URL)

	if err != nil {
		log.Println("InvoiceGenerator: unable to connect PrinterService")
		return err
	}
	log.Printf("InvoiceGenerator: created print job #%v via PrinterService", p.JobId)

	return nil
}

func main() {

	host := flag.String("host", "127.0.0.1", "Bind address")
	port := flag.String("port", "8081", "Bind port")
	flag.Parse()

	router := gin.Default()
	router.SetTrustedProxies([]string{"92.244.24.19"})

	router.POST("/invoices", func(c *gin.Context) {
		var iv Invoice
		if err := c.ShouldBindJSON(&iv); err != nil {
			c.JSON(400, gin.H{"error": "Invalid input!"})
			return
		}
		log.Println("InvoiceGenerator: creating new invoice...")
		rand.Seed(time.Now().UnixNano())
		iv.InvoiceId = rand.Intn(1000)
		log.Printf("InvoiceGenerator: created invoice #%v", iv.InvoiceId)
		var printerServiceAddress strings.Builder
		printerServiceAddress.Grow(90) // we will be writing 90 bytes
		_, err := printerServiceAddress.WriteString("http://")
		if err != nil {
			handleWriteStringError(err, c)
			return
		}

		_, err = printerServiceAddress.WriteString(*host)
		if err != nil {
			handleWriteStringError(err, c)
			return
		}

		_, err = printerServiceAddress.WriteString(":")
		if err != nil {
			handleWriteStringError(err, c)
			return
		}

		_, err = printerServiceAddress.WriteString("8080")
		if err != nil {
			handleWriteStringError(err, c)
			return
		}

		_, err = printerServiceAddress.WriteString("/print-jobs")
		if err != nil {
			handleWriteStringError(err, c)
			return
		}

		err = createPrintJob(iv.InvoiceId, printerServiceAddress.String()) // Ask PrinterService to create a print job
		if err != nil {
			log.Printf("PrinterService: %s", err.Error())
			c.JSON(500, gin.H{"error": "Unable to connect to PrinterService"})
		}
		c.JSON(200, iv)
	})

	var bindAddress strings.Builder
	bindAddress.Grow(90)
	_, err := bindAddress.WriteString(*host)
	if err != nil {
		log.Fatalf("InvoiceGenerator: %s", err.Error())
		return
	}
	_, err = bindAddress.WriteString(":")
	if err != nil {
		log.Fatalf("InvoiceGenerator: %s", err.Error())
		return
	}

	_, err = bindAddress.WriteString(*port)
	if err != nil {
		log.Fatalf("InvoiceGenerator: %s", err.Error())
		return
	}

	srv := getserver.NewServer(bindAddress.String(), router)

	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}
	}()

	servershutdown.Graceful(srv, 5*time.Second)

}
