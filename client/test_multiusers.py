import subprocess
import time

if __name__ == '__main__':
    processes = []
    for i in range(100):
        processes.append(subprocess.Popen(["go", "run", "client.go", "-D=img"+str(i), "4312", "http://localhost/test.jpg"]))
for p in processes:
    p.wait()
