import socket

ip_list = []
ais = socket.getaddrinfo("mynode",0,0,0,0)
for result in ais:
    ip_list.append(result[-1][0])
    ip_list = list(set(ip_list))

print(ip_list)