import os 
import time
from multiprocessing import Process

def launch(i):
    os.system("client.exe -D=img"+str(i)+" 5454 https://stsci-opo.org/STScI-01EVSVZHYGRC7KMTEJ6DJFXE4N.png")

if __name__ == '__main__':
    for i in range(100):
        p = Process(target=launch, args=(i,))
        p.start()
    p.join()