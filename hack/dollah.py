#!/usr/bin/python

import argparse
import shlex
import sys
import random
import time
import getpass
import tty
import termios
import os

DOLLAH_WAIT = True
DOLLAH_FAST = False
DOLLAH_PROMPT = '[demo]$ '

if os.environ.get("DOLLAH_WAIT") == "no":
    DOLLAH_WAIT = False

if os.environ.get("DOLLAH_FAST") == "yes":
    DOLLAH_FAST = True

if os.environ.get("DOLLAH_PROMPT", "") != "":
    DOLLAH_PROMPT = os.environ.get("DOLLAH_PROMPT", "")


def prompt(dest=sys.stdout):
    dest.write(DOLLAH_PROMPT)
    dest.flush()
    wait()


def wait():
    if not DOLLAH_WAIT:
        return
    fd = sys.stdin.fileno()
    before = termios.tcgetattr(fd)
    tty.setraw(fd)
    sys.stdin.read(1)
    termios.tcsetattr(fd, termios.TCSADRAIN, before)


def fast_sleep(seconds):
    time.sleep(seconds / 4)


def typewrite(s, dest=sys.stdout):
    if DOLLAH_FAST:
        _sleep = fast_sleep
    else:
        _sleep = time.sleep
    lslp = 0
    for c in s:
        _sleep(0.015 * random.randint(1, 10))
        if lslp > 5:
            _sleep(0.002 * random.randint(2, 20))
        if random.randint(1,10) == 1:
            _sleep(0.1)
        if c == " ":
            _sleep(0.2)
            lslp = 0
        else:
            lslp += 1
        dest.write(c)
        dest.flush()
        if c == " ":
            _sleep(0.03 * random.randint(1, 10))


def read(src):
    out = []
    for line in src:
        p = shlex.split(line, posix=True)
        cmd = list(p)
        if cmd:
            out.append(cmd)
    return out


def present(cmd):
    prompt()
    typewrite(shlex.join(cmd))
    sys.stdout.write('\n')

def main():
    #parser = argparse.ArgumentParser()
    #parser.add_argument('source')
    #cli = parser.parse_args()


    #with open(cli.source) as fh:
    #    cmds = read(fh)
    #for cmd in cmds:
    #    present(cmd)
    present(sys.argv[1:])


if __name__ == '__main__':
    main()
