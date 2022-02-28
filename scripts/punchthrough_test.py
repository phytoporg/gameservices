# Just for some basic E2E testing and iteration

import socket

ip = '127.0.0.1'
port = 9999

sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)

# MSG = hello 
# 0x01 + client name

msg = bytearray([0x01])
msg += bytearray("philjo5000".encode('utf-8'))
sock.sendto(msg, (ip, port))

data, addr = sock.recvfrom(512)
client_id = int.from_bytes(data, 'little')
print(f'Client ID: {client_id}')

# MSG = host room
# 0x03 + room name

msg = bytearray([0x03])
msg += client_id.to_bytes(4, 'little')
msg += bytearray("cool room for excellent people".encode('utf-8'))
print(msg)
sock.sendto(msg, (ip, port))

data, addr = sock.recvfrom(512)
room_id = int.from_bytes(data, 'little')
print(f'Hosting room ID: {room_id}')

# MSG = hello 
# 0x01 + client name

msg = bytearray([0x01])
msg += bytearray("philjo5001".encode('utf-8'))
sock.sendto(msg, (ip, port))

data, addr = sock.recvfrom(512)
second_client_id = int.from_bytes(data, 'little')
print(f'Second client ID: {second_client_id}')

# MSG = join room 
# 0x04 + client id + room id

msg = bytearray([0x04])
msg += second_client_id.to_bytes(4, 'little')
msg += room_id.to_bytes(4, 'little')
sock.sendto(msg, (ip, port))

print(f'Joined room!')
