package edgestacks

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	portainer "github.com/portainer/portainer/api"
	"github.com/stretchr/testify/assert"

	"github.com/segmentio/encoding/json"
)

// Delete
func TestDeleteAndInspect(t *testing.T) {
	handler, rawAPIKey := setupHandler(t)

	// Create
	endpoint := createEndpoint(t, handler.DataStore)
	edgeStack := createEdgeStack(t, handler.DataStore, endpoint.ID)

	// Inspect
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/edge_stacks/%d", edgeStack.ID), nil)
	if err != nil {
		t.Fatal("request error:", err)
	}

	req.Header.Add("x-api-key", rawAPIKey)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected a %d response, found: %d", http.StatusOK, rec.Code)
	}

	data := portainer.EdgeStack{}
	err = json.NewDecoder(rec.Body).Decode(&data)
	if err != nil {
		t.Fatal("error decoding response:", err)
	}

	if data.ID != edgeStack.ID {
		t.Fatalf("expected EdgeStackID %d, found %d", int(edgeStack.ID), data.ID)
	}

	// Delete
	req, err = http.NewRequest(http.MethodDelete, fmt.Sprintf("/edge_stacks/%d", edgeStack.ID), nil)
	if err != nil {
		t.Fatal("request error:", err)
	}

	req.Header.Add("x-api-key", rawAPIKey)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected a %d response, found: %d", http.StatusNoContent, rec.Code)
	}

	// Inspect
	req, err = http.NewRequest(http.MethodGet, fmt.Sprintf("/edge_stacks/%d", edgeStack.ID), nil)
	if err != nil {
		t.Fatal("request error:", err)
	}

	req.Header.Add("x-api-key", rawAPIKey)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected a %d response, found: %d", http.StatusNotFound, rec.Code)
	}
}

func TestDeleteInvalidEdgeStack(t *testing.T) {
	handler, rawAPIKey := setupHandler(t)

	cases := []struct {
		Name               string
		URL                string
		ExpectedStatusCode int
	}{
		{Name: "Non-existing EdgeStackID", URL: "/edge_stacks/-1", ExpectedStatusCode: http.StatusNotFound},
		{Name: "Invalid EdgeStackID", URL: "/edge_stacks/aaaaaaa", ExpectedStatusCode: http.StatusBadRequest},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodDelete, tc.URL, nil)
			if err != nil {
				t.Fatal("request error:", err)
			}

			req.Header.Add("x-api-key", rawAPIKey)
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			if rec.Code != tc.ExpectedStatusCode {
				t.Fatalf("expected a %d response, found: %d", tc.ExpectedStatusCode, rec.Code)
			}
		})
	}
}

func TestDeleteEdgeStack_RemoveProjectFolder(t *testing.T) {
	handler, rawAPIKey := setupHandler(t)

	edgeGroup := createEdgeGroup(t, handler.DataStore)

	payload := edgeStackFromStringPayload{
		Name:             "test-stack",
		DeploymentType:   portainer.EdgeStackDeploymentCompose,
		EdgeGroups:       []portainer.EdgeGroupID{edgeGroup.ID},
		StackFileContent: "version: '3.7'\nservices:\n  test:\n    image: test",
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(payload); err != nil {
		t.Fatal("error encoding payload:", err)
	}

	// Create
	req, err := http.NewRequest(http.MethodPost, "/edge_stacks/create/string", &buf)
	if err != nil {
		t.Fatal("request error:", err)
	}

	req.Header.Add("x-api-key", rawAPIKey)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected a %d response, found: %d", http.StatusNoContent, rec.Code)
	}

	assert.DirExists(t, handler.FileService.GetEdgeStackProjectPath("1"))

	// Delete
	if req, err = http.NewRequest(http.MethodDelete, "/edge_stacks/1", nil); err != nil {
		t.Fatal("request error:", err)
	}

	req.Header.Add("x-api-key", rawAPIKey)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected a %d response, found: %d", http.StatusNoContent, rec.Code)
	}

	assert.NoDirExists(t, handler.FileService.GetEdgeStackProjectPath("1"))
}
