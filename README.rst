Reading Password on the Command Line in Go
==========================================

Some of the code in this package is taken from
`code.google.com/p/go/crypto/ssh/terminal`_ and then has been modified
by Jochen Voss.  The original code is distributed under the following
license::

    Copyright 2011 The Go Authors. All rights reserved.
    Use of this source code is governed by a BSD-style
    license that can be found in the LICENSE file.

All added code and all changes to the original code are distributed
under the following license::

    Copyright 2013 Jochen Voss. All rights reserved.
    Use of this source code is governed by a BSD-style
    license that can be found in the LICENSE file.

.. _code.google.com/p/go/crypto/ssh/terminal: https://code.google.com/p/go/source/browse/?repo=crypto#hg%2Fssh%2Fterminal


Installation
------------

This package can be installed using the ``go get`` command::

    go get github.com/seehuhn/password

Usage
-----

The following command can be used to read a password from the command
line::

    input, err := password.Read("passwd: ")

This switches of echoing of input to the terminal, prints the given
prompt to the screen, reads input from standard input until the end of
line, and finally restores the original terminal settings.  The byte
slice returned does not include the terminating end-of-line character.
