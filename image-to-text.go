package main

import (
  "net/http"
  "fmt"
  "os"
  "path/filepath"
  "log"
  "encoding/json"
  "github.com/gorilla/mux"
  "github.com/jung-kurt/gofpdf"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/s3/s3manager"
  )

type Pdfinfo struct {
  Link string 
}


func generatepdf(w http.ResponseWriter, r *http.Request) {
  pdfdetails := Pdfinfo{}
  err := json.NewDecoder(r.Body).Decode(&pdfdetails)
  // log.Println(r)
  // log.Println(r.URL.Query().Get("text"))
  text := r.URL.Query().Get("text")
  pdf := gofpdf.New("P", "mm", "A4", "")
  pdf.AddPage()
  pdf.SetFont("Arial", "B", 16)
  pdf.Cell(40, 10, text)
  pdf.OutputFileAndClose("text.pdf")

  // Aws stuff

  bucket := "avanti-dev-resources"
  filename := "text.pdf"

  file, err := os.Open(filename)
  if err != nil {
    fmt.Println("Failed to open file", filename, err)
    os.Exit(1)
  }
  defer file.Close()


  conf := aws.Config{Region: aws.String("ap-southeast-1")}
  sess := session.New(&conf)
  svc := s3manager.NewUploader(sess)
  acl := "public-read" 


  fmt.Println("Uploading file to S3...")
  result, err := svc.Upload(&s3manager.UploadInput{
      Bucket: aws.String(bucket),
      Key:    aws.String(filepath.Base(filename)),
      Body:   file,
       ACL: &acl,
  })
  if err != nil {
    fmt.Println("error", err)
    os.Exit(1)
  }
  fmt.Printf("Successfully uploaded %s to %s\n", filename, result.Location)

  pdfdetails.Link = "https://s3-ap-southeast-1.amazonaws.com/avanti-dev-resources/text.pdf"



  pdfJson, err := json.Marshal(pdfdetails)
  
  // w.Header().Set("Content-type", "application/json")
  // w.WriteHeader(http.StatusOK)
  w.Write(pdfJson)

}



func main() {
  router := mux.NewRouter()
  router.HandleFunc("/generatepdf", generatepdf).Methods("POST")
  log.Fatal(http.ListenAndServe(":8000", router))
}



