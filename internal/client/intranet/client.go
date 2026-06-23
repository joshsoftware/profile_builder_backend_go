package intranet

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/specs"
	"go.uber.org/zap"
)

// IntranetClient is the interface for fetching employee data from the Intranet API.
type IntranetClient interface {
	GetEmployees(ctx context.Context) ([]specs.IntranetEmployee, error)
	GetEmployeeByID(ctx context.Context, employeeID string) (*specs.IntranetEmployee, error)
}

// httpClient is the concrete implementation of IntranetClient.
type httpClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// NewClient creates a new IntranetClient backed by an HTTP client.
func NewClient(baseURL, apiKey string) IntranetClient {
	return &httpClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetEmployees calls the Intranet API and returns a slice of IntranetEmployee records.
func (c *httpClient) GetEmployees(ctx context.Context) ([]specs.IntranetEmployee, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL, nil)
	if err != nil {
		zap.S().Error("Error creating intranet HTTP request: ", err)
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("X-API-Key", c.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		zap.S().Error("Error executing intranet HTTP request: ", err)
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		zap.S().Errorf("Intranet API returned non-2xx status %d: %s", resp.StatusCode, string(body))
		return nil, fmt.Errorf("intranet API returned status %d", resp.StatusCode)
	}

	var employees []specs.IntranetEmployee
	if err := json.NewDecoder(resp.Body).Decode(&employees); err != nil {
		zap.S().Error("Error decoding intranet API response: ", err)
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return employees, nil
}

// GetEmployeeByID calls the Intranet API and returns a single IntranetEmployee record.
func (c *httpClient) GetEmployeeByID(ctx context.Context, employeeID string) (*specs.IntranetEmployee, error) {
	url := fmt.Sprintf("%s/%s", c.baseURL, employeeID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		zap.S().Error("Error creating intranet HTTP request: ", err)
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("X-API-Key", c.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		zap.S().Error("Error executing intranet HTTP request: ", err)
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if resp.StatusCode == http.StatusNotFound {
			return nil, errors.ErrNoRecordFound
		}
		body, _ := io.ReadAll(resp.Body)
		zap.S().Errorf("Intranet API returned non-2xx status %d: %s", resp.StatusCode, string(body))
		return nil, fmt.Errorf("intranet API returned status %d", resp.StatusCode)
	}

	var employee specs.IntranetEmployee
	if err := json.NewDecoder(resp.Body).Decode(&employee); err != nil {
		zap.S().Error("Error decoding intranet API response: ", err)
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &employee, nil
}
