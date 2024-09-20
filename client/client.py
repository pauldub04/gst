import socket
import random
import struct
import sys
import hashlib

MIN_INT = -1000
MAX_INT = +1000

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

def main():
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
    print("SHA-256 Hash:", data_hash.hex())

    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as sock:
        sock.connect((host, port))

        data = [rows, cols] + [item for row in matrix for item in row] + vector
        packed_data = struct.pack(f'<ii{cols * rows}i{cols}i', *data)
        sock.sendall(packed_data)

        sock.sendall(data_hash)

        # result = struct.unpack(f'!{rows}I', sock.recv(rows * 4))
        # exec_time = struct.unpack('!d', sock.recv(8))[0]

        # print(f'Result: {result}')
        # print(f'Execution time: {exec_time:.6f} seconds')

if __name__ == '__main__':
    main()
