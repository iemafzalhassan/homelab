# Grafana MCP Server Setup Guide

Model Context Protocol (MCP) server for Grafana integration in Antigravity IDE.

---

## Installation & Token Setup

1. Generate a Service Account Token in Grafana Cloud (`Administration > Service Accounts`).
2. Install `uv` Python tool:
   ```bash
   curl -LsSf https://astral.sh/uv/install.sh | sh
   ```
3. Configure `mcp_config.json` in Antigravity IDE:
   ```json
   {
     "mcpServers": {
       "grafana": {
         "command": "uvx",
         "args": ["mcp-grafana"],
         "env": {
           "GRAFANA_URL": "https://<your-org>.grafana.net",
           "GRAFANA_SERVICE_ACCOUNT_TOKEN": "glsa_your_token_here"
         }
       }
     }
   }
   ```
