import os 
import time

if __name__ == '__main__':
    for i in range(100):
        os.popen("go run client.go -D=img"+str(i)+" 5454 https://stsci-opo.org/STScI-01EVSVZHYGRC7KMTEJ6DJFXE4N.png")
time.sleep(1000)