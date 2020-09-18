package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"intelligent-analysis-of-car-sensors-backend/models"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"strconv"
)

// Client represents the entity that store the configuration used for the request.
type Client struct {
	BaseURL string
}

// ClientRequest stores the entities used in the request.
type ClientRequest struct {
	Dataset          *models.Dataset
	ClustersNumber   string
	ComponentsNumber string
}

// ClassifySVMRequest stores the entities used in the request
// for ClassifySVM.
type ClassifySVMRequest struct {
	Dataset           *models.Dataset
	DatasetToClassify *multipart.FileHeader
}

// PredictionRequest stores the entities used in the LSTM request.
type PredictionRequest struct {
	Dataset *models.Dataset
	Feature string
	Epochs  int
}

// NewClient returns a new entity of Client.
func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
	}
}

// Result represents the entity that store the results from k-means and SVM.
type Result struct {
	models.Kmeans
	SVMPlot string `json:"svmPlot"`
}

// ResultSVM represents the entity that store the results from the SVM classification.
type ResultSVM struct {
	ClassificationList string `json:"classificationList"`
}

// ResultPrediction represents the entity that store the results from the LSTM prediction.
type ResultPrediction struct {
	models.Prediction
}

// Start performs a POST request to ia service sending a CSV file
// created from the Dataset entity and aplying k-means and SVM to the dataset.
func (c *Client) Start(clientRequest *ClientRequest) (*Result, error) {
	url := fmt.Sprintf(c.BaseURL + "/pca?components-number=" + clientRequest.ComponentsNumber + "&clusters-number=" + clientRequest.ClustersNumber)

	// Create CSV file
	csvFilename, err := models.CreateCSVFromDatasetEntity(clientRequest.Dataset)
	if err != nil {
		return nil, err
	}

	// Open file
	csvFile, err := os.Open(csvFilename)
	if err != nil {
		return nil, err
	}

	// Create form-data
	body := &bytes.Buffer{}
	writter := multipart.NewWriter(body)
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "csv", csvFile.Name()))
	h.Set("Content-Type", "text/csv")
	part, err := writter.CreatePart(h)
	if err != nil {
		return nil, err
	}

	// TODO: Check first return value from Copy method
	// Copy CSV file to form-data
	_, err = io.Copy(part, csvFile)
	if err != nil {
		return nil, err
	}
	writter.Close()

	// Create POST request
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", writter.FormDataContentType())

	// Make request
	data, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	// Get response
	result := &Result{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	// Delete csv file
	err = os.Remove(csvFilename)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// ClassifySVM performs a POST request to ia service sending a CSV file
// created from the Dataset entity and the dataset received in the POST request.
// It applies SMV using the dataset created as training data and the dataset
// received as data to classify (test).
func (c *Client) ClassifySVM(clientRequest *ClassifySVMRequest) (*ResultSVM, error) {
	url := fmt.Sprintf(c.BaseURL + "/svm?dataset-rows-number=" + strconv.Itoa(clientRequest.Dataset.RowsNumber))

	// Create CSV file
	csvFilename, err := models.CreateCSVFromDatasetEntitySVM(clientRequest.Dataset, clientRequest.DatasetToClassify)
	if err != nil {
		return nil, err
	}

	// Open file
	mergedCSVFile, err := os.Open(csvFilename)
	if err != nil {
		return nil, err
	}

	// Create form-data
	body := &bytes.Buffer{}
	writter := multipart.NewWriter(body)
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "csv", mergedCSVFile.Name()))
	h.Set("Content-Type", "text/csv")
	part, err := writter.CreatePart(h)
	if err != nil {
		return nil, err
	}

	// TODO: Check first return value from Copy method
	// Copy CSV file to form-data
	_, err = io.Copy(part, mergedCSVFile)
	if err != nil {
		return nil, err
	}
	writter.Close()

	// Create POST request
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", writter.FormDataContentType())

	// Make request
	data, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	// Get response
	result := &ResultSVM{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	// Delete csv file
	err = os.Remove(csvFilename)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Prediction performs a predictions over the LSTM neural network.
func (c *Client) Prediction(clientRequest *PredictionRequest) (*ResultPrediction, error) {
	url := fmt.Sprintf(c.BaseURL + "/predict?feature=" + clientRequest.Feature + "&epochs=" + strconv.Itoa(clientRequest.Epochs))

	// Create CSV file
	csvFilename, err := models.CreateCSVFromDatasetEntity(clientRequest.Dataset)
	if err != nil {
		return nil, err
	}

	// Open file
	csvFile, err := os.Open(csvFilename)
	if err != nil {
		return nil, err
	}

	// Create form-data
	body := &bytes.Buffer{}
	writter := multipart.NewWriter(body)
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "csv", csvFile.Name()))
	h.Set("Content-Type", "text/csv")
	part, err := writter.CreatePart(h)
	if err != nil {
		return nil, err
	}

	// TODO: Check first return value from Copy method
	// Copy CSV file to form-data
	_, err = io.Copy(part, csvFile)
	if err != nil {
		return nil, err
	}
	writter.Close()

	// Create POST request
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", writter.FormDataContentType())

	// Make request
	data, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	// Get response
	result := &ResultPrediction{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	// Delete csv file
	err = os.Remove(csvFilename)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if 200 != resp.StatusCode {
		return nil, fmt.Errorf("%s", body)
	}

	return body, nil
}
