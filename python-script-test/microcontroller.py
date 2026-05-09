import socket
import struct
import time
import random

UDP_IP = "127.0.0.1"
UDP_PORT = 9999

def send_seismic_data(device_id, is_trigger=False):
    timestamp = int(time.time() * 1000000) # Microseconds
    acc_x = random.uniform(-0.5, 0.5)
    acc_y = random.uniform(-0.5, 0.5)
    acc_z = 9.8 + random.uniform(-0.1, 0.1) # Gravity + noise
    pga = max(abs(acc_x), abs(acc_y), abs(acc_z))
    sta_lta = 5.2 if is_trigger else 1.1

    # Format Packing (Little Endian):
    # Q = uint64 (DeviceID)
    # q = int64 (Timestamp)
    # f = float32 (AccX, AccY, AccZ, PGA, STALTA)
    # ? = bool (IsTrigger)
    # Total: 8+8+4+4+4+4+4+1 = 37 bytes
    packer = struct.Struct('<Q q f f f f f ?')
    packed_data = packer.pack(
        device_id, 
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
    
    print(f"Sent: ID={device_id}, Trigger={is_trigger}, PGA={pga:.2f}")

if __name__ == "__main__":
    print(f"Sending data to {UDP_IP}:{UDP_PORT}...")
    try:
        while True:
            send_seismic_data(device_id=101, is_trigger=False)
            time.sleep(0.001)
    except KeyboardInterrupt:
        print("\nStopped.")
