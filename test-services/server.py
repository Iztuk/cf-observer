from http.server import BaseHTTPRequestHandler, HTTPServer
import json
import time
from urllib.parse import urlparse


USERS = [
    {"id": "1", "email": "john@example.com", "displayName": "John"},
    {"id": "2", "email": "jane@example.com", "displayName": "Jane"},
]


class Handler(BaseHTTPRequestHandler):
    def _send_json(self, status_code, payload):
        body = json.dumps(payload).encode("utf-8")

        self.send_response(status_code)
        self.send_header("Content-Type", "application/json")
        self.send_header("Content-Length", str(len(body)))
        self.end_headers()
        self.wfile.write(body)

    def do_GET(self):
        path = urlparse(self.path).path

        if path == "/health":
            self._send_json(200, {"status": "ok"})
            return

        if path == "/users":
            self._send_json(200, USERS)
            return

        if path.startswith("/users/"):
            user_id = path.removeprefix("/users/")

            for user in USERS:
                if user["id"] == user_id:
                    self._send_json(200, user)
                    return

            self._send_json(404, {"error": "user not found"})
            return

        # Should trigger cf-observer timeout handling if proxy timeout is < 10s.
        if path == "/timeout":
            time.sleep(10)
            self._send_json(200, {"message": "slow response completed"})
            return

        # Attempts to simulate an upstream connection/reset-style failure.
        if path == "/reset":
            self.close_connection = True
            return

        self._send_json(404, {"error": "route not found"})


def main():
    server = HTTPServer(("localhost", 8081), Handler)
    print("Test service listening on http://localhost:8081")
    server.serve_forever()


if __name__ == "__main__":
    main()
