# APIKeyValidator Traefik Plugin

The `APIKeyValidator` is a custom middleware plugin for Traefik, designed to provide API key validation for incoming HTTP requests. This plugin allows specifying valid API keys, defining ignore and block paths, and the option to remove specific headers from the request.

## Features

- **API Key Validation:** Validates incoming requests based on predefined API keys.
- **Ignore Paths:** Define paths that can bypass API key validation.
- **Block Paths:** Specify paths that should always be denied access.
- **Remove Headers:** Option to remove certain headers from the request.

## Configuration

Here's a breakdown of the configuration options for `APIKeyValidator`:

- `ValidAPIKeys`: A list of strings representing valid API keys.
- `APIKeyHeader`: The header name from which the API key will be extracted. Defaults to `"X-API-Key"`.
- `UseAuthorization`: If set to `true`, the plugin will look for the API key in the Authorization header (`Bearer` token).
- `IgnorePaths`: A list of paths (regex supported) to ignore for API key validation.
- `BlockPaths`: A list of paths (regex supported) that should always be blocked.
- `RemoveHeader`: If `true`, the API key header will be removed from the request after validation.

## Usage

To use the `APIKeyValidator` plugin in your Traefik instance:

1. Add the plugin to your Traefik configuration:

    ```yaml
    # Example Traefik configuration snippet
    experimental:
      plugins:
        apikey-validator:
          moduleName: "github.com/a8851625/traefik_plugin_api_keys"
          version: "v0.1.0"
    ```

2. Configure the middleware in your dynamic Traefik configuration (e.g., using labels in Docker or in Traefik's CRD in Kubernetes).

    ```yaml
    # Example Middleware Configuration
    http:
      middlewares:
        apikey-validator:
          plugin:
            apikey-validator:
              ValidAPIKeys:
                - "key1"
                - "key2"
              APIKeyHeader: "X-API-Key"
              UseAuthorization: false
              IgnorePaths:
                - "^/public"
                - "^/health"
              BlockPaths:
                - "^/admin"
              RemoveHeader: true
    ```

3. Apply the middleware to your routers.

## Development and Contributions

This plugin is open for contributions. Feel free to fork, modify, and make pull requests to enhance its functionalities.

## License

This plugin is released under the [MIT License](https://opensource.org/licenses/MIT).
