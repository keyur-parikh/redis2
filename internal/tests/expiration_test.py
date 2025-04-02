import socket
import time

def build_set_px_command(key, value, px_ms):
    return (
        f"*5\r\n"
        f"$3\r\nSET\r\n"
        f"${len(key)}\r\n{key}\r\n"
        f"${len(value)}\r\n{value}\r\n"
        f"$2\r\npx\r\n"
        f"${len(str(px_ms))}\r\n{px_ms}\r\n"
    )

def build_get_command(key):
    return f"*2\r\n$3\r\nGET\r\n${len(key)}\r\n{key}\r\n"

def recv_resp(sock):
    sock.settimeout(1)
    data = b""
    while True:
        try:
            chunk = sock.recv(4096)
            if not chunk:
                break
            data += chunk
        except socket.timeout:
            break
    return data.decode()

if __name__ == "__main__":
    host = "localhost"
    port = 6379

    with socket.create_connection((host, port)) as sock:
        print("Sending SET foo bar PX 100")
        sock.sendall(build_set_px_command("foo", "bar", 100).encode())
        time.sleep(0.05)
        print("Sending immediate GET foo")
        sock.sendall(build_get_command("foo").encode())
        time.sleep(0.2)
        print("Sending GET foo after 200ms")
        sock.sendall(build_get_command("foo").encode())

        # Receive all responses
        response = recv_resp(sock)
        print("Raw server response:\n", repr(response))
