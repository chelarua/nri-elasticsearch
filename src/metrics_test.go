package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/newrelic/infra-integrations-sdk/integration"

	"github.com/stretchr/testify/assert"
)

type testClient struct {
	endpointMapping    map[string]string
	ReturnRequestError bool
}

func (c *testClient) init(filename string, endpoint string) {
	c.endpointMapping = map[string]string{
		endpoint: filepath.Join("testdata", filename),
	}
}

func (c *testClient) Request(endpoint string, v interface{}) error {
	if c.ReturnRequestError {
		return errors.New("error")
	}

	jsonPath := c.endpointMapping[endpoint]

	jsonData, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonData, v)
}

func createNewTestClient() *testClient {
	return new(testClient)
}

func createGoldenFile(i *integration.Integration, sourceFile string) (string, []byte) {
	goldenFile := sourceFile + ".golden"
	actualContents, _ := i.Entities[0].Metrics[0].MarshalJSON()

	if *update {
		ioutil.WriteFile(goldenFile, actualContents, 0644)
	}
	return goldenFile, actualContents
}

func TestPopulateNodesMetrics(t *testing.T) {
	i := getTestingIntegration(t)
	client := createNewTestClient()
	client.init("nodeStatsMetricsResult.json", nodeStatsEndpoint)

	populateNodesMetrics(i, client)

	sourceFile := "testdata/nodeStatsMetricsResult.json"
	goldenFile, actualContents := createGoldenFile(i, sourceFile)
	expectedContents, _ := ioutil.ReadFile(goldenFile)

	actualLength := len(i.Entities[0].Metrics[0].Metrics)
	expectedLength := 166

	assert.Equal(t, 1, len(i.Entities))
	assert.Equal(t, 1, len(i.Entities[0].Metrics))
	assert.Equal(t, expectedContents, actualContents)
	assert.Equal(t, expectedLength, actualLength)
}

func TestPopulateNodesMetrics_Error(t *testing.T) {
	mockClient := createNewTestClient()
	mockClient.ReturnRequestError = true

	i := getTestingIntegration(t)
	err := populateNodesMetrics(i, mockClient)
	assert.Error(t, err, "should be an error")
}

func TestPopulateClusterMetrics(t *testing.T) {
	i := getTestingIntegration(t)
	client := createNewTestClient()
	client.init("clusterStatsMetricsResult.json", clusterEndpoint)

	populateClusterMetrics(i, client)

	sourceFile := "testData/clusterStatsMetricsResult.json"
	goldenFile, actualContents := createGoldenFile(i, sourceFile)
	expectedContents, _ := ioutil.ReadFile(goldenFile)

	actualLength := len(i.Entities[0].Metrics[0].Metrics)
	expectedLength := 11

	assert.Equal(t, expectedContents, actualContents)
	assert.Equal(t, expectedLength, actualLength)
}

func TestPopulateClusterMetrics_Error(t *testing.T) {
	mockClient := createNewTestClient()
	mockClient.ReturnRequestError = true

	i := getTestingIntegration(t)
	err := populateClusterMetrics(i, mockClient)
	assert.Error(t, err, "should be an error")
}

func TestPopulateCommonMetrics(t *testing.T) {
	i := getTestingIntegration(t)
	client := createNewTestClient()
	client.init("commonMetricsResult.json", commonStatsEndpoint)

	populateCommonMetrics(i, client)

	sourceFile := "testData/commonMetricsResult.json"
	goldenFile, actualContents := createGoldenFile(i, sourceFile)
	expectedContents, _ := ioutil.ReadFile(goldenFile)

	actualLength := len(i.Entities[0].Metrics[0].Metrics)
	expectedLength := 36

	assert.Equal(t, expectedContents, actualContents)
	assert.Equal(t, expectedLength, actualLength)
}

func TestPopulateCommonMetrics_Error(t *testing.T) {
	mockClient := createNewTestClient()
	mockClient.ReturnRequestError = true

	i := getTestingIntegration(t)
	err := populateCommonMetrics(i, mockClient)
	assert.Error(t, err, "should be an error")
}

func TestPopulateIndicesMetrics(t *testing.T) {
	i := getTestingIntegration(t)
	client := createNewTestClient()
	client.init("indicesMetricsResult.json", indicesStatsEndpoint)

	populateIndicesMetrics(i, client)

	sourceFile := "testData/indicesMetricsResult.json"
	goldenFile, actualContents := createGoldenFile(i, sourceFile)

	for j := range i.Entities {
		resultStruct := i.Entities[j].Metrics[0].Metrics
		actualLength := len(resultStruct)
		expectedLength := 10
		assert.Equal(t, expectedLength, actualLength)
	}

	expectedContents, _ := ioutil.ReadFile(goldenFile)
	assert.Equal(t, expectedContents, actualContents)
}

func TestPopulateIndicesMetrics_Error(t *testing.T) {
	mockClient := createNewTestClient()
	mockClient.ReturnRequestError = true

	i := getTestingIntegration(t)
	err := populateIndicesMetrics(i, mockClient)
	assert.Error(t, err, "should be an error")
}
