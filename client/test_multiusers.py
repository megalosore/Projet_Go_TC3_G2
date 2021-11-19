import subprocess
import time

if __name__ == '__main__':
    processes = []
    for i in range(100):
        processes.append(subprocess.Popen(["go", "run", "client.go", "-D=img"+str(i), "4312", "https://webypress.b-cdn.net/wp-content/uploads/2020/03/The-Ultimate-Coronavirus-Small-Business-Guide-WordPress-Tools.png"]))
for p in processes:
    p.wait()