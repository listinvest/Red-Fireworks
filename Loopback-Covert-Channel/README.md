 # "Covert Channel"
This is a small script/malware made for a Covert Channel course as a PoC for covert communications in a linux system. The idea is to send bash history to the loopback interface using ICMP and UDP, and another process will grab from there. 

 - I haven't implemented any kind of encryption to the data.
 - More than 2k lines in .history break it.
