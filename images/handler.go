package function

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	minio "github.com/minio/minio-go"
	"github.com/openfaas/openfaas-cloud/sdk"
)

const imgPath = "image/image.png"

func Handle(w http.ResponseWriter, r *http.Request) {
	region := regionName()

	bucketName := bucketName()

	minioClient, connectErr := connectToMinio(region)
	if connectErr != nil {
		log.Printf("S3/Minio connection error %s\n", connectErr.Error())
		os.Exit(1)
	}

	switch r.Method {
	case http.MethodPost:

		minioClient.MakeBucket(bucketName, region)

		defer r.Body.Close()
		_, err := minioClient.PutObject(bucketName,
			imgPath,
			r.Body,
			r.ContentLength,
			minio.PutObjectOptions{})

		if err != nil {
			log.Printf("error writing: %s, error: %s", imgPath, err.Error())
			w.Write([]byte(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return

	case http.MethodGet:
		obj, err := minioClient.GetObject(bucketName, imgPath, minio.GetObjectOptions{})

		if err != nil {
			log.Printf("error reading: %s, error: %s", imgPath, err.Error())
			os.Exit(1)
		}

		logBytes, _ := ioutil.ReadAll(obj)

		w.Header().Set("content-type", "image/png")
		w.Write(logBytes)
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Write([]byte(""))
	w.WriteHeader(http.StatusOK)
	return
}

func bucketName() string {
	bucketName, exist := os.LookupEnv("s3_bucket")
	if exist == false || len(bucketName) == 0 {
		bucketName = "pipeline"
		log.Printf("Bucket name not found, set to default: %v\n", bucketName)
	}
	return bucketName
}

func regionName() string {
	regionName, exist := os.LookupEnv("s3_region")
	if exist == false || len(regionName) == 0 {
		regionName = "us-east-1"
	}
	return regionName
}

func connectToMinio() (*minio.Client, error) {

	endpoint := os.Getenv("s3_url")

	tlsEnabled := tlsEnabled()

	secretKey, _ := sdk.ReadSecret("s3-secret-key")
	accessKey, _ := sdk.ReadSecret("s3-access-key")

	return minio.New(endpoint, accessKey, secretKey, tlsEnabled)
}

func tlsEnabled() bool {
	if connection := os.Getenv("s3_tls"); connection == "true" || connection == "1" {
		return true
	}
	return false
}
