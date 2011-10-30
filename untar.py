#!/usr/bin/python
# -----------------------------------------------------------------
# untar.py -- Untar tarballs while protecting against tarbombs.
# Copyright 2011 Michael Kelly (michael@michaelkelly.org)
#
# This program is released under the terms of the GNU General Public
# License as published by the Free Software Foundation; either version 2
# of the License, or (at your option) any later version.
#
# Sat Oct 29 22:01:23 EDT 2011
# -----------------------------------------------------------------

import os
import re
import subprocess
import sys

def tar_list(tar_file):
  """Returns a list of all files in the given tar file."""
  proc = subprocess.Popen(['/bin/tar', '-tf', tar_file],
                          stdout=subprocess.PIPE, stderr=subprocess.PIPE)
  stdout, stderr = proc.communicate()
  if proc.returncode:
    return None
  else:
    return stdout.strip().split('\n')


def base_directories(paths):
  """Given a list of paths, returns the first directory component (or the
  filename, if there is none) of each.
  """
  dirs = set()
  for path in paths:
    # We need to account for absolute paths
    if path.startswith(os.sep):
      dirs.add(os.sep + path.split(os.sep, 2)[1])
    else:
      dirs.add(path.split(os.sep, 1)[0])
  return dirs


def archive_name(archive_path):
  """Make a guess at the archive name based on its name without a suffix."""
  base = os.path.basename(archive_path)
  return re.sub(r'\.(tar|tar\.\w+|tgz|tbz2)$', '', base, 1)


def untar(tar_file, extra_args, directory):
  """Untars a file.

  Automatically 

  Args:
    extra_args: [str] Extra args to add, before the filename.
    directory: (str) Untar into the given dirctory. If None, use the
               CWD.
  """
  cmd = ['/bin/tar'] + extra_args + ['-xf', tar_file]
  if directory is not None:
    cmd += ['-C', directory]
  proc = subprocess.Popen(cmd)
  stdout, stderr = proc.communicate()
  return proc.returncode == 0


def usage():
  usage_str = ("Usage: %s [FLAGS] TARFILE\n\n"
               "Untars TARFILE to a subdirectory of the CWD. If the TARFILE\n"
               "will naturally expand only to a single subdirecory, that one\n"
               "is used. Otherwise, the name of the directory without a\n"
               "suffix is used.\n\n"
               "Any FLAGS are passed straight to tar." % sys.argv[0])
  print >>sys.stderr, usage_str


def main(argv):
  if len(argv) < 2:
    usage()
    return 2
  tar_file = argv[-1]
  flags = argv[1:-1]
  print "flags = %s" % flags

  files = tar_list(tar_file)
  if files is None:
    print >>sys.stderr, (
        'Could not parse tar file listing for %s. Aborting.' % tar_file)
  if len(base_directories(files)) > 0:
    base = archive_name(tar_file)
    os.mkdir(base)
  else:
    base = None
  untar(tar_file, extra_args=flags, directory=base)

sys.exit(main(sys.argv))
