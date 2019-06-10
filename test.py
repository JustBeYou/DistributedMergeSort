from pwn import *
context.log_level = 'error'

def parseTime(t):
    if t[-2:] == 'ms':
        t = t.replace('ms', '')
        return float(t)

    t = t.replace('s', '')
    return float(t) * 1000

def getTimes(CPUs):
    cmd = "go run main.go {}".format(CPUs).split(' ')
    p = process(cmd)
    p.recvuntil(':')
    t1 = p.recvline().strip()
    p.recvuntil(':')
    t2 = p.recvline().strip()
    p.close()

    t1 = parseTime(t1)
    t2 = parseTime(t2)

    return t1, t2

def calcImprovement(t1, t2):
    return (t1 - t2) / t1 * 100

best = 0
best_cpus = 0
for CPUs in range(4, 9):
    print ("Testing for {} threads".format(CPUs))
    print ("----------------------")

    total = 0
    for i in range(5):
        t1, t2 = getTimes(CPUs)
        imp = calcImprovement(t1, t2)
        total += imp
        print ("Time(1): {}\nTime(2): {}\nImprovement: {}%".format(t1, t2, imp))
    mean = total / 5
    print ("\nAverage {}%".format(mean))
    print ("----------------------\n")

    if mean > best:
        best = mean
        best_cpus = CPUs

print ("\n")
print ("Best score is {}% using {} threads".format(best, best_cpus))
