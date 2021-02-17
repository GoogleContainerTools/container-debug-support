#!/usr/bin/env python
#
# Copyright 2020 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# A `skaffold debug` launcher for Python.
#
# Supporting Python introduces some quirks as there are now three
# methods for hooking up a debugging backend (pydevd, ptvsd, and now
# debugpy, where ptvsd and debugpy themselves embed pydevd).
# Furthermore pydevd has pyx libraries which are specific to
# particular versions of Python.
#
# Further complicating matters is that a number of Python packages
# use launcher scripts (e.g., gunicorn), and so we can't simply run
# `python -m ptvsd -- gunicorn` as ptvsd/debugpy/etc don't look for
# the script file in the PATH.

import os, runpy, sys, argparse

def resolve(script):
  """
  Attempt to resolve a (python) script in the PATH and PYTHONUSERBASE's bin.
  The intention is to allow `skaffold debug` to simply prefix a command line,
  like `gunicorn ...`, with `python /dbg/debug.py ... --`.  These scripts,
  like `gunicorn`, are not normally resolved through Python's normal import
  mechanisms.

  Return the resolved script if found, or the provided name so that it
  can be used through the python import mechanisms.
  """
  if os.path.exists(script):
    return script

  for p in os.get_exec_path():
    if os.path.exists(os.path.join(p, script)):
      return os.path.join(p, script)

  userbase = os.getenv("PYTHONUSERBASE")
  if userbase is None:
    userbase = os.path.expanduser(os.path.join("~", ".local", "bin"))
  if os.path.exists(os.path.join(userbase, script)):
    print("found " + script + ": " + os.path.join(userbase, script))
    return os.path.join(userbase, script)

  return script

def updateSysPath():
  """
  Alter sys.path to add the appropriate skaffold-debug-support location.
  We append to sys.path in case the user installed a better version of
  our support libraries.
  """
  dbglib = "/dbg/lib/python{major}.{minor}/site-packages".format(major=sys.version_info.major, minor=sys.version_info.minor)
  if os.path.exists(dbglib):
    sys.path.append(dbglib)


if __name__ == '__main__':
  print("sys.argv=", sys.argv)
  print("sys.path=", sys.path)
  print("os.get_exec_path()=", os.get_exec_path())

  parser = argparse.ArgumentParser(description='Skaffold `debug` launcher for Python.')
  parser.add_argument('--mode', required=True, choices=['pydevd', 'debugpy', 'ptvsd'],
                      help='specify which debug backend to be used')
  parser.add_argument('--port', default='5678',
                      help='port to await for debug connections')
  parser.add_argument('--wait', default=False, action='store_true',
                      help='pause execution until a debug connection has been established')
  parser.add_argument('cmdline', metavar='...', nargs=argparse.REMAINDER,
                      help='application module or script with arguments to be debugged')
  args = parser.parse_args()
 
  if len(args.cmdline) == 0:
    parser.print_help()
    sys.exit(1)

  script = resolve(args.cmdline[0])

  print("sys.argv=", sys.argv)

  if args.mode == 'debugpy':
    cmdline = ['debugpy', '--listen', args.port]
    if args.wait:
      cmdline.append('--wait-for-client')

    if os.path.exists(script):
      cmdline.append(script)
    else:
      cmdline.append('-m')
      cmdline.append(script)
    cmdline.extend(args.cmdline[1:])

    sys.argv = cmdline
    print('debugpy mode: sys.argv=', sys.argv)
    runpy.run_module("debugpy", run_name="__main__")
    sys.exit(0)

  if args.mode == 'ptvsd':
    cmdline = ['ptvsd', '--host', 'localhost', '--port', args.port]
    if args.wait:
      cmdline.append('--wait')

    if os.path.exists(script):
      cmdline.append(script)
    else:
      cmdline.append('-m')
      cmdline.append(script)
    cmdline.extend(args.cmdline[1:])

    sys.argv = cmdline
    print('ptvsd mode: sys.argv=', sys.argv)
    runpy.run_module("ptvsd", run_name="__main__")
    sys.exit(0)

  # pydevd
  cmdline = ['--server', '--port', args.port]

  if os.path.exists(script):
    cmdline.append('--file')
    cmdline.append(script)
  else:
    cmdline.append('--module')
    cmdline.append('--file')
    cmdline.append(script)
  cmdline.extend(args.cmdline[1:])

  sys.argv = cmdline
  print('pydevd mode: sys.argv=', sys.argv)

  import pydevd
  pydevd.main()
  if args.wait:
    pydevd._wait_for_attach()
