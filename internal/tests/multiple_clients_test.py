import socket
import threading
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

def client_setter():
    with socket.create_connection(("localhost", 6379)) as sock:
        print("[SETTER] Sending SET foo bar PX 1000")
        sock.sendall(build_set_px_command("foo", "bar", 1000).encode())
        time.sleep(0.1)
        response = recv_resp(sock)
        print("[SETTER] Response:", repr(response))

def client_getter():
    time.sleep(0.2)  # small delay to let setter go first
    with socket.create_connection(("localhost", 6379)) as sock:
        print("[GETTER] Sending GET foo")
        sock.sendall(build_get_command("foo").encode())
        response = recv_resp(sock)
        print("[GETTER] Response:", repr(response))

if __name__ == "__main__":
    t1 = threading.Thread(target=client_setter)
    t2 = threading.Thread(target=client_getter)

    t1.start()
    t2.start()

    t1.join()
    t2.join()

    print("Finished concurrent test.")
