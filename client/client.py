import socket
import random
import struct
import sys
import hashlib
import time

RESULT_FILE = 'result.txt'
STATISTICS_FILE = 'statistics.txt'

MIN_INT = -10**2
MAX_INT = +10**2

def generate_matrix(rows, cols):
    return [[random.randint(MIN_INT, MAX_INT) for _ in range(cols)] for _ in range(rows)]

def generate_vector(size):
    return [random.randint(MIN_INT, MAX_INT) for _ in range(size)]

def calculate_hash(matrix, vector):
    hash_obj = hashlib.sha256()
    for row in matrix:
        for item in row:
            hash_obj.update(struct.pack('<i', item))
    for item in vector:
        hash_obj.update(struct.pack('<i', item))
    return hash_obj.digest()

def send_data(sock, rows, cols, matrix, vector, data_hash):
    data = [rows, cols] + [item for row in matrix for item in row] + vector
    packed_data = struct.pack(f'<ii{cols * rows}i{cols}i', *data)
    sock.sendall(packed_data)
    sock.sendall(data_hash)

def recv_data(sock, rows):
    result = sock.recv(rows * 4)
    result_unpacked = struct.unpack(f'<{rows}i', result)
    compute_time = struct.unpack('<d', sock.recv(8))[0]
    data_size = struct.unpack('<i', sock.recv(4))[0]
    return result_unpacked, compute_time, data_size

def multiply_matrix_vector(matrix, vector):
    rows = len(matrix)
    result = [0] * rows
    for i in range(rows):
        for j in range(len(vector)):
            result[i] += matrix[i][j] * vector[j]
    return result

def main():
    start_time = time.time()

    if len(sys.argv) != 4:
        print("Usage: python client.py <size_in_mb> <host> <port>")
        return
    size_in_mb = int(sys.argv[1])
    host = sys.argv[2]
    port = int(sys.argv[3])
    
    total_size = size_in_mb * (1024 * 1024)
    element_size = 4
    num_elements = total_size // element_size
    rows = cols = int(num_elements ** 0.5)

    matrix = generate_matrix(rows, cols)
    vector = generate_vector(cols)

    data_hash = calculate_hash(matrix, vector)
    print("generated SHA-256 hash:")
    print(data_hash.hex())

    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as sock:
        try:
            sock.connect((host, port))
            send_data(sock, rows, cols, matrix, vector, data_hash)
            result, compute_time, data_size = recv_data(sock, rows)
        except:
            print("Error connecting to server")
            return 0

        if multiply_matrix_vector(matrix, vector) == list(result):
            print("\nClient and server results match!")
        else:
            print("\nClient and server results do NOT match!")

        with open(RESULT_FILE, 'w') as f:
            for value in result:
                f.write(f"{value} ")

        compute_time_ms = compute_time * 1000.0
        data_size_mb = data_size / (1024 * 1024)
        with open(STATISTICS_FILE, 'w') as f:
            f.write(f"Client time: {time.time() - start_time:.3f} seconds\n")
            f.write(f"Compute time: {compute_time_ms:.3f} ms\n")
            f.write(f"Processed data size: {data_size_mb:.3f} mb\n")

    print(f"\n{RESULT_FILE}:")
    with open(RESULT_FILE, 'r') as f:
        print(f'{len(f.read().split())} elements saved')

    print(f"\n{STATISTICS_FILE}:")
    with open(STATISTICS_FILE, 'r') as f:
        for line in f.readlines():
            print(line, end='')

if __name__ == '__main__':
    main()
