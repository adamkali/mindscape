package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"time"

	"github.com/labstack/echo/v4"
)

type CoolifyClient struct {
	httpClient          *http.Client
	defaultTimeout      time.Duration
	personalAccessToken string
	baseUrl             string
}

type CoolifyClientOption func(*CoolifyClient)

type CoolifyMessage struct {
	Message string `json:"message"`
}
type CoolifyApplicationMessage struct {
	Message      string `json:"message"`
	DeploymentId string `json:"deployment_uuid"`
}

type CoolifyService struct {
	ID                              int    `json:"id"`
	UUID                            string `json:"uuid"`
	Name                            string `json:"name"`
	EnvironmentID                   int    `json:"environment_id"`
	ServerID                        int    `json:"server_id"`
	Description                     string `json:"description"`
	DockerComposeRaw                string `json:"docker_compose_raw"`
	DockerCompose                   string `json:"docker_compose"`
	DestinationType                 string `json:"destination_type"`
	DestinationID                   int    `json:"destination_id"`
	ConnectToDockerNetwork          bool   `json:"connect_to_docker_network"`
	IsContainerLabelEscapeEnabled   bool   `json:"is_container_label_escape_enabled"`
	IsContainerLabelReadonlyEnabled bool   `json:"is_container_label_readonly_enabled"`
	ConfigHash                      string `json:"config_hash"`
	ServiceType                     string `json:"service_type"`
	CreatedAt                       string `json:"created_at"`
	UpdatedAt                       string `json:"updated_at"`
	DeletedAt                       string `json:"deleted_at"`
}

type CoolifyApplication struct {
	ID                              int    `json:"id"`
	Description                     string `json:"description"`
	RepositoryProjectID             int    `json:"repository_project_id"`
	UUID                            string `json:"uuid"`
	Name                            string `json:"name"`
	Fqdn                            string `json:"fqdn"`
	ConfigHash                      string `json:"config_hash"`
	GitRepository                   string `json:"git_repository"`
	GitBranch                       string `json:"git_branch"`
	GitCommitSha                    string `json:"git_commit_sha"`
	GitFullURL                      string `json:"git_full_url"`
	DockerRegistryImageName         string `json:"docker_registry_image_name"`
	DockerRegistryImageTag          string `json:"docker_registry_image_tag"`
	BuildPack                       string `json:"build_pack"`
	StaticImage                     string `json:"static_image"`
	InstallCommand                  string `json:"install_command"`
	BuildCommand                    string `json:"build_command"`
	StartCommand                    string `json:"start_command"`
	PortsExposes                    string `json:"ports_exposes"`
	PortsMappings                   string `json:"ports_mappings"`
	CustomNetworkAliases            string `json:"custom_network_aliases"`
	BaseDirectory                   string `json:"base_directory"`
	PublishDirectory                string `json:"publish_directory"`
	HealthCheckEnabled              bool   `json:"health_check_enabled"`
	HealthCheckPath                 string `json:"health_check_path"`
	HealthCheckPort                 string `json:"health_check_port"`
	HealthCheckHost                 string `json:"health_check_host"`
	HealthCheckMethod               string `json:"health_check_method"`
	HealthCheckReturnCode           int    `json:"health_check_return_code"`
	HealthCheckScheme               string `json:"health_check_scheme"`
	HealthCheckResponseText         string `json:"health_check_response_text"`
	HealthCheckInterval             int    `json:"health_check_interval"`
	HealthCheckTimeout              int    `json:"health_check_timeout"`
	HealthCheckRetries              int    `json:"health_check_retries"`
	HealthCheckStartPeriod          int    `json:"health_check_start_period"`
	LimitsMemory                    string `json:"limits_memory"`
	LimitsMemorySwap                string `json:"limits_memory_swap"`
	LimitsMemorySwappiness          int    `json:"limits_memory_swappiness"`
	LimitsMemoryReservation         string `json:"limits_memory_reservation"`
	LimitsCpus                      string `json:"limits_cpus"`
	LimitsCpuset                    string `json:"limits_cpuset"`
	LimitsCPUShares                 int    `json:"limits_cpu_shares"`
	Status                          string `json:"status"`
	PreviewURLTemplate              string `json:"preview_url_template"`
	DestinationType                 string `json:"destination_type"`
	DestinationID                   int    `json:"destination_id"`
	SourceID                        int    `json:"source_id"`
	PrivateKeyID                    int    `json:"private_key_id"`
	EnvironmentID                   int    `json:"environment_id"`
	Dockerfile                      string `json:"dockerfile"`
	DockerfileLocation              string `json:"dockerfile_location"`
	CustomLabels                    string `json:"custom_labels"`
	DockerfileTargetBuild           string `json:"dockerfile_target_build"`
	ManualWebhookSecretGithub       string `json:"manual_webhook_secret_github"`
	ManualWebhookSecretGitlab       string `json:"manual_webhook_secret_gitlab"`
	ManualWebhookSecretBitbucket    string `json:"manual_webhook_secret_bitbucket"`
	ManualWebhookSecretGitea        string `json:"manual_webhook_secret_gitea"`
	DockerComposeLocation           string `json:"docker_compose_location"`
	DockerCompose                   string `json:"docker_compose"`
	DockerComposeRaw                string `json:"docker_compose_raw"`
	DockerComposeDomains            string `json:"docker_compose_domains"`
	DockerComposeCustomStartCommand string `json:"docker_compose_custom_start_command"`
	DockerComposeCustomBuildCommand string `json:"docker_compose_custom_build_command"`
	SwarmReplicas                   int    `json:"swarm_replicas"`
	SwarmPlacementConstraints       string `json:"swarm_placement_constraints"`
	CustomDockerRunOptions          string `json:"custom_docker_run_options"`
	PostDeploymentCommand           string `json:"post_deployment_command"`
	PostDeploymentCommandContainer  string `json:"post_deployment_command_container"`
	PreDeploymentCommand            string `json:"pre_deployment_command"`
	PreDeploymentCommandContainer   string `json:"pre_deployment_command_container"`
	WatchPaths                      string `json:"watch_paths"`
	CustomHealthcheckFound          bool   `json:"custom_healthcheck_found"`
	Redirect                        string `json:"redirect"`
	CreatedAt                       string `json:"created_at"`
	UpdatedAt                       string `json:"updated_at"`
	DeletedAt                       string `json:"deleted_at"`
	ComposeParsingVersion           string `json:"compose_parsing_version"`
	CustomNginxConfiguration        string `json:"custom_nginx_configuration"`
	IsHTTPBasicAuthEnabled          bool   `json:"is_http_basic_auth_enabled"`
	HTTPBasicAuthUsername           string `json:"http_basic_auth_username"`
	HTTPBasicAuthPassword           string `json:"http_basic_auth_password"`
}

type CoolifyFilter struct {
	QueryTarget string `json:"query_target"`
	Index       int    `json:"index"`
}

func SortApplications(apps []CoolifyApplication, by string) ([]CoolifyApplication, error) {
	switch by {
	case "name":
		sort.Slice(apps, func(i, j int) bool {
			return apps[i].Name < apps[j].Name 
		})
		break
	case "created_at":
		sort.Slice(apps, func(i, j int) bool {
			t1, err := time.Parse(time.RFC3339,apps[i].CreatedAt)
			if err != nil {
				return false	
			}
			t2, err := time.Parse(time.RFC3339, apps[j].CreatedAt)
			if err != nil {
				return false
			}
			return t1.Before(t2)
		})
		break
	case "updated_at":
		sort.Slice(apps, func(i, j int) bool {
			t1, err := time.Parse(time.RFC3339,apps[i].UpdatedAt)
			if err != nil {
				return false	
			}
			t2, err := time.Parse(time.RFC3339, apps[j].UpdatedAt)
			if err != nil {
				return false
			}
			return t1.Before(t2)
		})
		break
	default:
		return apps, nil
	}
	return apps, nil
}

func FilterApplications(apps []CoolifyApplication, filter ...string) {
	var filteredApps []CoolifyApplication
	for _, app := range apps {
		for _, filter := range filter {
			if app.UUID == filter {
				filteredApps = append(filteredApps, app)
			}
		}
	}
	apps = filteredApps
}

func FilterServices(services []CoolifyService, filter ...string) {
	var filteredServices []CoolifyService
	for _, service := range services {
		for _, filter := range filter {
			if service.UUID == filter {
				filteredServices = append(filteredServices, service)
			}
		}
	}
	services = filteredServices
}

func SortServices(services []CoolifyService, by string) ([]CoolifyService, error) {
	switch by {
	case "name":
		sort.Slice(services, func(i, j int) bool {
			return services[i].Name < services[j].Name 
		})
		break
	case "created_at":
		sort.Slice(services, func(i, j int) bool {
			t1, err := time.Parse(time.RFC3339,services[i].CreatedAt)
			if err != nil {
				return false	
			}
			t2, err := time.Parse(time.RFC3339, services[j].CreatedAt)
			if err != nil {
				return false
			}
			return t1.Before(t2)
		})
		break
	case "updated_at":
		sort.Slice(services, func(i, j int) bool {
			t1, err := time.Parse(time.RFC3339,services[i].UpdatedAt)
			if err != nil {
				return false	
			}
			t2, err := time.Parse(time.RFC3339, services[j].UpdatedAt)
			if err != nil {
				return false
			}
			return t1.Before(t2)
		})
		break
	default:	
		return services, nil
	}
	return services, nil
}

func NewCoolifyClient(pat string, baseUrl string, opts ...CoolifyClientOption) *CoolifyClient {
	client := &CoolifyClient{
		httpClient:          &http.Client{},
		personalAccessToken: pat,
		defaultTimeout:      30 * time.Second, // Default timeout
		baseUrl:             baseUrl,
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

func CoolifyWithHTTPClient(httpClient *http.Client) CoolifyClientOption {
	return func(c *CoolifyClient) {
		c.httpClient = httpClient
	}
}

func CoolifyWithTimeout(timeout time.Duration) CoolifyClientOption {
	return func(c *CoolifyClient) {
		c.defaultTimeout = timeout
	}
}

func (c *CoolifyClient) doRequest(ctx context.Context, url string, acceptHeader string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Accept", acceptHeader)
	if c.personalAccessToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.personalAccessToken)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("request timed out: %w", err)
		}
		if ctx.Err() == context.Canceled {
			return nil, fmt.Errorf("request canceled: %w", err)
		}
		return nil, fmt.Errorf("request failed: %w", err)

	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Coolify API error: %d", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}

func (c *CoolifyClient) CoolifyHealthCheck(ctx echo.Context) error {
	url := fmt.Sprintf("%s/health", c.baseUrl)

	body, err := c.doRequest(
		ctx.Request().Context(),
		url,
		"text/plain",
	)
	if err != nil {
		return err
	}

	if string(body) != "OK" {
		return fmt.Errorf("health check failed: %s", string(body))
	}

	return nil
}

// Service methods
func (c *CoolifyClient) GetCoolifyServices(ctx echo.Context) ([]CoolifyService, error) {
	url := fmt.Sprintf("%s/services", c.baseUrl)
	body, err := c.doRequest(ctx.Request().Context(), url, "text/plain")
	if err != nil {
		return nil, err
	}
	var services []CoolifyService
	if err := json.Unmarshal(body, &services); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return services, nil
}

func (c *CoolifyClient) getCoolifyServiceUrl(serviceId string) string {
	return fmt.Sprintf("%s/services/%s", c.baseUrl, serviceId)
}

func (c *CoolifyClient) GetCoolifyService(ctx echo.Context, serviceId string) (*CoolifyService, error) {
	url := c.getCoolifyServiceUrl(serviceId)
	body, err := c.doRequest(ctx.Request().Context(), url, "text/plain")
	if err != nil {
		return nil, err
	}

	var service CoolifyService
	if err := json.Unmarshal(body, &service); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &service, nil
}

func (c *CoolifyClient) restartCoolifyServiceUrl(serviceId string) string {
	return fmt.Sprintf("%s/services/%s/restart", c.baseUrl, serviceId)
}

func (c *CoolifyClient) RestartCoolifyService(ctx echo.Context, serviceId string) error {
	url := c.restartCoolifyServiceUrl(serviceId)
	body, err := c.doRequest(ctx.Request().Context(), url, "text/plain")
	if err != nil {
		return err
	}

	var service CoolifyMessage
	if err := json.Unmarshal(body, &service); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	if string(service.Message) != "Service restarting request queued." {
		return fmt.Errorf("Restart failed: %s", string(service.Message))
	}

	return nil
}

func (c *CoolifyClient) startCoolifyServiceUrl(serviceId string) string {
	return fmt.Sprintf("%s/services/%s/start", c.baseUrl, serviceId)
}

func (c *CoolifyClient) StartCoolifyService(ctx echo.Context, serviceId string) error {
	url := c.startCoolifyServiceUrl(serviceId)
	body, err := c.doRequest(ctx.Request().Context(), url, "text/plain")
	if err != nil {
		return err
	}

	var service CoolifyMessage
	if err := json.Unmarshal(body, &service); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	if string(service.Message) != "Service starting request queued." {
		return fmt.Errorf("Start failed: %s", string(service.Message))
	}
	return nil
}

func (c *CoolifyClient) stopCoolifyServiceUrl(serviceId string) string {
	return fmt.Sprintf("%s/services/%s/stop", c.baseUrl, serviceId)
}

func (c *CoolifyClient) StopCoolifyService(ctx echo.Context, serviceId string) error {
	url := c.stopCoolifyServiceUrl(serviceId)
	body, err := c.doRequest(ctx.Request().Context(), url, "text/plain")
	if err != nil {
		return err
	}

	var service CoolifyMessage
	if err := json.Unmarshal(body, &service); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	if string(service.Message) != "Service stopping request queued." {
		return fmt.Errorf("Stop failed: %s", string(service.Message))
	}
	return nil
}

// Application methods
func (c *CoolifyClient) GetCoolifyApplications(ctx echo.Context) ([]CoolifyApplication, error) {
	url := fmt.Sprintf("%s/applications", c.baseUrl)
	body, err := c.doRequest(ctx.Request().Context(), url, "text/plain")
	if err != nil {
		return nil, err
	}
	var applications []CoolifyApplication
	if err := json.Unmarshal(body, &applications); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return applications, nil

}

func (c *CoolifyClient) getCoolifyApplicationUrl(applicationId string) string {
	return fmt.Sprintf("%s/applications/%s", c.baseUrl, applicationId)
}

func (c *CoolifyClient) GetCoolifyApplication(ctx echo.Context, applicationId string) (*CoolifyApplication, error) {
	url := c.getCoolifyApplicationUrl(applicationId)
	body, err := c.doRequest(ctx.Request().Context(), url, "text/plain")
	if err != nil {
		return nil, err
	}
	var application CoolifyApplication
	if err := json.Unmarshal(body, &application); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return &application, nil
}

func (c *CoolifyClient) restartCoolifyApplicationUrl(applicationId string) string {
	return fmt.Sprintf("%s/applications/%s/restart", c.baseUrl, applicationId)
}

func (c *CoolifyClient) RestartCoolifyApplication(ctx echo.Context, applicationId string) (string, error) {
	url := c.restartCoolifyApplicationUrl(applicationId)
	body, err := c.doRequest(ctx.Request().Context(), url, "text/plain")
	if err != nil {
		return "", err
	}

	var application CoolifyApplicationMessage
	if err := json.Unmarshal(body, &application); err != nil {
		return "", fmt.Errorf("failed to parse JSON: %w", err)
	}

	if string(application.Message) != "Application restarting request queued." {
		return "", fmt.Errorf("Restart failed: %s", string(application.Message))
	}
	return application.DeploymentId, nil
}

func (c *CoolifyClient) startCoolifyApplicationUrl(applicationId string) string {
	return fmt.Sprintf("%s/applications/%s/start", c.baseUrl, applicationId)
}

func (c *CoolifyClient) StartCoolifyApplication(ctx echo.Context, applicationId string) (string, error) {
	url := c.startCoolifyApplicationUrl(applicationId)
	body, err := c.doRequest(ctx.Request().Context(), url, "text/plain")
	if err != nil {
		return "", err
	}
	var application CoolifyApplicationMessage
	if err := json.Unmarshal(body, &application); err != nil {
		return "", fmt.Errorf("failed to parse JSON: %w", err)
	}
	if string(application.Message) != "Application starting request queued." {
		return "", fmt.Errorf("Start failed: %s", string(application.Message))
	}
	return application.DeploymentId, nil
}

func (c *CoolifyClient) stopCoolifyApplicationUrl(applicationId string) string {
	return fmt.Sprintf("%s/applications/%s/stop", c.baseUrl, applicationId)
}

func (c *CoolifyClient) StopCoolifyApplication(ctx echo.Context, applicationId string) (string, error) {
	url := c.stopCoolifyApplicationUrl(applicationId)
	body, err := c.doRequest(ctx.Request().Context(), url, "text/plain")
	if err != nil {
		return "", err
	}
	var application CoolifyApplicationMessage
	if err := json.Unmarshal(body, &application); err != nil {
		return "", fmt.Errorf("failed to parse JSON: %w", err)
	}
	if string(application.Message) != "Application stopping request queued." {
		return "", fmt.Errorf("Stop failed: %s", string(application.Message))
	}
	return application.DeploymentId, nil
}
