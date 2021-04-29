#!/usr/bin/python

# A hacky little toy to print things in a fun way.

import argparse
import getpass
import os
import random
import shlex
import sys
import termios
import time
import tty


def get_color(color_name):
    cm = {
        "red": "\033[1;31;40m",
        "blue": "\033[1;34;40m",
        "green": "\033[1;32;40m",
        "yellow": "\033[1;33;40m",
        "white": "\033[1;37;40m",
    }
    return cm.get(color_name, None)


DOLLAH_WAIT = True
DOLLAH_FAST = False
DOLLAH_PROMPT = "[demo]$ "
DOLLAH_COLOR = None
DOLLAH_DISPLAY = False

if os.environ.get("DOLLAH_WAIT") == "no":
    DOLLAH_WAIT = False

if os.environ.get("DOLLAH_FAST") == "yes":
    DOLLAH_FAST = True

if os.environ.get("DOLLAH_PROMPT", "") != "":
    DOLLAH_PROMPT = os.environ.get("DOLLAH_PROMPT", "")

if os.environ.get("DOLLAH_COLOR", "") != "":
    DOLLAH_COLOR = get_color(os.environ.get("DOLLAH_COLOR"))

if os.environ.get("DOLLAH_DISPLAY") == "yes":
    DOLLAH_DISPLAY = True


def colorize(dest=sys.stdout):
    if DOLLAH_COLOR is None:
        return
    dest.write(DOLLAH_COLOR)
    dest.flush()


def normal(dest=sys.stdout):
    if DOLLAH_COLOR is None:
        return
    dest.write("\033[0;37;40m")
    dest.flush()


def prompt(dest=sys.stdout):
    colorize(dest)
    dest.write(DOLLAH_PROMPT)
    dest.flush()
    normal(dest)
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
        if random.randint(1, 10) == 1:
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
    sys.stdout.write("\n")


def display(cmd):
    colorize()
    sys.stdout.write(" ".join(cmd))
    sys.stdout.write("\n")
    sys.stdout.flush()
    normal()


def main():
    if DOLLAH_DISPLAY:
        display(sys.argv[1:])
    else:
        present(sys.argv[1:])


if __name__ == "__main__":
    main()
