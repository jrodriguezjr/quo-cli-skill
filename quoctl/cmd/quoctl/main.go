package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"golang.org/x/term"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 || args[0] == "help" || args[0] == "-h" || args[0] == "--help" {
		usage()
		return nil
	}

	switch args[0] {
	case "contacts":
		return contactsCmd(args[1:])
	case "messages":
		return messagesCmd(args[1:])
	case "phone-numbers":
		return phoneNumbersCmd(args[1:])
	case "users":
		return usersCmd(args[1:])
	case "api":
		return apiCmd(args[1:])
	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func usage() {
	fmt.Println(`quoctl - Quo/OpenPhone API CLI

Usage:
  quoctl <command> [subcommand] [flags]

Commands:
  contacts list|get
  messages list|send
  phone-numbers list
  users list
  api get|post|patch|delete

Auth:
  Set QUO_API_KEY in your environment.

Examples:
  quoctl contacts list --max-results 10
  quoctl messages send --from +12105551234 --to +12105559876 --content "Hello"
  quoctl api get /v1/conversations?maxResults=20
`)
}

type client struct {
	baseURL    string
	apiKey     string
	authScheme string
	httpClient *http.Client
}

type commonFlags struct {
	baseURL    *string
	authScheme *string
	apiKey     *string
	timeout    *time.Duration
}

func bindCommonFlags(fs *flag.FlagSet) commonFlags {
	return commonFlags{
		baseURL:    fs.String("base-url", getenv("QUO_BASE_URL", "https://api.openphone.com"), "API base URL"),
		authScheme: fs.String("auth-scheme", getenv("QUO_AUTH_SCHEME", "Bearer"), "Authorization scheme (Bearer, ApiKey, None)"),
		apiKey:     fs.String("api-key", "", "API key (overrides QUO_API_KEY)"),
		timeout:    fs.Duration("timeout", 30*time.Second, "request timeout"),
	}
}

func newClient(cf commonFlags) (*client, error) {
	apiKey := strings.TrimSpace(*cf.apiKey)
	if apiKey == "" {
		apiKey = strings.TrimSpace(os.Getenv("QUO_API_KEY"))
	}
	if apiKey == "" {
		apiKey = loadAPIKeyFromEnvFile(".quoctl.env")
	}
	if strings.ToLower(*cf.authScheme) != "none" && apiKey == "" {
		if isInteractiveTTY() {
			entered, err := promptForAPIKey()
			if err != nil {
				return nil, err
			}
			apiKey = entered
			if err := writeAPIKeyEnvFile(".quoctl.env", apiKey); err != nil {
				return nil, fmt.Errorf("failed to write .quoctl.env: %w", err)
			}
			fmt.Fprintln(os.Stderr, "Saved API key to .quoctl.env")
		} else {
			return nil, errors.New("missing API key: set QUO_API_KEY, pass --api-key, or create .quoctl.env")
		}
	}

	return &client{
		baseURL:    strings.TrimRight(*cf.baseURL, "/"),
		apiKey:     apiKey,
		authScheme: *cf.authScheme,
		httpClient: &http.Client{Timeout: *cf.timeout},
	}, nil
}

func (c *client) do(method, path string, query map[string]string, body any) error {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	u, err := url.Parse(c.baseURL + path)
	if err != nil {
		return err
	}
	q := u.Query()
	for k, v := range query {
		if strings.TrimSpace(v) != "" {
			q.Set(k, v)
		}
	}
	u.RawQuery = q.Encode()

	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reqBody = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, u.String(), reqBody)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	switch strings.ToLower(c.authScheme) {
	case "none":
	case "apikey":
		req.Header.Set("Authorization", c.apiKey)
	default:
		req.Header.Set("Authorization", c.authScheme+" "+c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if len(respBody) == 0 {
		respBody = []byte("{}")
	}

	var pretty bytes.Buffer
	if err := json.Indent(&pretty, respBody, "", "  "); err == nil {
		respBody = pretty.Bytes()
	}

	fmt.Println(string(respBody))

	if resp.StatusCode >= 400 {
		return fmt.Errorf("HTTP %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}
	return nil
}

func contactsCmd(args []string) error {
	if len(args) == 0 {
		return errors.New("contacts requires subcommand: list|get")
	}

	switch args[0] {
	case "list":
		fs := flag.NewFlagSet("contacts list", flag.ContinueOnError)
		cf := bindCommonFlags(fs)
		maxResults := fs.String("max-results", "20", "max results")
		pageToken := fs.String("page-token", "", "page token")
		if err := fs.Parse(args[1:]); err != nil {
			return err
		}
		c, err := newClient(cf)
		if err != nil {
			return err
		}
		return c.do(http.MethodGet, "/v1/contacts", map[string]string{"maxResults": *maxResults, "pageToken": *pageToken}, nil)
	case "get":
		if len(args) < 2 {
			return errors.New("usage: quoctl contacts get <id>")
		}
		fs := flag.NewFlagSet("contacts get", flag.ContinueOnError)
		cf := bindCommonFlags(fs)
		if err := fs.Parse(args[2:]); err != nil {
			return err
		}
		c, err := newClient(cf)
		if err != nil {
			return err
		}
		return c.do(http.MethodGet, "/v1/contacts/"+args[1], nil, nil)
	default:
		return fmt.Errorf("unknown contacts subcommand %q", args[0])
	}
}

func messagesCmd(args []string) error {
	if len(args) == 0 {
		return errors.New("messages requires subcommand: list|send")
	}
	switch args[0] {
	case "list":
		fs := flag.NewFlagSet("messages list", flag.ContinueOnError)
		cf := bindCommonFlags(fs)
		phoneNumberId := fs.String("phone-number-id", "", "OpenPhone number ID (required)")
		participants := fs.String("participants", "", "participant numbers, comma-separated (required)")
		maxResults := fs.String("max-results", "20", "max results")
		pageToken := fs.String("page-token", "", "page token")
		if err := fs.Parse(args[1:]); err != nil {
			return err
		}
		if *phoneNumberId == "" || *participants == "" {
			return errors.New("usage: quoctl messages list --phone-number-id <PN...> --participants <+1..., +1...>")
		}
		c, err := newClient(cf)
		if err != nil {
			return err
		}
		return c.do(http.MethodGet, "/v1/messages", map[string]string{
			"phoneNumberId": *phoneNumberId,
			"participants":  *participants,
			"maxResults":    *maxResults,
			"pageToken":     *pageToken,
		}, nil)
	case "send":
		fs := flag.NewFlagSet("messages send", flag.ContinueOnError)
		cf := bindCommonFlags(fs)
		from := fs.String("from", "", "from OpenPhone number (E.164)")
		to := fs.String("to", "", "to recipient number (E.164)")
		content := fs.String("content", "", "message body")
		userID := fs.String("user-id", "", "optional OpenPhone user ID")
		setInboxStatus := fs.String("set-inbox-status", "", "optional: done")
		if err := fs.Parse(args[1:]); err != nil {
			return err
		}
		if *from == "" || *to == "" || *content == "" {
			return errors.New("usage: quoctl messages send --from <num|PN...> --to <num> --content <msg>")
		}
		c, err := newClient(cf)
		if err != nil {
			return err
		}
		payload := map[string]any{
			"from":    *from,
			"to":      []string{*to},
			"content": *content,
		}
		if *userID != "" {
			payload["userId"] = *userID
		}
		if *setInboxStatus != "" {
			payload["setInboxStatus"] = *setInboxStatus
		}
		return c.do(http.MethodPost, "/v1/messages", nil, payload)
	default:
		return fmt.Errorf("unknown messages subcommand %q", args[0])
	}
}

func phoneNumbersCmd(args []string) error {
	if len(args) == 0 || args[0] != "list" {
		return errors.New("usage: quoctl phone-numbers list")
	}
	fs := flag.NewFlagSet("phone-numbers list", flag.ContinueOnError)
	cf := bindCommonFlags(fs)
	if err := fs.Parse(args[1:]); err != nil {
		return err
	}
	c, err := newClient(cf)
	if err != nil {
		return err
	}
	return c.do(http.MethodGet, "/v1/phone-numbers", nil, nil)
}

func usersCmd(args []string) error {
	if len(args) == 0 || args[0] != "list" {
		return errors.New("usage: quoctl users list")
	}
	fs := flag.NewFlagSet("users list", flag.ContinueOnError)
	cf := bindCommonFlags(fs)
	maxResults := fs.String("max-results", "20", "max results")
	pageToken := fs.String("page-token", "", "page token")
	if err := fs.Parse(args[1:]); err != nil {
		return err
	}
	c, err := newClient(cf)
	if err != nil {
		return err
	}
	return c.do(http.MethodGet, "/v1/users", map[string]string{"maxResults": *maxResults, "pageToken": *pageToken}, nil)
}

func apiCmd(args []string) error {
	if len(args) < 2 {
		return errors.New("usage: quoctl api <get|post|patch|delete> <path> [--data '{...}']")
	}
	method := strings.ToUpper(args[0])
	if method == "DEL" {
		method = "DELETE"
	}
	path := args[1]
	fs := flag.NewFlagSet("api", flag.ContinueOnError)
	cf := bindCommonFlags(fs)
	data := fs.String("data", "", "JSON payload")
	if err := fs.Parse(args[2:]); err != nil {
		return err
	}
	c, err := newClient(cf)
	if err != nil {
		return err
	}

	var payload any
	if strings.TrimSpace(*data) != "" {
		if err := json.Unmarshal([]byte(*data), &payload); err != nil {
			return fmt.Errorf("invalid --data JSON: %w", err)
		}
	}
	return c.do(method, path, nil, payload)
}

func loadAPIKeyFromEnvFile(path string) string {
	b, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(b), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "export ") {
			line = strings.TrimSpace(strings.TrimPrefix(line, "export "))
		}
		if strings.HasPrefix(line, "QUO_API_KEY=") {
			v := strings.TrimSpace(strings.TrimPrefix(line, "QUO_API_KEY="))
			v = strings.Trim(v, "\"'")
			if v != "" {
				return v
			}
		}
	}
	return ""
}

func writeAPIKeyEnvFile(path, apiKey string) error {
	content := "QUO_API_KEY=" + apiKey + "\n"
	return os.WriteFile(path, []byte(content), 0o600)
}

func isInteractiveTTY() bool {
	return term.IsTerminal(int(os.Stdin.Fd()))
}

func promptForAPIKey() (string, error) {
	fmt.Fprint(os.Stderr, "Enter QUO API key: ")
	r := bufio.NewReader(os.Stdin)
	v, err := r.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}
	v = strings.TrimSpace(v)
	if v == "" {
		return "", errors.New("empty API key")
	}
	return v, nil
}

func getenv(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}
