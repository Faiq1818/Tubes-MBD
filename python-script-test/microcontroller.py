import socket
import struct
import time
import random
import uuid

UDP_IP = "127.0.0.1"
UDP_PORT = 9999

DEVICE_IDS = [uuid.uuid4() for _ in range(20)]

def send_seismic_data(device_id, is_trigger=False):
    timestamp = int(time.time() * 1000000)  # Microseconds

    acc_x = random.uniform(-0.5, 0.5)
    acc_y = random.uniform(-0.5, 0.5)
    acc_z = 9.8 + random.uniform(-0.1, 0.1)

    pga = max(abs(acc_x), abs(acc_y), abs(acc_z))
    sta_lta = 5.2 if is_trigger else 1.1

    # UUID = 16 bytes
    uuid_bytes = device_id.bytes

    # Format:
    # 16s = UUID bytes (16 byte)
    # q   = int64 timestamp
    # f   = float32
    # ?   = bool
    #
    # Total:
    # 16 + 8 + 4 + 4 + 4 + 4 + 4 + 1 = 45 bytes

    packer = struct.Struct('<16s q f f f f f ?')

    packed_data = packer.pack(
        uuid_bytes,
        timestamp,
        acc_x,
        acc_y,
        acc_z,
        pga,
        sta_lta,
        is_trigger
    )

    sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    sock.sendto(packed_data, (UDP_IP, UDP_PORT))

    print(
        f"Sent: UUID={device_id}, "
        f"Trigger={is_trigger}, "
        f"PGA={pga:.2f}"
    )


if __name__ == "__main__":
    print(f"Sending data to {UDP_IP}:{UDP_PORT}...")
    print(f"Total device IDs: {len(DEVICE_IDS)}")

    try:
        while True:
            device_id = random.choice(DEVICE_IDS)

            send_seismic_data(
                device_id=device_id,
                is_trigger=False
            )

            time.sleep(0.001)

    except KeyboardInterrupt:
        print("\nStopped.")
