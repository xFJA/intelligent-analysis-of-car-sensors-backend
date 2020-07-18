package pca

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

// NewClient returns a new entity of Client.
func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
	}
}

// PCA performs a POST request to the PCA service sending a CSV file created from the Dataset entity.
func (c *Client) PCA(clientRequest *ClientRequest) (*models.PCAResult, error) {
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
	pcaResult := &models.PCAResult{}
	err = json.Unmarshal(data, &pcaResult)
	if err != nil {
		return nil, err
	}

	// Delete csv file
	err = os.Remove(csvFilename)
	if err != nil {
		return nil, err
	}

	return pcaResult, nil
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
