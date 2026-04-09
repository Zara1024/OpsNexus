import http.client
import os
import sys
import urllib.parse
from http.server import SimpleHTTPRequestHandler, ThreadingHTTPServer
from pathlib import Path


DIST_DIR = Path(sys.argv[1]).resolve() if len(sys.argv) > 1 else Path("web/dist").resolve()
API_HOST = sys.argv[2] if len(sys.argv) > 2 else "127.0.0.1"
API_PORT = int(sys.argv[3]) if len(sys.argv) > 3 else 8000
LISTEN_HOST = sys.argv[4] if len(sys.argv) > 4 else "0.0.0.0"
LISTEN_PORT = int(sys.argv[5]) if len(sys.argv) > 5 else 8080


class SPAProxyHandler(SimpleHTTPRequestHandler):
    def __init__(self, *args, **kwargs):
        super().__init__(*args, directory=str(DIST_DIR), **kwargs)

    def end_headers(self):
        self.send_header("Cache-Control", "no-cache")
        super().end_headers()

    def do_GET(self):
        if self.path.startswith("/api/"):
            self._proxy_request()
            return
        return super().do_GET()

    def do_POST(self):
        if self.path.startswith("/api/"):
            self._proxy_request()
            return
        self.send_error(405, "Method Not Allowed")

    def do_PUT(self):
        if self.path.startswith("/api/"):
            self._proxy_request()
            return
        self.send_error(405, "Method Not Allowed")

    def do_DELETE(self):
        if self.path.startswith("/api/"):
            self._proxy_request()
            return
        self.send_error(405, "Method Not Allowed")

    def do_OPTIONS(self):
        if self.path.startswith("/api/"):
            self._proxy_request()
            return
        self.send_response(204)
        self.send_header("Access-Control-Allow-Origin", "*")
        self.send_header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
        self.send_header("Access-Control-Allow-Headers", "Content-Type, Authorization")
        self.end_headers()

    def send_head(self):
        path = self.translate_path(self.path)
        if os.path.isdir(path):
            index_path = os.path.join(path, "index.html")
            if os.path.exists(index_path):
                path = index_path
        if not os.path.exists(path):
            path = str(DIST_DIR / "index.html")

        if path.endswith(os.sep):
            path = os.path.join(path, "index.html")

        ctype = self.guess_type(path)
        try:
            f = open(path, "rb")
        except OSError:
            self.send_error(404, "File not found")
            return None

        fs = os.fstat(f.fileno())
        self.send_response(200)
        self.send_header("Content-type", ctype)
        self.send_header("Content-Length", str(fs.st_size))
        self.end_headers()
        return f

    def _proxy_request(self):
        parsed = urllib.parse.urlsplit(self.path)
        conn = http.client.HTTPConnection(API_HOST, API_PORT, timeout=60)
        try:
            content_length = int(self.headers.get("Content-Length", "0"))
            body = self.rfile.read(content_length) if content_length > 0 else None
            headers = {}
            for key, value in self.headers.items():
                if key.lower() in {"host", "content-length", "connection"}:
                    continue
                headers[key] = value
            if body is not None:
                headers["Content-Length"] = str(len(body))

            target_path = parsed.path
            if parsed.query:
                target_path += "?" + parsed.query

            conn.request(self.command, target_path, body=body, headers=headers)
            resp = conn.getresponse()
            response_body = resp.read()

            self.send_response(resp.status)
            for key, value in resp.getheaders():
                if key.lower() in {"transfer-encoding", "connection", "content-length"}:
                    continue
                self.send_header(key, value)
            self.send_header("Content-Length", str(len(response_body)))
            self.end_headers()
            if response_body:
                self.wfile.write(response_body)
        finally:
            conn.close()


def main():
    if not DIST_DIR.exists():
        raise SystemExit(f"dist directory not found: {DIST_DIR}")
    server = ThreadingHTTPServer((LISTEN_HOST, LISTEN_PORT), SPAProxyHandler)
    print(f"serving dist={DIST_DIR} on http://{LISTEN_HOST}:{LISTEN_PORT}, proxy /api -> http://{API_HOST}:{API_PORT}")
    server.serve_forever()


if __name__ == "__main__":
    main()
